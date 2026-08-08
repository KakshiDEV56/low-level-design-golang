package countingsemaphore

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestSemaphore_AsMutex uses a semaphore initialized to 1 to protect a shared counter
// against concurrent increments, and checks the final value is exact — any lost update
// means Wait/Signal aren't providing real mutual exclusion.
func TestSemaphore_AsMutex(t *testing.T) {
	sem := NewSemaphore(1)
	count := 0

	const goroutines = 100
	const incrementsEach = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsEach; j++ {
				sem.Wait()
				count++ // deliberately unsynchronized read-modify-write, protected only by the semaphore
				sem.Signal()
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
		t.Fatal("timed out — Wait/Signal used as a mutex never let all goroutines finish (possible deadlock)")
	}

	want := goroutines * incrementsEach
	if count != want {
		t.Fatalf("count = %d, want %d — the semaphore did not provide real mutual exclusion", count, want)
	}
}

// TestSemaphore_BlocksWhenNegative checks that Wait() actually blocks a goroutine when
// the semaphore's value would go negative, and that it's the Signal() that unblocks it
// — not a spurious wakeup or a busy-wait that just happens to proceed anyway.
func TestSemaphore_BlocksWhenNegative(t *testing.T) {
	sem := NewSemaphore(0)
	proceeded := make(chan struct{})

	go func() {
		sem.Wait()
		close(proceeded)
	}()

	select {
	case <-proceeded:
		t.Fatal("Wait() returned before Signal() was ever called — it should have blocked")
	case <-time.After(200 * time.Millisecond):
		// Good: still blocked after a reasonable wait.
	}

	sem.Signal()

	select {
	case <-proceeded:
		// Good: Signal() unblocked the waiter.
	case <-time.After(2 * time.Second):
		t.Fatal("Wait() never returned after Signal() was called")
	}
}

// TestSemaphore_SignalWakesExactlyOne verifies the core semaphore invariant: each
// Signal() call unblocks AT MOST one waiting goroutine, not all of them.
func TestSemaphore_SignalWakesExactlyOne(t *testing.T) {
	sem := NewSemaphore(0)
	const waiters = 5

	var proceeded int32
	var wg sync.WaitGroup
	wg.Add(waiters)
	for i := 0; i < waiters; i++ {
		go func() {
			defer wg.Done()
			sem.Wait()
			atomic.AddInt32(&proceeded, 1)
		}()
	}

	// Give all goroutines a chance to reach Wait() and block.
	time.Sleep(200 * time.Millisecond)

	for i := 1; i <= waiters; i++ {
		sem.Signal()
		// After each individual signal, give the woken goroutine time to run, then
		// check that exactly `i` goroutines have proceeded so far — not more.
		time.Sleep(100 * time.Millisecond)
		if got := atomic.LoadInt32(&proceeded); got != int32(i) {
			t.Fatalf("after %d Signal() call(s), %d goroutines proceeded, want exactly %d", i, got, i)
		}
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("not all waiters finished even after signaling once per waiter")
	}
}

// TestSemaphore_SignalBeforeWaitIsRemembered checks that signaling a semaphore with no
// current waiters isn't lost — a subsequent Wait() should proceed immediately, since the
// semaphore's value (not just "is someone currently blocked") is what matters.
func TestSemaphore_SignalBeforeWaitIsRemembered(t *testing.T) {
	sem := NewSemaphore(0)
	sem.Signal()
	sem.Signal()
	sem.Signal()

	done := make(chan struct{})
	go func() {
		sem.Wait()
		sem.Wait()
		sem.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait() blocked even though Signal() had already been called 3 times before any Wait()")
	}
}

// TestSemaphore_Multiplex uses a semaphore initialized to n as a bounded-concurrency
// limiter and checks that at most n goroutines are ever in the "critical section" at
// once, using -race to also confirm the bookkeeping itself is race-free.
func TestSemaphore_Multiplex(t *testing.T) {
	const n = 4
	const goroutines = 50
	sem := NewSemaphore(n)

	var current int32
	var maxObserved int32
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			sem.Wait()
			cur := atomic.AddInt32(&current, 1)
			for {
				max := atomic.LoadInt32(&maxObserved)
				if cur <= max || atomic.CompareAndSwapInt32(&maxObserved, max, cur) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt32(&current, -1)
			sem.Signal()
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
		t.Fatal("timed out running the multiplex workload (possible deadlock)")
	}

	if maxObserved > n {
		t.Fatalf("observed %d goroutines concurrently past Wait(), want at most %d", maxObserved, n)
	}
	if maxObserved < n {
		t.Errorf("observed at most %d concurrent goroutines, want to see the full limit of %d reached at some point (test may be too fast/slow, or Wait/Signal are over-restrictive)", maxObserved, n)
	}
}

// TestSemaphore_StressManyGoroutines is a broader smoke test with many goroutines doing
// short Wait/Signal cycles, to help surface race conditions in the implementation.
// Run with `go test -race`.
func TestSemaphore_StressManyGoroutines(t *testing.T) {
	sem := NewSemaphore(3)
	const goroutines = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			sem.Wait()
			sem.Signal()
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
		t.Fatal("timed out — stress test did not complete (possible deadlock)")
	}
}
