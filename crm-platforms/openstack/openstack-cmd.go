package openstack

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	sh "github.com/codeskyblue/go-sh"
	"github.com/mobiledgex/edge-cloud-infra/mexos"
	"github.com/mobiledgex/edge-cloud/cloudcommon"
	"github.com/mobiledgex/edge-cloud/edgeproto"
	"github.com/mobiledgex/edge-cloud/log"
)

func (s *Platform) TimedOpenStackCommand(ctx context.Context, name string, a ...string) ([]byte, error) {
	parmstr := ""
	start := time.Now()
	for _, a := range a {
		parmstr += a + " "
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "OpenStack Command Start", "name", name, "parms", parmstr)
	newSh := sh.NewSession()
	for key, val := range s.openRCVars {
		newSh.SetEnv(key, val)
	}

	out, err := newSh.Command(name, a).CombinedOutput()
	if err != nil {
		log.InfoLog("Openstack command returned error", "parms", parmstr, "err", err, "out", string(out), "elapsed time", time.Since(start))
		return out, err
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "OpenStack Command Done", "parmstr", parmstr, "elapsed time", time.Since(start))
	return out, nil

}

//ListServers returns list of servers, KVM instances, running on the system
func (s *Platform) ListServers(ctx context.Context) ([]OSServer, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "server", "list", "-f", "json")

	if err != nil {
		err = fmt.Errorf("cannot get server list, %s, %v", out, err)
		return nil, err
	}
	var servers []OSServer
	err = json.Unmarshal(out, &servers)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return nil, err
	}
	return servers, nil
}

//ListServers returns list of servers, KVM instances, running on the system
func (s *Platform) ListPorts(ctx context.Context) ([]OSPort, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "port", "list", "-f", "json")

	if err != nil {
		err = fmt.Errorf("cannot get port list, %s, %v", out, err)
		return nil, err
	}
	var ports []OSPort
	err = json.Unmarshal(out, &ports)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return nil, err
	}
	return ports, nil
}

//ListPortsServerNetwork returns ports for a particular server on a given network
func (s *Platform) ListPortsServerNetwork(ctx context.Context, server, network string) ([]OSPort, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "port", "list", "--server", server, "--network", network, "-f", "json")
	if err != nil {
		err = fmt.Errorf("cannot get port list, %s, %v", out, err)
		return nil, err
	}
	var ports []OSPort
	err = json.Unmarshal(out, &ports)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return nil, err
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "list ports", "server", server, "network", network, "ports", ports)
	return ports, nil
}

//ListImages lists avilable images in glance
func (s *Platform) ListImages(ctx context.Context) ([]OSImage, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "image", "list", "-f", "json")
	if err != nil {
		err = fmt.Errorf("cannot get image list, %s, %v", out, err)
		return nil, err
	}
	var images []OSImage
	err = json.Unmarshal(out, &images)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return nil, err
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "list images", "images", images)
	return images, nil
}

//GetImageDetail show of a given image from Glance
func (s *Platform) GetImageDetail(ctx context.Context, name string) (*OSImageDetail, error) {
	out, err := s.TimedOpenStackCommand(
		ctx, "openstack", "image", "show", name, "-f", "json",
		"-c", "id",
		"-c", "status",
		"-c", "updated_at",
		"-c", "checksum",
		"-c", "disk_format",
	)
	if err != nil {
		err = fmt.Errorf("cannot get image Detail for %s, %s, %v", name, string(out), err)
		return nil, err
	}
	var imageDetail OSImageDetail
	err = json.Unmarshal(out, &imageDetail)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return nil, err
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "show image Detail", "Detail", imageDetail)
	return &imageDetail, nil
}

// fetch tags + properties etc of all images for resource mapping
func (s *Platform) ListImagesDetail(ctx context.Context) ([]OSImageDetail, error) {
	var img_details []OSImageDetail
	images, err := s.ListImages(ctx)
	if err != nil {
		return nil, err
	}

	for _, image := range images {
		details, err := s.GetImageDetail(ctx, image.Name)
		if err == nil {
			img_details = append(img_details, *details)
		}
	}
	return img_details, err
}

//
//ListNetworks lists networks known to the platform. Some created by the operator, some by users.
func (s *Platform) ListNetworks(ctx context.Context) ([]OSNetwork, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "network", "list", "-f", "json")
	if err != nil {
		log.SpanLog(ctx, log.DebugLevelMexos, "network list failed", "out", out)
		err = fmt.Errorf("cannot get network list, %s, %v", out, err)
		return nil, err
	}
	var networks []OSNetwork
	err = json.Unmarshal(out, &networks)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return nil, err
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "list networks", "networks", networks)
	return networks, nil
}

