// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: autoprovpolicy.proto

package ormctl

import (
	fmt "fmt"
	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	"github.com/mobiledgex/edge-cloud-infra/mc/ormapi"
	_ "github.com/mobiledgex/edge-cloud/d-match-engine/dme-proto"
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

var CreateAutoProvPolicyCmd = &ApiCommand{
	Name:         "CreateAutoProvPolicy",
	Use:          "create",
	Short:        "Create an Auto Provisioning Policy",
	RequiredArgs: "region " + strings.Join(CreateAutoProvPolicyRequiredArgs, " "),
	OptionalArgs: strings.Join(CreateAutoProvPolicyOptionalArgs, " "),
	AliasArgs:    strings.Join(AutoProvPolicyAliasArgs, " "),
	SpecialArgs:  &AutoProvPolicySpecialArgs,
	Comments:     addRegionComment(AutoProvPolicyComments),
	NoConfig:     "Cloudlets:#.Loc",
	ReqData:      &ormapi.RegionAutoProvPolicy{},
	ReplyData:    &edgeproto.Result{},
	Path:         "/auth/ctrl/CreateAutoProvPolicy",
	ProtobufApi:  true,
}

var DeleteAutoProvPolicyCmd = &ApiCommand{
	Name:         "DeleteAutoProvPolicy",
	Use:          "delete",
	Short:        "Delete an Auto Provisioning Policy",
	RequiredArgs: "region " + strings.Join(AutoProvPolicyRequiredArgs, " "),
	OptionalArgs: strings.Join(AutoProvPolicyOptionalArgs, " "),
	AliasArgs:    strings.Join(AutoProvPolicyAliasArgs, " "),
	SpecialArgs:  &AutoProvPolicySpecialArgs,
	Comments:     addRegionComment(AutoProvPolicyComments),
	NoConfig:     "Cloudlets:#.Loc",
	ReqData:      &ormapi.RegionAutoProvPolicy{},
	ReplyData:    &edgeproto.Result{},
	Path:         "/auth/ctrl/DeleteAutoProvPolicy",
	ProtobufApi:  true,
}

var UpdateAutoProvPolicyCmd = &ApiCommand{
	Name:         "UpdateAutoProvPolicy",
	Use:          "update",
	Short:        "Update an Auto Provisioning Policy",
	RequiredArgs: "region " + strings.Join(AutoProvPolicyRequiredArgs, " "),
	OptionalArgs: strings.Join(AutoProvPolicyOptionalArgs, " "),
	AliasArgs:    strings.Join(AutoProvPolicyAliasArgs, " "),
	SpecialArgs:  &AutoProvPolicySpecialArgs,
	Comments:     addRegionComment(AutoProvPolicyComments),
	NoConfig:     "Cloudlets:#.Loc",
	ReqData:      &ormapi.RegionAutoProvPolicy{},
	ReplyData:    &edgeproto.Result{},
	Path:         "/auth/ctrl/UpdateAutoProvPolicy",
	ProtobufApi:  true,
}

var ShowAutoProvPolicyCmd = &ApiCommand{
	Name:         "ShowAutoProvPolicy",
	Use:          "show",
	Short:        "Show Auto Provisioning Policies. Any fields specified will be used to filter results.",
	RequiredArgs: "region",
	OptionalArgs: strings.Join(append(AutoProvPolicyRequiredArgs, AutoProvPolicyOptionalArgs...), " "),
	AliasArgs:    strings.Join(AutoProvPolicyAliasArgs, " "),
	SpecialArgs:  &AutoProvPolicySpecialArgs,
	Comments:     addRegionComment(AutoProvPolicyComments),
	NoConfig:     "Cloudlets:#.Loc",
	ReqData:      &ormapi.RegionAutoProvPolicy{},
	ReplyData:    &edgeproto.AutoProvPolicy{},
	Path:         "/auth/ctrl/ShowAutoProvPolicy",
	StreamOut:    true,
	ProtobufApi:  true,
}

var AddAutoProvPolicyCloudletCmd = &ApiCommand{
	Name:         "AddAutoProvPolicyCloudlet",
	Use:          "addcloudlet",
	Short:        "Add a Cloudlet to the Auto Provisioning Policy",
	RequiredArgs: "region " + strings.Join(AutoProvPolicyCloudletRequiredArgs, " "),
	OptionalArgs: strings.Join(AutoProvPolicyCloudletOptionalArgs, " "),
	AliasArgs:    strings.Join(AutoProvPolicyCloudletAliasArgs, " "),
	SpecialArgs:  &AutoProvPolicyCloudletSpecialArgs,
	Comments:     addRegionComment(AutoProvPolicyCloudletComments),
	ReqData:      &ormapi.RegionAutoProvPolicyCloudlet{},
	ReplyData:    &edgeproto.Result{},
	Path:         "/auth/ctrl/AddAutoProvPolicyCloudlet",
	ProtobufApi:  true,
}

var RemoveAutoProvPolicyCloudletCmd = &ApiCommand{
	Name:         "RemoveAutoProvPolicyCloudlet",
	Use:          "removecloudlet",
	Short:        "Remove a Cloudlet from the Auto Provisioning Policy",
	RequiredArgs: "region " + strings.Join(AutoProvPolicyCloudletRequiredArgs, " "),
	OptionalArgs: strings.Join(AutoProvPolicyCloudletOptionalArgs, " "),
	AliasArgs:    strings.Join(AutoProvPolicyCloudletAliasArgs, " "),
	SpecialArgs:  &AutoProvPolicyCloudletSpecialArgs,
	Comments:     addRegionComment(AutoProvPolicyCloudletComments),
	ReqData:      &ormapi.RegionAutoProvPolicyCloudlet{},
	ReplyData:    &edgeproto.Result{},
	Path:         "/auth/ctrl/RemoveAutoProvPolicyCloudlet",
	ProtobufApi:  true,
}
var AutoProvPolicyApiCmds = []*ApiCommand{
	CreateAutoProvPolicyCmd,
	DeleteAutoProvPolicyCmd,
	UpdateAutoProvPolicyCmd,
	ShowAutoProvPolicyCmd,
	AddAutoProvPolicyCloudletCmd,
	RemoveAutoProvPolicyCloudletCmd,
}

const AutoProvPolicyGroup = "AutoProvPolicy"

func init() {
	AllApis.AddGroup(AutoProvPolicyGroup, "Manage AutoProvPolicys", AutoProvPolicyApiCmds)
}

var CreateAutoProvPolicyRequiredArgs = []string{
	"app-org",
	"name",
}
var CreateAutoProvPolicyOptionalArgs = []string{
	"deployclientcount",
	"deployintervalcount",
	"cloudlets:#.key.organization",
	"cloudlets:#.key.name",
	"minactiveinstances",
	"maxinstances",
	"undeployclientcount",
	"undeployintervalcount",
}
var AutoProvPolicyRequiredArgs = []string{
	"app-org",
	"name",
}
var AutoProvPolicyOptionalArgs = []string{
	"deployclientcount",
	"deployintervalcount",
	"cloudlets:empty",
	"cloudlets:#.key.organization",
	"cloudlets:#.key.name",
	"minactiveinstances",
	"maxinstances",
	"undeployclientcount",
	"undeployintervalcount",
}
var AutoProvPolicyAliasArgs = []string{
	"fields=autoprovpolicy.fields",
	"app-org=autoprovpolicy.key.organization",
	"name=autoprovpolicy.key.name",
	"deployclientcount=autoprovpolicy.deployclientcount",
	"deployintervalcount=autoprovpolicy.deployintervalcount",
	"cloudlets:empty=autoprovpolicy.cloudlets:empty",
	"cloudlets:#.key.organization=autoprovpolicy.cloudlets:#.key.organization",
	"cloudlets:#.key.name=autoprovpolicy.cloudlets:#.key.name",
	"cloudlets:#.loc.latitude=autoprovpolicy.cloudlets:#.loc.latitude",
	"cloudlets:#.loc.longitude=autoprovpolicy.cloudlets:#.loc.longitude",
	"cloudlets:#.loc.horizontalaccuracy=autoprovpolicy.cloudlets:#.loc.horizontalaccuracy",
	"cloudlets:#.loc.verticalaccuracy=autoprovpolicy.cloudlets:#.loc.verticalaccuracy",
	"cloudlets:#.loc.altitude=autoprovpolicy.cloudlets:#.loc.altitude",
	"cloudlets:#.loc.course=autoprovpolicy.cloudlets:#.loc.course",
	"cloudlets:#.loc.speed=autoprovpolicy.cloudlets:#.loc.speed",
	"cloudlets:#.loc.timestamp.seconds=autoprovpolicy.cloudlets:#.loc.timestamp.seconds",
	"cloudlets:#.loc.timestamp.nanos=autoprovpolicy.cloudlets:#.loc.timestamp.nanos",
	"minactiveinstances=autoprovpolicy.minactiveinstances",
	"maxinstances=autoprovpolicy.maxinstances",
	"undeployclientcount=autoprovpolicy.undeployclientcount",
	"undeployintervalcount=autoprovpolicy.undeployintervalcount",
}
var AutoProvPolicyComments = map[string]string{
	"fields":                             "Fields are used for the Update API to specify which fields to apply",
	"app-org":                            "Name of the organization for the cluster that this policy will apply to",
	"name":                               "Policy name",
	"deployclientcount":                  "Minimum number of clients within the auto deploy interval to trigger deployment",
	"deployintervalcount":                "Number of intervals to check before triggering deployment",
	"cloudlets:empty":                    "Allowed deployment locations, specify cloudlets:empty=true to clear",
	"cloudlets:#.key.organization":       "Organization of the cloudlet site",
	"cloudlets:#.key.name":               "Name of the cloudlet",
	"cloudlets:#.loc.latitude":           "Latitude in WGS 84 coordinates",
	"cloudlets:#.loc.longitude":          "Longitude in WGS 84 coordinates",
	"cloudlets:#.loc.horizontalaccuracy": "Horizontal accuracy (radius in meters)",
	"cloudlets:#.loc.verticalaccuracy":   "Vertical accuracy (meters)",
	"cloudlets:#.loc.altitude":           "On android only lat and long are guaranteed to be supplied Altitude in meters",
	"cloudlets:#.loc.course":             "Course (IOS) / bearing (Android) (degrees east relative to true north)",
	"cloudlets:#.loc.speed":              "Speed (IOS) / velocity (Android) (meters/sec)",
	"minactiveinstances":                 "Minimum number of active instances for High-Availability",
	"maxinstances":                       "Maximum number of instances (active or not)",
	"undeployclientcount":                "Number of active clients for the undeploy interval below which trigers undeployment, 0 (default) disables auto undeploy",
	"undeployintervalcount":              "Number of intervals to check before triggering undeployment",
}
var AutoProvPolicySpecialArgs = map[string]string{
	"autoprovpolicy.fields": "StringArray",
}
var AutoProvPolicyCloudletRequiredArgs = []string{
	"app-org",
	"name",
}
var AutoProvPolicyCloudletOptionalArgs = []string{
	"cloudlet-org",
	"cloudlet",
}
var AutoProvPolicyCloudletAliasArgs = []string{
	"app-org=autoprovpolicycloudlet.key.organization",
	"name=autoprovpolicycloudlet.key.name",
	"cloudlet-org=autoprovpolicycloudlet.cloudletkey.organization",
	"cloudlet=autoprovpolicycloudlet.cloudletkey.name",
}
var AutoProvPolicyCloudletComments = map[string]string{
	"app-org":      "Name of the organization for the cluster that this policy will apply to",
	"name":         "Policy name",
	"cloudlet-org": "Organization of the cloudlet site",
	"cloudlet":     "Name of the cloudlet",
}
var AutoProvPolicyCloudletSpecialArgs = map[string]string{}
