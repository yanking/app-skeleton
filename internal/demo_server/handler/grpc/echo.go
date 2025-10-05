package grpc

import (
	"context"

	v1 "github.com/yanking/app-skeleton/api/proto/gen/demo/v1"
)

func (h *Handler) Echo(ctx context.Context, req *v1.EchoRequest) (*v1.EchoResponse, error) {
	return &v1.EchoResponse{
		Value: req.Value,
	}, nil
}
