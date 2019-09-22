// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cloudletpool.proto

package ormctl

import edgeproto "github.com/mobiledgex/edge-cloud/edgeproto"
import "strings"
import "github.com/mobiledgex/edge-cloud-infra/mc/ormapi"
import "github.com/mobiledgex/edge-cloud/cli"
import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/googleapis/google/api"
import _ "github.com/mobiledgex/edge-cloud/protogen"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Auto-generated code: DO NOT EDIT

var CreateCloudletPoolCmd = &cli.Command{
	Use:          "CreateCloudletPool",
	RequiredArgs: strings.Join(append([]string{"region"}, CloudletPoolRequiredArgs...), " "),
	OptionalArgs: strings.Join(CloudletPoolOptionalArgs, " "),
	AliasArgs:    strings.Join(CloudletPoolAliasArgs, " "),
	SpecialArgs:  &CloudletPoolSpecialArgs,
	Comments:     addRegionComment(CloudletPoolComments),
	ReqData:      &ormapi.RegionCloudletPool{},
	ReplyData:    &edgeproto.Result{},
	Run:          runRest("/auth/ctrl/CreateCloudletPool"),
}

var DeleteCloudletPoolCmd = &cli.Command{
	Use:          "DeleteCloudletPool",
	RequiredArgs: strings.Join(append([]string{"region"}, CloudletPoolRequiredArgs...), " "),
	OptionalArgs: strings.Join(CloudletPoolOptionalArgs, " "),
	AliasArgs:    strings.Join(CloudletPoolAliasArgs, " "),
	SpecialArgs:  &CloudletPoolSpecialArgs,
	Comments:     addRegionComment(CloudletPoolComments),
	ReqData:      &ormapi.RegionCloudletPool{},
	ReplyData:    &edgeproto.Result{},
	Run:          runRest("/auth/ctrl/DeleteCloudletPool"),
}

var ShowCloudletPoolCmd = &cli.Command{
	Use:          "ShowCloudletPool",
	RequiredArgs: "region",
	OptionalArgs: strings.Join(append(CloudletPoolRequiredArgs, CloudletPoolOptionalArgs...), " "),
	AliasArgs:    strings.Join(CloudletPoolAliasArgs, " "),
	SpecialArgs:  &CloudletPoolSpecialArgs,
	Comments:     addRegionComment(CloudletPoolComments),
	ReqData:      &ormapi.RegionCloudletPool{},
	ReplyData:    &edgeproto.CloudletPool{},
	Run:          runRest("/auth/ctrl/ShowCloudletPool"),
	StreamOut:    true,
}

var CloudletPoolApiCmds = []*cli.Command{
	CreateCloudletPoolCmd,
	DeleteCloudletPoolCmd,
	ShowCloudletPoolCmd,
}

var CreateCloudletPoolMemberCmd = &cli.Command{
	Use:          "CreateCloudletPoolMember",
	RequiredArgs: strings.Join(append([]string{"region"}, CloudletPoolMemberRequiredArgs...), " "),
	OptionalArgs: strings.Join(CloudletPoolMemberOptionalArgs, " "),
	AliasArgs:    strings.Join(CloudletPoolMemberAliasArgs, " "),
	SpecialArgs:  &CloudletPoolMemberSpecialArgs,
	Comments:     addRegionComment(CloudletPoolMemberComments),
	ReqData:      &ormapi.RegionCloudletPoolMember{},
	ReplyData:    &edgeproto.Result{},
	Run:          runRest("/auth/ctrl/CreateCloudletPoolMember"),
}

var DeleteCloudletPoolMemberCmd = &cli.Command{
	Use:          "DeleteCloudletPoolMember",
	RequiredArgs: strings.Join(append([]string{"region"}, CloudletPoolMemberRequiredArgs...), " "),
	OptionalArgs: strings.Join(CloudletPoolMemberOptionalArgs, " "),
	AliasArgs:    strings.Join(CloudletPoolMemberAliasArgs, " "),
	SpecialArgs:  &CloudletPoolMemberSpecialArgs,
	Comments:     addRegionComment(CloudletPoolMemberComments),
	ReqData:      &ormapi.RegionCloudletPoolMember{},
	ReplyData:    &edgeproto.Result{},
	Run:          runRest("/auth/ctrl/DeleteCloudletPoolMember"),
}

var ShowCloudletPoolMemberCmd = &cli.Command{
	Use:          "ShowCloudletPoolMember",
	RequiredArgs: "region",
	OptionalArgs: strings.Join(append(CloudletPoolMemberRequiredArgs, CloudletPoolMemberOptionalArgs...), " "),
	AliasArgs:    strings.Join(CloudletPoolMemberAliasArgs, " "),
	SpecialArgs:  &CloudletPoolMemberSpecialArgs,
	Comments:     addRegionComment(CloudletPoolMemberComments),
	ReqData:      &ormapi.RegionCloudletPoolMember{},
	ReplyData:    &edgeproto.CloudletPoolMember{},
	Run:          runRest("/auth/ctrl/ShowCloudletPoolMember"),
	StreamOut:    true,
}

