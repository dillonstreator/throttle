package throttle

import (
	"context"
	"errors"
	"time"
)

var (
	ErrThrottled = errors.New("throttled")
)

//Throttler is the throttling interface
type Throttler interface {
	Do(ctx context.Context, key string, duration time.Duration, fn func(ctx context.Context) error) error
	New(key string, duration time.Duration, fn func(ctx context.Context) error) func(context.Context) error
}
