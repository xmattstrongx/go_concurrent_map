package go_concurrent_map

import (
	"context"
	"sync"
	"time"
)

const DefaultExpiration = time.Duration(30 * time.Second)

type Concurrentmap struct {
	internal          map[string]Entry
	mtx               sync.RWMutex
	defaultExpiration time.Duration
	purgeInterval     time.Duration
}

type ConcurrentMapBuilder interface {
	WithPurgeInterval(interval time.Duration) ConcurrentMapBuilder
	WithDefaultExpiration(interval time.Duration) ConcurrentMapBuilder
	Build() *Concurrentmap
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

func (c *concurrentMapBuilder) Build() *Concurrentmap {
	return &Concurrentmap{
		internal:          make(map[string]Entry),
		defaultExpiration: c.defaultExpiration,
		purgeInterval:     c.purgeInterval,
	}
}

type EntryBuilder interface {
	WithExpiration(expiration time.Duration) EntryBuilder
	WithDefaultExpiration() EntryBuilder
	WithValue([]byte) EntryBuilder
	Build() *Entry
}

type entryBuilder struct {
	Value      []byte
	Expiration time.Duration
}

type Entry struct {
	Expiration time.Duration
	Value      []byte
	setTime    time.Time
}

func NewEntry() entryBuilder {
	return entryBuilder{}
}

func (e *entryBuilder) WithExpiration(expiration time.Duration) EntryBuilder {
	e.Expiration = expiration
	return e
}

func (e *entryBuilder) WithDefaultExpiration() EntryBuilder {
	e.Expiration = DefaultExpiration
	return e
}

func (e *entryBuilder) WithValue(value []byte) EntryBuilder {
	e.Value = value
	return e
}

func (e *entryBuilder) Build() *Entry {
	entry := &Entry{
		Expiration: e.Expiration,
		Value:      e.Value,
	}
	return entry
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
		setTime: time.Now(),
		Value:   value,
	})
}

func (c *Concurrentmap) SetEntry(key string, e Entry) {
	c.mtx.Lock()
	c.internal[key] = e
	c.mtx.Unlock()
}

func (c *Concurrentmap) Delete(key string) {
	c.mtx.Lock()
	delete(c.internal, key)
	c.mtx.Unlock()
}

// func (c *Concurrentmap) PurgeExpiredEntries(ctx context.Context) {
// 	retry := time.After(0)

// 	for {
// 		select {
// 		case <-retry:
// 			retry = time.After(c.purgeInterval)
// 			c.mtx.RLock()
// 			for k, v := range c.internal {
// 				if v.Expiration > 0 && time.Since(v.setTime) > v.Expiration {
// 					c.mtx.RUnlock()
// 					c.Delete(k)
// 					c.mtx.RLock()
// 				}
// 			}
// 			c.mtx.RUnlock()

// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }

func (c *Concurrentmap) PurgeExpiredEntriesWithLockSpaghetti(ctx context.Context) {
	retry := time.After(0)

	for {
		select {
		case <-retry:
			retry = time.After(c.purgeInterval)
			c.mtx.RLock()
			for k, v := range c.internal {
				if v.Expiration > 0 && time.Since(v.setTime) > v.Expiration {
					c.mtx.RUnlock()
					c.Delete(k)
					c.mtx.RLock()
				}
			}
			c.mtx.RUnlock()

		case <-ctx.Done():
			return
		}
	}
}

func (c *Concurrentmap) PurgeExpiredEntriesWithExtraDeleteFuncCall(ctx context.Context) {
	retry := time.After(0)

	for {
		select {
		case <-retry:
			retry = time.After(c.purgeInterval)
			keysToDelete := make(map[string]struct{})
			c.mtx.RLock()
			for k, v := range c.internal {
				if v.Expiration > 0 && time.Since(v.setTime) > v.Expiration {
					keysToDelete[k] = struct{}{}
				}
			}
			c.mtx.RUnlock()
			for k := range keysToDelete {
				c.Delete(k)
			}

		case <-ctx.Done():
			return
		}
	}
}
