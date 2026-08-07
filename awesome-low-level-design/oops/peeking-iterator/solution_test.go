package peekingiterator

import "testing"

func TestPeekingIterator_Example1(t *testing.T) {
	p := NewPeekingIterator(NewIterator([]int{1, 2, 3}))
	if got := p.Next(); got != 1 {
		t.Fatalf("Next() = %d, want 1", got)
	}
	if got := p.Peek(); got != 2 {
		t.Fatalf("Peek() = %d, want 2", got)
	}
	if got := p.Next(); got != 2 {
		t.Fatalf("Next() = %d, want 2", got)
	}
	if got := p.Next(); got != 3 {
		t.Fatalf("Next() = %d, want 3", got)
	}
	if got := p.HasNext(); got {
		t.Fatalf("HasNext() = %v, want false", got)
	}
}

func TestPeekingIterator_PeekWithoutNextDoesNotAdvance(t *testing.T) {
	p := NewPeekingIterator(NewIterator([]int{5, 6, 7}))
	if got := p.Peek(); got != 5 {
		t.Fatalf("Peek() = %d, want 5", got)
	}
	if got := p.Peek(); got != 5 {
		t.Fatalf("Peek() second call = %d, want 5 (should not have advanced)", got)
	}
	if got := p.Peek(); got != 5 {
		t.Fatalf("Peek() third call = %d, want 5", got)
	}
	if got := p.Next(); got != 5 {
		t.Fatalf("Next() = %d, want 5", got)
	}
	if got := p.Next(); got != 6 {
		t.Fatalf("Next() = %d, want 6", got)
	}
}

func TestPeekingIterator_HasNextAfterFullyPeekedAndConsumed(t *testing.T) {
	p := NewPeekingIterator(NewIterator([]int{1}))
	if !p.HasNext() {
		t.Fatal("HasNext() = false, want true before consuming")
	}
	if got := p.Peek(); got != 1 {
		t.Fatalf("Peek() = %d, want 1", got)
	}
	if !p.HasNext() {
		t.Fatal("HasNext() = false, want true (peeked but not consumed)")
	}
	if got := p.Next(); got != 1 {
		t.Fatalf("Next() = %d, want 1", got)
	}
	if p.HasNext() {
		t.Fatal("HasNext() = true, want false after consuming the only element")
	}
}

func TestPeekingIterator_AllNextNoPeek(t *testing.T) {
	p := NewPeekingIterator(NewIterator([]int{10, 20, 30}))
	want := []int{10, 20, 30}
	for i, w := range want {
		if got := p.Next(); got != w {
			t.Fatalf("Next() call %d = %d, want %d", i, got, w)
		}
	}
	if p.HasNext() {
		t.Fatal("HasNext() = true, want false")
	}
}

func TestPeekingIterator_PeekBeforeEveryNext(t *testing.T) {
	p := NewPeekingIterator(NewIterator([]int{100, 200, 300}))
	for _, want := range []int{100, 200, 300} {
		if got := p.Peek(); got != want {
			t.Fatalf("Peek() = %d, want %d", got, want)
		}
		if got := p.Next(); got != want {
			t.Fatalf("Next() = %d, want %d", got, want)
		}
	}
	if p.HasNext() {
		t.Fatal("HasNext() = true, want false")
	}
}

func TestPeekingIterator_SingleElement(t *testing.T) {
	p := NewPeekingIterator(NewIterator([]int{42}))
	if !p.HasNext() {
		t.Fatal("HasNext() = false, want true")
	}
	if got := p.Peek(); got != 42 {
		t.Fatalf("Peek() = %d, want 42", got)
	}
	if got := p.Next(); got != 42 {
		t.Fatalf("Next() = %d, want 42", got)
	}
	if p.HasNext() {
		t.Fatal("HasNext() = true, want false")
	}
}

func TestPeekingIterator_LongSequenceMixedCalls(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	p := NewPeekingIterator(NewIterator(data))

	got := []int{}
	for p.HasNext() {
		got = append(got, p.Peek())
		got = append(got, p.Next())
	}
	if len(got) != 2*len(data) {
		t.Fatalf("collected %d values, want %d", len(got), 2*len(data))
	}
	for i, v := range data {
		if got[2*i] != v || got[2*i+1] != v {
			t.Fatalf("at position %d: peek/next = %d/%d, want both %d", i, got[2*i], got[2*i+1], v)
		}
	}
}
