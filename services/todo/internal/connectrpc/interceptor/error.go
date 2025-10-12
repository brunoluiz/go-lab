package interceptor

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
)

func Error(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err != nil {
				logger.ErrorContext(ctx, "gRPC error", "error", err, "method", req.Spec().Procedure)
			}
			return resp, err
		}
	}
}
