package reusablebarrier

import (
	"sync"
	"testing"
	"time"
)

// TestCyclicBarrier_SingleParticipant checks the trivial n=1 case: Await() should
// always return immediately since there's nobody else to wait for.
func TestCyclicBarrier_SingleParticipant(t *testing.T) {
	b := NewCyclicBarrier(1)
	done := make(chan struct{})
	go func() {
		for i := 0; i < 5; i++ {
			b.Await()
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out — a barrier of size 1 should never block")
	}
}

// TestCyclicBarrier_BlocksUntilAllArrive checks that with n=3, the first two arrivals
// block, and only once the third arrives do all three proceed.
func TestCyclicBarrier_BlocksUntilAllArrive(t *testing.T) {
	b := NewCyclicBarrier(3)
	returned := make(chan int, 3)

	go func() { b.Await(); returned <- 1 }()
	go func() { b.Await(); returned <- 2 }()

	select {
	case <-returned:
		t.Fatal("Await() returned before all 3 participants arrived")
	case <-time.After(200 * time.Millisecond):
		// Good: still blocked with only 2/3 arrived.
	}

	go func() { b.Await(); returned <- 3 }()

	for i := 0; i < 3; i++ {
		select {
		case <-returned:
		case <-time.After(2 * time.Second):
			t.Fatalf("only %d/3 participants returned from Await() after the 3rd arrived", i)
		}
	}
}

// event is one recorded arrival at, or departure from, a specific round of the barrier.
type event struct {
	goroutine int
	round     int
	phase     string // "arrive" or "depart"
}

// TestCyclicBarrier_NoLappingAcrossRounds is the core correctness test: it runs many
// goroutines through many rounds and verifies that ALL n departures from round r are
// logged before ANY departure from round r+1 — i.e. no goroutine can "lap" another by
// racing ahead into the next round before everyone has left the current one. This is
// exactly the bug the book calls out in "Reusable barrier non-solution #2".
func TestCyclicBarrier_NoLappingAcrossRounds(t *testing.T) {
	const goroutines = 10
	const rounds = 50
	b := NewCyclicBarrier(goroutines)

	var mu sync.Mutex
	var log []event

	record := func(e event) {
		mu.Lock()
		log = append(log, e)
		mu.Unlock()
	}

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		g := g
		go func() {
			defer wg.Done()
			for r := 0; r < rounds; r++ {
				record(event{g, r, "arrive"})
				b.Await()
				record(event{g, r, "depart"})
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
		t.Fatal("timed out — possible deadlock in the barrier")
	}

	// Post-hoc check: for each round r, the LAST "depart" event for round r must come
	// before the FIRST "depart" event for round r+1.
	lastDepartIndex := make(map[int]int)  // round -> log index of its last depart event
	firstDepartIndex := make(map[int]int) // round -> log index of its first depart event
	departCount := make(map[int]int)      // round -> how many departs seen so far

	for i, e := range log {
		if e.phase != "depart" {
			continue
		}
		if _, ok := firstDepartIndex[e.round]; !ok {
			firstDepartIndex[e.round] = i
		}
		lastDepartIndex[e.round] = i
		departCount[e.round]++
	}

	for r := 0; r < rounds; r++ {
		if departCount[r] != goroutines {
			t.Fatalf("round %d had %d depart events, want %d", r, departCount[r], goroutines)
		}
	}

	for r := 0; r < rounds-1; r++ {
		if firstDepartIndex[r+1] < lastDepartIndex[r] {
			t.Fatalf(
				"lapping detected: a goroutine departed round %d (log index %d) before all goroutines "+
					"had departed round %d (last depart at log index %d)",
				r+1, firstDepartIndex[r+1], r, lastDepartIndex[r])
		}
	}
}

// TestCyclicBarrier_ReusableImmediately checks that the barrier can be awaited again
// immediately in a tight loop without any extra reset step.
func TestCyclicBarrier_ReusableImmediately(t *testing.T) {
	const goroutines = 4
	const rounds = 100
	b := NewCyclicBarrier(goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for r := 0; r < rounds; r++ {
				b.Await()
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
		t.Fatal("timed out running many tight-loop rounds")
	}
}

// TestCyclicBarrier_StressLargeGroup is a broader smoke test with a bigger group, to
// help surface races. Run with `go test -race`.
func TestCyclicBarrier_StressLargeGroup(t *testing.T) {
	const goroutines = 50
	const rounds = 20
	b := NewCyclicBarrier(goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for r := 0; r < rounds; r++ {
				b.Await()
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
		t.Fatalf("timed out with %d goroutines over %d rounds", goroutines, rounds)
	}
}
