package get

import (
    "context"
    sonicpb "github.com/TeamNocsys/sonicpb/api/protobuf/sonic"
    "github.com/openconfig/gnmi/proto/gnmi"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "gnmi_server/pkg/gnmi/handler"
    handler_utils "gnmi_server/pkg/gnmi/handler/utils"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func GlobalIpRouteHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.State()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    kvs := handler.FetchPathKey(r)
    spec := "*"
    if v, ok := kvs["ip-prefix"]; ok {
        spec = v
    }

    sr := &sonicpb.SonicRoute{
        Route: &sonicpb.SonicRoute_Route{},
    }
    if hkeys, err := conn.GetKeys(swsssdk.APPL_DB, []string{"ROUTE_TABLE", spec}); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(swsssdk.APPL_DB, hkey)
            if len(keys) != 1 {
                continue
            }
            c := helper.GlobalIpRoute{
                Key: keys[0],
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sr.Route.GlobalRouteList = append(sr.Route.GlobalRouteList,
                &sonicpb.SonicRoute_Route_GlobalRouteListKey{
                    IpPrefix: keys[0],
                    GlobalRouteList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sr)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}


func IpRouteHandler(ctx context.Context, r *gnmi.GetRequest, db command.Client) (*gnmi.GetResponse, error) {
    conn := db.State()
    if conn == nil {
        return nil, status.Error(codes.Internal, "")
    }
    kvs := handler.FetchPathKey(r)
    spec := []string{}
    if v, ok := kvs["vrf-name"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }
    if v, ok := kvs["ip-prefix"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    sr := &sonicpb.SonicRoute{
        Route: &sonicpb.SonicRoute_Route{},
    }
    if hkeys, err := conn.GetKeys("ROUTE_TABLE", spec); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    } else {
        for _, hkey := range hkeys {
            keys := conn.SplitKeys(swsssdk.APPL_DB, hkey)
            if len(keys) != 2 {
                continue
            }
            c := helper.IpRoute{
                Keys: keys,
                Client: db,
                Data: nil,
            }
            if err := c.LoadFromDB(); err != nil {
                return nil, status.Errorf(codes.Internal, err.Error())
            }
            sr.Route.RouteList = append(sr.Route.RouteList,
                &sonicpb.SonicRoute_Route_RouteListKey{
                    VrfName: keys[0],
                    IpPrefix: keys[1],
                    RouteList: c.Data,
                })
        }
    }

    response, err := handler_utils.CreateGetResponse(ctx, r, sr)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return response, nil
}


