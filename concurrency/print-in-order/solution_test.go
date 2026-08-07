package printinorder

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// runFoo launches First/Second/Third in the given goroutine-launch order (which must NOT
// affect the outcome — that's the whole point of the problem) and returns the concatenated
// output in the order the print callbacks actually fired.
func runFoo(t *testing.T, launchOrder [3]int) string {
	t.Helper()

	f := NewFoo()

	var mu sync.Mutex
	var sb strings.Builder

	calls := [3]func(){
		func() { f.First(func() { mu.Lock(); sb.WriteString("first"); mu.Unlock() }) },
		func() { f.Second(func() { mu.Lock(); sb.WriteString("second"); mu.Unlock() }) },
		func() { f.Third(func() { mu.Lock(); sb.WriteString("third"); mu.Unlock() }) },
	}

	var wg sync.WaitGroup
	wg.Add(3)
	for _, idx := range launchOrder {
		idx := idx
		go func() {
			defer wg.Done()
			calls[idx]()
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("Foo: timed out for launch order %v (possible deadlock)", launchOrder)
	}

	mu.Lock()
	defer mu.Unlock()
	return sb.String()
}

// TestFoo_AllLaunchOrders tries every possible order in which the First/Second/Third
// goroutines could get scheduled first, and repeats each order a few times, since the
// bug this problem guards against (missing synchronization) is scheduler-dependent and
// won't necessarily show up on the first try.
func TestFoo_AllLaunchOrders(t *testing.T) {
	const want = "firstsecondthird"

	// 0 = First, 1 = Second, 2 = Third. Each triple is the order the three goroutines
	// are launched in — the synchronization inside Foo must make the output identical
	// regardless of this order.
	orders := [][3]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	}

	for _, order := range orders {
		for run := 1; run <= 3; run++ {
			name := fmt.Sprintf("order_%v_run%d", order, run)
			t.Run(name, func(t *testing.T) {
				got := runFoo(t, order)
				if got != want {
					t.Errorf("Foo() launched in order %v = %q, want %q", order, got, want)
				}
			})
		}
	}
}

// TestFoo_Repeated stress-tests a fixed launch order many times in a row to help surface
// flaky synchronization bugs. Run with `go test -race` to also catch data races.
func TestFoo_Repeated(t *testing.T) {
	const want = "firstsecondthird"
	order := [3]int{2, 1, 0}

	for i := 0; i < 15; i++ {
		got := runFoo(t, order)
		if got != want {
			t.Fatalf("iteration %d: Foo() = %q, want %q", i, got, want)
		}
	}
}
