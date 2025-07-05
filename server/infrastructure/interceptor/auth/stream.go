package interceptor

import "google.golang.org/grpc"

func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {

		if !interceptor.isRestricted(info.FullMethod) {
			return handler(srv, stream)
		}

		claims, err := interceptor.claimsToken(stream.Context())
		if err != nil {
			return err
		}

		err = interceptor.authorize(stream.Context(), claims, info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}
