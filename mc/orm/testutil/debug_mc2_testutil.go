// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: debug.proto

package testutil

import edgeproto "github.com/mobiledgex/edge-cloud/edgeproto"
import "github.com/mobiledgex/edge-cloud-infra/mc/ormclient"
import "github.com/mobiledgex/edge-cloud-infra/mc/ormapi"
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

func TestEnableDebugLevels(mcClient *ormclient.Client, uri, token, region string, in *edgeproto.DebugRequest) ([]edgeproto.DebugReply, int, error) {
	dat := &ormapi.RegionDebugRequest{}
	dat.Region = region
	dat.DebugRequest = *in
	return mcClient.EnableDebugLevels(uri, token, dat)
}
func TestPermEnableDebugLevels(mcClient *ormclient.Client, uri, token, region, org string) ([]edgeproto.DebugReply, int, error) {
	in := &edgeproto.DebugRequest{}
	return TestEnableDebugLevels(mcClient, uri, token, region, in)
}

func TestDisableDebugLevels(mcClient *ormclient.Client, uri, token, region string, in *edgeproto.DebugRequest) ([]edgeproto.DebugReply, int, error) {
	dat := &ormapi.RegionDebugRequest{}
	dat.Region = region
	dat.DebugRequest = *in
	return mcClient.DisableDebugLevels(uri, token, dat)
}
func TestPermDisableDebugLevels(mcClient *ormclient.Client, uri, token, region, org string) ([]edgeproto.DebugReply, int, error) {
	in := &edgeproto.DebugRequest{}
	return TestDisableDebugLevels(mcClient, uri, token, region, in)
}

func TestShowDebugLevels(mcClient *ormclient.Client, uri, token, region string, in *edgeproto.DebugRequest) ([]edgeproto.DebugReply, int, error) {
	dat := &ormapi.RegionDebugRequest{}
	dat.Region = region
	dat.DebugRequest = *in
	return mcClient.ShowDebugLevels(uri, token, dat)
}
func TestPermShowDebugLevels(mcClient *ormclient.Client, uri, token, region, org string) ([]edgeproto.DebugReply, int, error) {
	in := &edgeproto.DebugRequest{}
	return TestShowDebugLevels(mcClient, uri, token, region, in)
}

func TestRunDebug(mcClient *ormclient.Client, uri, token, region string, in *edgeproto.DebugRequest) ([]edgeproto.DebugReply, int, error) {
	dat := &ormapi.RegionDebugRequest{}
	dat.Region = region
	dat.DebugRequest = *in
	return mcClient.RunDebug(uri, token, dat)
}
func TestPermRunDebug(mcClient *ormclient.Client, uri, token, region, org string) ([]edgeproto.DebugReply, int, error) {
	in := &edgeproto.DebugRequest{}
	return TestRunDebug(mcClient, uri, token, region, in)
}