package shufflearray

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func sortedCopy(nums []int) []int {
	out := make([]int, len(nums))
	copy(out, nums)
	sort.Ints(out)
	return out
}

func TestSolution_ResetReturnsOriginal(t *testing.T) {
	original := []int{1, 2, 3}
	s := NewSolution(original)
	s.Shuffle()
	got := s.Reset()
	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Reset() = %v, want [1 2 3]", got)
	}
}

func TestSolution_ShuffleIsAPermutation(t *testing.T) {
	original := []int{1, 2, 3, 4, 5}
	s := NewSolution(original)
	for i := 0; i < 20; i++ {
		got := s.Shuffle()
		if len(got) != len(original) {
			t.Fatalf("Shuffle() returned length %d, want %d", len(got), len(original))
		}
		if !reflect.DeepEqual(sortedCopy(got), []int{1, 2, 3, 4, 5}) {
			t.Fatalf("Shuffle() = %v, is not a permutation of %v", got, original)
		}
	}
}

func TestSolution_MultipleResetsAlwaysMatchOriginal(t *testing.T) {
	original := []int{-5, 0, 10, 3}
	s := NewSolution(original)
	for i := 0; i < 5; i++ {
		s.Shuffle()
		got := s.Reset()
		if !reflect.DeepEqual(got, []int{-5, 0, 10, 3}) {
			t.Fatalf("Reset() call %d = %v, want [-5 0 10 3]", i, got)
		}
	}
}

func TestSolution_ResetDoesNotMutateOriginalReferenceAcrossCalls(t *testing.T) {
	original := []int{7, 8, 9}
	s := NewSolution(original)
	first := s.Reset()
	if len(first) != 3 {
		t.Fatalf("Reset() = %v, want length 3", first)
	}
	first[0] = 999 // mutate the returned slice
	second := s.Reset()
	if len(second) != 3 || second[0] != 7 {
		t.Fatalf("Reset() after mutating a prior Reset() result = %v, want [7 8 9] (internal original must not alias caller's slice)", second)
	}
}

func TestSolution_SingleElementAlwaysSame(t *testing.T) {
	s := NewSolution([]int{42})
	for i := 0; i < 5; i++ {
		got := s.Shuffle()
		if !reflect.DeepEqual(got, []int{42}) {
			t.Fatalf("Shuffle() on single-element array = %v, want [42]", got)
		}
	}
	got := s.Reset()
	if !reflect.DeepEqual(got, []int{42}) {
		t.Fatalf("Reset() = %v, want [42]", got)
	}
}

func TestSolution_NegativeAndLargeValues(t *testing.T) {
	original := []int{-1000000, 0, 1000000}
	s := NewSolution(original)
	got := s.Shuffle()
	if !reflect.DeepEqual(sortedCopy(got), []int{-1000000, 0, 1000000}) {
		t.Fatalf("Shuffle() = %v, is not a permutation of %v", got, original)
	}
}

func TestSolution_ShuffleProducesVariedOutcomes(t *testing.T) {
	// Statistical smoke test: over enough shuffles of a reasonably sized array, we
	// should see more than one distinct permutation. (Not a rigorous uniformity test —
	// just guards against a broken "shuffle" that always returns the same order.)
	original := []int{1, 2, 3, 4, 5, 6}
	s := NewSolution(original)
	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		got := s.Shuffle()
		seen[keyOf(got)] = true
	}
	if len(seen) < 2 {
		t.Fatalf("Shuffle() produced only %d distinct permutation(s) over 50 calls, want more variety", len(seen))
	}
}

func keyOf(nums []int) string {
	return fmt.Sprint(nums)
}

func TestSolution_LargerArrayEachValueAppearsExactlyOnce(t *testing.T) {
	original := make([]int, 50)
	for i := range original {
		original[i] = i
	}
	s := NewSolution(original)
	got := s.Shuffle()
	seen := make(map[int]bool)
	for _, v := range got {
		if seen[v] {
			t.Fatalf("value %d appeared more than once in shuffled output %v", v, got)
		}
		seen[v] = true
	}
	if len(seen) != 50 {
		t.Fatalf("shuffled output has %d distinct values, want 50", len(seen))
	}
}