//ShowFlavor returns the details of a given flavor.
func (s *Platform) ShowFlavor(ctx context.Context, flavor string) (details OSFlavorDetail, err error) {

	var flav OSFlavorDetail
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "flavor", "show", flavor, "-f", "json")
	if err != nil {
		log.SpanLog(ctx, log.DebugLevelMexos, "flavor show failed", "out", out)
		return flav, err
	}

	err = json.Unmarshal(out, &flav)
	if err != nil {
		return flav, err
	}
	return flav, nil
}

//ListFlavors lists flavors known to the platform.   The ones matching the flavorMatchPattern are returned
func (s *Platform) ListFlavors(ctx context.Context) ([]OSFlavorDetail, error) {
	flavorMatchPattern := s.GetCloudletFlavorMatchPattern()
	r, err := regexp.Compile(flavorMatchPattern)
	if err != nil {
		return nil, fmt.Errorf("Cannot compile flavor match pattern")
	}
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "flavor", "list", "--long", "-f", "json")
	if err != nil {
		err = fmt.Errorf("cannot get flavor list, %s, %v", out, err)
		return nil, err
	}
	var flavors []OSFlavorDetail
	var flavorsMatched []OSFlavorDetail

	err = json.Unmarshal(out, &flavors)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return nil, err
	}
	for _, f := range flavors {
		if r.MatchString(f.Name) {
			flavorsMatched = append(flavorsMatched, f)
		}
	}
	return flavorsMatched, nil
}

func (s *Platform) ListAZones(ctx context.Context) ([]OSAZone, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "availability", "zone", "list", "-f", "json")
	if err != nil {
		err = fmt.Errorf("cannot get availability zone list, %s, %v", out, err)
		return nil, err
	}
	var zones []OSAZone
	err = json.Unmarshal(out, &zones)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return nil, err
	}
	return zones, nil
}

func (s *Platform) ListFloatingIPs(ctx context.Context) ([]OSFloatingIP, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "floating", "ip", "list", "-f", "json")
	if err != nil {
		err = fmt.Errorf("cannot get floating ip list, %s, %v", out, err)
		return nil, err
	}
	var fips []OSFloatingIP
	err = json.Unmarshal(out, &fips)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return nil, err
	}
	return fips, nil
}

//CreateServer instantiates a new server instance, which is a KVM instance based on a qcow2 image from glance
func (s *Platform) CreateServer(ctx context.Context, opts *OSServerOpt) error {
	args := []string{
		"server", "create",
		"--config-drive", "true", //XXX always
		"--image", opts.Image, "--flavor", opts.Flavor,
	}
	if opts.UserData != "" {
		args = append(args, "--user-data", opts.UserData)
	}
	for _, p := range opts.Properties {
		args = append(args, "--property", p)
		// `p` should be like: "key=value"
	}
	for _, n := range opts.NetIDs {
		args = append(args, "--nic", "net-id="+n)
		// `n` should be like: "public,v4-fixed-ip=172.24.4.201"
	}
	args = append(args, opts.Name)
	//TODO additional args
	iargs := make([]string, len(args))
	for i, v := range args {
		iargs[i] = v
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "creating server with args", "iargs", iargs)

	//log.SpanLog(ctx,log.DebugLevelMexos, "openstack create server", "opts", opts, "iargs", iargs)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", iargs...)
	if err != nil {
		err = fmt.Errorf("cannot create server, %v, '%s'", err, out)
		return err
	}
	return nil
}

// GetActiveServerDetails returns details of the KVM instance waiting for it to be ACTIVE
func (s *Platform) GetActiveServerDetails(ctx context.Context, name string) (*OSServerDetail, error) {
	active := false
	srvDetail := &OSServerDetail{}
	for i := 0; i < 10; i++ {
		out, err := s.TimedOpenStackCommand(ctx, "openstack", "server", "show", "-f", "json", name)
		if err != nil {
			err = fmt.Errorf("can't show server %s, %s, %v", name, out, err)
			return nil, err
		}
		//fmt.Printf("%s\n", out)
		err = json.Unmarshal(out, srvDetail)
		if err != nil {
			err = fmt.Errorf("cannot unmarshal while getting server detail, %v", err)
			return nil, err
		}
		if srvDetail.Status == "ACTIVE" {
			active = true
			break
		}
		log.SpanLog(ctx, log.DebugLevelMexos, "wait for server to become ACTIVE", "server detail", srvDetail)
		time.Sleep(30 * time.Second)
	}
	if !active {
		return nil, fmt.Errorf("while getting server detail, waited but server %s is too slow getting to active state", name)
	}
	//log.SpanLog(ctx,log.DebugLevelMexos, "server detail", "server detail", srvDetail)
	return srvDetail, nil
}

// GetServerDetails returns details of the KVM instance
func (s *Platform) GetServerDetails(ctx context.Context, name string) (*OSServerDetail, error) {
	srvDetail := &OSServerDetail{}
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "server", "show", "-f", "json", name)
	if err != nil {
		err = fmt.Errorf("can't show server %s, %s, %v", name, out, err)
		return nil, err
	}
	//fmt.Printf("%s\n", out)
	err = json.Unmarshal(out, srvDetail)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal while getting server detail, %v", err)
		return nil, err
	}
	return srvDetail, nil
}

