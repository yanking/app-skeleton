package config

import (
	"github.com/yanking/app-skeleton/pkg/conf"
	"github.com/yanking/app-skeleton/pkg/log"
)

var config *Config

func Init(configFile string, fs ...func()) error {
	config = &Config{}
	return conf.Parse(configFile, config, fs...)
}

func Get() *Config {
	if config == nil {
		panic("configs is nil, please call configs.Init() first")
	}
	return config
}

func Set(conf *Config) {
	config = conf
}

type Config struct {
	AppName       string       `mapstructure:"app-name" yaml:"app-name" json:"app-name"`
	EnableMetrics bool         `mapstructure:"enable-metrics" yaml:"enable-metrics" json:"enable-metrics"`
	ServerMode    string       `mapstructure:"server-mode" yaml:"server-mode" json:"server-mode"`
	JwtKey        string       `mapstructure:"jwt-key" yaml:"jwt-key" json:"jwt-key"`
	Expiration    string       `mapstructure:"expiration" yaml:"expiration" json:"expiration"`
	HTTP          HTTPConfig   `mapstructure:"http" yaml:"http" json:"http"`
	Grpc          GrpcConfig   `mapstructure:"grpc" yaml:"grpc" json:"grpc"`
	Log           *log.Options `mapstructure:"log" yaml:"log" json:"log"`
	Mysql         MysqlConfig  `mapstructure:"mysql" yaml:"mysql" json:"mysql"`
	Redis         RedisConfig  `mapstructure:"redis" yaml:"redis" json:"redis"`
	Jaeger        JaegerConfig `mapstructure:"jaeger" yaml:"jaeger" json:"jaeger"`
}

// HTTPConfig 对应 HTTP 相关配置 (主要用于 gRPC-Gateway)
type HTTPConfig struct {
	Addr    string `mapstructure:"addr" yaml:"addr" json:"addr"`
	Timeout string `mapstructure:"timeout" yaml:"timeout" json:"timeout"`
}

// GrpcConfig 对应 gRPC 相关配置
type GrpcConfig struct {
	Addr string `mapstructure:"addr" yaml:"addr" json:"addr"`
	// Gateway 用于配置 gRPC-Gateway 相关选项
	Gateway GatewayConfig `mapstructure:"gateway" yaml:"gateway" json:"gateway"`
}

// GatewayConfig 对应 gRPC-Gateway 相关配置
type GatewayConfig struct {
	// Enabled 控制是否启用 gRPC-Gateway
	Enabled bool `mapstructure:"enabled" yaml:"enabled" json:"enabled"`
}

// MysqlConfig 对应 MySQL 数据库相关配置
type MysqlConfig struct {
	Addr                  string `mapstructure:"addr" yaml:"addr" json:"addr"`
	Username              string `mapstructure:"username" yaml:"username" json:"username"`
	Password              string `mapstructure:"password" yaml:"password" json:"password"`
	Database              string `mapstructure:"database" yaml:"database" json:"database"`
	MaxIdleConnections    int    `mapstructure:"max-idle-connections" yaml:"max-idle-connections" json:"max-idle-connections"`
	MaxOpenConnections    int    `mapstructure:"max-open-connections" yaml:"max-open-connections" json:"max-open-connections"`
	MaxConnectionLifeTime string `mapstructure:"max-connection-life-time" yaml:"max-connection-life-time" json:"max-connection-life-time"`
	LogLevel              int    `mapstructure:"log-level" yaml:"log-level" json:"log-level"`
}

// RedisConfig 对应 Redis 相关配置
type RedisConfig struct {
	Addr         string   `mapstructure:"addr" yaml:"addr" json:"addr"`
	Password     string   `mapstructure:"password" yaml:"password" json:"password"`
	Db           int      `mapstructure:"db" yaml:"db" json:"db"`
	Addrs        []string `mapstructure:"addrs" yaml:"addrs" json:"addrs"`
	DialTimeout  int      `mapstructure:"dialTimeout" yaml:"dialTimeout" json:"dialTimeout"`
	ReadTimeout  int      `mapstructure:"readTimeout" yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout int      `mapstructure:"writeTimeout" yaml:"writeTimeout" json:"writeTimeout"`
}

// JaegerConfig 对应 Jaeger 追踪相关配置
type JaegerConfig struct {
	AgentHost string `mapstructure:"agentHost" yaml:"agentHost" json:"agentHost"`
	AgentPort int    `mapstructure:"agentPort" yaml:"agentPort" json:"agentPort"`
}