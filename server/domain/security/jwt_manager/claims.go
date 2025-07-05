package jwtmanager

import (
	"github.com/dgrijalva/jwt-go"
	auth_pb "github.com/kurnhyalcantara/koer-tax-service/protogen/auth-service"
)

type UserClaims struct {
	jwt.StandardClaims
	UserType     string                      `json:"user_type"`
	ProductRoles []*auth_pb.ProductAuthority `json:"product_roles"`
	Authorities  []string                    `json:"authorities"`
}
