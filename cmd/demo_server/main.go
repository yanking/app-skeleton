// Package main 演示了如何使用 demo_server.yaml 配置文件
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yanking/app-skeleton/internal/config"
	"github.com/yanking/app-skeleton/pkg/app"
	"github.com/yanking/app-skeleton/pkg/grpc"
)

func main() {
	// 初始化配置
	if err := config.Init("configs/demo_server.yaml"); err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	cfg := config.Get()

	fmt.Println(cfg)
	
	var components []app.IComponent

	rpcServer := grpc.NewServer(
		grpc.WithAddress(cfg.Grpc.Addr),
		grpc.WithMetrics(true),
		grpc.WithTimeout(time.Second*5),
	)
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
