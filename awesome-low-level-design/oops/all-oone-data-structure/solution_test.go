package alloonedatastructure

import "testing"

func isOneOf(got string, candidates ...string) bool {
	for _, c := range candidates {
		if got == c {
			return true
		}
	}
	return false
}

func TestAllOne_Example1(t *testing.T) {
	a := NewAllOne()
	a.Inc("hello")
	a.Inc("hello")
	if got := a.GetMaxKey(); got != "hello" {
		t.Fatalf("GetMaxKey() = %q, want hello", got)
	}
	if got := a.GetMinKey(); got != "hello" {
		t.Fatalf("GetMinKey() = %q, want hello", got)
	}
	a.Inc("leet")
	if got := a.GetMaxKey(); got != "hello" {
		t.Fatalf("GetMaxKey() = %q, want hello", got)
	}
	if got := a.GetMinKey(); got != "leet" {
		t.Fatalf("GetMinKey() = %q, want leet", got)
	}
}

func TestAllOne_EmptyReturnsEmptyString(t *testing.T) {
	a := NewAllOne()
	if got := a.GetMaxKey(); got != "" {
		t.Fatalf("GetMaxKey() on empty structure = %q, want \"\"", got)
	}
	if got := a.GetMinKey(); got != "" {
		t.Fatalf("GetMinKey() on empty structure = %q, want \"\"", got)
	}
}

func TestAllOne_DecToZeroRemovesKey(t *testing.T) {
	a := NewAllOne()
	a.Inc("a")
	a.Inc("b")
	a.Dec("a")
	if got := a.GetMinKey(); got != "b" {
		t.Fatalf("GetMinKey() = %q, want b", got)
	}
	if got := a.GetMaxKey(); got != "b" {
		t.Fatalf("GetMaxKey() = %q, want b", got)
	}
}

func TestAllOne_MultipleKeysSameCountTieAllowed(t *testing.T) {
	a := NewAllOne()
	a.Inc("x")
	a.Inc("y")
	a.Inc("z")
	// All three have count 1; both max and min can be any of them.
	if got := a.GetMaxKey(); !isOneOf(got, "x", "y", "z") {
		t.Fatalf("GetMaxKey() = %q, want one of x/y/z", got)
	}
	if got := a.GetMinKey(); !isOneOf(got, "x", "y", "z") {
		t.Fatalf("GetMinKey() = %q, want one of x/y/z", got)
	}
}

func TestAllOne_IncDecInterleaved(t *testing.T) {
	a := NewAllOne()
	a.Inc("a") // a:1
	a.Inc("a") // a:2
	a.Inc("b") // b:1
	a.Inc("b") // b:2
	a.Inc("b") // b:3
	if got := a.GetMaxKey(); got != "b" {
		t.Fatalf("GetMaxKey() = %q, want b", got)
	}
	if got := a.GetMinKey(); got != "a" {
		t.Fatalf("GetMinKey() = %q, want a", got)
	}
	a.Dec("b") // b:2
	a.Dec("b") // b:1
	// Now a:2, b:1
	if got := a.GetMaxKey(); got != "a" {
		t.Fatalf("GetMaxKey() = %q, want a", got)
	}
	if got := a.GetMinKey(); got != "b" {
		t.Fatalf("GetMinKey() = %q, want b", got)
	}
}

func TestAllOne_ReinsertAfterFullRemoval(t *testing.T) {
	a := NewAllOne()
	a.Inc("k")
	a.Dec("k")
	if got := a.GetMaxKey(); got != "" {
		t.Fatalf("GetMaxKey() after full removal = %q, want \"\"", got)
	}
	a.Inc("k")
	if got := a.GetMaxKey(); got != "k" {
		t.Fatalf("GetMaxKey() after reinsert = %q, want k", got)
	}
	if got := a.GetMinKey(); got != "k" {
		t.Fatalf("GetMinKey() after reinsert = %q, want k", got)
	}
}

func TestAllOne_MaxMinDistinctWithMiddleKeys(t *testing.T) {
	a := NewAllOne()
	a.Inc("low")
	for i := 0; i < 5; i++ {
		a.Inc("high")
	}
	for i := 0; i < 3; i++ {
		a.Inc("mid")
	}
	if got := a.GetMaxKey(); got != "high" {
		t.Fatalf("GetMaxKey() = %q, want high", got)
	}
	if got := a.GetMinKey(); got != "low" {
		t.Fatalf("GetMinKey() = %q, want low", got)
	}
}

func TestAllOne_DecBucketMergingDoesNotLoseOtherKeys(t *testing.T) {
	a := NewAllOne()
	a.Inc("p") // p:1
	a.Inc("q") // q:1
	a.Inc("q") // q:2
	a.Dec("q") // q:1 -- must merge back into the count=1 bucket alongside p, not clobber it
	if got := a.GetMaxKey(); !isOneOf(got, "p", "q") {
		t.Fatalf("GetMaxKey() = %q, want p or q", got)
	}
	if got := a.GetMinKey(); !isOneOf(got, "p", "q") {
		t.Fatalf("GetMinKey() = %q, want p or q", got)
	}
}

func TestAllOne_StressManyKeys(t *testing.T) {
	a := NewAllOne()
	keys := []string{"aa", "bb", "cc", "dd", "ee"}
	for i, k := range keys {
		for c := 0; c <= i; c++ {
			a.Inc(k)
		}
	}
	// counts: aa=1, bb=2, cc=3, dd=4, ee=5
	if got := a.GetMaxKey(); got != "ee" {
		t.Fatalf("GetMaxKey() = %q, want ee", got)
	}
	if got := a.GetMinKey(); got != "aa" {
		t.Fatalf("GetMinKey() = %q, want aa", got)
	}
	// Bring ee all the way down to 0 (removed) and bump aa up past everything.
	for c := 0; c < 5; c++ {
		a.Dec("ee")
	}
	for c := 0; c < 10; c++ {
		a.Inc("aa")
	}
	if got := a.GetMaxKey(); got != "aa" {
		t.Fatalf("GetMaxKey() = %q, want aa", got)
	}
	if got := a.GetMinKey(); got != "bb" {
		t.Fatalf("GetMinKey() = %q, want bb", got)
	}
}
