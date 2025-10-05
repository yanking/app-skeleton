package serverinterceptors

import (
	"context"

	"github.com/google/uuid"
	grpcmetadata "github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// ContextRequestIDKey request id key for context
	ContextRequestIDKey = "request_id"
)

// KV key value
type KV struct {
	Key string
	Val interface{}
}

// WrapServerCtx wrap context, used in grpc server-side
func WrapServerCtx(ctx context.Context, kvs ...KV) context.Context {
	ctx = context.WithValue(ctx, ContextRequestIDKey, grpcmetadata.ExtractIncoming(ctx).Get(ContextRequestIDKey)) //nolint
	for _, kv := range kvs {
		ctx = context.WithValue(ctx, kv.Key, kv.Val) //nolint
	}
	return ctx
}

// ServerCtxRequestID get request id from rpc server context.Context
func ServerCtxRequestID(ctx context.Context) string {
	return grpcmetadata.ExtractIncoming(ctx).Get(ContextRequestIDKey)
}

// ServerCtxRequestIDField get request id field from rpc server context.Context
func ServerCtxRequestIDField(ctx context.Context) zap.Field {
	return zap.String(ContextRequestIDKey, grpcmetadata.ExtractIncoming(ctx).Get(ContextRequestIDKey))
}

// UnaryServerRequestID server-side request_id unary interceptor
func UnaryServerRequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestID := ServerCtxRequestID(ctx)
		if requestID == "" {
			requestID = uuid.NewString()
			ctx = grpcmetadata.ExtractIncoming(ctx).Add(ContextRequestIDKey, requestID).ToIncoming(ctx)
		}

		return handler(ctx, req)
	}
}

// StreamServerRequestID server-side request id stream interceptor
func StreamServerRequestID() grpc.StreamServerInterceptor {
	// todo
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		//ctx := stream.Context()
		//requestID := ServerCtxRequestID(ctx)
		//if requestID == "" {
		//	requestID = krand.String(krand.R_All, 10)
		//	ctx = grpc_metadata.ExtractIncoming(ctx).Add(ContextRequestIDKey, requestID).ToIncoming(ctx)
		//}
		return handler(srv, stream)
	}
}
