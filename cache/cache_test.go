package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

// reading mirrors the kind of value the CLI memoizes (a fetched temperature),
// exercising the generic Cache with a struct payload.
type reading struct {
	tempC float64
	unit  string
}

// newClocked returns a cache whose clock is driven by *cur, so tests advance
// time deterministically instead of sleeping on the wall clock.
func newClocked[T any](ttl time.Duration, cur *time.Time) *Cache[T] {
	c := New[T](ttl)
	c.now = func() time.Time { return *cur }
	return c
}

// UT-030: a value stored in the cache is returned on a lookup within the TTL.
func TestGetReturnsValueWithinTTL(t *testing.T) {
	now := time.Unix(0, 0)
	c := newClocked[reading](100*time.Millisecond, &now)

	want := reading{tempC: 21.3, unit: "°C"}
	c.Set("38.72,-9.14", want)

	// Advance partway through the TTL: still a hit.
	now = now.Add(50 * time.Millisecond)
	got, ok := c.Get("38.72,-9.14")
	if !ok {
		t.Fatalf("Get within TTL: got miss, want hit")
	}
	if got != want {
		t.Fatalf("Get within TTL: got %+v, want %+v", got, want)
	}
}

// UT-031: a lookup after the TTL has elapsed reports a miss.
func TestGetReportsMissAfterTTL(t *testing.T) {
	now := time.Unix(0, 0)
	c := newClocked[reading](100*time.Millisecond, &now)

	c.Set("38.72,-9.14", reading{tempC: 21.3, unit: "°C"})

	// Advance past the TTL: the entry has expired.
	now = now.Add(150 * time.Millisecond)
	if got, ok := c.Get("38.72,-9.14"); ok {
		t.Fatalf("Get after TTL: got hit %+v, want miss", got)
	}
}

func TestGetBoundaryIsMiss(t *testing.T) {
	now := time.Unix(0, 0)
	c := newClocked[reading](100*time.Millisecond, &now)

	c.Set("k", reading{tempC: 1, unit: "°C"})

	// Exactly at the expiry instant counts as elapsed -> miss.
	now = now.Add(100 * time.Millisecond)
	if _, ok := c.Get("k"); ok {
		t.Fatalf("Get at expiry boundary: got hit, want miss")
	}
}

func TestGetUnknownKeyIsMiss(t *testing.T) {
	c := New[reading](time.Minute)
	if got, ok := c.Get("never-set"); ok {
		t.Fatalf("Get unknown key: got hit %+v, want miss", got)
	}
}

func TestSetRefreshesExpiry(t *testing.T) {
	now := time.Unix(0, 0)
	c := newClocked[int](100*time.Millisecond, &now)

	c.Set("k", 1)
	now = now.Add(80 * time.Millisecond)
	c.Set("k", 2) // re-Set resets the TTL window from this instant

	now = now.Add(80 * time.Millisecond) // 160ms since first Set, 80ms since re-Set
	got, ok := c.Get("k")
	if !ok {
		t.Fatalf("Get after re-Set within TTL: got miss, want hit")
	}
	if got != 2 {
		t.Fatalf("Get after re-Set: got %d, want 2", got)
	}
}

// Exercises concurrent Get/Set so `go test -race` can flag data races (R3).
func TestConcurrentAccess(t *testing.T) {
	c := New[int](time.Minute)
	const workers = 8

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := "coord-" + strconv.Itoa(id)
			for i := 0; i < 1000; i++ {
				c.Set(key, i)
				c.Get(key)
			}
		}(w)
	}
	wg.Wait()
}
