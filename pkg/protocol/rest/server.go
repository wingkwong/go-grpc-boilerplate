package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/wingkwong/go-grpc-boilerplate/pkg/api/v1"
	"github.com/wingkwong/go-grpc-boilerplate/pkg/logger"
	"github.com/wingkwong/go-grpc-boilerplate/pkg/protocol/rest/middleware"
)

func RunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterFooServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		logger.Log.Fatal("Failed to start HTTP gateway", zap.String("reason", err.Error()))
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: middleware.AddRequestID(middleware.AddLogger(logger.Log, mux)),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// handle ^C
		}
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	logger.Log.Info("Starting HTTP/REST Gateway...")
	return srv.ListenAndServe()
}
