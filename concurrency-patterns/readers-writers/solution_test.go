package readerswriters

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ---- RWLock (basic) ----

func TestRWLock_MultipleReadersConcurrently(t *testing.T) {
	l := NewRWLock()
	const readers = 10

	var current int32
	var maxObserved int32
	var wg sync.WaitGroup
	wg.Add(readers)

	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			l.RLock()
			cur := atomic.AddInt32(&current, 1)
			for {
				max := atomic.LoadInt32(&maxObserved)
				if cur <= max || atomic.CompareAndSwapInt32(&maxObserved, max, cur) {
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
			atomic.AddInt32(&current, -1)
			l.RUnlock()
		}()
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out")
	}

	if maxObserved < 2 {
		t.Fatalf("max concurrent readers observed = %d, want readers to overlap (>= 2)", maxObserved)
	}
}

func TestRWLock_WriterExcludesReaders(t *testing.T) {
	l := NewRWLock()
	l.Lock() // writer holds the lock

	readerProceeded := make(chan struct{})
	go func() {
		l.RLock()
		close(readerProceeded)
		l.RUnlock()
	}()

	select {
	case <-readerProceeded:
		t.Fatal("a reader acquired RLock() while a writer held the lock")
	case <-time.After(200 * time.Millisecond):
	}

	l.Unlock()

	select {
	case <-readerProceeded:
	case <-time.After(2 * time.Second):
		t.Fatal("reader never proceeded after the writer released the lock")
	}
}

func TestRWLock_ReadersExcludeWriter(t *testing.T) {
	l := NewRWLock()
	l.RLock() // one reader holds the lock

	writerProceeded := make(chan struct{})
	go func() {
		l.Lock()
		close(writerProceeded)
		l.Unlock()
	}()

	select {
	case <-writerProceeded:
		t.Fatal("a writer acquired Lock() while a reader held the lock")
	case <-time.After(200 * time.Millisecond):
	}

	l.RUnlock()

	select {
	case <-writerProceeded:
	case <-time.After(2 * time.Second):
		t.Fatal("writer never proceeded after the last reader released the lock")
	}
}

func TestRWLock_WritersMutuallyExclusive(t *testing.T) {
	l := NewRWLock()
	count := 0
	const goroutines = 50
	const perGoroutine = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				l.Lock()
				count++
				l.Unlock()
			}
		}()
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out")
	}

	want := goroutines * perGoroutine
	if count != want {
		t.Fatalf("count = %d, want %d — writers were not mutually exclusive", count, want)
	}
}

// ---- FairRWLock (starvation-free) ----

func TestFairRWLock_MultipleReadersConcurrently(t *testing.T) {
	l := NewFairRWLock()
	const readers = 10

	var current int32
	var maxObserved int32
	var wg sync.WaitGroup
	wg.Add(readers)

	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			l.RLock()
			cur := atomic.AddInt32(&current, 1)
			for {
				max := atomic.LoadInt32(&maxObserved)
				if cur <= max || atomic.CompareAndSwapInt32(&maxObserved, max, cur) {
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
			atomic.AddInt32(&current, -1)
			l.RUnlock()
		}()
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out")
	}

	if maxObserved < 2 {
		t.Fatalf("max concurrent readers observed = %d, want readers to overlap (>= 2)", maxObserved)
	}
}

func TestFairRWLock_WritersMutuallyExclusiveAndExcludeReaders(t *testing.T) {
	l := NewFairRWLock()
	count := 0
	const goroutines = 30
	const perGoroutine = 30

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				l.Lock()
				count++
				l.Unlock()
			}
		}()
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out")
	}

	want := goroutines * perGoroutine
	if count != want {
		t.Fatalf("count = %d, want %d", count, want)
	}
}

// TestFairRWLock_NewReadersCannotCutInFrontOfAWaitingWriter is the core no-starve test.
// Sequence, deterministically orchestrated:
//  1. reader1 acquires RLock.
//  2. A writer calls Lock() and (correctly) blocks, since reader1 is inside.
//  3. reader2 arrives AFTER the writer is already waiting, and calls RLock(). In a
//     naive (non-fair) implementation, reader2 could join reader1 immediately since
//     concurrent readers are normally allowed — but the whole point of the fair
//     version is that reader2 must NOT be allowed to cut in front of the already-
//     waiting writer. So reader2 must also block.
//  4. reader1 releases. Now the writer (not reader2) must be the one to proceed next.
//  5. The writer releases. Only now may reader2 proceed.
func TestFairRWLock_NewReadersCannotCutInFrontOfAWaitingWriter(t *testing.T) {
	l := NewFairRWLock()

	l.RLock() // step 1: reader1 inside

	writerProceeded := make(chan struct{})
	go func() {
		l.Lock()
		close(writerProceeded)
	}()

	// Give the writer time to reach and block on Lock().
	time.Sleep(150 * time.Millisecond)
	select {
	case <-writerProceeded:
		t.Fatal("writer acquired Lock() while reader1 was still inside")
	default:
	}

	reader2Proceeded := make(chan struct{})
	go func() {
		l.RLock()
		close(reader2Proceeded)
		l.RUnlock()
	}()

	// step 3: reader2 must NOT be able to proceed while the writer is queued, even
	// though reader1 is still holding its read lock (a naive readers-writers lock
	// would let reader2 in here).
	select {
	case <-reader2Proceeded:
		t.Fatal("reader2 was allowed to acquire RLock() ahead of an already-waiting writer")
	case <-time.After(200 * time.Millisecond):
	}

	l.RUnlock() // step 4: reader1 leaves

	// The writer must be the next one to proceed — not reader2.
	select {
	case <-writerProceeded:
	case <-time.After(2 * time.Second):
		t.Fatal("writer never proceeded after reader1 released")
	}
	select {
	case <-reader2Proceeded:
		t.Fatal("reader2 proceeded before the writer got its turn")
	default:
	}

	l.Unlock() // step 5: writer leaves

	select {
	case <-reader2Proceeded:
	case <-time.After(2 * time.Second):
		t.Fatal("reader2 never proceeded after the writer released")
	}
}

// TestFairRWLock_StressManyGoroutines is a broader smoke test mixing readers and
// writers to help surface races. Run with `go test -race`.
func TestFairRWLock_StressManyGoroutines(t *testing.T) {
	l := NewFairRWLock()
	var shared int64
	const readers = 50
	const writers = 10
	const opsEach = 50

	var wg sync.WaitGroup
	wg.Add(readers + writers)

	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsEach; j++ {
				l.RLock()
				_ = atomic.LoadInt64(&shared)
				l.RUnlock()
			}
		}()
	}
	for i := 0; i < writers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsEach; j++ {
				l.Lock()
				atomic.AddInt64(&shared, 1)
				l.Unlock()
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

	if shared != int64(writers*opsEach) {
		t.Fatalf("shared = %d, want %d", shared, writers*opsEach)
	}
}
