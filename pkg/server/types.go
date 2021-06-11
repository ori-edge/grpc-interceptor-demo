package server

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ori-edge/grpc-interceptor-demo/pkg/api"
)

// EdgeLocationsServer is a struct that contains our edge location store and
// represents our server
type EdgeLocationsServer struct {
	api.UnimplementedEdgeLocationsServer

	LocationStore map[string]EdgeLocation
}

// EdgeLocation is a struct that represents the structure and data held in
// memory for an EdgeLocation
type EdgeLocation struct {
	Id              string
	IpAddress       string
	OperatingSystem string
	UpdatedAt       time.Time
}

// hydrateType takes a gRPC message representation of an EdgeLocation and
// converts it to our EdgeLocation type. We need to do this due to a mutex on
// the original object not allowing it to be manipulated
func hydrateType(el *api.EdgeLocation) EdgeLocation {
	return EdgeLocation{
		Id:              el.Id,
		IpAddress:       el.IpAddress,
		OperatingSystem: el.OperatingSystem,
		UpdatedAt:       el.UpdatedAt.AsTime(),
	}
}

// hydrateResponse takes our EdgeLocation type and converts it back into our
// gRPC defined message. This allows us to send it back to the client in a way
// that it understands
func hydrateResponse(el EdgeLocation) *api.EdgeLocation {
	return &api.EdgeLocation{
		Id:              el.Id,
		IpAddress:       el.IpAddress,
		OperatingSystem: el.OperatingSystem,
		UpdatedAt:       timestamppb.New(el.UpdatedAt),
	}
}
