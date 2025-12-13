package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	devicesv1 "github.com/yaninyzwitty/golang-monorepo-md/gen/devices/v1"
	"github.com/yaninyzwitty/golang-monorepo-md/packages/shared/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			if !strings.Contains(err.Error(), "The handle is invalid") {
				logger.Error("failed to flush logs", zap.Error(err))
			}
		}
	}()

	// important for kubernates
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	var cfg config.Config
	if err := cfg.Load(logger, *configPath); err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// gRPC server address
	// devicesAddr := fmt.Sprintf(":%d", cfg.DevicesPort)

	// TODO-CONSIDER -- KUBERNATES but above still work
	devicesAddr := fmt.Sprintf(
		"device-service.testing:%d",
		cfg.DevicesPort,
	)

	// Connect to gRPC server
	grpcConn, err := grpc.NewClient(
		devicesAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("failed to connect to grpc server", zap.Error(err))
	}
	defer grpcConn.Close()

	// Create the gRPC client
	devicesClient := devicesv1.NewCloudServiceClient(grpcConn)

	// HTTP mux
	mux := http.NewServeMux()

	// GET /devices → gRPC GetDevices
	mux.HandleFunc("GET /devices", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		resp, err := devicesClient.GetDevices(ctx, &devicesv1.GetDevicesRequest{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get devices: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp.Devices)
	})

	mux.HandleFunc("POST /devices/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Parse JSON body
		var body struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		// Call gRPC
		resp, err := devicesClient.CreateDevice(r.Context(), &devicesv1.CreateDeviceRequest{
			Device: &devicesv1.Device{
				Name: body.Name,
				Type: body.Type,
			},
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create device: %v", err), http.StatusInternalServerError)
			return
		}

		// Return created device
		json.NewEncoder(w).Encode(resp.Device)
	})

	// Start HTTP server
	httpAddr := fmt.Sprintf(":%d", cfg.GatewayPort)
	logger.Info("HTTP server started", zap.String("address", httpAddr))

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		logger.Fatal("failed to start http server", zap.Error(err))
	}
}
