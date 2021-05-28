package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ori-edge/grpc-interceptor-demo/pkg/api"
)

const (
	// Register is the constant for the sub-command to register an edge location with the server
	Register = "register"
	// Get is the constant for the sub-command to get an individual edge location
	// already registered with the server
	Get = "get"
	// List allow users to list all registered edge locations, allows users to
	// limit this list through the use of the limit flag
	List = "list"
)

func main() {
	var idFlag string
	var limitFlag int
	var regionFlag string

	// Set up our flags
	fs := flag.NewFlagSet("flagset", flag.ContinueOnError)
	fs.StringVar(&idFlag, "id", "", "the id of an edge location to lookup")
	fs.StringVar(&regionFlag, "region", "undefined", "the region that the edge location resides in")
	fs.IntVar(&limitFlag, "limit", 10, "the number of edge locations to return from a stream")

	// Parse flags from the second onwards - the first is the sub-command addressed lower
	err := fs.Parse(os.Args[2:])
	if err != nil {
		log.Fatalf("couldn't parse flags: %v", err)
		return
	}

	// Dial the gRPC server in insecure mode (you should use SSL in production)
	conn, err := grpc.Dial("localhost:5565", grpc.WithInsecure())
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
		// Register a location with a user defined region
		_, err = client.Register(context.Background(), &api.EdgeLocation{
			Id:        uuid.New().String(),
			Region:    regionFlag,
			UpdatedAt: timestamppb.New(time.Now()),
		})
		if err != nil {
			log.Fatal(err)
		}
		return
	case Get:
		// Get an individual edge location based on it's ID, we require this
		// flag to progress
		if idFlag == "" {
			fs.PrintDefaults()
			os.Exit(1)
		}
		location, err := client.Get(context.Background(), &api.EdgeLocation{Id: idFlag})
		if err != nil {
			log.Fatal(err)
		}

		// Print the edge location to the client
		log.Println(location)
		return
	case List:
		// List based on the limit flag all of the registered edge locations
		stream, err := client.List(context.Background(), &api.ListEdgeLocationParams{Limit: int32(limitFlag)})
		if err != nil {
			log.Fatal(err)
		}

		done := make(chan bool)
		go func() {
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					done <- true
					return
				}
				if err != nil {
					log.Fatal(err)
				}

				log.Println(resp)
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
