package throttle

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisThrottler redis throttler
type RedisThrottler struct {
	redisClient redis.UniversalClient
}

var _ Throttler = (*RedisThrottler)(nil)

// NewRedisThrottler creates a new RedisThrottler
func NewRedisThrottler(redisClient redis.UniversalClient) *RedisThrottler {
	return &RedisThrottler{
		redisClient: redisClient,
	}
}

// New curried redis throttling
func (r *RedisThrottler) New(key string, duration time.Duration, fn func(ctx context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		return r.do(ctx, key, duration, fn)
	}
}

// New inline redis throttling
func (r *RedisThrottler) Do(ctx context.Context, key string, duration time.Duration, fn func(ctx context.Context) error) error {
	return r.do(ctx, key, duration, fn)
}

func (r *RedisThrottler) do(ctx context.Context, key string, duration time.Duration, fn func(ctx context.Context) error) error {
	err := r.redisClient.SetArgs(ctx, key, "", redis.SetArgs{Get: true, KeepTTL: true}).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("getting key: %w", err)
	}

	if err == nil {
		return ErrThrottled
	}

	err = r.redisClient.SetEX(ctx, key, "", duration).Err()
	if err != nil {
		return fmt.Errorf("setting expiration: %w", err)
	}

	return fn(ctx)
}
