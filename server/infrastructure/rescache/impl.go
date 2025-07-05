package rescache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/kurnhyalcantara/koer-tax-service/pkg/log"
	"github.com/kurnhyalcantara/koer-tax-service/server/domain/rescache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	redisV8 "github.com/go-redis/redis/v8"
)

// DeleteResponse implements ResCacheCore.
func (r *ResCache) DeleteResponse(ctx context.Context, code string) error {
	panic("unimplemented")
}

// GetResponse implements ResCacheCore.
func (r *ResCache) GetResponse(ctx context.Context, code string, userId uint64) (*rescache.CacheData, error) {

	userIdStr := strconv.FormatUint(userId, 10)

	key, err := r.KeyGen(code, userIdStr)
	if err != nil {
		r.serviceLogger.Error(log.LogPayload{Message: "[ResCache] failed to generate key", Metadata: err})
		return nil, status.Error(codes.Internal, "Unexpected Error")
	}

	found := &rescache.CacheData{}

	val, err := r.redis.Get(ctx, key)
	switch {
	case err == redisV8.Nil:
		return nil, nil
	case err != nil:
		r.serviceLogger.Error(log.LogPayload{Message: "[ResCache] failed to get data from redis", Metadata: err})
		return nil, status.Error(codes.Internal, "Unexpected Error")
	}

	if err := json.Unmarshal([]byte(val), found); err != nil {
		r.serviceLogger.Error(log.LogPayload{Message: "[ResCache] failed unmarshal data from redis", Metadata: err})
		return nil, status.Error(codes.Internal, "Unexpected Error")
	}

	if found.Progress == "OnProgress" {
		return nil, nil
	}

	if found.Progress == "Done" {
		err := r.redis.Del(context.Background(), key)
		if err != nil {
			r.serviceLogger.Error(log.LogPayload{Message: "[ResCache] failed to del cache data from redis", Metadata: err})
			return nil, err
		}
	}

	return found, nil
}

// GetResponseWithoutUser implements ResCacheCore.
func (r *ResCache) GetResponseWithoutUser(ctx context.Context, code string) (*rescache.CacheData, error) {
	panic("unimplemented")
}

// KeyGen implements ResCacheCore.
func (r *ResCache) KeyGen(rcode string, userId string) (string, error) {
	panic("unimplemented")
}

// Start implements ResCacheCore.
func (r *ResCache) Start(ctx context.Context, code string, userId uint64) (*rescache.CacheData, error) {

	userIdStr := strconv.FormatUint(userId, 10)

	data := &rescache.CacheData{
		Code:         code,
		User:         userIdStr,
		Response:     "",
		ResponseCode: "200",
		Error:        "",
		StartTime:    time.Now(),
		Progress:     "OnProgress",
	}

	key, err := r.KeyGen(code, userIdStr)
	if err != nil {
		r.serviceLogger.Error(log.LogPayload{Message: "[ResCache] failed to generate key", Metadata: err})
		return nil, err
	}

	if err := r.redis.Set(
		ctx,
		key,
		data,
		r.onProgressDuration,
	); err != nil {
		r.serviceLogger.Error(log.LogPayload{Message: "[ResCache] failed to set redis", Metadata: err})
		return nil, err
	}

	r.serviceLogger.Info(log.LogPayload{Message: "[ResCache] Start Response Cache"})

	return data, nil
}

// StartWithoutUser implements ResCacheCore.
func (r *ResCache) StartWithoutUser(ctx context.Context, code string) *rescache.CacheData {
	panic("unimplemented")
}

// StoreResponse implements ResCacheCore.
func (r *ResCache) StoreResponse(ctx context.Context, cacheData *rescache.CacheData) error {

	cacheData.Progress = "Stored"

	key, err := r.KeyGen(cacheData.Code, cacheData.User)
	if err != nil {
		r.serviceLogger.Error(log.LogPayload{Message: "[ResCache] failed to generate key", Metadata: err})
		return err
	}

	if err := r.redis.Set(
		ctx,
		key,
		cacheData,
		r.doneDuration,
	); err != nil {
		r.serviceLogger.Error(log.LogPayload{Message: "[ResCache] failed to set redis", Metadata: err})
		return err
	}

	r.serviceLogger.Info(log.LogPayload{Message: "[ResCache] Store Response Cache"})

	return nil
}
