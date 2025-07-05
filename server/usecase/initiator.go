package usecase

import (
	"github.com/kurnhyalcantara/koer-tax-service/pkg/log"
	"github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/event/rabbitmq"
	grpcclient "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/grpc_client"
	"github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/rescache"
	jwtmanager "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/security/jwt_manager"
	"github.com/kurnhyalcantara/koer-tax-service/server/repository"
	"github.com/kurnhyalcantara/koer-tax-service/server/usecase/account"
	"github.com/kurnhyalcantara/koer-tax-service/server/usecase/tax"
)

type UseCases struct {
	Account account.UseCase
	Tax     tax.UseCase
}

type Dependencies struct {
	Logger log.LoggerCore
	Repo   *repository.Repositories

	Manager   jwtmanager.JwtManagerCore
	Rcache    rescache.ResCacheCore
	Publisher *rabbitmq.Publisher

	AccountService grpcclient.AccountService
}

func InitUseCases(dep Dependencies) *UseCases {
	return &UseCases{
		Account: account.NewAccountUseCase(dep.AccountService, dep.Logger),
		Tax:     tax.NewTaxUseCase(dep.Logger, dep.Repo.TaxRepo),
	}
}
