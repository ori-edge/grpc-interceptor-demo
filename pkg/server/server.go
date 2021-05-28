package server

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/ori-edge/grpc-interceptor-demo/pkg/api"
)

// New is a constructor function that creates our new server with an in-memory
// store for edge locations. This means that when the server finishes or exits
// unexpectedly, we lose our edge locations. In a real environment, this store
// would likely be a database of some kind
func New() *EdgeLocationsServer {
	return &EdgeLocationsServer{
		LocationStore: make(map[string]EdgeLocation),
	}
}

// Register takes a client supplied edge location and stores it in-memory
func (s EdgeLocationsServer) Register(ctx context.Context, el *api.EdgeLocation) (*empty.Empty, error) {
	s.LocationStore[el.Id] = hydrateType(el)
	log.Printf("registering client success")

	return &empty.Empty{}, nil
}

// Get retrieves a location from the in-memory store and returns it to the client
func (s EdgeLocationsServer) Get(ctx context.Context, el *api.EdgeLocation) (*api.EdgeLocation, error) {
	log.Println("retrieving edge location...")

	return hydrateResponse(s.LocationStore[el.Id]), nil
}

// List iterates through all items within the in-memory store and streams them
// one by one, back to the client. This function finishes either when the client
// supplied limit is hit - or there are no more objects in the store. The
// default limit supplied by the client is 10
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
