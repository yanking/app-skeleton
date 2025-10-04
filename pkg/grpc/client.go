package grpc

import (
	"context"

	"time"

	"github.com/yanking/app-skeleton/pkg/grpc/clientinterceptors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	endpoint      string
	timeout       time.Duration
	unaryInts     []grpc.UnaryClientInterceptor
	streamInts    []grpc.StreamClientInterceptor
	rpcOpts       []grpc.DialOption
	enableTracing bool
	enableMetrics bool
}

func WithEnableTracing(enable bool) ClientOption {
	return func(o *clientOptions) {
		o.enableTracing = enable
	}
}

// WithEndpoint 设置地址
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithClientTimeout 设置超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// WithClientUnaryInterceptor 设置拦截器
func WithClientUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.unaryInts = in
	}
}

// WithClientStreamInterceptor 设置stream拦截器
func WithClientStreamInterceptor(in ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.streamInts = in
	}
}

// WithClientOptions 设置grpc的dial选项
func WithClientOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.rpcOpts = opts
	}
}

func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout:       2000 * time.Millisecond,
		enableTracing: true,
	}

	for _, o := range opts {
		o(&options)
	}

	//TODO 客户端默认拦截器
	ints := []grpc.UnaryClientInterceptor{
		clientinterceptors.UnaryClientTimeoutInterceptor(options.timeout),
	}

	if options.enableMetrics {
		ints = append(ints, clientinterceptors.PrometheusInterceptor())
	}

	var streamInts []grpc.StreamClientInterceptor

	if len(options.unaryInts) > 0 {
		ints = append(ints, options.unaryInts...)
	}
	if len(options.streamInts) > 0 {
		streamInts = append(streamInts, options.streamInts...)
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(ints...),
		grpc.WithChainStreamInterceptor(streamInts...),
	}

	if options.enableTracing {
		grpcOpts = append(grpcOpts, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	}

	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}

	if len(options.rpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.rpcOpts...)
	}

	return grpc.NewClient(options.endpoint, grpcOpts...)
}
