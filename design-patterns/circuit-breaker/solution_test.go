package circuitbreaker

import (
	"errors"
	"sync"
	"testing"
	"time"
)

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock(start time.Time) *fakeClock { return &fakeClock{now: start} }

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
var errBoom = errors.New("boom")

func failFn() error    { return errBoom }
func succeedFn() error { return nil }

func TestCircuitBreaker_StartsClosed(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second, newFakeClock(epoch))
	if cb.State() != Closed {
		t.Fatalf("State() at start = %v, want Closed", cb.State())
	}
}

func TestCircuitBreaker_ClosedCallsThroughAndPropagatesResult(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second, newFakeClock(epoch))
	if err := cb.Call(succeedFn); err != nil {
		t.Fatalf("Call(succeed) = %v, want nil", err)
	}
	if err := cb.Call(failFn); !errors.Is(err, errBoom) {
		t.Fatalf("Call(fail) = %v, want errBoom", err)
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second, newFakeClock(epoch))
	cb.Call(failFn)
	cb.Call(failFn)
	if cb.State() != Closed {
		t.Fatalf("State() after 2/3 failures = %v, want Closed", cb.State())
	}
	cb.Call(failFn)
	if cb.State() != Open {
		t.Fatalf("State() after 3/3 failures = %v, want Open", cb.State())
	}
}

func TestCircuitBreaker_SuccessResetsFailureCount(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second, newFakeClock(epoch))
	cb.Call(failFn)
	cb.Call(failFn)
	cb.Call(succeedFn) // resets the count
	cb.Call(failFn)
	cb.Call(failFn)
	if cb.State() != Closed {
		t.Fatalf("State() = %v, want Closed (a success in between should have reset the failure streak)", cb.State())
	}
}

func TestCircuitBreaker_OpenRejectsWithoutCallingFn(t *testing.T) {
	clock := newFakeClock(epoch)
	cb := NewCircuitBreaker(1, time.Second, clock)
	cb.Call(failFn) // opens the breaker

	calls := 0
	err := cb.Call(func() error { calls++; return nil })
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("Call() while Open = %v, want ErrCircuitOpen", err)
	}
	if calls != 0 {
		t.Fatalf("fn was called %d times while circuit was Open, want 0", calls)
	}
}

func TestCircuitBreaker_TransitionsToHalfOpenAfterCooldownAndProbeSucceeds(t *testing.T) {
	clock := newFakeClock(epoch)
	cb := NewCircuitBreaker(1, time.Second, clock)
	cb.Call(failFn) // opens

	clock.Advance(1100 * time.Millisecond) // cooldown elapsed

	calls := 0
	err := cb.Call(func() error { calls++; return nil }) // this call IS the probe
	if err != nil {
		t.Fatalf("probe Call() = %v, want nil (probe succeeded)", err)
	}
	if calls != 1 {
		t.Fatalf("fn was called %d times for the probe, want 1", calls)
	}
	if cb.State() != Closed {
		t.Fatalf("State() after a successful probe = %v, want Closed", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenProbeFailureReopens(t *testing.T) {
	clock := newFakeClock(epoch)
	cb := NewCircuitBreaker(1, time.Second, clock)
	cb.Call(failFn) // opens
	clock.Advance(1100 * time.Millisecond)

	err := cb.Call(failFn) // probe fails
	if !errors.Is(err, errBoom) {
		t.Fatalf("probe Call() = %v, want errBoom", err)
	}
	if cb.State() != Open {
		t.Fatalf("State() after a failed probe = %v, want Open", cb.State())
	}
}

func TestCircuitBreaker_OnlyOneProbeAllowedPerCooldown(t *testing.T) {
	clock := newFakeClock(epoch)
	cb := NewCircuitBreaker(1, time.Second, clock)
	cb.Call(failFn) // opens
	clock.Advance(1100 * time.Millisecond)
	cb.Call(failFn) // probe fails, reopens (cooldown timer resets to now)

	// Immediately after the failed probe, with no further time advance, the breaker
	// must be freshly Open again — the very next call must be rejected without
	// touching fn.
	calls := 0
	err := cb.Call(func() error { calls++; return nil })
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("Call() immediately after a failed probe = %v, want ErrCircuitOpen", err)
	}
	if calls != 0 {
		t.Fatalf("fn was called %d times right after a failed probe, want 0 (fresh cooldown must apply)", calls)
	}
}

func TestCircuitBreaker_RejectsBeforeCooldownElapses(t *testing.T) {
	clock := newFakeClock(epoch)
	cb := NewCircuitBreaker(1, time.Second, clock)
	cb.Call(failFn) // opens

	clock.Advance(500 * time.Millisecond) // still within the 1s cooldown
	if err := cb.Call(succeedFn); !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("Call() before cooldown elapsed = %v, want ErrCircuitOpen", err)
	}
	if cb.State() != Open {
		t.Fatalf("State() before cooldown elapsed = %v, want Open", cb.State())
	}
}

func TestCircuitBreaker_FullCycleClosedOpenHalfOpenClosed(t *testing.T) {
	clock := newFakeClock(epoch)
	cb := NewCircuitBreaker(2, 5*time.Second, clock)

	if cb.State() != Closed {
		t.Fatal("expected Closed at start")
	}
	cb.Call(failFn)
	cb.Call(failFn)
	if cb.State() != Open {
		t.Fatal("expected Open after threshold failures")
	}

	clock.Advance(6 * time.Second)
	cb.Call(succeedFn) // probe succeeds
	if cb.State() != Closed {
		t.Fatal("expected Closed after a successful probe")
	}

	// And the breaker should require a fresh threshold's worth of failures again.
	cb.Call(failFn)
	if cb.State() != Closed {
		t.Fatal("expected Closed after a single failure post-recovery (threshold not yet reached)")
	}
}
