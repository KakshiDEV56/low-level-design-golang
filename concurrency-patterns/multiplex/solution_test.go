package multiplex

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestMultiplex_NeverExceedsLimit spawns far more goroutines than the limit and checks
// that the observed concurrency never exceeds n, while also confirming the limit is
// actually reached (not overly conservative).
func TestMultiplex_NeverExceedsLimit(t *testing.T) {
	const n = 5
	const goroutines = 100
	m := NewMultiplex(n)

	var current int32
	var maxObserved int32
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			m.Acquire()
			cur := atomic.AddInt32(&current, 1)
			for {
				max := atomic.LoadInt32(&maxObserved)
				if cur <= max || atomic.CompareAndSwapInt32(&maxObserved, max, cur) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt32(&current, -1)
			m.Release()
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out — possible deadlock in Acquire/Release")
	}

	if maxObserved > n {
		t.Fatalf("observed %d concurrent holders, want at most %d", maxObserved, n)
	}
	if maxObserved < n {
		t.Errorf("observed at most %d concurrent holders, want the full limit of %d to be reached at some point", maxObserved, n)
	}
}

// TestMultiplex_LimitOne behaves like a plain mutex when n=1.
func TestMultiplex_LimitOne(t *testing.T) {
	m := NewMultiplex(1)
	count := 0
	const goroutines = 50
	const perGoroutine = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				m.Acquire()
				count++
				m.Release()
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out")
	}

	want := goroutines * perGoroutine
	if count != want {
		t.Fatalf("count = %d, want %d — n=1 multiplex should behave exactly like a mutex", count, want)
	}
}

// TestMultiplex_AcquireBlocksAtCapacity checks that once the limit is reached, the next
// Acquire() genuinely blocks rather than proceeding.
func TestMultiplex_AcquireBlocksAtCapacity(t *testing.T) {
	m := NewMultiplex(2)
	m.Acquire()
	m.Acquire() // capacity now fully used

	proceeded := make(chan struct{})
	go func() {
		m.Acquire()
		close(proceeded)
	}()

	select {
	case <-proceeded:
		t.Fatal("Acquire() returned even though the multiplex was already at capacity")
	case <-time.After(200 * time.Millisecond):
		// Good: still blocked.
	}

	m.Release() // frees one slot

	select {
	case <-proceeded:
		// Good.
	case <-time.After(2 * time.Second):
		t.Fatal("Acquire() never returned after Release() freed a slot")
	}
}

// TestMultiplex_ReleaseWakesOnlyOneWaiter checks that releasing one slot unblocks
// exactly one waiting goroutine, not all of them.
func TestMultiplex_ReleaseWakesOnlyOneWaiter(t *testing.T) {
	m := NewMultiplex(1)
	m.Acquire() // hold the only slot

	const waiters = 4
	var proceeded int32
	var wg sync.WaitGroup
	wg.Add(waiters)
	for i := 0; i < waiters; i++ {
		go func() {
			defer wg.Done()
			m.Acquire()
			atomic.AddInt32(&proceeded, 1)
		}()
	}

	time.Sleep(200 * time.Millisecond)
	if got := atomic.LoadInt32(&proceeded); got != 0 {
		t.Fatalf("before any Release(), %d waiter(s) already proceeded, want 0", got)
	}

	m.Release()
	time.Sleep(200 * time.Millisecond)
	if got := atomic.LoadInt32(&proceeded); got != 1 {
		t.Fatalf("after one Release(), %d waiter(s) proceeded, want exactly 1", got)
	}

	// Drain the rest so the test doesn't leak goroutines.
	for i := 0; i < waiters-1; i++ {
		m.Release()
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("not all waiters finished after releasing all slots")
	}
}

// TestMultiplex_StressManyGoroutines is a broader smoke test to help surface races.
// Run with `go test -race`.
func TestMultiplex_StressManyGoroutines(t *testing.T) {
	m := NewMultiplex(8)
	const goroutines = 300

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			m.Acquire()
			m.Release()
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out — possible deadlock")
	}
}
