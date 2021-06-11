package server

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ori-edge/grpc-interceptor-demo/pkg/api"
)

type EdgeLocationsServer struct {
	api.UnimplementedEdgeLocationsServer

	LocationStore map[string]EdgeLocation
}

type EdgeLocation struct {
	Id              string
	IpAddress       string
	OperatingSystem string
	UpdatedAt       time.Time
}

func hydrateType(el *api.EdgeLocation) EdgeLocation {
	return EdgeLocation{
		Id:              el.Id,
		IpAddress:       el.IpAddress,
		OperatingSystem: el.OperatingSystem,
		UpdatedAt:       el.UpdatedAt.AsTime(),
	}
}

func hydrateResponse(el EdgeLocation) *api.EdgeLocation {
	return &api.EdgeLocation{
		Id:              el.Id,
		IpAddress:       el.IpAddress,
		OperatingSystem: el.OperatingSystem,
		UpdatedAt:       timestamppb.New(el.UpdatedAt),
	}
}
