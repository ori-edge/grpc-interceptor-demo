package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/ori-edge/grpc-interceptor-demo/pkg/api"
)

const (
	Register = "register"
	List     = "list"
)

func main() {
	var listFlag string

	fs := flag.NewFlagSet("flagset", flag.ContinueOnError)
	fs.StringVar(&listFlag, "list", "", "a comma seperated list of ids to fetch from the server")

	conn, err := grpc.Dial("localhost:5565", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("couldn't connect to server: %v", err)
		return
	}
	defer conn.Close()

	client := api.NewEdgeLocationsClient(conn)

	switch os.Args[1] {
	case Register:
		_, err = client.Register(context.Background(), &api.EdgeLocation{Id: uuid.New().String()})
		if err != nil {
			log.Fatal(err)
		}
		return
	case List:
		// Seperate the user flafs into an array of ids to send to the server
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
	default:
		log.Fatal("argument not supported")
		return
	}
}
