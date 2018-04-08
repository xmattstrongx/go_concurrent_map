package go_concurrent_map

import (
	"testing"
	"time"
)

func Test_concurrentmap_Set(t *testing.T) {
	// c := NewConcurrentMap()
	c := New().
		WithDefaultExpiration(30 * time.Second).
		WithPurgeInterval(30 * time.Second).
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
			name: "tets1",
			args: args{
				key:   "DeckardCain",
				value: []byte("Hello my friend. Stay awhile, and listen.."),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := c.Get(tt.args.key)
			if ok {
				t.Fatalf("Key should not be present before being set")
			}

			c.Set(tt.args.key, tt.args.value)
			val, ok := c.Get(tt.args.key)
			if !ok {
				t.Fatalf("Key not found after being set")
			}
			if string(val) != string(tt.args.value) {
				t.Fatalf("Expected %s got %s", string(tt.args.value), string(val))
			}

			c.Delete(tt.args.key)
			_, ok = c.Get(tt.args.key)
			if ok {
				t.Fatalf("Key should not be present after deletion")
			}
		})
	}
}
