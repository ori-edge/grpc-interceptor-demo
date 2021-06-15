package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ori-edge/grpc-interceptor-demo/pkg/api"
	"github.com/ori-edge/grpc-interceptor-demo/pkg/interceptor"
)

const (
	// Register is the constant for the sub-command to register an edge location with the server
	Register = "register"
	List     = "list"
)

func main() {
	var listFlag string

	fs := flag.NewFlagSet("flagset", flag.ContinueOnError)
	fs.StringVar(&listFlag, "list", "", "a comma seperated list of ids to fetch from the server")

	// Parse flags from the second onwards - the first is the sub-command addressed lower
	err := fs.Parse(os.Args[2:])
	if err != nil {
		log.Fatalf("couldn't parse flags: %v", err)
		return
	}

	// Dial the gRPC server in insecure mode (you should use SSL in production)
	conn, err := grpc.Dial(
		"localhost:5565",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(interceptor.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("couldn't connect to server: %v", err)
		return
	}

	// Make sure the connection is closed when we are done
	defer conn.Close()

	// Create a client that can connect to the server
	client := api.NewEdgeLocationsClient(conn)

	// Switch through our sub-commands, register, get, list
	switch os.Args[1] {
	case Register:
		_, err = client.Register(context.Background(), &api.EdgeLocation{
			Id:        uuid.New().String(),
			UpdatedAt: timestamppb.New(time.Now()),
		})
		if err != nil {
			log.Fatal(err)
		}
		return
	case List:
		// Seperate the user flags into an array of ids to send to the server
		list := strings.Split(listFlag, ",")
		// Open a streaming connection to the server
		stream, err := client.List(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			// For every id provided, send them on the bidi stream
			for i := 0; i < len(list); i++ {
				el := api.EdgeLocation{
					Id: list[i],
				}
				if err := stream.Send(&el); err != nil {
					log.Fatalf("can not send %v", err)
				}
			}
			if err := stream.CloseSend(); err != nil {
				log.Println(err)
			}
		}()

		done := make(chan bool)
		go func() {
			for {
				locations, err := stream.Recv()
				if err == io.EOF {
					done <- true
					return
				}
				if err != nil {
					log.Fatal(err)
				}

				// Display the returned locations
				log.Println(locations)
			}
		}()
		<-done
		return
	// Fail gracefully when the user supplies an invalid sub-command
	default:
		log.Fatal("argument not supported")
		return
	}
}
