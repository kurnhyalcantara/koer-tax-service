package inmemory

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

//go:generate mockery --name=RedisCore --output=../../tests/mocks --structname=MockRedisCore

type RedisCore interface {
	GetClient() *redis.Client
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type Redis struct {
	client *redis.Client
}

// Del implements RedisCore.
func (r *Redis) Del(ctx context.Context, key string) error {
	panic("unimplemented")
}

// Get implements RedisCore.
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	panic("unimplemented")
}

// GetClient implements RedisCore.
func (r *Redis) GetClient() *redis.Client {
	return r.client
}

// Set implements RedisCore.
func (r *Redis) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	panic("unimplemented")
}

// Exists implements RedisCore.
func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func New(addr string, pass string, db int) RedisCore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	return &Redis{
		client: client,
	}
}
