package rescache

import (
	"context"
	"time"

	"github.com/kurnhyalcantara/koer-tax-service/pkg/log"
	"github.com/kurnhyalcantara/koer-tax-service/server/domain/rescache"
	inmemory "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/db/in_memory"
	jwtmanager "github.com/kurnhyalcantara/koer-tax-service/server/infrastructure/security/jwt_manager"
)

//go:generate mockery --name=ResCacheCore --output=../../tests/mocks --structname=MockResCacheCore

type ResCacheCore interface {
	Start(ctx context.Context, code string, userId uint64) (*rescache.CacheData, error)
	StartWithoutUser(ctx context.Context, code string) *rescache.CacheData
	StoreResponse(ctx context.Context, cacheData *rescache.CacheData) error
	GetResponse(ctx context.Context, code string, userId uint64) (*rescache.CacheData, error)
	GetResponseWithoutUser(ctx context.Context, code string) (*rescache.CacheData, error)
	DeleteResponse(ctx context.Context, code string) error
	KeyGen(rcode string, userId string) (string, error)
}

type ResCache struct {
	redis              inmemory.RedisCore
	manager            jwtmanager.JwtManagerCore
	service            string
	minTimeout         time.Duration
	onProgressDuration time.Duration
	doneDuration       time.Duration
	serviceLogger      log.LoggerCore
}

func NewResCache(rdb inmemory.RedisCore, jwtManager jwtmanager.JwtManagerCore, service string, minTimeout, onProgressDuration, doneDuration time.Duration, logger log.LoggerCore) ResCacheCore {
	return &ResCache{
		redis:              rdb,
		manager:            jwtManager,
		service:            service,
		minTimeout:         minTimeout,
		onProgressDuration: onProgressDuration,
		doneDuration:       doneDuration,
		serviceLogger:      logger,
	}
}
