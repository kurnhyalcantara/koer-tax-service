package interceptor

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	auth_interceptor "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/interceptor/auth"
	jwtmanager "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/security/jwt_manager"
	"google.golang.org/grpc"
)

// Interceptors implements the grpc.UnaryServerInteceptor function to add
// interceptors around all gRPC unary calls
func UnaryInterceptors(jwtManager jwtmanager.JwtManagerCore) grpc.UnaryServerInterceptor {
	var authInterceptor = auth_interceptor.NewAuthInterceptor(jwtManager)
	return grpc_middleware.ChainUnaryServer(
		authInterceptor.Unary(),
	)
}

// Interceptors implements the grpc.StreamServerInteceptor function to add
// interceptors around all gRPC stream calls
func StreamInterceptors(jwtManager jwtmanager.JwtManagerCore) grpc.StreamServerInterceptor {
	var authInterceptor = auth_interceptor.NewAuthInterceptor(jwtManager)
	return grpc_middleware.ChainStreamServer(
		authInterceptor.Stream(),
	)
}
