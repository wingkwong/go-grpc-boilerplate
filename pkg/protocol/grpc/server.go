package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/wingkwong/go-grpc-boilerplate/pkg/api/v1"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, v1API v1.FooServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterFooServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("[Info] Shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	log.Println("[Info] Starting gRPC server...")
	return server.Serve(listen)
}
