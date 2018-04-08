package go_concurrent_map

import (
	"context"
	"sync"
	"testing"
	"time"
)

func Test_concurrentmap_Set(t *testing.T) {
	c := New().
		Build()

	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				key:   "DeckardCain",
				value: []byte("Hello my friend. Stay awhile, and listen.."),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := c.Get(tc.args.key)
			if ok {
				t.Fatalf("Key should not be present before being set")
			}

			c.Set(tc.args.key, tc.args.value)
			val, ok := c.Get(tc.args.key)
			if !ok {
				t.Fatalf("Key not found after being set")
			}
			if string(val) != string(tc.args.value) {
				t.Fatalf("Expected %s got %s", string(tc.args.value), string(val))
			}

			c.Delete(tc.args.key)
			_, ok = c.Get(tc.args.key)
			if ok {
				t.Fatalf("Key should not be present after deletion")
			}
		})
	}
}

func TestConcurrentmap_PurgeExpiredEntries(t *testing.T) {
	c := New().
		WithDefaultExpiration(30 * time.Second).
		WithPurgeInterval(20 * time.Millisecond).
		Build()

	tests := []struct {
		name            string
		key             string
		existAfterPurge bool
		entry           Entry
	}{
		{
			name:            "test1",
			key:             "DeckardCain",
			existAfterPurge: false,
			entry: Entry{
				Expiration: time.Duration(1 * time.Millisecond),
				Value:      []byte("Hello my friend. Stay awhile, and listen.."),
				setTime:    time.Now(),
			},
		},
		{
			name:            "test2",
			key:             "Butcher",
			existAfterPurge: true,
			entry: Entry{
				Expiration: time.Duration(20 * time.Second),
				Value:      []byte("Rarrgh Rarghhh!"),
				setTime:    time.Now(),
			},
		},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1*time.Second))
		c.PurgeExpiredEntries(ctx)
		defer cancel()
		defer wg.Done()
	}()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c.SetEntry(tc.key, tc.entry)
			if _, ok := c.GetEntry(tc.key); !ok {
				t.Fatalf("Key %s not found after being set", tc.key)
			}
		})
	}

	wg.Wait()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := c.GetEntry(tc.key); ok != tc.existAfterPurge {
				t.Errorf("After PurgeExpiredEntries have %t but want %t", ok, tc.existAfterPurge)
			}
		})
	}
}
