package orm

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/mobiledgex/edge-cloud-infra/mc/ormapi"
	"github.com/mobiledgex/edge-cloud-infra/mc/ormclient"
	"github.com/mobiledgex/edge-cloud/log"
	"github.com/stretchr/testify/require"
)

type entry struct {
	Org     string            // Organization/Developer
	OrgType string            // Organization Type
	Users   map[string]string // User:UserType
}

var (
	testEntries []entry = []entry{
		entry{
			Org:     "bigorg1",
			OrgType: OrgTypeDeveloper,
			Users: map[string]string{
				"orgman1":   RoleDeveloperManager,
				"worker1":   RoleDeveloperContributor,
				"worKer1.1": RoleDeveloperViewer,
			},
		},
		entry{
			Org:     "bigOrg2",
			OrgType: OrgTypeDeveloper,
			Users: map[string]string{
				"orgMan2":   RoleDeveloperManager,
				"worker2":   RoleDeveloperContributor,
				"wOrKer2.1": RoleDeveloperViewer,
			},
		},
		entry{
			Org:     "operatorOrg",
			OrgType: OrgTypeOperator,
			Users: map[string]string{
				"oper1": RoleOperatorManager,
			},
		},
	}

	// Extra entries only present in Artifactory/Gitlab but not in MC
	extraEntries []entry = []entry{
		entry{
			Org:     "extraOrg1",
			OrgType: OrgTypeDeveloper,
			Users: map[string]string{
				"extraUser1":   RoleDeveloperManager,
				"extraWorker1": RoleDeveloperContributor,
			},
		},
	}

	// Missing entries only present in MC but not in Artifactory/Gitlab
	missingEntries []entry = []entry{
		entry{
			Org:     "missingOrg1",
			OrgType: OrgTypeDeveloper,
			Users: map[string]string{
				"missingUser1": RoleDeveloperManager,
				"missingUser2": RoleDeveloperViewer,
			},
		},
		entry{
			Org:     "missingOperOrg",
			OrgType: OrgTypeOperator,
			Users: map[string]string{
				"missingOperUser1": RoleOperatorManager,
				"missingOperUser2": RoleOperatorViewer,
			},
		},
	}

	// create operator entries present in both MC and Artifactory/Gitlab,
	// should be removed by sync thread.
	operatorEntries []entry = []entry{
		entry{
			Org:     "oldOperOrg",
			OrgType: OrgTypeOperator,
		},
		entry{
			Org:     "oldOperOrg2",
			OrgType: OrgTypeOperator,
		},
	}
)

const (
	ExtraObj   string = "extra"
	MCObj      string = "mc"
	OldOperObj string = "oldoper"
)