// GetPortDetails gets details of the specified port
func (s *Platform) GetPortDetails(ctx context.Context, name string) (*OSPortDetail, error) {
	log.SpanLog(ctx, log.DebugLevelMexos, "get port details", "name", name)
	portDetail := &OSPortDetail{}

	out, err := s.TimedOpenStackCommand(ctx, "openstack", "port", "show", name, "-f", "json")
	if err != nil {
		err = fmt.Errorf("can't get port detail for port: %s, %s, %v", name, out, err)
		return nil, err
	}
	err = json.Unmarshal(out, &portDetail)
	if err != nil {
		log.SpanLog(ctx, log.DebugLevelMexos, "port unmarshal failed", "err", err)
		err = fmt.Errorf("can't unmarshal port, %v", err)
		return nil, err
	}
	return portDetail, nil
}

// AttachPortToServer attaches a port to a server
func (s *Platform) AttachPortToServer(ctx context.Context, serverName, portName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "AttachPortToServer", "serverName", serverName, "portName", portName)

	out, err := s.TimedOpenStackCommand(ctx, "openstack", "server", "add", "port", serverName, portName)
	if err != nil {
		if strings.Contains(string(out), "still in use") {
			// port already attached
			log.SpanLog(ctx, log.DebugLevelMexos, "port already attached", "serverName", serverName, "portName", portName, "out", out, "err", err)
			return nil
		}
		log.SpanLog(ctx, log.DebugLevelMexos, "can't attach port", "serverName", serverName, "portName", portName, "out", out, "err", err)
		err = fmt.Errorf("can't attach port: %s, %s, %v", portName, out, err)
		return err
	}
	return nil
}

// DetachPortFromServer removes a port from a server
func (s *Platform) DetachPortFromServer(ctx context.Context, serverName, portName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "DetachPortFromServer", "serverName", serverName, "portName", portName)

	out, err := s.TimedOpenStackCommand(ctx, "openstack", "server", "remove", "port", serverName, portName)
	if err != nil {
		log.SpanLog(ctx, log.DebugLevelMexos, "can't remove port", "serverName", serverName, "portName", portName, "out", out, "err", err)
		if strings.Contains(string(out), "No Port found") {
			// when ports are removed they are detached from any server they are connected to.
			log.SpanLog(ctx, log.DebugLevelMexos, "port is gone", "portName", portName)
			err = nil
		} else {
			log.SpanLog(ctx, log.DebugLevelMexos, "can't remove port", "serverName", serverName, "portName", portName, "out", out, "err", err)
		}
		err = fmt.Errorf("can't detach port %s from server %s: %s, %v", portName, serverName, out, err)
		return err
	}
	return nil
}

//DeleteServer destroys a KVM instance
//  sometimes it is not possible to destroy. Like most things in Openstack, try again.
func (s *Platform) DeleteServer(ctx context.Context, id string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "deleting server", "id", id)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "server", "delete", id)
	if err != nil {
		err = fmt.Errorf("can't delete server %s, %s, %v", id, out, err)
		return err
	}
	return nil
}

// CreateNetwork creates a network with a name.
func (s *Platform) CreateNetwork(ctx context.Context, name string, netType string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "creating network", "network", name)
	args := []string{"network", "create"}
	if netType != "" {
		args = append(args, []string{"--provider-network-type", netType}...)
	}
	args = append(args, name)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", args...)
	if err != nil {
		err = fmt.Errorf("can't create network %s, %s, %v", name, out, err)
		return err
	}
	return nil
}

//DeleteNetwork destroys a named network
//  Sometimes it will fail. Openstack will refuse if there are resources attached.
func (s *Platform) DeleteNetwork(ctx context.Context, name string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "deleting network", "network", name)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "network", "delete", name)
	if err != nil {
		err = fmt.Errorf("can't delete network %s, %s, %v", name, out, err)
		return err
	}
	return nil
}

//CreateSubnet creates a subnet within a network. A subnet is assigned ranges. Optionally DHCP can be enabled.
func (s *Platform) CreateSubnet(ctx context.Context, netRange, networkName, gatewayAddr, subnetName string, dhcpEnable bool) error {
	var dhcpFlag string
	if dhcpEnable {
		dhcpFlag = "--dhcp"
	} else {
		dhcpFlag = "--no-dhcp"
	}
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "subnet", "create",
		"--subnet-range", netRange, // e.g. 10.101.101.0/24
		"--network", networkName, // mex-k8s-net-1
		dhcpFlag,
		"--gateway", gatewayAddr, // e.g. 10.101.101.1
		subnetName) // e.g. mex-k8s-subnet-1
	if err != nil {
		nerr := &NeutronErrorType{}
		if ix := strings.Index(string(out), `{"NeutronError":`); ix > 0 {
			neutronErr := out[ix:]
			if jerr := json.Unmarshal(neutronErr, nerr); jerr != nil {
				err = fmt.Errorf("can't create subnet %s, %s, %v, error while parsing neutron error, %v", subnetName, out, err, jerr)
				return err
			}
			if strings.Index(nerr.NeutronError.Message, "overlap") > 0 {
				sd, serr := s.GetSubnetDetail(ctx, subnetName)
				if serr != nil {
					return fmt.Errorf("cannot get subnet detail for %s, while fixing overlap error, %v", subnetName, serr)
				}
				log.SpanLog(ctx, log.DebugLevelMexos, "create subnet, existing subnet detail", "subnet detail", sd)

				//XXX do more validation

				log.SpanLog(ctx, log.DebugLevelMexos, "create subnet, reusing existing subnet", "result", out, "error", err)
				return nil
			}
		}
		err = fmt.Errorf("can't create subnet %s, %s, %v", subnetName, out, err)
		return err
	}
	return nil
}

