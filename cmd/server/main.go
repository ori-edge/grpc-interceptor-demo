package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ori-edge/grpc-interceptor-demo/pkg/api"
	"github.com/ori-edge/grpc-interceptor-demo/pkg/interceptor"
	"github.com/ori-edge/grpc-interceptor-demo/pkg/server"
)

func main() {
	// The server is set up to listen on localhost, on port 5565
	lis, err := net.Listen("tcp", "localhost:5565")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()),
		grpc.StreamInterceptor(interceptor.StreamServerInterceptor()),
	)

	// Register our server
	api.RegisterEdgeLocationsServer(s, server.New())
	log.Println("starting server...")

	// Start serving!
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
