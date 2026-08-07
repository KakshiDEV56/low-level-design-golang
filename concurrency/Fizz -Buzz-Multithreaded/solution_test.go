package fizzBuzzMultithreaded

import (
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

// expectedSequence builds the reference [1,2,"fizz",4,"buzz",...] sequence for a given n,
// with numbers rendered as strings so it can be compared against the collected output.
func expectedSequence(n int) []string {
	seq := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			seq = append(seq, "fizzbuzz")
		case i%3 == 0:
			seq = append(seq, "fizz")
		case i%5 == 0:
			seq = append(seq, "buzz")
		default:
			seq = append(seq, strconv.Itoa(i))
		}
	}
	return seq
}

// runFizzBuzz wires up a FizzBuzz instance the same way the problem statement describes:
// four goroutines (Fizz, Buzz, Fizzbuzz, Number), each calling its print callback which
// appends to a shared, mutex-protected slice. If the four goroutines are correctly
// synchronized, the append order matches 1..n regardless of OS scheduling.
func runFizzBuzz(t *testing.T, n int) []string {
	t.Helper()

	fb := NewFizzBuzz(n)

	var mu sync.Mutex
	results := make([]string, 0, n)

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		fb.Fizz(func() {
			mu.Lock()
			results = append(results, "fizz")
			mu.Unlock()
		})
	}()
	go func() {
		defer wg.Done()
		fb.Buzz(func() {
			mu.Lock()
			results = append(results, "buzz")
			mu.Unlock()
		})
	}()
	go func() {
		defer wg.Done()
		fb.Fizzbuzz(func() {
			mu.Lock()
			results = append(results, "fizzbuzz")
			mu.Unlock()
		})
	}()
	go func() {
		defer wg.Done()
		fb.Number(func(x int) {
			mu.Lock()
			results = append(results, strconv.Itoa(x))
			mu.Unlock()
		})
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("FizzBuzz(%d): timed out waiting for all 4 goroutines to finish (possible deadlock)", n)
	}

	mu.Lock()
	defer mu.Unlock()
	out := make([]string, len(results))
	copy(out, results)
	return out
}

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		// Problem statement examples.
		{"Example1_n15", 15},
		{"Example2_n5", 5},

		// Edge cases.

		{"n1_onlyNumber", 1},
		{"n2_onlyNumbers", 2},

		// Boundary of first fizz/buzz.
		{"n3_firstFizz", 3},
		{"n4", 4},
		{"n5_firstBuzz", 5},
		{"n6_secondFizz", 6},

		// First fizzbuzz boundary.
		{"n14_justBeforeFizzbuzz", 14},
		{"n15_firstFizzbuzz", 15},
		{"n16_justAfterFizzbuzz", 16},

		// Larger inputs with multiple repeats of fizz/buzz/fizzbuzz.
		{"n20", 20},
		{"n30_twoFizzbuzz", 30},
		{"n45_threeFizzbuzz", 45},
		{"n50", 50},
		{"n100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runFizzBuzz(t, tt.n)
			want := expectedSequence(tt.n)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("FizzBuzz(%d) = %v, want %v", tt.n, got, want)
			}
		})
	}
}

// TestFizzBuzz_Repeated runs the same n multiple times back-to-back to help surface
// flaky synchronization bugs that don't show up on every run. Run with `go test -race`
// to also catch data races on the shared results slice/counter.
func TestFizzBuzz_Repeated(t *testing.T) {
	const n = 15
	want := expectedSequence(n)

	for i := 0; i < 20; i++ {
		got := runFizzBuzz(t, n)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("iteration %d: FizzBuzz(%d) = %v, want %v", i, n, got, want)
		}
	}
}

// TestFizzBuzz_OutputLength checks that every number 1..n produces exactly one token,
// which would catch bugs where a goroutine double-prints or a number gets dropped.
func TestFizzBuzz_OutputLength(t *testing.T) {
	for _, n := range []int{0, 1, 5, 15, 33} {
		got := runFizzBuzz(t, n)
		if len(got) != n {
			t.Errorf("FizzBuzz(%d) produced %d tokens, want %d", n, len(got), n)
		}
	}
}
