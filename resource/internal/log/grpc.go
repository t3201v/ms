package log

import (
	"context"

	"google.golang.org/grpc"
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	Info(ctx, "grpc", "method", info.FullMethod, "req", req)
	return handler(ctx, req)
}
