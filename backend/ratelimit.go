package main

import (
	"sync"
	"time"
)

// limiter is a simple in-memory fixed-window rate limiter, keyed by client IP.
// State lives in memory only: it resets on restart and is NOT shared across
// replicas. Fine for a single-instance demo — not real distributed protection.
type limiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	max    int
	window time.Duration
}

func newLimiter(maxPerMinute int) *limiter {
	return &limiter{
		hits:   make(map[string][]time.Time),
		max:    maxPerMinute,
		window: time.Minute,
	}
}

// allow reports whether ip may make another request at time now.
func (l *limiter) allow(ip string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := now.Add(-l.window)
	kept := l.hits[ip][:0]
	for _, t := range l.hits[ip] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= l.max {
		l.hits[ip] = kept
		return false
	}
	l.hits[ip] = append(kept, now)
	return true
}
