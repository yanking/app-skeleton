package clientinterceptors

import (
	"context"

	"github.com/google/uuid"
	grpcmetadata "github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	// ContextRequestIDKey request id key for context
	ContextRequestIDKey = "request_id"
)

// CtxKeyString for context.WithValue key type
type CtxKeyString string

// RequestIDKey request_id
var RequestIDKey = CtxKeyString(ContextRequestIDKey)

// CtxRequestIDField get request id field from context.Context
func CtxRequestIDField(ctx context.Context) zap.Field {
	return zap.String(ContextRequestIDKey, grpcmetadata.ExtractOutgoing(ctx).Get(ContextRequestIDKey))
}

// ClientCtxRequestID get request id from rpc client context.Context
func ClientCtxRequestID(ctx context.Context) string {
	return grpcmetadata.ExtractOutgoing(ctx).Get(ContextRequestIDKey)
}

// ClientCtxRequestIDField get request id field from rpc client context.Context
func ClientCtxRequestIDField(ctx context.Context) zap.Field {
	return zap.String(ContextRequestIDKey, grpcmetadata.ExtractOutgoing(ctx).Get(ContextRequestIDKey))
}

// UnaryClientRequestID client-side request_id unary interceptor
func UnaryClientRequestID() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requestID := ClientCtxRequestID(ctx)
		if requestID == "" {
			requestID = uuid.NewString()
			ctx = metadata.AppendToOutgoingContext(ctx, ContextRequestIDKey, requestID)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamClientRequestID client request id stream interceptor
func StreamClientRequestID() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		requestID := ClientCtxRequestID(ctx)
		if requestID == "" {
			requestID = uuid.NewString()
			ctx = metadata.AppendToOutgoingContext(ctx, ContextRequestIDKey, requestID)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}
