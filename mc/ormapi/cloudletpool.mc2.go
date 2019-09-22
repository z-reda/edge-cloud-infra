// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cloudletpool.proto

package ormapi

import edgeproto "github.com/mobiledgex/edge-cloud/edgeproto"
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

type RegionCloudletPool struct {
	Region       string
	CloudletPool edgeproto.CloudletPool
}

type RegionCloudletPoolMember struct {
	Region             string
	CloudletPoolMember edgeproto.CloudletPoolMember
}

type RegionCloudletKey struct {
	Region      string
	CloudletKey edgeproto.CloudletKey
}

type RegionCloudletPoolKey struct {
	Region          string
	CloudletPoolKey edgeproto.CloudletPoolKey
}

type RegionCloudletPoolList struct {
	Region           string
	CloudletPoolList edgeproto.CloudletPoolList
}