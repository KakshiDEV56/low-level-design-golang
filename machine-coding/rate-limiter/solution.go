package ratelimiter

import "time"

/*
See README.md for the full problem statement, requirements, and rubric.

Clock and RealClock below are already complete — do not modify them. They exist so
your limiters never call time.Now() directly, which is what makes them testable
without real sleeps.
*/

// Clock is a source of the current time. Do not modify.
type Clock interface {
	Now() time.Time
}

// RealClock is the production Clock, backed by the real wall clock. Do not modify.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// Limiter is the common interface both rate-limiting algorithms must implement.
type Limiter interface {
	// Allow reports whether a request for key is allowed right now, and if so,
	// counts it against that key's quota.
	Allow(key string) bool
}

// TODO: implement TokenBucketLimiter.
//   - NewTokenBucketLimiter(capacity int, refillPerSecond float64, clock Clock) *TokenBucketLimiter
//   - Allow(key string) bool
//
// Each key gets its own bucket of `capacity` tokens, refilling continuously at
// `refillPerSecond` tokens/second (based on elapsed time since that key's last
// refill, via the injected Clock — never time.Now() directly). Allow consumes one
// token if available and returns true; otherwise returns false without blocking.
// Must be safe for concurrent use.
type TokenBucketLimiter struct {
	// TODO: add your fields here
}

func NewTokenBucketLimiter(capacity int, refillPerSecond float64, clock Clock) *TokenBucketLimiter {
	// TODO: implement
	return &TokenBucketLimiter{}
}

func (l *TokenBucketLimiter) Allow(key string) bool {
	// TODO: implement
	return false
}

// TODO: implement SlidingWindowLimiter.
//   - NewSlidingWindowLimiter(maxRequests int, window time.Duration, clock Clock) *SlidingWindowLimiter
//   - Allow(key string) bool
//
// Each key may make at most `maxRequests` requests in any trailing window of duration
// `window`, measured via the injected Clock. Must be safe for concurrent use.
type SlidingWindowLimiter struct {
	// TODO: add your fields here
}

func NewSlidingWindowLimiter(maxRequests int, window time.Duration, clock Clock) *SlidingWindowLimiter {
	// TODO: implement
	return &SlidingWindowLimiter{}
}

func (l *SlidingWindowLimiter) Allow(key string) bool {
	// TODO: implement
	return false
}
