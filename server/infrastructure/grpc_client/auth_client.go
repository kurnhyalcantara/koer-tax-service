package grpcclient

import (
	"context"

	auth_pb "github.com/kurnhyalcantara/koer-tax-service/protogen/auth-service"

	"google.golang.org/grpc"
)

type AuthGrpcClient struct {
	client auth_pb.ApiServiceClient
}

// SetMe implements AuthService.
func (a *AuthGrpcClient) SetMe(ctx context.Context, opts ...grpc.CallOption) (*auth_pb.SetMeRes, error) {
	return a.client.SetMe(ctx, &auth_pb.VerifyTokenReq{}, opts...)
}

// VerifyToken implements AuthService.
func (a *AuthGrpcClient) VerifyToken(ctx context.Context, accessToken string) (*auth_pb.VerifyTokenRes, error) {
	return a.client.VerifyToken(ctx, &auth_pb.VerifyTokenReq{
		AccessToken: accessToken,
	})
}

func NewAuthGrpcClient(conn *grpc.ClientConn) AuthService {
	return &AuthGrpcClient{
		client: auth_pb.NewApiServiceClient(conn),
	}
}
