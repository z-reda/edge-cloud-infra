package openstack

import (
	"context"
	"fmt"
	"os"

	"github.com/mobiledgex/edge-cloud-infra/mexos"
	"github.com/mobiledgex/edge-cloud/cloud-resource-manager/nginx"
	"github.com/mobiledgex/edge-cloud/cloud-resource-manager/platform"
	"github.com/mobiledgex/edge-cloud/cloud-resource-manager/platform/pc"
	"github.com/mobiledgex/edge-cloud/cloudcommon"
	"github.com/mobiledgex/edge-cloud/edgeproto"
	"github.com/mobiledgex/edge-cloud/log"
	"github.com/mobiledgex/edge-cloud/vmspec"
)

const MINIMUM_DISK_SIZE uint64 = 20

type Platform struct {
	rootLBName  string
	rootLB      *mexos.MEXRootLB
	cloudletKey *edgeproto.CloudletKey
	flavorList  []*edgeproto.FlavorInfo
	config      platform.PlatformConfig
}

func (s *Platform) GetType() string {
	return "openstack"
}

func (s *Platform) Init(ctx context.Context, platformConfig *platform.PlatformConfig, updateCallback edgeproto.CacheUpdateCallback) error {
	rootLBName := cloudcommon.GetRootLBFQDN(platformConfig.CloudletKey)
	s.cloudletKey = platformConfig.CloudletKey
	s.config = *platformConfig
	log.SpanLog(ctx,
		log.DebugLevelMexos, "init openstack",
		"rootLB", rootLBName,
		"physicalName", platformConfig.PhysicalName,
		"vaultAddr", platformConfig.VaultAddr)

	updateCallback(edgeproto.UpdateTask, "Initializing Openstack platform")

	updateCallback(edgeproto.UpdateTask, "Fetching Openstack access credentials")
	if err := mexos.InitInfraCommon(ctx, platformConfig.VaultAddr); err != nil {
		return err
	}
	if err := mexos.InitOpenstackProps(ctx, platformConfig.CloudletKey.OperatorKey.Name, platformConfig.PhysicalName, platformConfig.VaultAddr); err != nil {
		return err
	}
	mexos.CloudletInfraCommon.NetworkScheme = os.Getenv("MEX_NETWORK_SCHEME")
	if mexos.CloudletInfraCommon.NetworkScheme == "" {
		mexos.CloudletInfraCommon.NetworkScheme = "name=mex-k8s-net-1,cidr=10.101.X.0/24"
	}
	var err error
	s.flavorList, err = mexos.GetFlavorInfo(ctx)
	if err != nil {
		return err
	}

	// create rootLB
	updateCallback(edgeproto.UpdateTask, "Creating RootLB")
	crmRootLB, cerr := mexos.NewRootLB(ctx, rootLBName)
	if cerr != nil {
		return cerr
	}
	if crmRootLB == nil {
		return fmt.Errorf("rootLB is not initialized")
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "created rootLB", "rootlb", crmRootLB.Name)
	s.rootLB = crmRootLB
	s.rootLBName = rootLBName

	var sharedRootLBFlavor edgeproto.Flavor
	err = mexos.GetCloudletSharedRootLBFlavor(&sharedRootLBFlavor)
	if err != nil {
		return fmt.Errorf("unable to get Shared RootLB Flavor: %v", err)
	}
	vmspec, err := vmspec.GetVMSpec(s.flavorList, sharedRootLBFlavor)
	if err != nil {
		return fmt.Errorf("unable to find VM spec for Shared RootLB: %v", err)
	}

	log.SpanLog(ctx, log.DebugLevelMexos, "calling SetupRootLB")
	updateCallback(edgeproto.UpdateTask, "Setting up RootLB")
	err = mexos.SetupRootLB(ctx, rootLBName, vmspec, edgeproto.DummyUpdateCallback)
	if err != nil {
		return err
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "ok, SetupRootLB")

	// set up L7 load balancer
	client, err := s.GetPlatformClientRootLB(ctx, rootLBName)
	if err != nil {
		return err
	}
	updateCallback(edgeproto.UpdateTask, "Setting up Nginx Proxy")
	err = nginx.InitL7Proxy(client, nginx.WithDockerNetwork("host"))
	if err != nil {
		return err
	}
	return nil
}

func (s *Platform) GatherCloudletInfo(ctx context.Context, info *edgeproto.CloudletInfo) error {
	return mexos.OSGetLimits(ctx, info)
}

func (s *Platform) GetPlatformClientRootLB(ctx context.Context, rootLBName string) (pc.PlatformClient, error) {
	log.SpanLog(ctx, log.DebugLevelMexos, "GetPlatformClientRootLB", "rootLBName", rootLBName)

	if rootLBName == "" {
		return nil, fmt.Errorf("cannot GetPlatformClientRootLB, rootLB is empty")
	}
	if mexos.GetCloudletExternalNetwork() == "" {
		return nil, fmt.Errorf("GetPlatformClientRootLB, missing external network in platform config")
	}
	return mexos.GetSSHClient(ctx, rootLBName, mexos.GetCloudletExternalNetwork(), mexos.SSHUser)
}

func (s *Platform) GetPlatformClient(ctx context.Context, clusterInst *edgeproto.ClusterInst) (pc.PlatformClient, error) {
	rootLBName := s.rootLBName
	if clusterInst.IpAccess == edgeproto.IpAccess_IP_ACCESS_DEDICATED {
		rootLBName = cloudcommon.GetDedicatedLBFQDN(s.cloudletKey, &clusterInst.Key.ClusterKey)
	}
	return s.GetPlatformClientRootLB(ctx, rootLBName)
}
