// Package main 演示了如何使用 demo_server.yaml 配置文件
package main

import (
	"fmt"
	"log"

	"github.com/yanking/app-skeleton/internal/config"
)

func main() {
	// 初始化配置
	if err := config.Init("configs/demo_server.yaml"); err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	// 获取配置实例
	cfg := config.Get()

	// 打印部分配置项验证
	fmt.Printf("App Name: %s\n", cfg.AppName)
	fmt.Printf("Server Mode: %s\n", cfg.ServerMode)
	fmt.Printf("JWT Key: %s\n", cfg.JwtKey)
	fmt.Printf("Expiration: %s\n", cfg.Expiration)
	fmt.Printf("HTTP Addr: %s\n", cfg.HTTP.Addr)
	fmt.Printf("HTTP Timeout: %s\n", cfg.HTTP.Timeout)
	fmt.Printf("GRPC Addr: %s\n", cfg.Grpc.Addr)
	fmt.Printf("Log Level: %s\n", cfg.Log.Level)
	fmt.Printf("Log Format: %s\n", cfg.Log.Format)
	fmt.Printf("MySQL Addr: %s\n", cfg.Mysql.Addr)
	fmt.Printf("MySQL Username: %s\n", cfg.Mysql.Username)
	fmt.Printf("MySQL Database: %s\n", cfg.Mysql.Database)
	fmt.Printf("Redis Addr: %s\n", cfg.Redis.Addr)
	fmt.Printf("Jaeger Host: %s\n", cfg.Jaeger.AgentHost)
	fmt.Printf("Jaeger Port: %d\n", cfg.Jaeger.AgentPort)
}
