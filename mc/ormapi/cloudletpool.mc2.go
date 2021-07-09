// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cloudletpool.proto

package ormapi

import (
	fmt "fmt"
	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/mobiledgex/edge-cloud/d-match-engine/dme-proto"
	edgeproto "github.com/mobiledgex/edge-cloud/edgeproto"
	_ "github.com/mobiledgex/edge-cloud/protogen"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Auto-generated code: DO NOT EDIT

// Request summary for CreateCloudletPool
// swagger:parameters CreateCloudletPool
type swaggerCreateCloudletPool struct {
	// in: body
	Body RegionCloudletPool
}

type RegionCloudletPool struct {
	// required: true
	// Region name
	Region       string
	CloudletPool edgeproto.CloudletPool
}

func (s *RegionCloudletPool) GetRegion() string {
	return s.Region
}

func (s *RegionCloudletPool) GetObj() interface{} {
	return &s.CloudletPool
}

func (s *RegionCloudletPool) GetObjName() string {
	return "CloudletPool"
}
func (s *RegionCloudletPool) GetObjFields() []string {
	return s.CloudletPool.Fields
}

func (s *RegionCloudletPool) SetObjFields(fields []string) {
	s.CloudletPool.Fields = fields
}

// Request summary for DeleteCloudletPool
// swagger:parameters DeleteCloudletPool
type swaggerDeleteCloudletPool struct {
	// in: body
	Body RegionCloudletPool
}

// Request summary for UpdateCloudletPool
// swagger:parameters UpdateCloudletPool
type swaggerUpdateCloudletPool struct {
	// in: body
	Body RegionCloudletPool
}

// Request summary for ShowCloudletPool
// swagger:parameters ShowCloudletPool
type swaggerShowCloudletPool struct {
	// in: body
	Body RegionCloudletPool
}

// Request summary for AddCloudletPoolMember
// swagger:parameters AddCloudletPoolMember
type swaggerAddCloudletPoolMember struct {
	// in: body
	Body RegionCloudletPoolMember
}

type RegionCloudletPoolMember struct {
	// required: true
	// Region name
	Region             string
	CloudletPoolMember edgeproto.CloudletPoolMember
}

func (s *RegionCloudletPoolMember) GetRegion() string {
	return s.Region
}

func (s *RegionCloudletPoolMember) GetObj() interface{} {
	return &s.CloudletPoolMember
}

func (s *RegionCloudletPoolMember) GetObjName() string {
	return "CloudletPoolMember"
}

// Request summary for RemoveCloudletPoolMember
// swagger:parameters RemoveCloudletPoolMember
type swaggerRemoveCloudletPoolMember struct {
	// in: body
	Body RegionCloudletPoolMember
}
