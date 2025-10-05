package grpc

import "github.com/yanking/app-skeleton/api/proto/gen/demo/v1"

type Handler struct {
	v1.UnimplementedDemoServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}
