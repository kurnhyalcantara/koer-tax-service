package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if !interceptor.isRestricted(info.FullMethod) {
			return handler(ctx, req)
		}

		claims, err := interceptor.claimsToken(ctx)
		if err != nil {
			return nil, err
		}

		err = interceptor.authorize(ctx, claims, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
