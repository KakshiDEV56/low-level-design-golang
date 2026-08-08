package nostarvemutex

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestNoStarveMutex_MutualExclusion is the basic sanity check: however fairness is
// implemented, it must still be a real mutex (no lost updates).
func TestNoStarveMutex_MutualExclusion(t *testing.T) {
	m := NewNoStarveMutex()
	count := 0
	const goroutines = 50
	const perGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				m.Lock()
				count++
				m.Unlock()
			}
		}()
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("timed out — possible deadlock")
	}

	want := goroutines * perGoroutine
	if count != want {
		t.Fatalf("count = %d, want %d — not a real mutex", count, want)
	}
}

// TestNoStarveMutex_StrictFIFOOrder is the core fairness test. It first grabs the
// mutex from the main goroutine (so every subsequent Lock() call is GUARANTEED to
// block, not spuriously succeed), then starts n goroutines one at a time with a pause
// between each start so each one has time to actually reach and block on Lock() before
// the next is started — establishing a known Lock()-call order of 0, 1, 2, ..., n-1.
// It then releases the mutex and checks the goroutines enter the critical section in
// EXACTLY that order, proving no goroutine can be skipped over by one that queued later.
func TestNoStarveMutex_StrictFIFOOrder(t *testing.T) {
	m := NewNoStarveMutex()
	const n = 8

	m.Lock() // hold the mutex so every goroutine below is guaranteed to block on Lock()

	var mu sync.Mutex
	var order []int

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			m.Lock()
			mu.Lock()
			order = append(order, i)
			mu.Unlock()
			m.Unlock()
		}()
		// Give this goroutine time to actually reach and block on m.Lock() before the
		// next one is started, establishing a known call order.
		time.Sleep(50 * time.Millisecond)
	}

	m.Unlock() // release the initial hold; the queued goroutines should now drain in order

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for all goroutines to pass through the mutex")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(order) != n {
		t.Fatalf("recorded %d critical-section entries, want %d", len(order), n)
	}
	for i, v := range order {
		if v != i {
			t.Fatalf("critical section entry order = %v, want [0,1,2,...,%d] (goroutine %d got skipped ahead of/behind schedule)", order, n-1, v)
		}
	}
}

// TestNoStarveMutex_LateArrivalEventuallyGetsIn is the book's own motivating scenario
// (page 82): two goroutines (A, B) hammer the mutex in a tight loop. A third goroutine
// (C) arrives late and calls Lock() exactly once. With a naive (unfair) mutex this CAN
// starve C forever in the worst case; with a fair mutex, C is guaranteed to get in
// within a bounded amount of time. This is a liveness smoke test alongside the strict
// ordering test above.
func TestNoStarveMutex_LateArrivalEventuallyGetsIn(t *testing.T) {
	m := NewNoStarveMutex()
	stop := make(chan struct{})
	var aCount, bCount int64

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
			}
			m.Lock()
			atomic.AddInt64(&aCount, 1)
			m.Unlock()
		}
	}()
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
			}
			m.Lock()
			atomic.AddInt64(&bCount, 1)
			m.Unlock()
		}
	}()

	// Let A and B contend for a while before C ever shows up.
	time.Sleep(200 * time.Millisecond)

	cDone := make(chan struct{})
	go func() {
		m.Lock()
		m.Unlock()
		close(cDone)
	}()

	select {
	case <-cDone:
		// Good: C got in within a bounded time despite A and B's continuous contention.
	case <-time.After(5 * time.Second):
		close(stop)
		t.Fatal("goroutine C never acquired the mutex while A and B were looping — looks like starvation")
	}

	close(stop)
}

// TestNoStarveMutex_StressManyGoroutines is a broader smoke test to help surface races.
// Run with `go test -race`.
func TestNoStarveMutex_StressManyGoroutines(t *testing.T) {
	m := NewNoStarveMutex()
	const goroutines = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			m.Lock()
			m.Unlock()
		}()
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("timed out — possible deadlock")
	}
}
