package throttle

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestMemoryThrottler_New(t *testing.T) {
	throttler := NewMemoryThrottler()

	var calls int32
	fn := throttler.New("test:new", time.Millisecond*100, func(ctx context.Context) error {
		atomic.AddInt32(&calls, 1)
		return nil
	})

	ctx := context.Background()

	err := fn(ctx)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	time.Sleep(time.Millisecond * 50)

	err = fn(ctx)
	if !errors.Is(err, ErrThrottled) {
		t.Fatalf("unexpected error %v", err)
	}

	time.Sleep(time.Millisecond * 50)

	err = fn(ctx)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if calls != 2 {
		t.Fatalf("expected 2 calls but got %d", calls)
	}
}

func TestMemoryThrottler_Do(t *testing.T) {
	throttler := NewMemoryThrottler()

	var calls int32
	fn := func(ctx context.Context) error {
		atomic.AddInt32(&calls, 1)
		return nil
	}

	ctx := context.Background()

	err := throttler.Do(ctx, "test:do", time.Millisecond*100, fn)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	time.Sleep(time.Millisecond * 50)

	err = throttler.Do(ctx, "test:do", time.Millisecond*100, fn)
	if !errors.Is(err, ErrThrottled) {
		t.Fatalf("unexpected error %v", err)
	}

	time.Sleep(time.Millisecond * 50)

	err = throttler.Do(ctx, "test:do", time.Millisecond*100, fn)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if calls != 2 {
		t.Fatalf("expected 2 calls but got %d", calls)
	}
}
