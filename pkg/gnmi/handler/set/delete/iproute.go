package delete

import (
    "context"
    "gnmi_server/cmd/command"
    "gnmi_server/internal/pkg/swsssdk/helper"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func GlobalIpRouteHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
    if v, ok := kvs["ip-prefix"]; !ok {
        return status.Error(codes.Internal, ErrNoKey)
    } else {
        c := &helper.GlobalIpRoute{
            Key:    v,
            Client: db,
            Data:   nil,
        }
        if err := c.RemoveFromDB(); err != nil {
            return err
        }
    }

    return nil
}

func IpRouteHandler(ctx context.Context, kvs map[string]string, db command.Client) error {
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
    if v, ok := kvs["nexthop"]; ok {
        spec = append(spec, v)
    } else {
        spec = append(spec, "*")
    }

    c := &helper.IpRoute{
        Keys:   spec,
        Client: db,
        Data:   nil,
    }

    if err := c.RemoveFromDB(); err != nil {
        return err
    }
    return nil
}