//DeleteSubnet deletes the subnet. If this fails, remove any attached resources, like router, and try again.
func (s *Platform) DeleteSubnet(ctx context.Context, subnetName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "deleting subnet", "name", subnetName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "subnet", "delete", subnetName)
	if err != nil {
		err = fmt.Errorf("can't delete subnet %s, %s, %v", subnetName, out, err)
		return err
	}
	return nil
}

//CreateRouter creates new router. A router can be attached to network and subnets.
func (s *Platform) CreateRouter(ctx context.Context, routerName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "creating router", "name", routerName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "router", "create", routerName)
	if err != nil {
		err = fmt.Errorf("can't create router %s, %s, %v", routerName, out, err)
		return err
	}
	return nil
}

//DeleteRouter removes the named router. The router needs to not be in use at the time of deletion.
func (s *Platform) DeleteRouter(ctx context.Context, routerName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "deleting router", "name", routerName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "router", "delete", routerName)
	if err != nil {
		err = fmt.Errorf("can't delete router %s, %s, %v", routerName, out, err)
		return err
	}
	return nil
}

//SetRouter assigns the router to a particular network. The network needs to be attached to
// a real external network. This is intended only for routing to external network for now. No internal routers.
// Sometimes, oftentimes, it will fail if the network is not external.
func (s *Platform) SetRouter(ctx context.Context, routerName, networkName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "setting router to network", "router", routerName, "network", networkName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "router", "set", routerName, "--external-gateway", networkName)
	if err != nil {
		err = fmt.Errorf("can't set router %s to %s, %s, %v", routerName, networkName, out, err)
		return err
	}
	return nil
}

//AddRouterSubnet will connect subnet to another network, possibly external, via a router
func (s *Platform) AddRouterSubnet(ctx context.Context, routerName, subnetName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "adding router to subnet", "router", routerName, "network", subnetName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "router", "add", "subnet", routerName, subnetName)
	if err != nil {
		err = fmt.Errorf("can't add router %s to subnet %s, %s, %v", routerName, subnetName, out, err)
		return err
	}
	return nil
}

//RemoveRouterSubnet is useful to remove the router from the subnet before deletion. Otherwise subnet cannot
//  be deleted.
func (s *Platform) RemoveRouterSubnet(ctx context.Context, routerName, subnetName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "removing subnet from router", "router", routerName, "subnet", subnetName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "router", "remove", "subnet", routerName, subnetName)
	if err != nil {
		err = fmt.Errorf("can't remove router %s from subnet %s, %s, %v", routerName, subnetName, out, err)
		return err
	}
	return nil
}

//ListSubnets returns a list of subnets available
func (s *Platform) ListSubnets(ctx context.Context, netName string) ([]OSSubnet, error) {
	var err error
	var out []byte
	if netName != "" {
		out, err = s.TimedOpenStackCommand(ctx, "openstack", "subnet", "list", "--network", netName, "-f", "json")
	} else {
		out, err = s.TimedOpenStackCommand(ctx, "openstack", "subnet", "list", "-f", "json")
	}
	if err != nil {
		err = fmt.Errorf("can't get a list of subnets, %s, %v", out, err)
		return nil, err
	}
	subnets := []OSSubnet{}
	err = json.Unmarshal(out, &subnets)
	if err != nil {
		err = fmt.Errorf("can't unmarshal subnets, %v", err)
		return nil, err
	}
	//log.SpanLog(ctx,log.DebugLevelMexos, "list subnets", "subnets", subnets)
	return subnets, nil
}

//ListProjects returns a list of projects we can see
func (s *Platform) ListProjects(ctx context.Context) ([]OSProject, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "project", "list", "-f", "json")
	if err != nil {
		err = fmt.Errorf("can't get a list of projects, %s, %v", out, err)
		return nil, err
	}
	projects := []OSProject{}
	err = json.Unmarshal(out, &projects)
	if err != nil {
		err = fmt.Errorf("can't unmarshal projects, %v", err)
		return nil, err
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "list projects", "projects", projects)
	return projects, nil
}

