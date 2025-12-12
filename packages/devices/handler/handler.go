package handler

import (
	"context"
	"fmt"
	"sync"

	devicesv1 "github.com/yaninyzwitty/golang-monorepo-md/gen/devices/v1"
)

// devicesServiceHandler implements devicesv1.CloudServiceServer
type devicesServiceHandler struct {
	devicesv1.UnimplementedCloudServiceServer
	mu      sync.Mutex
	devices []*devicesv1.Device
}

// NewDevicesServiceHandler creates a new instance
func NewDevicesServiceHandler() *devicesServiceHandler {
	return &devicesServiceHandler{
		devices: []*devicesv1.Device{
			{Id: "1", Name: "Thermostat", Type: "Sensor"},
			{Id: "2", Name: "Smart Light", Type: "Actuator"},
			{Id: "3", Name: "Security Camera", Type: "Sensor"},
		},
	}
}

// GetDevices returns all devices
func (s *devicesServiceHandler) GetDevices(ctx context.Context, req *devicesv1.GetDevicesRequest) (*devicesv1.GetDevicesResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &devicesv1.GetDevicesResponse{
		Devices: s.devices,
	}, nil
}

// CreateDevice adds a new device
func (s *devicesServiceHandler) CreateDevice(ctx context.Context, req *devicesv1.CreateDeviceRequest) (*devicesv1.CreateDeviceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate a simple string ID (could use UUID in production)
	newID := len(s.devices) + 1
	device := &devicesv1.Device{
		Id:   fmt.Sprintf("%d", newID), // convert int to string
		Name: req.Device.Name,
		Type: req.Device.Type,
	}

	s.devices = append(s.devices, device)

	return &devicesv1.CreateDeviceResponse{
		Device: device,
	}, nil
}