var ShowPoolsForCloudletCmd = &cli.Command{
	Use:          "ShowPoolsForCloudlet",
	RequiredArgs: "region",
	OptionalArgs: strings.Join(append(CloudletKeyRequiredArgs, CloudletKeyOptionalArgs...), " "),
	AliasArgs:    strings.Join(CloudletKeyAliasArgs, " "),
	SpecialArgs:  &CloudletKeySpecialArgs,
	Comments:     addRegionComment(CloudletKeyComments),
	ReqData:      &ormapi.RegionCloudletKey{},
	ReplyData:    &edgeproto.CloudletPool{},
	Run:          runRest("/auth/ctrl/ShowPoolsForCloudlet"),
	StreamOut:    true,
}

var ShowCloudletsForPoolCmd = &cli.Command{
	Use:          "ShowCloudletsForPool",
	RequiredArgs: "region",
	OptionalArgs: strings.Join(append(CloudletPoolKeyRequiredArgs, CloudletPoolKeyOptionalArgs...), " "),
	AliasArgs:    strings.Join(CloudletPoolKeyAliasArgs, " "),
	SpecialArgs:  &CloudletPoolKeySpecialArgs,
	Comments:     addRegionComment(CloudletPoolKeyComments),
	ReqData:      &ormapi.RegionCloudletPoolKey{},
	ReplyData:    &edgeproto.Cloudlet{},
	Run:          runRest("/auth/ctrl/ShowCloudletsForPool"),
	StreamOut:    true,
}

var ShowCloudletsForPoolListCmd = &cli.Command{
	Use:          "ShowCloudletsForPoolList",
	RequiredArgs: "region",
	OptionalArgs: strings.Join(append(CloudletPoolListRequiredArgs, CloudletPoolListOptionalArgs...), " "),
	AliasArgs:    strings.Join(CloudletPoolListAliasArgs, " "),
	SpecialArgs:  &CloudletPoolListSpecialArgs,
	Comments:     addRegionComment(CloudletPoolListComments),
	ReqData:      &ormapi.RegionCloudletPoolList{},
	ReplyData:    &edgeproto.Cloudlet{},
	Run:          runRest("/auth/ctrl/ShowCloudletsForPoolList"),
	StreamOut:    true,
}

var CloudletPoolMemberApiCmds = []*cli.Command{
	CreateCloudletPoolMemberCmd,
	DeleteCloudletPoolMemberCmd,
	ShowCloudletPoolMemberCmd,
	ShowPoolsForCloudletCmd,
	ShowCloudletsForPoolCmd,
	ShowCloudletsForPoolListCmd,
}

var CloudletPoolKeyRequiredArgs = []string{}
var CloudletPoolKeyOptionalArgs = []string{
	"name",
}
var CloudletPoolKeyAliasArgs = []string{
	"name=cloudletpoolkey.name",
}
var CloudletPoolKeyComments = map[string]string{
	"name": "CloudletPool Name",
}
var CloudletPoolKeySpecialArgs = map[string]string{}
var CloudletPoolRequiredArgs = []string{
	"name",
}
var CloudletPoolOptionalArgs = []string{}
var CloudletPoolAliasArgs = []string{
	"name=cloudletpool.key.name",
}
var CloudletPoolComments = map[string]string{
	"name": "CloudletPool Name",
}
var CloudletPoolSpecialArgs = map[string]string{}
var CloudletPoolMemberRequiredArgs = []string{
	"pool",
	"operator",
	"cloudlet",
}
var CloudletPoolMemberOptionalArgs = []string{}
var CloudletPoolMemberAliasArgs = []string{
	"pool=cloudletpoolmember.poolkey.name",
	"operator=cloudletpoolmember.cloudletkey.operatorkey.name",
	"cloudlet=cloudletpoolmember.cloudletkey.name",
}
var CloudletPoolMemberComments = map[string]string{
	"pool":     "CloudletPool Name",
	"operator": "Company or Organization name of the operator",
	"cloudlet": "Name of the cloudlet",
}
var CloudletPoolMemberSpecialArgs = map[string]string{}
var CloudletPoolListRequiredArgs = []string{}
var CloudletPoolListOptionalArgs = []string{
	"poolname",
}
var CloudletPoolListAliasArgs = []string{
	"poolname=cloudletpoollist.poolname",
}
var CloudletPoolListComments = map[string]string{
	"poolname": "Name of Cloudlet Pool (may be repeated)",
}
var CloudletPoolListSpecialArgs = map[string]string{
	"poolname": "StringArray",
}