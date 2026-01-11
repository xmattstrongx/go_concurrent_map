package go_concurrent_map

import (
	"context"
	"sync"
	"testing"
	"time"
)

func buildTestMap(t *testing.T, builder ConcurrentMapBuilder) *Concurrentmap {
	t.Helper()
	m, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %s", err)
	}
	return m
}

func Test_concurrentmap_Set(t *testing.T) {
	c := buildTestMap(t, New().
		WithPurgeInterval(1*time.Second))

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
			t.Cleanup(func() {
				c.Delete(tc.args.key)
			})

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
	c := buildTestMap(t, New().
		WithDefaultExpiration(30 * time.Second).
		WithPurgeInterval(1*time.Second))

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
				KeyExpiration: time.Duration(1 * time.Millisecond),
				Value:         []byte("Hello my friend. Stay awhile, and listen.."),
				timeCreated:   time.Now(),
			},
		},
		{
			name:            "test2",
			key:             "Butcher",
			existAfterPurge: true,
			entry: Entry{
				KeyExpiration: time.Duration(20 * time.Second),
				Value:         []byte("Rarrgh Rarghhh!"),
				timeCreated:   time.Now(),
			},
		},
	}

	t.Cleanup(func() {
		for _, tc := range tests {
			c.Delete(tc.key)
		}
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

func TestConcurrentmap_DefaultExpirationApplied(t *testing.T) {
	defaultExp := 50 * time.Millisecond
	c := buildTestMap(t, New().
		WithDefaultExpiration(defaultExp).
		WithPurgeInterval(1*time.Second))

	c.Set("default-exp", []byte("value"))
	t.Cleanup(func() {
		c.Delete("default-exp")
	})
	entry, ok := c.GetEntry("default-exp")
	if !ok {
		t.Fatalf("Key not found after being set")
	}
	if entry.KeyExpiration != defaultExp {
		t.Fatalf("Expected default expiration %s got %s", defaultExp, entry.KeyExpiration)
	}
	if entry.timeCreated.IsZero() {
		t.Fatalf("Expected timeCreated to be populated")
	}
}

func TestConcurrentmap_SetEntryAndGetEntry(t *testing.T) {
	c := buildTestMap(t, New().
		WithPurgeInterval(1*time.Second))
	entry := Entry{
		KeyExpiration: 10 * time.Second,
		Value:         []byte("payload"),
		timeCreated:   time.Now(),
	}

	c.SetEntry("k1", entry)
	t.Cleanup(func() {
		c.Delete("k1")
	})
	got, ok := c.GetEntry("k1")
	if !ok {
		t.Fatalf("Key not found after SetEntry")
	}
	if string(got.Value) != string(entry.Value) {
		t.Fatalf("Expected %s got %s", string(entry.Value), string(got.Value))
	}
	if got.KeyExpiration != entry.KeyExpiration {
		t.Fatalf("Expected expiration %s got %s", entry.KeyExpiration, got.KeyExpiration)
	}
}


func TestConcurrentmap_DeleteMissing(t *testing.T) {
	c := buildTestMap(t, New().
		WithPurgeInterval(1*time.Second))
	c.Delete("missing")
	if _, ok := c.Get("missing"); ok {
		t.Fatalf("Expected missing key to remain absent after delete")
	}
}

func TestConcurrentmap_PurgeExpiredEntries_DefaultExpiration(t *testing.T) {
	c := buildTestMap(t, New().
		WithDefaultExpiration(10 * time.Millisecond).
		WithPurgeInterval(1*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		c.PurgeExpiredEntries(ctx)
		close(done)
	}()

	c.Set("expiring", []byte("value"))
	t.Cleanup(func() {
		c.Delete("expiring")
	})

	deadline := time.Now().Add(1500 * time.Millisecond)
	for {
		if _, ok := c.Get("expiring"); !ok {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("Entry did not expire within expected time")
		}
		time.Sleep(20 * time.Millisecond)
	}

	<-done
}

func TestConcurrentmap_PurgeExpiredEntries_ContextCancel(t *testing.T) {
	c := buildTestMap(t, New().
		WithPurgeInterval(1*time.Second))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		c.PurgeExpiredEntries(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("PurgeExpiredEntries did not stop after context cancel")
	}
}

func TestConcurrentmap_SetEntry_NoExpiration(t *testing.T) {
	c := buildTestMap(t, New().
		WithDefaultExpiration(10 * time.Millisecond).
		WithPurgeInterval(1*time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		c.PurgeExpiredEntries(ctx)
		close(done)
	}()

	c.SetEntry("permanent", Entry{
		NeverExpire: true,
		Value:         []byte("persist"),
	})
	t.Cleanup(func() {
		c.Delete("permanent")
	})

	time.Sleep(1200 * time.Millisecond)
	if _, ok := c.Get("permanent"); !ok {
		t.Fatalf("Expected entry to remain without expiration")
	}

	<-done
}
