package interceptor

import (
	"context"
	"strings"

	"github.com/kurnhyalcantara/koer-tax-service/pkg/utils"
	jwtmanagerDom "github.com/kurnhyalcantara/koer-tax-service/server/domain/security/jwt_manager"
	jwtmanager "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/security/jwt_manager"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager jwtmanager.JwtManagerCore
}

func NewAuthInterceptor(jwtManager jwtmanager.JwtManagerCore) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

func (interceptor *AuthInterceptor) claimsToken(ctx context.Context) (*jwtmanagerDom.UserClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	split := strings.Split(values[0], " ")
	accessToken := split[0]
	if len(split) > 1 {
		accessToken = split[1]
	}

	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	getUser, err := interceptor.jwtManager.GetMeFromAuthService(ctx, accessToken)

	claims.ProductRoles = getUser.ProductRoles
	claims.UserType = getUser.UserType

	if err != nil {
		return nil, err
	}
	if getUser.IsExpired && !getUser.IsValid {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return claims, nil
}

func (interceptor *AuthInterceptor) isRestricted(infoMethod string) bool {
	_, restricted := accessibleRoles[infoMethod]
	return restricted
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, claims *jwtmanagerDom.UserClaims, method string) error {
	featureRoles := []string{}
	for _, v := range claims.ProductRoles {
		// TODO: must be check for final product name
		if utils.Contains([]string{"MPN"}, v.ProductName) {
			featureRoles = append(featureRoles, v.Authorities...)
		}
	}

	accessibleRoles, ok := accessibleRoles[method]
	if !ok {
		// everyone can access
		return nil
	}

	if len(accessibleRoles) < 1 {
		return nil
	}

	for _, role := range accessibleRoles {
		for _, exist := range featureRoles {
			if role == exist {
				return nil
			}
		}
	}

	return status.Error(codes.PermissionDenied, "Access denied")
}
