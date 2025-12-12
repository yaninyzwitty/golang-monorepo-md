package main

import (
	"fmt"
	"net"
	"strings"

	devicesv1 "github.com/yaninyzwitty/golang-monorepo-md/gen/devices/v1"
	"github.com/yaninyzwitty/golang-monorepo-md/packages/devices/handler"
	"github.com/yaninyzwitty/golang-monorepo-md/packages/shared/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			// ignore the "handle is invalid" error on Windows
			if !strings.Contains(err.Error(), "The handle is invalid") {
				logger.Error("failed to flush logs", zap.Error(err))
			}
		}
	}()

	var cfg config.Config
	if err := cfg.Load(logger, "config.yaml"); err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	devicesAddr := fmt.Sprintf(":%d", cfg.DevicesPort)

	lis, err := net.Listen("tcp", devicesAddr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()

	// Create handler and register it
	devicesHandler := handler.NewDevicesServiceHandler()
	devicesv1.RegisterCloudServiceServer(grpcServer, devicesHandler)

	logger.Info("Devices gRPC server started", zap.String("address", devicesAddr))

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve gRPC", zap.Error(err))
	}
}
