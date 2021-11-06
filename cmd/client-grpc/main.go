package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	v1 "github.com/wingkwong/go-grpc-boilerplate/pkg/api/v1"
)

const (
	apiVersion = "v1"
)

func main() {
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[Error] Failed to connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewFooServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create
	req1 := v1.CreateRequest{
		ApiVersion: apiVersion,
		Foo: &v1.Foo{
			Title: "Dummy Title",
			Desc:  "Dummy Description",
		},
	}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("[Error] Failed to create : %v", err)
	}
	log.Printf("[Info] Create result: <%+v>\n\n", res1)

	id := res1.Id

	// Read
	req2 := v1.ReadRequest{
		ApiVersion: apiVersion,
		Id:         id,
	}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatalf("[Error] Failed to read : %v", err)
	}
	log.Printf("[Info] Read result: <%+v>\n\n", res2)

	// ReadAll
	req3 := v1.ReadAllRequest{
		ApiVersion: apiVersion,
	}
	res3, err := c.ReadAll(ctx, &req3)
	if err != nil {
		log.Fatalf("[Error] Failed to read all : %v", err)
	}
	log.Printf("[INFO] ReadAll result: <%+v>\n\n", res3)

	// Update
	req4 := v1.UpdateRequest{
		ApiVersion: apiVersion,
		Foo: &v1.Foo{
			Id:    res2.Foo.Id,
			Title: res2.Foo.Title,
			Desc:  res2.Foo.Desc + " + (Updated)",
		},
	}
	res4, err := c.Update(ctx, &req4)
	if err != nil {
		log.Fatalf("[Error] Failed to update : %v", err)
	}
	log.Printf("[Info] Update result: <%+v>\n\n", res4)

	// Delete
	req5 := v1.DeleteRequest{
		ApiVersion: apiVersion,
		Id:         id,
	}
	res5, err := c.Delete(ctx, &req5)
	if err != nil {
		log.Fatalf("[Error] Failed to delete : %v", err)
	}
	log.Printf("[Info] Delete result: <%+v>\n\n", res5)
}