//ListSecurityGroups returns a list of security groups
func (s *Platform) ListSecurityGroups(ctx context.Context) ([]OSSecurityGroup, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "security", "group", "list", "-f", "json")
	if err != nil {
		err = fmt.Errorf("can't get a list of security groups, %s, %v", out, err)
		return nil, err
	}
	secgrps := []OSSecurityGroup{}
	err = json.Unmarshal(out, &secgrps)
	if err != nil {
		err = fmt.Errorf("can't unmarshal security groups, %v", err)
		return nil, err
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "list security groups", "security groups", secgrps)
	return secgrps, nil
}

//ListSecurityGroups returns a list of security groups
func (s *Platform) ListSecurityGroupRules(ctx context.Context, secGrp string) ([]OSSecurityGroupRule, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "security", "group", "rule", "list", secGrp, "-f", "json")
	if err != nil {
		err = fmt.Errorf("can't get a list of security group rules, %s, %v", out, err)
		return nil, err
	}
	rules := []OSSecurityGroupRule{}
	err = json.Unmarshal(out, &rules)
	if err != nil {
		err = fmt.Errorf("can't unmarshal security group rules, %v", err)
		return nil, err
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "list security group rules", "security groups", rules)
	return rules, nil
}

func (s *Platform) CreateSecurityGroup(ctx context.Context, groupName string) error {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "security", "group", "create", groupName)
	if err != nil {
		err = fmt.Errorf("can't create security group, %s, %v", out, err)
		return err
	}
	return nil
}

func (s *Platform) AddSecurityGroupToPort(ctx context.Context, portID, groupName string) error {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "port", "set", "--security-group", groupName, portID)
	if err != nil {
		err = fmt.Errorf("can't add security group to port, %s, %v", out, err)
		return err
	}
	return nil
}

//ListRouters returns a list of routers available
func (s *Platform) ListRouters(ctx context.Context) ([]OSRouter, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "router", "list", "-f", "json")
	if err != nil {
		err = fmt.Errorf("can't get a list of routers, %s, %v", out, err)
		return nil, err
	}
	routers := []OSRouter{}
	err = json.Unmarshal(out, &routers)
	if err != nil {
		err = fmt.Errorf("can't unmarshal routers, %v", err)
		return nil, err
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "list routers", "routers", routers)
	return routers, nil
}

//GetRouterDetail returns details per router
func (s *Platform) GetRouterDetail(ctx context.Context, routerName string) (*OSRouterDetail, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "router", "show", "-f", "json", routerName)
	if err != nil {
		err = fmt.Errorf("can't get router details for %s, %s, %v", routerName, out, err)
		return nil, err
	}
	routerDetail := &OSRouterDetail{}
	err = json.Unmarshal(out, routerDetail)
	if err != nil {
		err = fmt.Errorf("can't unmarshal router detail, %v", err)
		return nil, err
	}
	//log.SpanLog(ctx,log.DebugLevelMexos, "router detail", "router detail", routerDetail)
	return routerDetail, nil
}

//CreateServerImage snapshots running service into a qcow2 image
func (s *Platform) CreateServerImage(ctx context.Context, serverName, imageName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "creating image snapshot from server", "server", serverName, "image", imageName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "server", "image", "create", serverName, "--name", imageName)
	if err != nil {
		err = fmt.Errorf("can't create image from %s into %s, %s, %v", serverName, imageName, out, err)
		return err
	}
	return nil
}

//CreateImage puts images into glance
func (s *Platform) CreateImage(ctx context.Context, imageName, fileName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "creating image in glance", "image", imageName, "fileName", fileName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "image", "create",
		imageName,
		"--disk-format", s.GetCloudletImageDiskFormat(),
		"--container-format", "bare",
		"--file", fileName)
	if err != nil {
		err = fmt.Errorf("can't create image in glance, %s, %s, %s, %v", imageName, fileName, out, err)
		return err
	}
	return nil
}

//CreateImageFromUrl downloads image from URL and then puts into glance
func (s *Platform) CreateImageFromUrl(ctx context.Context, imageName, imageUrl, md5Sum string) error {
	fileExt, err := cloudcommon.GetFileNameWithExt(imageUrl)
	if err != nil {
		return err
	}
	filePath := "/tmp/" + fileExt
	defer func() {
		// Stale file might be present if download fails/succeeds, deleting it
		if delerr := mexos.DeleteFile(filePath); delerr != nil {
			log.SpanLog(ctx, log.DebugLevelMexos, "delete file failed", "filePath", filePath)
		}
	}()
	err = cloudcommon.DownloadFile(ctx, s.vaultConfig, imageUrl, filePath, nil)
	if err != nil {
		return fmt.Errorf("error downloading image from %s, %v", imageUrl, err)
	}
	// Verify checksum
	if md5Sum != "" {
		fileMd5Sum, err := mexos.Md5SumFile(filePath)
		if err != nil {
			return err
		}
		log.SpanLog(ctx, log.DebugLevelMexos, "verify md5sum", "downloaded-md5sum", fileMd5Sum, "actual-md5sum", md5Sum)
		if fileMd5Sum != md5Sum {
			return fmt.Errorf("mismatch in md5sum for downloaded image: %s", imageName)
		}
	}

	err = s.CreateImage(ctx, imageName, filePath)
	if err != nil {
		return fmt.Errorf("error creating image %v", err)
	}
	return err
}

