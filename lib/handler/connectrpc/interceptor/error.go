package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"github.com/samber/oops"
)

func Error() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if oopsErr, ok := oops.AsOops(err); ok {
				switch oopsErr.Code() {
				case "NOT_FOUND":
					return resp, connect.NewError(connect.CodeNotFound, err)
				case "VALIDATION":
					return resp, connect.NewError(connect.CodeInvalidArgument, err)
				case "CONFLICT":
					return resp, connect.NewError(connect.CodeAlreadyExists, err)
				case "UNKNOWN":
					return resp, connect.NewError(connect.CodeInternal, err)
				default:
					return resp, connect.NewError(connect.CodeInternal, err)
				}
			}

			return resp, err
		}
	}
}