func TestAppStoreApi(t *testing.T) {
	artifactoryAddr := "https://dummy-artifactory.mobiledgex.net"
	artifactoryApiKey := "dummyKey"

	gitlabAddr := "https://dummy-gitlab.mobiledgex.net"
	gitlabApiKey := "dummyKey"

	var status int

	log.SetDebugLevel(log.DebugLevelApi)
	log.InitTracer("")
	defer log.FinishTracer()

	ctx := log.StartTestSpan(context.Background())

	os.Setenv("gitlab_token", gitlabApiKey)

	// mock http to redirect requests
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	// any requests that don't have a registered URL will be fetched normally
	httpmock.RegisterNoResponder(httpmock.InitialTransport.RoundTrip)

	// master controller
	addr := "127.0.0.1:9999"
	uri := "http://" + addr + "/api/v1"

	config := ServerConfig{
		ServAddr:        addr,
		SqlAddr:         "127.0.0.1:5445",
		RunLocal:        true,
		InitLocal:       true,
		IgnoreEnv:       true,
		ArtifactoryAddr: artifactoryAddr,
		GitlabAddr:      gitlabAddr,
		SkipVerifyEmail: true,
		LocalVault:      true,
	}

	server, err := RunServer(&config)
	require.Nil(t, err, "run server")
	defer server.Stop()

	os.Setenv("VAULT_ROLE_ID", roleID)
	os.Setenv("VAULT_SECRET_ID", secretID)

	rtfuri, err := url.ParseRequestURI(artifactoryAddr)
	require.Nil(t, err, "parse artifactory url")

	path := "secret/registry/" + rtfuri.Host
	server.vault.Run("vault", fmt.Sprintf("kv put %s apikey=%s", path, artifactoryApiKey), &err)
	require.Nil(t, err, "added secret to vault")

	err = server.WaitUntilReady()
	require.Nil(t, err, "server online")

	mcClient := &ormclient.Client{}

	// login as super user
	tokenAdmin, err := mcClient.DoLogin(uri, DefaultSuperuser, DefaultSuperpass)
	require.Nil(t, err, "login as superuser")

	// Before Artifactory/Gitlab are hooked into mock, create "missing" data
	// so it doesn't automatically get populated into Art/Gitlab
	for _, v := range missingEntries {
		mcClientCreate(t, v, mcClient, uri)
	}

	// mock artifactory
	rtf := NewArtifactoryMock(artifactoryAddr)
	// mock gitlab
	gm := NewGitlabMock(gitlabAddr)

	// Create new users & orgs from MC
	for _, v := range testEntries {
		mcClientCreate(t, v, mcClient, uri)
		rtf.verify(t, v, MCObj)
		gm.verify(t, v, MCObj)
	}

	// Create users & orgs which are not present in MC
	for _, v := range extraEntries {
		org := ormapi.Organization{
			Name: v.Org,
			Type: v.OrgType,
		}
		artifactoryCreateGroupObjects(ctx, v.Org, v.OrgType)
		gitlabCreateGroup(ctx, &org)
		for user, userType := range v.Users {
			userObj := ormapi.User{
				Name: user,
			}
			artifactoryCreateUser(ctx, &userObj, nil, nil)
			gitlabCreateLDAPUser(ctx, &userObj)

			roleArg := ormapi.Role{
				Username: user,
				Org:      v.Org,
				Role:     userType,
			}
			gitlabAddGroupMember(ctx, &roleArg, org.Type)
			artifactoryCreateGroupObjects(ctx, org.Name, org.Type)
			artifactoryAddUserToGroup(ctx, &roleArg, org.Type)
		}
		rtf.verify(t, v, ExtraObj)
		gm.verify(t, v, ExtraObj)
	}

	// Create operator entries in MC and then force populate them
	// in artifactory/gitlab to test that sync will remove them.
	for _, v := range operatorEntries {
		testCreateOrg(t, mcClient, uri, tokenAdmin, v.OrgType, v.Org)
		// leave Type empty so artifactory/gitlab funcs will push it
		org := ormapi.Organization{
			Name: v.Org,
		}
		// verify not in artifactory/gitlab
		rtf.verify(t, v, MCObj)
		gm.verify(t, v, MCObj)

		artifactoryCreateGroupObjects(ctx, org.Name, org.Type)
		gitlabCreateGroup(ctx, &org)

		// verify now in artifactory/gitlab
		rtf.verify(t, v, OldOperObj)
		gm.verify(t, v, OldOperObj)
	}

	// Trigger resync to delete extra objects and create missing ones
	status, err = mcClient.ArtifactoryResync(uri, tokenAdmin)
	require.Nil(t, err, "artifactory resync")
	require.Equal(t, http.StatusOK, status, "artifactory resync status")
	status, err = mcClient.GitlabResync(uri, tokenAdmin)
	require.Nil(t, err, "gitlab resync")
	require.Equal(t, http.StatusOK, status, "gitlab resync status")

	// Delete MC created Objects
	for _, v := range testEntries {
		mcClientDelete(t, v, mcClient, uri, tokenAdmin)
	}

	waitSyncCount(t, gitlabSync, 2)
	waitSyncCount(t, artifactorySync, 2)

	// verify missing entries are there
	for _, v := range missingEntries {
		rtf.verify(t, v, MCObj)
		gm.verify(t, v, MCObj)
		// delete them
		mcClientDelete(t, v, mcClient, uri, tokenAdmin)
	}

	// By now, appstore Sync threads should delete all extra objects as well
	rtf.verifyEmpty(t)
	gm.verifyEmpty(t)
}

func mcClientCreate(t *testing.T, v entry, mcClient *ormclient.Client, uri string) {
	token := ""
	for user, userType := range v.Users {
		if userType == RoleDeveloperManager || userType == RoleOperatorManager {
			_, token = testCreateUser(t, mcClient, uri, user)
			testCreateOrg(t, mcClient, uri, token, v.OrgType, v.Org)
			break
		}
	}
	for user, userType := range v.Users {
		if userType != RoleDeveloperManager && userType != RoleOperatorManager {
			worker, _ := testCreateUser(t, mcClient, uri, user)
			testAddUserRole(t, mcClient, uri, token, v.Org, userType, worker.Name, Success)
		}
	}
}

func mcClientDelete(t *testing.T, v entry, mcClient *ormclient.Client, uri, tokenAdmin string) {
	for user, userType := range v.Users {
		if userType == RoleDeveloperManager {
			continue
		}
		roleArg := ormapi.Role{
			Username: user,
			Org:      v.Org,
			Role:     userType,
		}
		// admin user can remove role
		status, err := mcClient.RemoveUserRole(uri, tokenAdmin, &roleArg)
		require.Nil(t, err, "remove user role")
		require.Equal(t, http.StatusOK, status)
	}

	// delete org
	org := ormapi.Organization{
		Name: v.Org,
	}
	status, err := mcClient.DeleteOrg(uri, tokenAdmin, &org)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, status)

	// delete all users
	for user, _ := range v.Users {
		userObj := ormapi.User{
			Name: user,
		}
		status, err = mcClient.DeleteUser(uri, tokenAdmin, &userObj)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, status)
	}
}

func waitSyncCount(t *testing.T, sync *AppStoreSync, count int64) {
	for ii := 0; ii < 10; ii++ {
		if sync.count >= count {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	require.True(t, sync.count == count, fmt.Sprintf("sync count %d != expected %d", sync.count, count))
}