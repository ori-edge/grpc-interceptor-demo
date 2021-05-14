package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ori-edge/grpc-interceptor-demo/pkg/api"
	"github.com/ori-edge/grpc-interceptor-demo/pkg/server"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:5565")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterEdgeLocationsServer(s, server.New())
	log.Println("starting server...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
