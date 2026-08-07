package designhashset

import "testing"

func TestMyHashSet_Example1(t *testing.T) {
	s := NewMyHashSet()
	s.Add(1)
	s.Add(2)
	if !s.Contains(1) {
		t.Fatal("Contains(1) = false, want true")
	}
	if s.Contains(3) {
		t.Fatal("Contains(3) = true, want false")
	}
	s.Add(2)
	if !s.Contains(2) {
		t.Fatal("Contains(2) = false, want true")
	}
	s.Remove(2)
	if s.Contains(2) {
		t.Fatal("Contains(2) after Remove = true, want false")
	}
}

func TestMyHashSet_ContainsOnEmptySet(t *testing.T) {
	s := NewMyHashSet()
	if s.Contains(0) {
		t.Fatal("Contains(0) on empty set = true, want false")
	}
}

func TestMyHashSet_RemoveMissingKeyIsNoOp(t *testing.T) {
	s := NewMyHashSet()
	s.Add(1)
	s.Remove(999)
	if !s.Contains(1) {
		t.Fatal("Contains(1) after removing unrelated key = false, want true")
	}
}

func TestMyHashSet_AddDuplicateIsIdempotent(t *testing.T) {
	s := NewMyHashSet()
	for i := 0; i < 5; i++ {
		s.Add(42)
	}
	if !s.Contains(42) {
		t.Fatal("Contains(42) = false, want true")
	}
	s.Remove(42)
	if s.Contains(42) {
		t.Fatal("Contains(42) after single Remove following repeated Add = true, want false")
	}
}

func TestMyHashSet_CollisionsWithinSameBucket(t *testing.T) {
	s := NewMyHashSet()
	s.Add(1)
	s.Add(1001)
	s.Add(2001)
	if !s.Contains(1) || !s.Contains(1001) || !s.Contains(2001) {
		t.Fatal("expected all of 1, 1001, 2001 to be present")
	}
	s.Remove(1001)
	if s.Contains(1001) {
		t.Fatal("Contains(1001) after Remove = true, want false")
	}
	if !s.Contains(1) || !s.Contains(2001) {
		t.Fatal("removing 1001 must not disturb 1 or 2001 in the same bucket chain")
	}
}

func TestMyHashSet_ReAddAfterRemove(t *testing.T) {
	s := NewMyHashSet()
	s.Add(7)
	s.Remove(7)
	if s.Contains(7) {
		t.Fatal("Contains(7) after Remove = true, want false")
	}
	s.Add(7)
	if !s.Contains(7) {
		t.Fatal("Contains(7) after re-Add = false, want true")
	}
}

func TestMyHashSet_BoundaryKeys(t *testing.T) {
	s := NewMyHashSet()
	s.Add(0)
	s.Add(1000000)
	if !s.Contains(0) {
		t.Fatal("Contains(0) = false, want true")
	}
	if !s.Contains(1000000) {
		t.Fatal("Contains(1000000) = false, want true")
	}
}

func TestMyHashSet_StressManyKeys(t *testing.T) {
	s := NewMyHashSet()
	const n = 5000

	for i := 0; i < n; i++ {
		s.Add(i)
	}
	for i := 0; i < n; i++ {
		if !s.Contains(i) {
			t.Fatalf("Contains(%d) = false, want true", i)
		}
	}
	for i := 0; i < n; i += 2 {
		s.Remove(i)
	}
	for i := 0; i < n; i++ {
		want := i%2 != 0
		if got := s.Contains(i); got != want {
			t.Fatalf("Contains(%d) after removing evens = %v, want %v", i, got, want)
		}
	}
}
