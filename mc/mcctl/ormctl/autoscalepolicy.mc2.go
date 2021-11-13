// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: autoscalepolicy.proto

package ormctl

import (
	fmt "fmt"
	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	"github.com/mobiledgex/edge-cloud-infra/mc/ormapi"
	edgeproto "github.com/mobiledgex/edge-cloud/edgeproto"
	_ "github.com/mobiledgex/edge-cloud/protogen"
	math "math"
	"strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Auto-generated code: DO NOT EDIT

var CreateAutoScalePolicyCmd = &ApiCommand{
	Name:         "CreateAutoScalePolicy",
	Use:          "create",
	Short:        "Create an Auto Scale Policy",
	RequiredArgs: "region " + strings.Join(CreateAutoScalePolicyRequiredArgs, " "),
	OptionalArgs: strings.Join(CreateAutoScalePolicyOptionalArgs, " "),
	AliasArgs:    strings.Join(AutoScalePolicyAliasArgs, " "),
	SpecialArgs:  &AutoScalePolicySpecialArgs,
	Comments:     addRegionComment(AutoScalePolicyComments),
	NoConfig:     "DeletePrepare",
	ReqData:      &ormapi.RegionAutoScalePolicy{},
	ReplyData:    &edgeproto.Result{},
	Path:         "/auth/ctrl/CreateAutoScalePolicy",
	ProtobufApi:  true,
}

var DeleteAutoScalePolicyCmd = &ApiCommand{
	Name:         "DeleteAutoScalePolicy",
	Use:          "delete",
	Short:        "Delete an Auto Scale Policy",
	RequiredArgs: "region " + strings.Join(AutoScalePolicyRequiredArgs, " "),
	OptionalArgs: strings.Join(AutoScalePolicyOptionalArgs, " "),
	AliasArgs:    strings.Join(AutoScalePolicyAliasArgs, " "),
	SpecialArgs:  &AutoScalePolicySpecialArgs,
	Comments:     addRegionComment(AutoScalePolicyComments),
	NoConfig:     "DeletePrepare",
	ReqData:      &ormapi.RegionAutoScalePolicy{},
	ReplyData:    &edgeproto.Result{},
	Path:         "/auth/ctrl/DeleteAutoScalePolicy",
	ProtobufApi:  true,
}

var UpdateAutoScalePolicyCmd = &ApiCommand{
	Name:         "UpdateAutoScalePolicy",
	Use:          "update",
	Short:        "Update an Auto Scale Policy",
	RequiredArgs: "region " + strings.Join(AutoScalePolicyRequiredArgs, " "),
	OptionalArgs: strings.Join(AutoScalePolicyOptionalArgs, " "),
	AliasArgs:    strings.Join(AutoScalePolicyAliasArgs, " "),
	SpecialArgs:  &AutoScalePolicySpecialArgs,
	Comments:     addRegionComment(AutoScalePolicyComments),
	NoConfig:     "DeletePrepare",
	ReqData:      &ormapi.RegionAutoScalePolicy{},
	ReplyData:    &edgeproto.Result{},
	Path:         "/auth/ctrl/UpdateAutoScalePolicy",
	ProtobufApi:  true,
}

var ShowAutoScalePolicyCmd = &ApiCommand{
	Name:         "ShowAutoScalePolicy",
	Use:          "show",
	Short:        "Show Auto Scale Policies. Any fields specified will be used to filter results.",
	RequiredArgs: "region",
	OptionalArgs: strings.Join(append(AutoScalePolicyRequiredArgs, AutoScalePolicyOptionalArgs...), " "),
	AliasArgs:    strings.Join(AutoScalePolicyAliasArgs, " "),
	SpecialArgs:  &AutoScalePolicySpecialArgs,
	Comments:     addRegionComment(AutoScalePolicyComments),
	NoConfig:     "DeletePrepare",
	ReqData:      &ormapi.RegionAutoScalePolicy{},
	ReplyData:    &edgeproto.AutoScalePolicy{},
	Path:         "/auth/ctrl/ShowAutoScalePolicy",
	StreamOut:    true,
	ProtobufApi:  true,
}
var AutoScalePolicyApiCmds = []*ApiCommand{
	CreateAutoScalePolicyCmd,
	DeleteAutoScalePolicyCmd,
	UpdateAutoScalePolicyCmd,
	ShowAutoScalePolicyCmd,
}

const AutoScalePolicyGroup = "AutoScalePolicy"

func init() {
	AllApis.AddGroup(AutoScalePolicyGroup, "Manage AutoScalePolicys", AutoScalePolicyApiCmds)
}

var CreateAutoScalePolicyRequiredArgs = []string{
	"cluster-org",
	"name",
	"minnodes",
	"maxnodes",
}
var CreateAutoScalePolicyOptionalArgs = []string{
	"scaleupcputhresh",
	"scaledowncputhresh",
	"triggertimesec",
	"stabilizationwindowsec",
	"targetcpu",
	"targetmem",
	"targetactiveconnections",
}
var AutoScalePolicyRequiredArgs = []string{
	"cluster-org",
	"name",
}
var AutoScalePolicyOptionalArgs = []string{
	"minnodes",
	"maxnodes",
	"scaleupcputhresh",
	"scaledowncputhresh",
	"triggertimesec",
	"stabilizationwindowsec",
	"targetcpu",
	"targetmem",
	"targetactiveconnections",
}
var AutoScalePolicyAliasArgs = []string{
	"fields=autoscalepolicy.fields",
	"cluster-org=autoscalepolicy.key.organization",
	"name=autoscalepolicy.key.name",
	"minnodes=autoscalepolicy.minnodes",
	"maxnodes=autoscalepolicy.maxnodes",
	"scaleupcputhresh=autoscalepolicy.scaleupcputhresh",
	"scaledowncputhresh=autoscalepolicy.scaledowncputhresh",
	"triggertimesec=autoscalepolicy.triggertimesec",
	"stabilizationwindowsec=autoscalepolicy.stabilizationwindowsec",
	"targetcpu=autoscalepolicy.targetcpu",
	"targetmem=autoscalepolicy.targetmem",
	"targetactiveconnections=autoscalepolicy.targetactiveconnections",
	"deleteprepare=autoscalepolicy.deleteprepare",
}
var AutoScalePolicyComments = map[string]string{
	"fields":                  "Fields are used for the Update API to specify which fields to apply",
	"cluster-org":             "Name of the organization for the cluster that this policy will apply to",
	"name":                    "Policy name",
	"minnodes":                "Minimum number of cluster nodes",
	"maxnodes":                "Maximum number of cluster nodes",
	"scaleupcputhresh":        "(Deprecated) Scale up cpu threshold (percentage 1 to 100), 0 means disabled",
	"scaledowncputhresh":      "(Deprecated) Scale down cpu threshold (percentage 1 to 100), 0 means disabled",
	"triggertimesec":          "(Deprecated) Trigger time defines how long the target must be satified in seconds before acting upon it.",
	"stabilizationwindowsec":  "Stabilization window is the time for which past triggers are considered; the largest scale factor is always taken.",
	"targetcpu":               "Target per-node cpu utilization (percentage 1 to 100), 0 means disabled",
	"targetmem":               "Target per-node memory utilization (percentage 1 to 100), 0 means disabled",
	"targetactiveconnections": "Target per-node number of active connections, 0 means disabled",
	"deleteprepare":           "Preparing to be deleted",
}
var AutoScalePolicySpecialArgs = map[string]string{
	"autoscalepolicy.fields": "StringArray",
}
