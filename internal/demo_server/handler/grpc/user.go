package grpc

import (
	"context"

	v1 "github.com/yanking/app-skeleton/api/proto/gen/demo/v1"
)

type UserHandler struct {
	v1.UnimplementedUserServiceServer
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	// 模拟获取用户信息
	return &v1.GetUserResponse{
		Id:    req.Id,
		Name:  "User " + req.Id,
		Email: "user" + req.Id + "@example.com",
	}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	// 模拟获取用户列表
	users := []*v1.GetUserResponse{
		{
			Id:    "1",
			Name:  "User 1",
			Email: "user1@example.com",
		},
		{
			Id:    "2",
			Name:  "User 2",
			Email: "user2@example.com",
		},
	}

	return &v1.ListUsersResponse{
		Users: users,
	}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	// 模拟创建用户
	return &v1.CreateUserResponse{
		Id:    "3",
		Name:  req.Name,
		Email: req.Email,
	}, nil
}