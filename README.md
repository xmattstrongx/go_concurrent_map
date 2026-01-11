# go_concurrent_map

Go concurrent map lib

## Usage

Note: the purge interval is required and must be greater than 0.
Note: per-entry expiration is optional. If `KeyExpiration` is 0 or `NeverExpire: true`, the entry never expires.
Note: if your Go build cache is not writable, set `GOCACHE` to a writable path when running tests (example below).

```sh
GOCACHE=/path/to/writable/cache go test ./...
```

```go
package main

import (
	"context"
	"time"

	cmap "github.com/xmattstrongx/go_concurrent_map"
)

func main() {
	m, err := cmap.New().
		WithDefaultExpiration(10 * time.Second).
		WithPurgeInterval(1 * time.Second).
		Build()
	if err != nil {
		panic(err)
	}

	m.Set("greeting", []byte("hello"))
	val, ok := m.Get("greeting")
	if ok {
		_ = val
	}

	// Override expiration for a specific entry.
	m.SetEntry("short-lived", cmap.Entry{
		KeyExpiration: 1 * time.Second,
		Value:      []byte("bye"),
	})

	// A non-expiring entry.
	m.SetEntry("permanent", cmap.Entry{
		NeverExpire: true,
		Value:       []byte("forever"),
	})

	ctx, cancel := context.WithCancel(context.Background())
	go m.PurgeExpiredEntries(ctx)
	defer cancel()
}
```