//SaveImage takes the image name available in glance, as a result of for example the above create image.
// It will then save that into a local file. The image transfer happens from glance into your own laptop
// or whatever.
// This can take a while, transferring all the data.
func (s *Platform) SaveImage(ctx context.Context, saveName, imageName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "saving image", "save name", saveName, "image name", imageName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "image", "save", "--file", saveName, imageName)
	if err != nil {
		err = fmt.Errorf("can't save image from %s to file %s, %s, %v", imageName, saveName, out, err)
		return err
	}
	return nil
}

//DeleteImage deletes the named image from glance. Sometimes backing store is still busy and
// will refuse to honor the request. Like most things in Openstack, wait for a while and try
// again.
func (s *Platform) DeleteImage(ctx context.Context, imageName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "deleting image", "name", imageName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "image", "delete", imageName)
	if err != nil {
		err = fmt.Errorf("can't delete image %s, %s, %v", imageName, out, err)
		return err
	}
	return nil
}

//GetSubnetDetail returns details for the subnet. This is useful when getting router/gateway
//  IP for a given subnet.  The gateway info is used for creating a server.
//  Also useful in general, like other `detail` functions, to get the ID map for the name of subnet.
func (s *Platform) GetSubnetDetail(ctx context.Context, subnetName string) (*OSSubnetDetail, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "subnet", "show", "-f", "json", subnetName)
	if err != nil {
		err = fmt.Errorf("can't get subnet details for %s, %s, %v", subnetName, out, err)
		return nil, err
	}
	subnetDetail := &OSSubnetDetail{}
	err = json.Unmarshal(out, subnetDetail)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal subnet detail, %v", err)
	}
	//log.SpanLog(ctx,log.DebugLevelMexos, "get subnet detail", "subnet detail", subnetDetail)
	return subnetDetail, nil
}

//GetNetworkDetail returns details about a network.  It is used, for example, by GetExternalGateway.
func (s *Platform) GetNetworkDetail(ctx context.Context, networkName string) (*OSNetworkDetail, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "network", "show", "-f", "json", networkName)
	if err != nil {
		err = fmt.Errorf("can't get details for network %s, %s, %v", networkName, out, err)
		return nil, err
	}
	networkDetail := &OSNetworkDetail{}
	err = json.Unmarshal(out, networkDetail)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal network detail, %v", err)
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "get network detail", "network detail", networkDetail)
	return networkDetail, nil
}

//SetServerProperty sets properties for the server
func (s *Platform) SetServerProperty(ctx context.Context, name, property string) error {
	if name == "" {
		return fmt.Errorf("empty name")
	}
	if property == "" {
		return fmt.Errorf("empty property")
	}
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "server", "set", "--property", property, name)
	if err != nil {
		return fmt.Errorf("can't set property %s on server %s, %s, %v", property, name, out, err)
	}
	log.SpanLog(ctx, log.DebugLevelMexos, "set server property", "name", name, "property", property)
	return nil
}

// createHeatStack creates a stack with the given template
func (s *Platform) createHeatStack(ctx context.Context, templateFile string, stackName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "create heat stack", "template", templateFile, "stackName", stackName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "stack", "create", "--template", templateFile, stackName)
	if err != nil {
		return fmt.Errorf("error creating heat stack: %s, %s -- %v", templateFile, string(out), err)
	}
	return nil
}

func (s *Platform) updateHeatStack(ctx context.Context, templateFile string, stackName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "update heat stack", "template", templateFile, "stackName", stackName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "stack", "update", "--template", templateFile, stackName)
	if err != nil {
		return fmt.Errorf("error udpating heat stack: %s -- %s, %v", templateFile, out, err)
	}
	return nil
}

// deleteHeatStack delete a stack with the given name
func (s *Platform) deleteHeatStack(ctx context.Context, stackName string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "delete heat stack", "stackName", stackName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "stack", "delete", stackName)
	if err != nil {
		if strings.Contains("Stack not found", string(out)) {
			log.SpanLog(ctx, log.DebugLevelMexos, "stack not found")
			return nil
		}
		log.SpanLog(ctx, log.DebugLevelMexos, "stack deletion failed", "stackName", stackName, "out", string(out), "err", err)
		if strings.Contains(string(out), "Stack not found") {
			log.SpanLog(ctx, log.DebugLevelMexos, "stack already deleted", "stackName", stackName)
			return nil
		}
		return fmt.Errorf("stack deletion failed: %s, %s %v", stackName, out, err)
	}
	return nil
}

