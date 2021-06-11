package server

import (
	"context"
	"fmt"
	"io"
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
	log.Printf("registering client success, id: %v", el.Id)
	return &empty.Empty{}, nil
}

func (s EdgeLocationsServer) List(stream api.EdgeLocations_ListServer) error {
	log.Println("streaming edge locations...")
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		el, err := stream.Recv()
		if err == io.EOF {
			log.Println("end of stream")
			return nil
		}
		if err != nil {
			return fmt.Errorf("receive error %v", err)
		}

		for _, location := range s.LocationStore {
			if location.Id == el.Id {
				if err := stream.Send(hydrateResponse(location)); err != nil {
					return err
				}
			}
		}
	}
}
