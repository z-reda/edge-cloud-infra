// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: flavor.proto

package ormclient

import edgeproto "github.com/mobiledgex/edge-cloud/edgeproto"
import "github.com/mobiledgex/edge-cloud-infra/mc/ormapi"
import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/googleapis/google/api"
import _ "github.com/mobiledgex/edge-cloud/protogen"
import _ "github.com/mobiledgex/edge-cloud/protoc-gen-cmd/protocmd"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Auto-generated code: DO NOT EDIT

func (s *Client) CreateFlavor(uri, token string, in *ormapi.RegionFlavor) (edgeproto.Result, int, error) {
	out := edgeproto.Result{}
	status, err := s.PostJson(uri+"/auth/ctrl/CreateFlavor", token, in, &out)
	return out, status, err
}

func (s *Client) DeleteFlavor(uri, token string, in *ormapi.RegionFlavor) (edgeproto.Result, int, error) {
	out := edgeproto.Result{}
	status, err := s.PostJson(uri+"/auth/ctrl/DeleteFlavor", token, in, &out)
	return out, status, err
}

func (s *Client) UpdateFlavor(uri, token string, in *ormapi.RegionFlavor) (edgeproto.Result, int, error) {
	out := edgeproto.Result{}
	status, err := s.PostJson(uri+"/auth/ctrl/UpdateFlavor", token, in, &out)
	return out, status, err
}

func (s *Client) ShowFlavor(uri, token string, in *ormapi.RegionFlavor) ([]edgeproto.Flavor, int, error) {
	out := edgeproto.Flavor{}
	outlist := []edgeproto.Flavor{}
	status, err := s.PostJsonStreamOut(uri+"/auth/ctrl/ShowFlavor", token, in, &out, func() {
		outlist = append(outlist, out)
	})
	return outlist, status, err
}

type FlavorApiClient interface {
	CreateFlavor(uri, token string, in *ormapi.RegionFlavor) (edgeproto.Result, int, error)
	DeleteFlavor(uri, token string, in *ormapi.RegionFlavor) (edgeproto.Result, int, error)
	UpdateFlavor(uri, token string, in *ormapi.RegionFlavor) (edgeproto.Result, int, error)
	ShowFlavor(uri, token string, in *ormapi.RegionFlavor) ([]edgeproto.Flavor, int, error)
}