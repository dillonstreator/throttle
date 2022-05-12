package throttle

import (
	"context"
	"sync"
	"time"
)

type MemoryThrottler struct {
	mus     map[string]*sync.Mutex
	history map[string]time.Time

	mu sync.Mutex
}

var _ Throttler = (*MemoryThrottler)(nil)

func NewMemoryThrottler() *MemoryThrottler {
	return &MemoryThrottler{
		mus:     map[string]*sync.Mutex{},
		history: map[string]time.Time{},
	}
}

var DefaultMemoryThrottler = NewMemoryThrottler()

func (m *MemoryThrottler) New(key string, duration time.Duration, fn func(ctx context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		return m.do(ctx, key, duration, fn)
	}
}

func (m *MemoryThrottler) Do(ctx context.Context, key string, duration time.Duration, fn func(ctx context.Context) error) error {
	return m.do(ctx, key, duration, fn)
}

func (m *MemoryThrottler) do(ctx context.Context, key string, duration time.Duration, fn func(ctx context.Context) error) error {
	mu := m.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	lastCall, ok := m.history[key]
	if !ok || lastCall.Add(duration).Before(time.Now()) {
		m.history[key] = time.Now()

		return fn(ctx)
	}

	return ErrThrottled
}

func (m *MemoryThrottler) getMutex(key string) *sync.Mutex {
	m.mu.Lock()
	defer m.mu.Unlock()

	mu, ok := m.mus[key]
	if ok {
		return mu
	}

	mu = &sync.Mutex{}
	m.mus[key] = mu

	return mu
}
