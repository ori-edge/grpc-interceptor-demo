package server

import (
	"context"
	"fmt"
	"io"
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
