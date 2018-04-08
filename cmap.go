package go_concurrent_map

import (
	"sync"
	"time"
)

type concurrentmap struct {
	internal map[string]Entry
	mtx      sync.RWMutex
}

type Entry struct {
	Expiration time.Duration
	Value      []byte
	setTime    time.Time
}

func NewConcurrentMap() *concurrentmap {
	return &concurrentmap{
		internal: make(map[string]Entry),
	}
}

func (c *concurrentmap) Get(key string) (value []byte, ok bool) {
	e, ok := c.GetEntry(key)
	return e.Value, ok
}

func (c *concurrentmap) GetEntry(key string) (entry Entry, ok bool) {
	c.mtx.RLock()
	entry, ok = c.internal[key]
	c.mtx.RUnlock()
	return entry, ok
}

func (c *concurrentmap) Set(key string, value []byte) {
	c.SetEntry(key, Entry{
		setTime: time.Now(),
		Value:   value,
	})
}

func (c *concurrentmap) SetEntry(key string, e Entry) {
	c.mtx.Lock()
	c.internal[key] = e
	c.mtx.Unlock()
}

func (c *concurrentmap) Delete(key string) {
	c.mtx.Lock()
	delete(c.internal, key)
	c.mtx.Unlock()
}