// getHeatStackDetail gets details of the provided stack
func (s *Platform) getHeatStackDetail(ctx context.Context, stackName string) (*OSHeatStackDetail, error) {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "stack", "show", "-f", "json", stackName)
	if err != nil {
		err = fmt.Errorf("can't get stack details for %s, %s, %v", stackName, out, err)
		return nil, err
	}
	stackDetail := &OSHeatStackDetail{}
	err = json.Unmarshal(out, stackDetail)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal stack detail, %v", err)
	}
	return stackDetail, nil
}

// Get resource limits
func (s *Platform) OSGetLimits(ctx context.Context, info *edgeproto.CloudletInfo) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "GetLimits (Openstack) - Resources info & Supported flavors")
	var limits []OSLimit
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "limits", "show", "--absolute", "-f", "json")
	if err != nil {
		err = fmt.Errorf("cannot get limits from openstack, %s, %v", out, err)
		return err
	}
	err = json.Unmarshal(out, &limits)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return err
	}
	for _, l := range limits {
		if l.Name == "maxTotalRAMSize" {
			info.OsMaxRam = uint64(l.Value)
		} else if l.Name == "maxTotalCores" {
			info.OsMaxVcores = uint64(l.Value)
		} else if l.Name == "maxTotalVolumeGigabytes" {
			info.OsMaxVolGb = uint64(l.Value)
		}
	}

	finfo, _, _, err := s.GetFlavorInfo(ctx)
	if err != nil {
		return err
	}
	info.Flavors = finfo
	return nil
}

func (s *Platform) OSGetAllLimits(ctx context.Context) ([]OSLimit, error) {
	log.SpanLog(ctx, log.DebugLevelMexos, "GetLimits (Openstack) - Resources info and usage")
	var limits []OSLimit
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "limits", "show", "--absolute", "-f", "json")
	if err != nil {
		err = fmt.Errorf("cannot get limits from openstack, %s, %v", out, err)
		return nil, err
	}
	err = json.Unmarshal(out, &limits)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal, %v", err)
		return nil, err
	}
	return limits, nil
}

func (s *Platform) GetFlavorInfo(ctx context.Context) ([]*edgeproto.FlavorInfo, []OSAZone, []OSImage, error) {

	var props map[string]string

	osflavors, err := s.ListFlavors(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get flavors, %v", err.Error())
	}
	if len(osflavors) == 0 {
		return nil, nil, nil, fmt.Errorf("no flavors found")
	}
	var finfo []*edgeproto.FlavorInfo
	for _, f := range osflavors {
		if f.Properties != "" {
			props = ParseFlavorProperties(f)
		}

		finfo = append(
			finfo,
			&edgeproto.FlavorInfo{
				Name:    f.Name,
				Vcpus:   uint64(f.VCPUs),
				Ram:     uint64(f.RAM),
				Disk:    uint64(f.Disk),
				PropMap: props},
		)
	}
	zones, err := s.ListAZones(ctx)
	images, err := s.ListImages(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	return finfo, zones, images, nil
}

func (s *Platform) GetSecurityGroupIDForProject(ctx context.Context, grpname string, projectID string) (string, error) {
	grps, err := s.ListSecurityGroups(ctx)
	if err != nil {
		return "", err
	}
	for _, g := range grps {
		if g.Name == grpname {
			if g.Project == projectID {
				log.SpanLog(ctx, log.DebugLevelMexos, "GetSecurityGroupIDForProject", "projectID", projectID, "group", grpname)
				return g.ID, nil
			}
			if g.Project == "" {
				// This is an openstack bug in some environments in which it may not show the project ids when listing the group
				// all we can do is hope for no conflicts in this case
				log.SpanLog(ctx, log.DebugLevelMexos, "Warning: no project id returned for security group", "group", grpname)
				return g.ID, nil
			}
		}
	}
	return "", fmt.Errorf("unable to find security group %s project %s", grpname, projectID)
}

func (s *Platform) OSGetConsoleUrl(ctx context.Context, serverName string) (*OSConsoleUrl, error) {
	log.SpanLog(ctx, log.DebugLevelMexos, "get console url", "server", serverName)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "console", "url", "show", "-f", "json", "-c", "url", "--novnc", serverName)
	if err != nil {
		err = fmt.Errorf("can't get console url details for %s, %s, %v", serverName, out, err)
		return nil, err
	}
	consoleUrl := &OSConsoleUrl{}
	err = json.Unmarshal(out, consoleUrl)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal console url output, %v", err)
	}
	return consoleUrl, nil
}

