package randompickwithweight

import "testing"

func TestSolution_SingleIndexAlwaysZero(t *testing.T) {
	s := NewSolution([]int{1})
	for i := 0; i < 50; i++ {
		if got := s.PickIndex(); got != 0 {
			t.Fatalf("PickIndex() = %d, want 0 (only one index exists)", got)
		}
	}
}

func TestSolution_AllPicksAreInBounds(t *testing.T) {
	s := NewSolution([]int{3, 1, 4, 1, 5, 9, 2, 6})
	for i := 0; i < 500; i++ {
		got := s.PickIndex()
		if got < 0 || got >= 8 {
			t.Fatalf("PickIndex() = %d, out of bounds [0,8)", got)
		}
	}
}

func TestSolution_ZeroWeightIndexNeverPicked(t *testing.T) {
	// Note: constraints say w[i] >= 1, but a robust implementation naturally handles a
	// weight of effectively "never" via a very small weight relative to a large one; here
	// we just check a very lopsided distribution behaves sanely, not literal zero-weight.
	s := NewSolution([]int{1, 1000000})
	countZero := 0
	const trials = 2000
	for i := 0; i < trials; i++ {
		if s.PickIndex() == 0 {
			countZero++
		}
	}
	// Index 0 has probability ~1/1000001; over 2000 trials it should be picked rarely.
	if countZero > trials/10 {
		t.Fatalf("index 0 picked %d/%d times, expected it to be rare given its tiny weight", countZero, trials)
	}
}

func TestSolution_TwoWeightsRoughDistribution(t *testing.T) {
	// w = [1, 3]: index 0 should be picked ~25% of the time, index 1 ~75%.
	s := NewSolution([]int{1, 3})
	const trials = 10000
	counts := [2]int{}
	for i := 0; i < trials; i++ {
		got := s.PickIndex()
		if got != 0 && got != 1 {
			t.Fatalf("PickIndex() = %d, want 0 or 1", got)
		}
		counts[got]++
	}
	ratio0 := float64(counts[0]) / float64(trials)
	// Loose bound: expect roughly 0.25, allow generous slack to avoid test flakiness.
	if ratio0 < 0.15 || ratio0 > 0.35 {
		t.Errorf("index 0 picked %.3f of the time, want roughly 0.25", ratio0)
	}
}

func TestSolution_EqualWeightsRoughlyUniform(t *testing.T) {
	s := NewSolution([]int{1, 1, 1, 1})
	const trials = 8000
	counts := [4]int{}
	for i := 0; i < trials; i++ {
		got := s.PickIndex()
		if got < 0 || got >= 4 {
			t.Fatalf("PickIndex() = %d, out of bounds", got)
		}
		counts[got]++
	}
	for i, c := range counts {
		ratio := float64(c) / float64(trials)
		if ratio < 0.15 || ratio > 0.35 {
			t.Errorf("index %d picked %.3f of the time, want roughly 0.25", i, ratio)
		}
	}
}

func TestSolution_LargeWeightsArray(t *testing.T) {
	weights := make([]int, 100)
	for i := range weights {
		weights[i] = i + 1 // weights 1..100
	}
	s := NewSolution(weights)
	seen := make(map[int]bool)
	for i := 0; i < 3000; i++ {
		got := s.PickIndex()
		if got < 0 || got >= 100 {
			t.Fatalf("PickIndex() = %d, out of bounds", got)
		}
		seen[got] = true
	}
	if len(seen) < 50 {
		t.Fatalf("only observed %d distinct indices over 3000 trials, distribution looks suspiciously narrow", len(seen))
	}
}

func TestSolution_TwoEqualWeightsBothAppear(t *testing.T) {
	s := NewSolution([]int{5, 5})
	seen := map[int]bool{}
	for i := 0; i < 200; i++ {
		seen[s.PickIndex()] = true
	}
	if len(seen) != 2 {
		t.Fatalf("expected to see both indices over 200 trials, saw %v", seen)
	}
}
