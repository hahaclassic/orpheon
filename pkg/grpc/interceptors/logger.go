package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

func UnaryLoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Info("incoming gRPC request", "method", info.FullMethod)

		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error("gRPC request failed", "method", info.FullMethod, "error", err)
		} else {
			logger.Info("gRPC request completed", "method", info.FullMethod)
		}

		return resp, err
	}
}
