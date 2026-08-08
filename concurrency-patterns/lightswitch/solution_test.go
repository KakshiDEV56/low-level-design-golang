package lightswitch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// countingLocker wraps a real sync.Mutex but also counts how many times Lock/Unlock
// were actually invoked, so tests can verify the Lightswitch calls them exactly once
// per "occupied" transition — not once per goroutine that passes through.
type countingLocker struct {
	mu          sync.Mutex
	lockCalls   int32
	unlockCalls int32
}

func (c *countingLocker) Lock() {
	atomic.AddInt32(&c.lockCalls, 1)
	c.mu.Lock()
}

func (c *countingLocker) Unlock() {
	c.mu.Unlock()
	atomic.AddInt32(&c.unlockCalls, 1)
}

// TestLightswitch_SequentialUsageLocksAndUnlocksEachTime checks the simple case: with
// no overlap, every Lock/Unlock pair on the Lightswitch really does lock/unlock the
// underlying room, since each caller is both first-in and last-out for its own turn.
func TestLightswitch_SequentialUsageLocksAndUnlocksEachTime(t *testing.T) {
	ls := NewLightswitch()
	room := &countingLocker{}

	const n = 10
	for i := 0; i < n; i++ {
		ls.Lock(room)
		ls.Unlock(room)
	}

	if got := atomic.LoadInt32(&room.lockCalls); got != n {
		t.Fatalf("room.Lock() was called %d times, want %d", got, n)
	}
	if got := atomic.LoadInt32(&room.unlockCalls); got != n {
		t.Fatalf("room.Unlock() was called %d times, want %d", got, n)
	}
}

// TestLightswitch_OverlappingEntryLocksOnlyOnce is the core correctness test: it forces
// two goroutines to overlap inside the Lightswitch (the second enters while the first
// is still inside, and leaves before the first does) using explicit channel handshakes
// for determinism, and checks that the underlying room was locked exactly once and
// unlocked exactly once for that whole overlapping session — proving only the true
// first-in and true last-out actually touch the room.
func TestLightswitch_OverlappingEntryLocksOnlyOnce(t *testing.T) {
	ls := NewLightswitch()
	room := &countingLocker{}

	entered1 := make(chan struct{})
	canLeave1 := make(chan struct{})
	goroutine1Done := make(chan struct{})

	go func() {
		ls.Lock(room)
		close(entered1)
		<-canLeave1
		ls.Unlock(room)
		close(goroutine1Done)
	}()

	select {
	case <-entered1:
	case <-time.After(2 * time.Second):
		t.Fatal("goroutine 1 never entered the Lightswitch")
	}

	// Goroutine 2 enters and leaves entirely while goroutine 1 is still inside.
	ls.Lock(room)
	ls.Unlock(room)

	close(canLeave1)
	select {
	case <-goroutine1Done:
	case <-time.After(2 * time.Second):
		t.Fatal("goroutine 1 never left the Lightswitch")
	}

	if got := atomic.LoadInt32(&room.lockCalls); got != 1 {
		t.Fatalf("room.Lock() was called %d times during an overlapping session, want exactly 1", got)
	}
	if got := atomic.LoadInt32(&room.unlockCalls); got != 1 {
		t.Fatalf("room.Unlock() was called %d times during an overlapping session, want exactly 1", got)
	}
}

// TestLightswitch_ExcludesIndependentLocker checks the whole point of the pattern: an
// independent goroutine calling room.Lock() directly (playing the role of a "writer" in
// readers-writers) must be blocked for as long as ANY Lightswitch user is inside, and
// must be able to proceed once the last one leaves.
func TestLightswitch_ExcludesIndependentLocker(t *testing.T) {
	ls := NewLightswitch()
	room := &sync.Mutex{}

	ls.Lock(room) // one "reader" is now inside; room should be locked

	writerDone := make(chan struct{})
	go func() {
		room.Lock()
		close(writerDone)
		room.Unlock()
	}()

	select {
	case <-writerDone:
		t.Fatal("an independent writer acquired the room while a Lightswitch user was still inside")
	case <-time.After(200 * time.Millisecond):
		// Good: writer is blocked.
	}

	ls.Unlock(room) // last (and only) reader leaves; room should be unlocked

	select {
	case <-writerDone:
		// Good.
	case <-time.After(2 * time.Second):
		t.Fatal("writer never acquired the room after the last Lightswitch user left")
	}
}

// TestLightswitch_ManyConcurrentUsersStress spawns many goroutines entering and leaving
// concurrently and checks the run completes without deadlocking. Run with `go test
// -race` to also catch races in the internal counter bookkeeping.
func TestLightswitch_ManyConcurrentUsersStress(t *testing.T) {
	ls := NewLightswitch()
	room := &countingLocker{}

	const goroutines = 100
	const cyclesEach = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for c := 0; c < cyclesEach; c++ {
				ls.Lock(room)
				ls.Unlock(room)
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
	case <-time.After(15 * time.Second):
		t.Fatal("timed out — possible deadlock in Lightswitch")
	}

	// The room must never have been left in a locked state, and every lock must have
	// been matched by an unlock.
	if got := atomic.LoadInt32(&room.lockCalls); got != atomic.LoadInt32(&room.unlockCalls) {
		t.Fatalf("room.Lock() called %d times but room.Unlock() called %d times — mismatched",
			atomic.LoadInt32(&room.lockCalls), atomic.LoadInt32(&room.unlockCalls))
	}
}
