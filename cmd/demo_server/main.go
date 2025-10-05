// Package main 演示了如何使用 demo_server.yaml 配置文件
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/yanking/app-skeleton/api/proto/gen/demo/v1"
	"github.com/yanking/app-skeleton/internal/config"
	"github.com/yanking/app-skeleton/pkg/app"
	pkgGrpc "github.com/yanking/app-skeleton/pkg/grpc"
	"github.com/yanking/app-skeleton/pkg/log"
	"github.com/yanking/app-skeleton/pkg/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	demoHandler "github.com/yanking/app-skeleton/internal/demo_server/handler/grpc"
)

func main() {
	// 初始化配置
	if err := config.Init("configs/demo_server.yaml"); err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}
	cfg := config.Get()

	log.Init(cfg.Log)

	fmt.Println(cfg)

	var components []app.IComponent

	// 创建 gRPC 服务器（带 gRPC-Gateway）
	grpcOptions := []pkgGrpc.ServerOption{
		pkgGrpc.WithAddress(cfg.Grpc.Addr),
		pkgGrpc.WithTimeout(time.Second * 5),
		pkgGrpc.WithGateway(cfg.Grpc.Gateway.Enabled, cfg.HTTP.Addr), // 根据配置启用 gRPC-Gateway
	}
	if cfg.EnableMetrics {
		metric.Handler()
		grpcOptions = append(grpcOptions, pkgGrpc.WithMetrics(true))
	}
	rpcServer := pkgGrpc.NewServer(grpcOptions...)
	
	// 注册 gRPC 服务
	v1.RegisterDemoServiceServer(rpcServer.Server, demoHandler.NewHandler())
	userHandler := demoHandler.NewUserHandler()
	v1.RegisterUserServiceServer(rpcServer.Server, userHandler)

	// 注册 gRPC-Gateway 处理程序（直接使用生成的函数）
	if cfg.Grpc.Gateway.Enabled {
		ctx := context.Background()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		
		// 注册DemoService
		if err := v1.RegisterDemoServiceHandlerFromEndpoint(ctx, rpcServer.GetGatewayMux(), cfg.Grpc.Addr, opts); err != nil {
			log.Fatalf("failed to register demo service gateway: %v", err)
		}
		
		// 注册UserService
		if err := v1.RegisterUserServiceHandlerFromEndpoint(ctx, rpcServer.GetGatewayMux(), cfg.Grpc.Addr, opts); err != nil {
			log.Fatalf("failed to register user service gateway: %v", err)
		}
	}

	components = append(components, rpcServer)

	a, err := app.New(cfg.AppName, components)
	if err != nil {
		panic(err)
	}

	err = a.Run()
	if err != nil {
		panic(err)
	}
}