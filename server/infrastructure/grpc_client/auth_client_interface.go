package grpcclient

import (
	"context"

	auth_pb "github.com/kurnhyalcantara/koer-tax-service/protogen/auth-service"
	"google.golang.org/grpc"
)

//go:generate mockery --name=AuthService --output=../../tests/mocks --structname=MockAuthService

type AuthService interface {
	VerifyToken(ctx context.Context, accessToken string) (*auth_pb.VerifyTokenRes, error)
	SetMe(ctx context.Context, opts ...grpc.CallOption) (*auth_pb.SetMeRes, error)
}
