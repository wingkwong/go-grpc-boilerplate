package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wingkwong/go-grpc-boilerplate/pkg/logger"
	"github.com/wingkwong/go-grpc-boilerplate/pkg/protocol/grpc"
	"github.com/wingkwong/go-grpc-boilerplate/pkg/protocol/rest"
	v1 "github.com/wingkwong/go-grpc-boilerplate/pkg/service/v1"
)

type Config struct {
	GRPCPort            string
	HTTPPort            string
	DatastoreDBHost     string
	DatastoreDBUser     string
	DatastoreDBPassword string
	DatastoreDBSchema   string
	LogLevel            int
	LogTimeFormat       string
}

func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "HTTP port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
	flag.IntVar(&cfg.LogLevel, "log-level", 0, "Global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "", "Print time format for logger e.g. 2019-07-21T23:20:00Z08:00")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("[ERROR] Invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("[ERROR] Invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("[ERROR] Failed to initialize logger: %v", err)
	}

	param := "parseTime=true"

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBSchema,
		param,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to open database: %v", err)
	}

	defer db.Close()

	v1API := v1.NewFooServiceServer(db)

	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
