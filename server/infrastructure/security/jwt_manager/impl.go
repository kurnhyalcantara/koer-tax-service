package jwtmanager

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	auth_pb "github.com/kurnhyalcantara/koer-tax-service/protogen/auth-service"

	jwtmanager "github.com/kurnhyalcantara/koer-tax-service/server/domain/security/jwt_manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (j *JwtManager) Verify(accessToken string) (*jwtmanager.UserClaims, error) {

	token, err := jwt.ParseWithClaims(
		accessToken,
		&jwtmanager.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwtmanager.UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil

}

// GetMeFromAuthService implements jwtmanager.JwtManagerInterface.
func (j *JwtManager) GetMeFromAuthService(ctx context.Context, accessToken string) (*auth_pb.VerifyTokenRes, error) {
	userData, err := j.AuthClient.VerifyToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed verify token: %v", err)
	}

	if userData == nil {
		return nil, fmt.Errorf("failed verify user data")
	}

	return userData, nil
}

// GetMeFromMD implements JwtManagerCore.
func (j *JwtManager) GetMeFromMD(ctx context.Context) (*jwtmanager.UserData, metadata.MD, error) {
	var mdResult metadata.MD
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if len(md) > 0 && len(md["user-userid"]) > 0 {
			mdResult = md
		}
		ctx = metadata.NewOutgoingContext(context.Background(), md)
	}

	if mdResult == nil {

		var trailer metadata.MD

		_, err := j.AuthClient.SetMe(ctx, grpc.Trailer(&trailer))
		if err != nil {
			return nil, nil, err
		}

		mdResult = metadata.Join(mdResult, trailer)

	}

	user := &jwtmanager.UserData{}
	var err error

	user.UserID, err = strconv.ParseUint(mdResult["user-userid"][0], 10, 64)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "Internal Error")
	}

	user.CompanyID, err = strconv.ParseUint(mdResult["user-companyid"][0], 10, 64)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "Internal Error")
	}

	user.Username = mdResult["user-username"][0]
	user.CompanyName = mdResult["user-companyname"][0]
	user.UserType = mdResult["user-usertype"][0]

	user.Authorities = strings.Split(mdResult["user-authorities"][0], ",")

	ids := strings.Split(mdResult["user-groupids"][0], ",")
	for _, v := range ids {
		if len(v) > 0 {
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, nil, status.Errorf(codes.Internal, "Internal Error")
			}
			user.GroupIDs = append(user.GroupIDs, id)
		}
	}

	ids = strings.Split(mdResult["user-roleids"][0], ",")
	for _, v := range ids {
		if len(v) > 0 {
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, nil, status.Errorf(codes.Internal, "Internal Error")
			}
			user.RoleIDs = append(user.RoleIDs, id)
		}
	}

	user.SessionID = mdResult["user-sessionid"][0]
	user.DateTime = mdResult["user-datetime"][0]
	user.TokenCreatedAt = mdResult["user-tokencreatedat"][0]
	user.IdToken = mdResult["user-idtoken"][0]

	return user, mdResult, nil
}
