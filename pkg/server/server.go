package server

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/ori-edge/grpc-interceptor-demo/pkg/api"
)

func New() *EdgeLocationsServer {
	return &EdgeLocationsServer{
		LocationStore: make(map[string]EdgeLocation),
	}
}

func (s EdgeLocationsServer) Register(ctx context.Context, el *api.EdgeLocation) (*empty.Empty, error) {
	s.LocationStore[el.Id] = hydrateType(el)
	log.Printf("registering client success")
	return &empty.Empty{}, nil
}

func (s EdgeLocationsServer) Get(ctx context.Context, el *api.EdgeLocation) (*api.EdgeLocation, error) {
	log.Println("retrieving edge location...")
	return hydrateResponse(s.LocationStore[el.Id]), nil
}

func (s EdgeLocationsServer) List(param *api.ListEdgeLocationParams, stream api.EdgeLocations_ListServer) error {
	log.Println("streaming edge locations...")

	i := 0
	for _, location := range s.LocationStore {
		if i == int(param.Limit) {
			log.Println("limit reached - stopping stream")
			return nil
		}
		if err := stream.Send(hydrateResponse(location)); err != nil {
			return err
		}
		i++
	}

	return nil
}
