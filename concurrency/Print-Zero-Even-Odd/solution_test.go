package printzeroevenodd

import (
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

// expectedZeroEvenOdd builds the reference "0" <number> "0" <number> ... sequence
// (e.g. n=2 -> ["0","1","0","2"], which joins to "0102") for a given n.
func expectedZeroEvenOdd(n int) []string {
	seq := make([]string, 0, 2*n)
	for i := 1; i <= n; i++ {
		seq = append(seq, "0", strconv.Itoa(i))
	}
	return seq
}

// runZeroEvenOdd wires up Zero/Even/Odd on three goroutines, each calling printNumber,
// which appends to a shared mutex-protected slice. If synchronized correctly, the append
// order matches "0","1","0","2","0","3",... regardless of scheduling.
func runZeroEvenOdd(t *testing.T, n int) []string {
	t.Helper()

	z := NewZeroEvenOdd(n)

	var mu sync.Mutex
	results := make([]string, 0, 2*n)
	record := func(x int) {
		mu.Lock()
		results = append(results, strconv.Itoa(x))
		mu.Unlock()
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() { defer wg.Done(); z.Zero(record) }()
	go func() { defer wg.Done(); z.Even(record) }()
	go func() { defer wg.Done(); z.Odd(record) }()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("ZeroEvenOdd(%d): timed out waiting for Zero/Even/Odd to finish (possible deadlock)", n)
	}

	mu.Lock()
	defer mu.Unlock()
	out := make([]string, len(results))
	copy(out, results)
	return out
}

func TestZeroEvenOdd(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		// Problem statement examples.
		{"Example1_n2", 2},
		{"Example2_n5", 5},

		// Small boundary values.
		{"n1", 1},
		{"n3", 3},
		{"n4", 4},
		{"n6", 6},
		{"n7", 7},

		// Larger values.
		{"n10", 10},
		{"n15", 15},
		{"n20", 20},
		{"n50", 50},
		{"n100", 100},
		{"n500", 500},

		// Max constraint (1 <= n <= 1000).
		{"n1000_maxConstraint", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runZeroEvenOdd(t, tt.n)
			want := expectedZeroEvenOdd(tt.n)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ZeroEvenOdd(%d) = %v, want %v", tt.n, got, want)
			}
		})
	}
}

// TestZeroEvenOdd_Repeated runs the same n multiple times back-to-back to help surface
// flaky synchronization bugs. Run with `go test -race` to also catch data races.
func TestZeroEvenOdd_Repeated(t *testing.T) {
	const n = 15
	want := expectedZeroEvenOdd(n)

	for i := 0; i < 20; i++ {
		got := runZeroEvenOdd(t, n)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("iteration %d: ZeroEvenOdd(%d) = %v, want %v", i, n, got, want)
		}
	}
}
