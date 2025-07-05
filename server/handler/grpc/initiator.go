package grpc

import (
	"github.com/kurnhyalcantara/koer-tax-service/pkg/validator"
	pb "github.com/kurnhyalcantara/koer-tax-service/protogen/koer-tax-service"
	jwtmanager "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/security/jwt_manager"
	"github.com/kurnhyalcantara/koer-tax-service/server/usecase"
)

type Handler struct {
	validator validator.Validator
	uc        *usecase.UseCases
	manager   jwtmanager.JwtManagerCore

	pb.KoerTaxServiceServer
}

func NewHandler(validator validator.Validator, uc *usecase.UseCases, manager jwtmanager.JwtManagerCore) *Handler {
	return &Handler{
		validator: validator,
		uc:        uc,
		manager:   manager,
	}
}
