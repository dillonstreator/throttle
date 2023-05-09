# [throttle](https://github.com/dillonstreator/throttle)

A simple golang throttle utility for redis and in memory, allowing inline and curried throttling.

## Installation

```sh
go get github.com/dillonstreator/throttle
```

## Usage

Both the [RedisThrottler](./redis.go) and [MemoryThrottler](./memory.go) adhere to the [Throttler](./throttler.go) interface which exposes a `Do` method for inline throttling and a `New` method for curried throttling.

The error `throttle.ErrThrottled` is returned in the event that a function call is throttled, allowing handling of this event if necessary.

### Inline throttling with `.Do`

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dillonstreator/throttle"
)

func main() {
	throttler := throttle.NewMemoryThrottler()

	key := "example:do"
	fn := func(ctx context.Context) error {
		fmt.Printf("called %s\n", key)
		return nil
	}

	err := throttler.Do(ctx, key, time.Second, fn)
	fmt.Println(err) // nil

	time.Sleep(time.Millisecond * 500)

	err = throttler.Do(ctx, key, time.Second, fn)
	fmt.Println(err) // throttle.ErrThrottled

	time.Sleep(time.Millisecond * 500)

	err = throttler.Do(ctx, key, time.Second, fn)
	fmt.Println(err) // nil
}
```

### Curried throttling with `.New`

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dillonstreator/throttle"
)

func main() {
	throttler := throttle.NewMemoryThrottler()

	key := "example:new"
	fn := throttler.New(key, time.Second, func(ctx context.Context) error {
		fmt.Printf("called %s\n", key)
		return nil
	})

	ctx := context.Background()

	err := fn(ctx)
	fmt.Println(err) // nil

	time.Sleep(time.Millisecond * 500)

	err = fn(ctx)
	fmt.Println(err) // throttle.ErrThrottled

	time.Sleep(time.Millisecond * 500)

	err = fn(ctx)
	fmt.Println(err) // nil
}
```

## Note

Trailing edge function invocation is not supported by this library. Only the leading edge is considered.
