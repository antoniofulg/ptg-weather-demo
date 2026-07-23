// Package cache provides a tiny, concurrency-safe in-memory TTL cache.
//
// A Cache[T] memoizes values by string key for a fixed time-to-live. A Get
// within the TTL returns the stored value; once the TTL has elapsed for a key,
// Get reports a miss so the caller can refresh it (for example, by re-fetching
// from the network). It is safe for concurrent use by multiple goroutines and
// depends only on the standard library.
package cache

import (
	"sync"
	"time"
)

// Cache is a concurrency-safe in-memory cache whose entries expire after a
// fixed TTL. The zero value is not usable; construct one with New.
type Cache[T any] struct {
	ttl time.Duration

	mu      sync.Mutex
	entries map[string]entry[T]

	// now returns the current time. It defaults to time.Now and exists so
	// tests can drive expiry deterministically without wall-clock sleeps.
	now func() time.Time
}

type entry[T any] struct {
	value   T
	expires time.Time
}

// New returns an empty Cache whose entries expire ttl after they are stored.
func New[T any](ttl time.Duration) *Cache[T] {
	return &Cache[T]{
		ttl:     ttl,
		entries: make(map[string]entry[T]),
		now:     time.Now,
	}
}

// Set stores value under key, (re)setting its expiry to ttl from now.
func (c *Cache[T]) Set(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry[T]{value: value, expires: c.now().Add(c.ttl)}
}

// Get returns the value stored under key and true when the entry exists and its
// TTL has not elapsed. It returns the zero value of T and false on a miss or
// once the entry has expired.
func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok || !c.now().Before(e.expires) {
		var zero T
		return zero, false
	}
	return e.value, true
}
