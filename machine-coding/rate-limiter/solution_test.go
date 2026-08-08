package ratelimiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeClock is a controllable Clock for deterministic tests — no real sleeps.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock(start time.Time) *fakeClock {
	return &fakeClock{now: start}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	c.mu.Unlock()
}

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

// ---- TokenBucketLimiter ----

func TestTokenBucket_AllowsUpToCapacityThenBlocks(t *testing.T) {
	clock := newFakeClock(epoch)
	l := NewTokenBucketLimiter(3, 0, clock) // no refill, isolate burst behavior

	for i := 0; i < 3; i++ {
		if !l.Allow("user1") {
			t.Fatalf("Allow() call %d = false, want true (within capacity)", i)
		}
	}
	if l.Allow("user1") {
		t.Fatal("Allow() after exhausting capacity = true, want false")
	}
}

func TestTokenBucket_RefillsOverTime(t *testing.T) {
	clock := newFakeClock(epoch)
	l := NewTokenBucketLimiter(2, 1, clock) // refill 1 token/sec

	if !l.Allow("user1") || !l.Allow("user1") {
		t.Fatal("expected first two Allow() calls to succeed")
	}
	if l.Allow("user1") {
		t.Fatal("Allow() with an empty bucket = true, want false")
	}

	clock.Advance(1 * time.Second) // refills exactly 1 token
	if !l.Allow("user1") {
		t.Fatal("Allow() after 1s refill = false, want true")
	}
	if l.Allow("user1") {
		t.Fatal("Allow() immediately after consuming the refilled token = true, want false")
	}
}

func TestTokenBucket_RefillCappedAtCapacity(t *testing.T) {
	clock := newFakeClock(epoch)
	l := NewTokenBucketLimiter(2, 1, clock)

	l.Allow("user1")
	l.Allow("user1") // bucket now empty

	clock.Advance(100 * time.Second) // would overflow capacity if uncapped

	successes := 0
	for i := 0; i < 5; i++ {
		if l.Allow("user1") {
			successes++
		}
	}
	if successes != 2 {
		t.Fatalf("successes after huge time advance = %d, want 2 (capped at capacity)", successes)
	}
}

func TestTokenBucket_KeysAreIndependent(t *testing.T) {
	clock := newFakeClock(epoch)
	l := NewTokenBucketLimiter(1, 0, clock)

	if !l.Allow("user1") {
		t.Fatal("Allow(user1) first call = false, want true")
	}
	if !l.Allow("user2") {
		t.Fatal("Allow(user2) = false, want true — different key, should have its own bucket")
	}
	if l.Allow("user1") {
		t.Fatal("Allow(user1) second call = true, want false — its own bucket should be empty")
	}
}

func TestTokenBucket_ConcurrentAllowNeverOverAdmits(t *testing.T) {
	clock := newFakeClock(epoch)
	const capacity = 20
	l := NewTokenBucketLimiter(capacity, 0, clock)

	var admitted int32
	var wg sync.WaitGroup
	const goroutines = 200
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if l.Allow("shared-key") {
				atomic.AddInt32(&admitted, 1)
			}
		}()
	}
	wg.Wait()

	if admitted != capacity {
		t.Fatalf("admitted = %d, want exactly %d (capacity), out of %d concurrent attempts", admitted, capacity, goroutines)
	}
}

// ---- SlidingWindowLimiter ----

func TestSlidingWindow_AllowsUpToMaxThenBlocks(t *testing.T) {
	clock := newFakeClock(epoch)
	l := NewSlidingWindowLimiter(3, time.Second, clock)

	for i := 0; i < 3; i++ {
		if !l.Allow("user1") {
			t.Fatalf("Allow() call %d = false, want true", i)
		}
	}
	if l.Allow("user1") {
		t.Fatal("Allow() after hitting the limit = true, want false")
	}
}

func TestSlidingWindow_OldRequestsExpireOutOfWindow(t *testing.T) {
	clock := newFakeClock(epoch)
	l := NewSlidingWindowLimiter(2, time.Second, clock)

	l.Allow("user1")
	l.Allow("user1")
	if l.Allow("user1") {
		t.Fatal("Allow() at capacity = true, want false")
	}

	clock.Advance(1100 * time.Millisecond) // both prior requests now outside the window
	if !l.Allow("user1") {
		t.Fatal("Allow() after the window fully elapsed = false, want true")
	}
}

func TestSlidingWindow_PartialExpiryKeepsRecentRequests(t *testing.T) {
	clock := newFakeClock(epoch)
	l := NewSlidingWindowLimiter(2, time.Second, clock)

	l.Allow("user1")                     // t=0
	clock.Advance(600 * time.Millisecond) // t=600ms
	l.Allow("user1")                      // t=600ms; both requests within [{-1000ms},600ms]
	if l.Allow("user1") {
		t.Fatal("3rd Allow() within the window = true, want false")
	}

	clock.Advance(500 * time.Millisecond) // t=1100ms; first request (t=0) now expired, second (t=600ms) still within window
	if !l.Allow("user1") {
		t.Fatal("Allow() after the first request expired = false, want true")
	}
	if l.Allow("user1") {
		t.Fatal("Allow() immediately after = true, want false (two requests — t=600ms and the new one — still within window)")
	}
}

func TestSlidingWindow_KeysAreIndependent(t *testing.T) {
	clock := newFakeClock(epoch)
	l := NewSlidingWindowLimiter(1, time.Second, clock)

	if !l.Allow("user1") {
		t.Fatal("Allow(user1) = false, want true")
	}
	if !l.Allow("user2") {
		t.Fatal("Allow(user2) = false, want true — independent key")
	}
}

func TestSlidingWindow_ConcurrentAllowNeverOverAdmits(t *testing.T) {
	clock := newFakeClock(epoch)
	const maxRequests = 15
	l := NewSlidingWindowLimiter(maxRequests, time.Second, clock)

	var admitted int32
	var wg sync.WaitGroup
	const goroutines = 200
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if l.Allow("shared-key") {
				atomic.AddInt32(&admitted, 1)
			}
		}()
	}
	wg.Wait()

	if admitted != maxRequests {
		t.Fatalf("admitted = %d, want exactly %d, out of %d concurrent attempts", admitted, maxRequests, goroutines)
	}
}

// Both implementations satisfy the same interface.
var (
	_ Limiter = (*TokenBucketLimiter)(nil)
	_ Limiter = (*SlidingWindowLimiter)(nil)
)