// Finds a resource by name by instance id.
// There are resources that are metered for instance-id, which are resources of their own
// The examples are instance_network_interface and instance_disk
// Openstack example call:
//   <openstack metric resource search --type instance_network_interface instance_id=dc32daa6-0d0a-4512-a9fa-2b989e913014>
// We only use the the first found result
func (s *Platform) OSFindResourceByInstId(ctx context.Context, resourceType string, instId string) (*OSMetricResource, error) {
	log.SpanLog(ctx, log.DebugLevelMexos, "find resource for instance Id", "id", instId,
		"resource", resourceType)
	osRes := []OSMetricResource{}
	instArg := fmt.Sprintf("instance_id=%s", instId)
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "metric", "resource", "search",
		"-f", "json", "--type", resourceType, instArg)
	if err != nil {
		err = fmt.Errorf("can't find resource %s, for %s, %s %v", resourceType, instId, out, err)
		return nil, err
	}
	err = json.Unmarshal(out, &osRes)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal Metric Resource, %v", err)
		return nil, err
	}
	if len(osRes) != 1 {
		return nil, fmt.Errorf("Unexpected Number of Meters found")
	}
	return &osRes[0], nil
}

// Get openstack metrics from ceilometer tsdb
// Example openstack call:
//   <openstack metric measures show --resource-id a9bf10cf-a709-5a47-8b69-da920b8f65cd network.incoming.bytes>
// This will return a range of measurements from the startTime
func (s *Platform) OSGetMetricsRangeForId(ctx context.Context, resId string, metric string, startTime time.Time) ([]OSMetricMeasurement, error) {
	log.SpanLog(ctx, log.DebugLevelMexos, "get measure for Id", "id", resId, "metric", metric)
	measurements := []OSMetricMeasurement{}

	startStr := startTime.Format(time.RFC3339)

	out, err := s.TimedOpenStackCommand(ctx, "openstack", "metric", "measures", "show",
		"-f", "json", "--start", startStr, "--resource-id", resId, metric)
	if err != nil {
		err = fmt.Errorf("can't get measurements %s, for %s, %s %v", metric, resId, out, err)
		return []OSMetricMeasurement{}, err
	}
	err = json.Unmarshal(out, &measurements)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal measurements, %v", err)
		return []OSMetricMeasurement{}, err
	}
	// No value, means we don't need to write it
	if len(measurements) == 0 {
		return []OSMetricMeasurement{}, fmt.Errorf("No values for the metric")
	}
	return measurements, nil
}

func (s *Platform) AddImageIfNotPresent(ctx context.Context, imgPathPrefix, imgVersion string, updateCallback edgeproto.CacheUpdateCallback) (string, error) {
	imgPath := mexos.GetCloudletVMImagePath(imgPathPrefix, imgVersion)

	// Fetch platform base image name
	pfImageName, err := cloudcommon.GetFileName(imgPath)
	if err != nil {
		return "", err
	}
	// Use PlatformBaseImage, if not present then fetch it from MobiledgeX VM registry
	imageDetail, err := s.GetImageDetail(ctx, pfImageName)
	if err == nil && imageDetail.Status != "active" {
		return "", fmt.Errorf("image %s is not active", pfImageName)
	}
	if err != nil {
		// Validate if pfImageName is same as we expected
		_, md5Sum, err := mexos.GetUrlInfo(ctx, s.vaultConfig, imgPath)
		if err != nil {
			return "", err
		}
		// Download platform base image and Add to Openstack Glance
		updateCallback(edgeproto.UpdateTask, "Downloading platform base image: "+pfImageName)
		err = s.CreateImageFromUrl(ctx, pfImageName, imgPath, md5Sum)
		if err != nil {
			return "", fmt.Errorf("Error downloading platform base image %s: %v", pfImageName, err)
		}
	}
	return pfImageName, nil
}

func (s *Platform) OSSetPowerState(ctx context.Context, serverName, serverAction string) error {
	log.SpanLog(ctx, log.DebugLevelMexos, "setting server state", "serverName", serverName, "serverAction", serverAction)

	out, err := s.TimedOpenStackCommand(ctx, "openstack", "server", serverAction, serverName)
	if err != nil {
		err = fmt.Errorf("unable to %s server %s, %s, %v", serverAction, serverName, out, err)
		return err
	}

	return nil
}

func (s *Platform) AddSecurityRuleCIDR(ctx context.Context, cidr string, proto string, groupName string, port string) error {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "security", "group", "rule", "create", "--remote-ip", cidr, "--proto", proto, "--dst-port", port, "--ingress", groupName)
	if err != nil {
		if strings.Contains(string(out), "Security group rule already exists") {
			log.SpanLog(ctx, log.DebugLevelMexos, "security group rule already exists, proceeding")
		} else {
			return fmt.Errorf("can't add security group rule for port %s to %s,%s,%v", port, groupName, string(out), err)
		}
	}
	return nil
}

func (s *Platform) DeleteSecurityGroupRule(ctx context.Context, ruleID string) error {
	out, err := s.TimedOpenStackCommand(ctx, "openstack", "security", "group", "rule", "delete", ruleID)
	if err != nil {
		return fmt.Errorf("can't delete security group rule %s,%s,%v", ruleID, string(out), err)
	}
	return nil
}