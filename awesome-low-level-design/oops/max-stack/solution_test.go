package maxstack

import "testing"

func TestMaxStack_Example1(t *testing.T) {
	s := NewMaxStack()
	s.Push(5)
	s.Push(1)
	s.Push(5)
	if got := s.Top(); got != 5 {
		t.Fatalf("Top() = %d, want 5", got)
	}
	if got := s.PopMax(); got != 5 {
		t.Fatalf("PopMax() = %d, want 5", got)
	}
	if got := s.Top(); got != 1 {
		t.Fatalf("Top() = %d, want 1", got)
	}
	if got := s.PeekMax(); got != 5 {
		t.Fatalf("PeekMax() = %d, want 5", got)
	}
	if got := s.Pop(); got != 1 {
		t.Fatalf("Pop() = %d, want 1", got)
	}
	if got := s.Top(); got != 5 {
		t.Fatalf("Top() = %d, want 5", got)
	}
}

func TestMaxStack_SingleElement(t *testing.T) {
	s := NewMaxStack()
	s.Push(42)
	if got := s.Top(); got != 42 {
		t.Fatalf("Top() = %d, want 42", got)
	}
	if got := s.PeekMax(); got != 42 {
		t.Fatalf("PeekMax() = %d, want 42", got)
	}
	if got := s.Pop(); got != 42 {
		t.Fatalf("Pop() = %d, want 42", got)
	}
}

func TestMaxStack_PopMaxRemovesOnlyTopmostDuplicate(t *testing.T) {
	s := NewMaxStack()
	s.Push(1)
	s.Push(5)
	s.Push(5)
	s.Push(3)
	if got := s.PopMax(); got != 5 {
		t.Fatalf("PopMax() = %d, want 5", got)
	}
	// Remaining stack bottom->top: [1, 5, 3]; the other 5 must still be there.
	if got := s.PeekMax(); got != 5 {
		t.Fatalf("PeekMax() after first PopMax = %d, want 5", got)
	}
	if got := s.Pop(); got != 3 {
		t.Fatalf("Pop() = %d, want 3", got)
	}
	if got := s.Top(); got != 5 {
		t.Fatalf("Top() = %d, want 5", got)
	}
}

func TestMaxStack_MaxUpdatesAfterPop(t *testing.T) {
	s := NewMaxStack()
	s.Push(3)
	s.Push(9)
	s.Push(7)
	if got := s.PeekMax(); got != 9 {
		t.Fatalf("PeekMax() = %d, want 9", got)
	}
	s.Pop() // removes 7
	if got := s.PeekMax(); got != 9 {
		t.Fatalf("PeekMax() = %d, want 9", got)
	}
	s.Pop() // removes 9
	if got := s.PeekMax(); got != 3 {
		t.Fatalf("PeekMax() = %d, want 3", got)
	}
}

func TestMaxStack_MaxUpdatesAfterPopMax(t *testing.T) {
	s := NewMaxStack()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(2)
	if got := s.PopMax(); got != 3 {
		t.Fatalf("PopMax() = %d, want 3", got)
	}
	if got := s.PeekMax(); got != 2 {
		t.Fatalf("PeekMax() = %d, want 2", got)
	}
}

func TestMaxStack_PopMaxFromMiddleThenPushKeepsOrder(t *testing.T) {
	s := NewMaxStack()
	s.Push(1)
	s.Push(10)
	s.Push(2)
	s.PopMax() // removes 10; stack bottom->top now [1, 2]
	s.Push(4)  // stack bottom->top now [1, 2, 4]
	if got := s.Top(); got != 4 {
		t.Fatalf("Top() = %d, want 4", got)
	}
	if got := s.Pop(); got != 4 {
		t.Fatalf("Pop() = %d, want 4", got)
	}
	if got := s.Pop(); got != 2 {
		t.Fatalf("Pop() = %d, want 2 (order below the removed max must be preserved)", got)
	}
	if got := s.Pop(); got != 1 {
		t.Fatalf("Pop() = %d, want 1", got)
	}
}

func TestMaxStack_AllSameValue(t *testing.T) {
	s := NewMaxStack()
	for i := 0; i < 5; i++ {
		s.Push(7)
	}
	for i := 0; i < 5; i++ {
		if got := s.PopMax(); got != 7 {
			t.Fatalf("PopMax() iteration %d = %d, want 7", i, got)
		}
	}
}

func TestMaxStack_NegativeValues(t *testing.T) {
	s := NewMaxStack()
	s.Push(-5)
	s.Push(-1)
	s.Push(-10)
	if got := s.PeekMax(); got != -1 {
		t.Fatalf("PeekMax() = %d, want -1", got)
	}
	if got := s.PopMax(); got != -1 {
		t.Fatalf("PopMax() = %d, want -1", got)
	}
	if got := s.PeekMax(); got != -5 {
		t.Fatalf("PeekMax() = %d, want -5", got)
	}
}

func TestMaxStack_StressManyOperations(t *testing.T) {
	s := NewMaxStack()
	const n = 500

	for i := 1; i <= n; i++ {
		s.Push(i)
	}
	if got := s.PeekMax(); got != n {
		t.Fatalf("PeekMax() = %d, want %d", got, n)
	}
	if got := s.PopMax(); got != n {
		t.Fatalf("PopMax() = %d, want %d", got, n)
	}
	if got := s.Top(); got != n-1 {
		t.Fatalf("Top() = %d, want %d", got, n-1)
	}
	if got := s.PeekMax(); got != n-1 {
		t.Fatalf("PeekMax() = %d, want %d", got, n-1)
	}
}
