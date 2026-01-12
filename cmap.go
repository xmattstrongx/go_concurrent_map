package go_concurrent_map

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Concurrentmap struct {
	internal          map[string]Entry
	mtx               sync.RWMutex
	defaultExpiration time.Duration
	purgeInterval     time.Duration
}

type Entry struct {
	KeyExpiration time.Duration
	Value         []byte
	timeCreated   time.Time
	NeverExpire   bool
}

type ConcurrentMapBuilder interface {
	WithPurgeInterval(interval time.Duration) ConcurrentMapBuilder
	WithDefaultExpiration(interval time.Duration) ConcurrentMapBuilder
	Build() (*Concurrentmap, error)
}

type concurrentMapBuilder struct {
	defaultExpiration time.Duration
	purgeInterval     time.Duration
}

func New() ConcurrentMapBuilder {
	return &concurrentMapBuilder{}
}

func (c *concurrentMapBuilder) WithPurgeInterval(interval time.Duration) ConcurrentMapBuilder {
	c.purgeInterval = interval
	return c
}

func (c *concurrentMapBuilder) WithDefaultExpiration(interval time.Duration) ConcurrentMapBuilder {
	c.defaultExpiration = interval
	return c
}

func (c *concurrentMapBuilder) Build() (*Concurrentmap, error) {
	if c.purgeInterval <= 0 {
		return nil, fmt.Errorf("map purge interval time must be set")
	}

	return &Concurrentmap{
		internal:          make(map[string]Entry),
		defaultExpiration: c.defaultExpiration,
		purgeInterval:     c.purgeInterval,
	}, nil
}

func (c *Concurrentmap) Get(key string) (value []byte, ok bool) {
	e, ok := c.GetEntry(key)
	return e.Value, ok
}

func (c *Concurrentmap) GetEntry(key string) (entry Entry, ok bool) {
	c.mtx.RLock()
	entry, ok = c.internal[key]
	c.mtx.RUnlock()
	return entry, ok
}

func (c *Concurrentmap) Set(key string, value []byte) {
	c.SetEntry(key, Entry{
		KeyExpiration: c.defaultExpiration,
		timeCreated:   time.Now(),
		Value:         value,
	})
}

func (c *Concurrentmap) SetEntry(key string, e Entry) {
	if e.NeverExpire || e.KeyExpiration == 0 {
		e.KeyExpiration = 0
	}

	e.timeCreated = time.Now()

	c.mtx.Lock()
	c.internal[key] = e
	c.mtx.Unlock()
}

func (c *Concurrentmap) Delete(key string) {
	c.mtx.Lock()
	delete(c.internal, key)
	c.mtx.Unlock()
}

func (c *Concurrentmap) PurgeExpiredEntries(ctx context.Context) {
	retry := time.After(c.purgeInterval)

	for {
		select {
		case <-retry:
			retry = time.After(c.purgeInterval)
			keysToBeDelete := make(map[string]struct{})
			c.mtx.RLock()
			for k, v := range c.internal {
				if v.KeyExpiration > 0 && time.Since(v.timeCreated) > v.KeyExpiration {
					keysToBeDelete[k] = struct{}{}
				}
			}
			c.mtx.RUnlock()
			for k := range keysToBeDelete {
				c.Delete(k)
			}

		case <-ctx.Done():
			return
		}
	}
}
