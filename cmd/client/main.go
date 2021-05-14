package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/ori-edge/grpc-interceptor-demo/pkg/api"
)

const (
	Register = "register"
	Get      = "get"
	List     = "list"
)

func main() {
	var idFlag string
	var limitFlag int

	flag.StringVar(&idFlag, "id", "", "the id of an edge location to lookup")
	flag.IntVar(&limitFlag, "limit", 10, "the number of edge locations to return from a stream")
	flag.Parse()

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
	case Get:
		if idFlag == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
		location, err := client.Get(context.Background(), &api.EdgeLocation{Id: idFlag})
		if err != nil {
			log.Fatal(err)
		}

		log.Println(location)
		return
	case List:
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
	default:
		log.Fatal("argument not supported")
		return
	}
}
