package buildingh2o

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// runH2O fires one goroutine per character of water ('H' -> Hydrogen, 'O' -> Oxygen),
// each releasing into a shared mutex-protected builder, and returns the bonding order
// the release callbacks actually fired in.
func runH2O(t *testing.T, water string) string {
	t.Helper()

	h := NewH2O()

	var mu sync.Mutex
	var sb strings.Builder

	var wg sync.WaitGroup
	wg.Add(len(water))

	for _, c := range water {
		c := c
		go func() {
			defer wg.Done()
			switch c {
			case 'H':
				h.Hydrogen(func() {
					mu.Lock()
					sb.WriteByte('H')
					mu.Unlock()
				})
			case 'O':
				h.Oxygen(func() {
					mu.Lock()
					sb.WriteByte('O')
					mu.Unlock()
				})
			default:
				t.Fatalf("unexpected character %q in water string %q", c, water)
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
	case <-time.After(5 * time.Second):
		t.Fatalf("H2O(%q): timed out waiting for all threads to finish (possible deadlock)", water)
	}

	mu.Lock()
	defer mu.Unlock()
	return sb.String()
}

// validateWater checks that the output is a permutation-preserving bonding of the input:
// every consecutive group of 3 released atoms must contain exactly 2 H's and 1 O, which
// is only possible if molecules bonded as complete, non-interleaved groups of three.
func validateWater(t *testing.T, input, output string) {
	t.Helper()

	if len(output) != len(input) {
		t.Fatalf("output length = %d, want %d (output=%q)", len(output), len(input), output)
	}

	var wantH, wantO int
	for _, c := range input {
		if c == 'H' {
			wantH++
		} else {
			wantO++
		}
	}
	var gotH, gotO int
	for _, c := range output {
		if c == 'H' {
			gotH++
		} else {
			gotO++
		}
	}
	if gotH != wantH || gotO != wantO {
		t.Fatalf("output %q has %d H/%d O, want %d H/%d O", output, gotH, gotO, wantH, wantO)
	}

	for i := 0; i+3 <= len(output); i += 3 {
		group := output[i : i+3]
		h, o := 0, 0
		for _, c := range group {
			if c == 'H' {
				h++
			} else {
				o++
			}
		}
		if h != 2 || o != 1 {
			t.Errorf("molecule %d (%q) is not a valid water molecule (want 2 H + 1 O)", i/3, group)
		}
	}
}

func allHydrogenFirst(n int) string { return strings.Repeat("H", 2*n) + strings.Repeat("O", n) }
func allOxygenFirst(n int) string   { return strings.Repeat("O", n) + strings.Repeat("H", 2*n) }
func interleavedHOH(n int) string   { return strings.Repeat("HOH", n) }

func TestH2O(t *testing.T) {
	tests := []struct {
		name  string
		water string
	}{
		// Problem statement examples.
		{"Example1", "HOH"},
		{"Example2", "OOHHHH"},

		// Small hand-picked orderings.
		{"n1_OHH", "OHH"},
		{"n2_HHOHOH", "HHOHOH"},
		{"n3_interleaved", interleavedHOH(3)},

		// All of one kind arriving before the other.
		{"n1_allOxygenFirst", allOxygenFirst(1)},
		{"n2_allOxygenFirst", allOxygenFirst(2)},
		{"n5_allOxygenFirst", allOxygenFirst(5)},
		{"n1_allHydrogenFirst", allHydrogenFirst(1)},
		{"n2_allHydrogenFirst", allHydrogenFirst(2)},
		{"n10_allHydrogenFirst", allHydrogenFirst(10)},

		// Larger inputs, including the max constraint (1 <= n <= 20).
		{"n10_interleaved", interleavedHOH(10)},
		{"n20_maxConstraint_interleaved", interleavedHOH(20)},
		{"n20_maxConstraint_allOxygenFirst", allOxygenFirst(20)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runH2O(t, tt.water)
			validateWater(t, tt.water, got)
		})
	}
}

// TestH2O_Repeated stress-tests a fixed input many times in a row to help surface
// flaky synchronization bugs. Run with `go test -race` to also catch data races.
func TestH2O_Repeated(t *testing.T) {
	water := interleavedHOH(5)
	for i := 0; i < 15; i++ {
		got := runH2O(t, water)
		validateWater(t, water, got)
	}
}
