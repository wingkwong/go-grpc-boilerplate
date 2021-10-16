package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	v1 "github.com/wingkwong/go-grpc-boilerplate/pkg/api/v1"
	"github.com/wingkwong/go-grpc-boilerplate/pkg/logger"
	"github.com/wingkwong/go-grpc-boilerplate/pkg/protocol/grpc/middleware"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, v1API v1.FooServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	opts = middleware.AddLogging(logger.Log, opts)

	server := grpc.NewServer(opts...)
	v1.RegisterFooServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logger.Log.Warn("Shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	logger.Log.Info("Starting gRPC server...")
	return server.Serve(listen)
}
