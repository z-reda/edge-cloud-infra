// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: app.proto

package orm

import edgeproto "github.com/mobiledgex/edge-cloud/edgeproto"
import "github.com/labstack/echo"
import "net/http"
import "context"
import "io"
import "github.com/mobiledgex/edge-cloud/log"
import "github.com/mobiledgex/edge-cloud-infra/mc/ormapi"
import "google.golang.org/grpc/status"
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

func CreateApp(c echo.Context) error {
	ctx := GetContext(c)
	rc := &RegionContext{}
	claims, err := getClaims(c)
	if err != nil {
		return err
	}
	rc.username = claims.Username

	in := ormapi.RegionApp{}
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, Msg("Invalid POST data"))
	}
	rc.region = in.Region
	span := log.SpanFromContext(ctx)
	span.SetTag("org", in.App.Key.DeveloperKey.Name)
	resp, err := CreateAppObj(ctx, rc, &in.App)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			err = fmt.Errorf("%s", st.Message())
		}
	}
	return setReply(c, err, resp)
}

func CreateAppObj(ctx context.Context, rc *RegionContext, obj *edgeproto.App) (*edgeproto.Result, error) {
	if !rc.skipAuthz {
		if err := authzCreateApp(ctx, rc.region, rc.username, obj,
			ResourceApps, ActionManage); err != nil {
			return nil, err
		}
	}
	if rc.conn == nil {
		conn, err := connectController(ctx, rc.region)
		if err != nil {
			return nil, err
		}
		rc.conn = conn
		defer func() {
			rc.conn.Close()
			rc.conn = nil
		}()
	}
	api := edgeproto.NewAppApiClient(rc.conn)
	return api.CreateApp(ctx, obj)
}

func DeleteApp(c echo.Context) error {
	ctx := GetContext(c)
	rc := &RegionContext{}
	claims, err := getClaims(c)
	if err != nil {
		return err
	}
	rc.username = claims.Username

	in := ormapi.RegionApp{}
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, Msg("Invalid POST data"))
	}
	rc.region = in.Region
	span := log.SpanFromContext(ctx)
	span.SetTag("org", in.App.Key.DeveloperKey.Name)
	resp, err := DeleteAppObj(ctx, rc, &in.App)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			err = fmt.Errorf("%s", st.Message())
		}
	}
	return setReply(c, err, resp)
}

func DeleteAppObj(ctx context.Context, rc *RegionContext, obj *edgeproto.App) (*edgeproto.Result, error) {
	if !rc.skipAuthz && !authorized(ctx, rc.username, obj.Key.DeveloperKey.Name,
		ResourceApps, ActionManage) {
		return nil, echo.ErrForbidden
	}
	if rc.conn == nil {
		conn, err := connectController(ctx, rc.region)
		if err != nil {
			return nil, err
		}
		rc.conn = conn
		defer func() {
			rc.conn.Close()
			rc.conn = nil
		}()
	}
	api := edgeproto.NewAppApiClient(rc.conn)
	return api.DeleteApp(ctx, obj)
}

func UpdateApp(c echo.Context) error {
	ctx := GetContext(c)
	rc := &RegionContext{}
	claims, err := getClaims(c)
	if err != nil {
		return err
	}
	rc.username = claims.Username

	in := ormapi.RegionApp{}
	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, Msg("Invalid POST data"))
	}
	rc.region = in.Region
	span := log.SpanFromContext(ctx)
	span.SetTag("org", in.App.Key.DeveloperKey.Name)
	resp, err := UpdateAppObj(ctx, rc, &in.App)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			err = fmt.Errorf("%s", st.Message())
		}
	}
	return setReply(c, err, resp)
}

func UpdateAppObj(ctx context.Context, rc *RegionContext, obj *edgeproto.App) (*edgeproto.Result, error) {
	if !rc.skipAuthz {
		if err := authzUpdateApp(ctx, rc.region, rc.username, obj,
			ResourceApps, ActionManage); err != nil {
			return nil, err
		}
	}
	if rc.conn == nil {
		conn, err := connectController(ctx, rc.region)
		if err != nil {
			return nil, err
		}
		rc.conn = conn
		defer func() {
			rc.conn.Close()
			rc.conn = nil
		}()
	}
	api := edgeproto.NewAppApiClient(rc.conn)
	return api.UpdateApp(ctx, obj)
}

func ShowApp(c echo.Context) error {
	ctx := GetContext(c)
	rc := &RegionContext{}
	claims, err := getClaims(c)
	if err != nil {
		return err
	}
	rc.username = claims.Username

	in := ormapi.RegionApp{}
	success, err := ReadConn(c, &in)
	if !success {
		return err
	}
	defer CloseConn(c)
	rc.region = in.Region
	span := log.SpanFromContext(ctx)
	span.SetTag("org", in.App.Key.DeveloperKey.Name)

	err = ShowAppStream(ctx, rc, &in.App, func(res *edgeproto.App) {
		payload := ormapi.StreamPayload{}
		payload.Data = res
		WriteStream(c, &payload)
	})
	if err != nil {
		WriteError(c, err)
	}
	return nil
}

func ShowAppStream(ctx context.Context, rc *RegionContext, obj *edgeproto.App, cb func(res *edgeproto.App)) error {
	var authz *ShowAuthz
	var err error
	if !rc.skipAuthz {
		authz, err = NewShowAuthz(ctx, rc.region, rc.username, ResourceApps, ActionView)
		if err == echo.ErrForbidden {
			return nil
		}
		if err != nil {
			return err
		}
	}
	if rc.conn == nil {
		conn, err := connectController(ctx, rc.region)
		if err != nil {
			return err
		}
		rc.conn = conn
		defer func() {
			rc.conn.Close()
			rc.conn = nil
		}()
	}
	api := edgeproto.NewAppApiClient(rc.conn)
	stream, err := api.ShowApp(ctx, obj)
	if err != nil {
		return err
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return err
		}
		if !rc.skipAuthz {
			if !authz.Ok(res.Key.DeveloperKey.Name) {
				continue
			}
		}
		cb(res)
	}
	return nil
}

func ShowAppObj(ctx context.Context, rc *RegionContext, obj *edgeproto.App) ([]edgeproto.App, error) {
	arr := []edgeproto.App{}
	err := ShowAppStream(ctx, rc, obj, func(res *edgeproto.App) {
		arr = append(arr, *res)
	})
	return arr, err
}