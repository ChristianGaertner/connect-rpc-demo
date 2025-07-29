package interceptor

import (
	"context"
	"log"

	"connectrpc.com/connect"
)

func Logging() connect.Interceptor {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			log.Println("req", req.Spec().Procedure)
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
