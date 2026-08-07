package insertdeletegetrandom

import "testing"

func TestRandomizedSet_Example1(t *testing.T) {
	s := NewRandomizedSet()
	if got := s.Insert(1); !got {
		t.Fatal("Insert(1) = false, want true")
	}
	if got := s.Remove(2); got {
		t.Fatal("Remove(2) = true, want false (not present)")
	}
	if got := s.Insert(2); !got {
		t.Fatal("Insert(2) = false, want true")
	}
	if got := s.GetRandom(); got != 1 && got != 2 {
		t.Fatalf("GetRandom() = %d, want 1 or 2", got)
	}
	if got := s.Remove(1); !got {
		t.Fatal("Remove(1) = false, want true")
	}
	if got := s.Insert(2); got {
		t.Fatal("Insert(2) = true, want false (already present)")
	}
	if got := s.GetRandom(); got != 2 {
		t.Fatalf("GetRandom() = %d, want 2", got)
	}
}

func TestRandomizedSet_InsertDuplicateReturnsFalse(t *testing.T) {
	s := NewRandomizedSet()
	s.Insert(5)
	if got := s.Insert(5); got {
		t.Fatal("Insert(5) second time = true, want false")
	}
}

func TestRandomizedSet_RemoveMissingReturnsFalse(t *testing.T) {
	s := NewRandomizedSet()
	if got := s.Remove(1); got {
		t.Fatal("Remove(1) on empty set = true, want false")
	}
}

func TestRandomizedSet_GetRandomAfterSwapRemove(t *testing.T) {
	s := NewRandomizedSet()
	s.Insert(1)
	s.Insert(2)
	s.Insert(3)
	s.Remove(1) // if implemented via swap-with-last, this must not corrupt 2 or 3
	seen := map[int]bool{}
	for i := 0; i < 100; i++ {
		v := s.GetRandom()
		if v != 2 && v != 3 {
			t.Fatalf("GetRandom() = %d, want 2 or 3", v)
		}
		seen[v] = true
	}
	if len(seen) != 2 {
		t.Fatalf("GetRandom() only ever returned %v over 100 calls, want to see both 2 and 3", seen)
	}
}

func TestRandomizedSet_RemoveLastInsertedElement(t *testing.T) {
	s := NewRandomizedSet()
	s.Insert(1)
	s.Insert(2)
	s.Insert(3)
	s.Remove(3) // removing the most recently inserted (last in slice) element
	for i := 0; i < 20; i++ {
		v := s.GetRandom()
		if v != 1 && v != 2 {
			t.Fatalf("GetRandom() = %d, want 1 or 2", v)
		}
	}
}

func TestRandomizedSet_ReinsertAfterRemove(t *testing.T) {
	s := NewRandomizedSet()
	s.Insert(1)
	s.Remove(1)
	if got := s.Insert(1); !got {
		t.Fatal("Insert(1) after Remove = false, want true (should be re-insertable)")
	}
	if got := s.GetRandom(); got != 1 {
		t.Fatalf("GetRandom() = %d, want 1", got)
	}
}

func TestRandomizedSet_SingleElementAlwaysReturnsIt(t *testing.T) {
	s := NewRandomizedSet()
	s.Insert(42)
	for i := 0; i < 10; i++ {
		if got := s.GetRandom(); got != 42 {
			t.Fatalf("GetRandom() = %d, want 42", got)
		}
	}
}

func TestRandomizedSet_NegativeAndBoundaryValues(t *testing.T) {
	s := NewRandomizedSet()
	s.Insert(-2147483648)
	s.Insert(2147483647)
	seen := map[int]bool{}
	for i := 0; i < 50; i++ {
		seen[s.GetRandom()] = true
	}
	if len(seen) != 2 {
		t.Fatalf("expected to observe both boundary values, saw %v", seen)
	}
}

func TestRandomizedSet_UniformDistributionRoughCheck(t *testing.T) {
	s := NewRandomizedSet()
	for i := 0; i < 5; i++ {
		s.Insert(i)
	}
	counts := make(map[int]int)
	const trials = 5000
	for i := 0; i < trials; i++ {
		counts[s.GetRandom()]++
	}
	if len(counts) != 5 {
		t.Fatalf("expected all 5 values to appear over %d trials, got %d distinct values: %v", trials, len(counts), counts)
	}
	// Loose sanity bound: each value should appear at least a few times (not a strict
	// uniformity test, just a smoke test against wildly broken randomization).
	for v, c := range counts {
		if c < trials/20 {
			t.Errorf("value %d only appeared %d/%d times, distribution looks broken", v, c, trials)
		}
	}
}

func TestRandomizedSet_StressManyInsertsAndRemoves(t *testing.T) {
	s := NewRandomizedSet()
	const n = 1000

	for i := 0; i < n; i++ {
		if !s.Insert(i) {
			t.Fatalf("Insert(%d) = false, want true", i)
		}
	}
	for i := 0; i < n; i += 2 {
		if !s.Remove(i) {
			t.Fatalf("Remove(%d) = false, want true", i)
		}
	}
	for i := 0; i < 200; i++ {
		v := s.GetRandom()
		if v%2 == 0 {
			t.Fatalf("GetRandom() = %d, want an odd (still-present) value", v)
		}
	}
}
