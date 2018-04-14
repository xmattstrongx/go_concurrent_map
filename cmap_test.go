package go_concurrent_map

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var tests1 []testcase = generateTestCases(1)
var tests10 []testcase = generateTestCases(10)
var tests100 []testcase = generateTestCases(100)
var tests1000 []testcase = generateTestCases(1000)
var tests10000 []testcase = generateTestCases(10000)
var tests100000 []testcase = generateTestCases(100000)
var tests1000000 []testcase = generateTestCases(1000000)

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
		c.PurgeExpiredEntriesWithLockSpaghetti(ctx)
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

func benchmarkPurgeWithLockSpaghetti(tests []testcase, b *testing.B) {
	for n := 0; n < b.N; n++ {
		doThePurgeWithLockSpaghetti(tests)
	}
}

func BenchmarkPurgeWithLockSpaghetti1(b *testing.B) {
	benchmarkPurgeWithLockSpaghetti(tests1, b)
}
func BenchmarkPurgeWithLockSpaghetti10(b *testing.B) {
	benchmarkPurgeWithLockSpaghetti(tests10, b)
}
func BenchmarkPurgeWithLockSpaghetti100(b *testing.B) {
	benchmarkPurgeWithLockSpaghetti(tests100, b)
}
func BenchmarkPurgeWithLockSpaghetti1000(b *testing.B) {
	benchmarkPurgeWithLockSpaghetti(tests1000, b)
}
func BenchmarkPurgeWithLockSpaghetti10000(b *testing.B) {
	benchmarkPurgeWithLockSpaghetti(tests100000, b)
}
func BenchmarkPurgeWithLockSpaghetti1000000(b *testing.B) {
	benchmarkPurgeWithLockSpaghetti(tests1000000, b)
}

func doThePurgeWithLockSpaghetti(tests []testcase) {
	c := New().
		WithDefaultExpiration(30 * time.Second).
		WithPurgeInterval(200 * time.Millisecond).
		Build()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1*time.Second))
		c.PurgeExpiredEntriesWithLockSpaghetti(ctx)
		defer cancel()
		defer wg.Done()
	}()
	for _, tc := range tests {
		c.SetEntry(tc.key, tc.entry)
		if _, ok := c.GetEntry(tc.key); !ok {
			log.Fatalf("Key %s not found after being set", tc.key)
		}
	}

	wg.Wait()

	for _, tc := range tests {
		if _, ok := c.GetEntry(tc.key); ok != tc.existAfterPurge {
			log.Fatalf("After PurgeExpiredEntries have %t but want %t", ok, tc.existAfterPurge)
		}
	}
}

func benchmarkPurgeWithExtraDelete(tests []testcase, b *testing.B) {
	for n := 0; n < b.N; n++ {
		doThePurgeWithExtraDelete(tests)
	}
}

func BenchmarkPurgeWithExtraDelete1(b *testing.B) {
	benchmarkPurgeWithExtraDelete(tests1, b)
}
func BenchmarkPurgeWithExtraDelete100(b *testing.B) {
	benchmarkPurgeWithExtraDelete(tests10, b)
}
func BenchmarkPurgeWithExtraDelete1000(b *testing.B) {
	benchmarkPurgeWithExtraDelete(tests100, b)
}
func BenchmarkPurgeWithExtraDelete10000(b *testing.B) {
	benchmarkPurgeWithExtraDelete(tests1000, b)
}
func BenchmarkPurgeWithExtraDelete100000(b *testing.B) {
	benchmarkPurgeWithExtraDelete(tests10000, b)
}
func BenchmarkPurgeWithExtraDelete1000000(b *testing.B) {
	benchmarkPurgeWithExtraDelete(tests100000, b)
}

func doThePurgeWithExtraDelete(tests []testcase) {
	c := New().
		WithDefaultExpiration(30 * time.Second).
		WithPurgeInterval(200 * time.Millisecond).
		Build()

	// tests := generateTestCases(mapEntries)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1*time.Second))
		c.PurgeExpiredEntriesWithExtraDeleteFuncCall(ctx)
		defer cancel()
		defer wg.Done()
	}()
	for _, tc := range tests {
		c.SetEntry(tc.key, tc.entry)
		if _, ok := c.GetEntry(tc.key); !ok {
			log.Fatalf("Key %s not found after being set", tc.key)
		}
	}

	wg.Wait()

	for _, tc := range tests {
		if _, ok := c.GetEntry(tc.key); ok != tc.existAfterPurge {
			log.Fatalf("After PurgeExpiredEntries have %t but want %t", ok, tc.existAfterPurge)
		}
	}
}

type testcase struct {
	name            string
	key             string
	existAfterPurge bool
	entry           Entry
}

func generateTestCases(amountToGenerate int) []testcase {
	testcases := []testcase{}
	t1 := testcase{
		name:            "test1",
		key:             "placeholder",
		existAfterPurge: false,
		entry: Entry{
			Expiration: time.Duration(1 * time.Millisecond),
			Value:      []byte("Hello my friend. Stay awhile, and listen.."),
			setTime:    time.Now(),
		},
	}

	for i := 0; i < amountToGenerate; i++ {
		t1.key = RandomString(12)
		testcases = append(testcases, t1)
	}
	return testcases
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
