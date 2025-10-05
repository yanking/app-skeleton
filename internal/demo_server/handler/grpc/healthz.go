package grpc

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	v1 "github.com/yanking/app-skeleton/api/proto/gen/demo/v1"
)

func (h *Handler) Healthz(ctx context.Context, req *empty.Empty) (*v1.HealthzResponse, error) {
	return &v1.HealthzResponse{
		Status:    v1.ServiceStatus_HEALTHY,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   "",
	}, nil
}
