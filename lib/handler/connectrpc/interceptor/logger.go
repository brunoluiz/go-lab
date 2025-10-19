package interceptor

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
)

func Logger(logger *slog.Logger) connect.UnaryInterceptorFunc {
	l := logger
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			startedAt := time.Now()
			resp, err := next(ctx, req)
			finishedAt := time.Now()

			level := slog.LevelInfo
			switch connect.CodeOf(err) {
			case connect.CodeNotFound, connect.CodeInvalidArgument, connect.CodeAlreadyExists:
				level = slog.LevelWarn
				l.With("error", err, "method", req.Spec().Procedure)
			case connect.CodeInternal:
				level = slog.LevelError
				l.With("error", err, "method", req.Spec().Procedure)
			default:
			}

			l.Log(ctx, level, "connectrpc request",
				slog.String("method", req.Spec().Procedure),
				slog.Duration("duration", (finishedAt.Sub(startedAt))/time.Millisecond),
			)
			return resp, err
		}
	}
}
