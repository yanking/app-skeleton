package grpc

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	srvintc "github.com/yanking/app-skeleton/pkg/grpc/serverinterceptors"
	"github.com/yanking/app-skeleton/pkg/log"

	apimd "github.com/go-kratos/kratos/v2/api/metadata"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type ServerOption func(o *Server)

type Server struct {
	*grpc.Server

	address    string
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
	grpcOpts   []grpc.ServerOption
	lis        net.Listener
	timeout    time.Duration

	health   *health.Server
	metadata *apimd.Server
	endpoint *url.URL

	enableMetrics bool
	enableTracing bool

	// gRPC-Gateway 相关字段
	enableGateway bool
	gatewayAddr   string
	gatewayServer *http.Server
	gatewayMux    *runtime.ServeMux
}

func (s *Server) Name() string {
	if s.enableGateway {
		return "grpc+gatewayServer"
	}
	return "grpcServer"
}

func (s *Server) Endpoint() *url.URL {
	return s.endpoint
}

func (s *Server) Address() string {
	return s.address
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address: ":0",
		health:  health.NewServer(),
	}

	for _, opt := range opts {
		opt(srv)
	}

	unaryInts := []grpc.UnaryServerInterceptor{
		srvintc.UnaryCrashInterceptor,
	}

	streamInts := []grpc.StreamServerInterceptor{
		srvintc.StreamCrashInterceptor,
	}

	if srv.enableMetrics {
		unaryInts = append(unaryInts, srvintc.UnaryPrometheusInterceptor)
	}

	if srv.timeout > 0 {
		unaryInts = append(unaryInts, srvintc.UnaryTimeoutInterceptor(srv.timeout))
	}

	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}

	//把我们传入的拦截器转换成grpc的ServerOption
	grpcOpts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(unaryInts...)}

	//把用户自己传入的grpc.ServerOption放在一起
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	if len(srv.streamInts) > 0 {
		streamInts = append(streamInts, srv.streamInts...)
	}

	//处理流式拦截器
	if len(streamInts) > 0 {
		grpcOpts = append(grpcOpts, grpc.ChainStreamInterceptor(streamInts...))
	}

	//处理tracing拦截器
	if srv.enableTracing {
		grpcOpts = append(grpcOpts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}

	srv.Server = grpc.NewServer(grpcOpts...)

	//注册metadata的Server
	srv.metadata = apimd.NewServer(srv.Server)

	//解析address
	err := srv.listenAndEndpoint()
	if err != nil {
		panic(err)
	}

	//注册health
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	apimd.RegisterMetadataServer(srv.Server, srv.metadata)

	reflection.Register(srv.Server)
	//可以支持用户直接通过grpc的一个接口查看当前支持的所有的rpc服务

	// 初始化 gRPC-Gateway
	if srv.enableGateway {
		srv.gatewayMux = runtime.NewServeMux()
		srv.gatewayServer = &http.Server{
			Addr:    srv.gatewayAddr,
			Handler: srv.gatewayMux,
		}
	}

	return srv
}

func WithAddress(address string) ServerOption {
	return func(o *Server) {
		o.address = address
	}
}

func WithMetrics(metric bool) ServerOption {
	return func(o *Server) {
		o.enableMetrics = metric
	}
}

func WithTracing(tracing bool) ServerOption {
	return func(o *Server) {
		o.enableTracing = tracing
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(o *Server) {
		o.timeout = timeout
	}
}

func WithLis(lis net.Listener) ServerOption {
	return func(o *Server) {
		o.lis = lis
	}
}

func WithUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(o *Server) {
		o.unaryInts = in
	}
}

func WithStreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInts = in
	}
}

func WithOptions(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

// WithGateway 启用 gRPC-Gateway
func WithGateway(enable bool, addr string) ServerOption {
	return func(s *Server) {
		s.enableGateway = enable
		s.gatewayAddr = addr
	}
}

// GetGatewayMux 返回 gRPC-Gateway 的多路复用器，用于注册 HTTP 处理程序
func (s *Server) GetGatewayMux() *runtime.ServeMux {
	return s.gatewayMux
}

// 完成ip和端口的提取
func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}

	s.endpoint = &url.URL{Scheme: "grpc", Host: s.address}
	return nil
}

// Start 启动grpc的服务
func (s *Server) Start(ctx context.Context) error {
	// 启动 gRPC 服务器
	log.Infof("[grpc] server listening on: %s", s.lis.Addr().String())
	s.health.Resume()

	// 在单独的goroutine中启动服务器，以便可以监听上下文取消
	go func() {
		<-ctx.Done()
		s.GracefulStop()
		if s.gatewayServer != nil {
			_ = s.gatewayServer.Shutdown(context.Background())
		}
	}()

	// 如果启用了 gRPC-Gateway，则同时启动 HTTP 服务器
	if s.enableGateway {
		go func() {
			log.Infof("[gateway] server listening on: %s", s.gatewayAddr)
			if err := s.gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Errorf("[gateway] server error: %v", err)
			}
		}()
	}

	return s.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	//设置服务的状态为not_serving，防止接收新的请求过来
	s.health.Shutdown()

	if s.gatewayServer != nil {
		_ = s.gatewayServer.Shutdown(ctx)
		log.Infof("[gateway] server stopped")
	}

	s.GracefulStop()
	log.Infof("[grpc] server stopped")
	return nil
}
