package jwtmanager

import (
	"context"
	"log"
	"time"

	"github.com/kurnhyalcantara/koer-tax-service/config"
	auth_pb "github.com/kurnhyalcantara/koer-tax-service/protogen/auth-service"
	jwtmanager "github.com/kurnhyalcantara/koer-tax-service/server/domain/security/jwt_manager"
	grpcClient "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/grpc_client"

	"google.golang.org/grpc/metadata"
)

//go:generate mockery --name=JwtManagerCore --output=../../tests/mocks --structname=MockJwtManagerCore

type JwtManagerCore interface {
	Verify(accessToken string) (*jwtmanager.UserClaims, error)
	GetMeFromAuthService(ctx context.Context, accessToken string) (*auth_pb.VerifyTokenRes, error)
	GetMeFromMD(ctx context.Context) (*jwtmanager.UserData, metadata.MD, error)
}

type JwtManager struct {
	SecretKey     string
	TokenDuration time.Duration
	AuthClient    grpcClient.AuthService
}

func NewJwtManager(cfg *config.Config, authClient grpcClient.AuthService) JwtManagerCore {
	tokenDuration, err := time.ParseDuration(cfg.JwtDuration)
	if err != nil {
		log.Fatalf("failed to parse duration for jwt manager: %v", err)
	}

	return &JwtManager{
		SecretKey:     cfg.JwtSecret,
		TokenDuration: tokenDuration,
		AuthClient:    authClient,
	}
}
