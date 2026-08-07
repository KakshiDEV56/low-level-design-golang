package minstack

import "testing"

func TestMinStack_Example1(t *testing.T) {
	s := NewMinStack()
	s.Push(-2)
	s.Push(0)
	s.Push(-3)
	if got := s.GetMin(); got != -3 {
		t.Fatalf("GetMin() = %d, want -3", got)
	}
	s.Pop()
	if got := s.Top(); got != 0 {
		t.Fatalf("Top() = %d, want 0", got)
	}
	if got := s.GetMin(); got != -2 {
		t.Fatalf("GetMin() = %d, want -2", got)
	}
}

func TestMinStack_SingleElement(t *testing.T) {
	s := NewMinStack()
	s.Push(5)
	if got := s.Top(); got != 5 {
		t.Fatalf("Top() = %d, want 5", got)
	}
	if got := s.GetMin(); got != 5 {
		t.Fatalf("GetMin() = %d, want 5", got)
	}
}

func TestMinStack_MinUpdatesOnPushAndPop(t *testing.T) {
	s := NewMinStack()
	s.Push(3)
	s.Push(1)
	s.Push(2)
	if got := s.GetMin(); got != 1 {
		t.Fatalf("GetMin() = %d, want 1", got)
	}
	s.Pop() // removes 2
	if got := s.GetMin(); got != 1 {
		t.Fatalf("GetMin() = %d, want 1", got)
	}
	s.Pop() // removes 1
	if got := s.GetMin(); got != 3 {
		t.Fatalf("GetMin() = %d, want 3", got)
	}
}

func TestMinStack_DuplicateMinimums(t *testing.T) {
	s := NewMinStack()
	s.Push(0)
	s.Push(1)
	s.Push(0)
	if got := s.GetMin(); got != 0 {
		t.Fatalf("GetMin() = %d, want 0", got)
	}
	s.Pop() // removes second 0
	if got := s.GetMin(); got != 0 {
		t.Fatalf("GetMin() = %d, want 0 (first 0 is still on the stack)", got)
	}
	if got := s.Top(); got != 1 {
		t.Fatalf("Top() = %d, want 1", got)
	}
}

func TestMinStack_DescendingPushes(t *testing.T) {
	s := NewMinStack()
	for _, v := range []int{10, 9, 8, 7, 6} {
		s.Push(v)
	}
	if got := s.GetMin(); got != 6 {
		t.Fatalf("GetMin() = %d, want 6", got)
	}
	s.Pop()
	if got := s.GetMin(); got != 7 {
		t.Fatalf("GetMin() = %d, want 7", got)
	}
}

func TestMinStack_AscendingPushes(t *testing.T) {
	s := NewMinStack()
	for _, v := range []int{1, 2, 3, 4, 5} {
		s.Push(v)
	}
	if got := s.GetMin(); got != 1 {
		t.Fatalf("GetMin() = %d, want 1", got)
	}
	for i := 0; i < 4; i++ {
		s.Pop()
	}
	if got := s.GetMin(); got != 1 {
		t.Fatalf("GetMin() = %d, want 1 (only the first push remains)", got)
	}
	if got := s.Top(); got != 1 {
		t.Fatalf("Top() = %d, want 1", got)
	}
}

func TestMinStack_NegativeAndBoundaryValues(t *testing.T) {
	s := NewMinStack()
	s.Push(2147483647)
	s.Push(-2147483648)
	if got := s.GetMin(); got != -2147483648 {
		t.Fatalf("GetMin() = %d, want -2147483648", got)
	}
	s.Pop()
	if got := s.GetMin(); got != 2147483647 {
		t.Fatalf("GetMin() = %d, want 2147483647", got)
	}
}

func TestMinStack_PushPopInterleaved(t *testing.T) {
	s := NewMinStack()
	s.Push(5)
	s.Push(3)
	s.Pop()
	s.Push(7)
	s.Push(1)
	if got := s.GetMin(); got != 1 {
		t.Fatalf("GetMin() = %d, want 1", got)
	}
	s.Pop()
	if got := s.GetMin(); got != 5 {
		t.Fatalf("GetMin() = %d, want 5", got)
	}
	if got := s.Top(); got != 7 {
		t.Fatalf("Top() = %d, want 7", got)
	}
}

func TestMinStack_StressManyOperations(t *testing.T) {
	s := NewMinStack()
	const n = 1000

	for i := n; i >= 1; i-- {
		s.Push(i)
	}
	if got := s.GetMin(); got != 1 {
		t.Fatalf("GetMin() = %d, want 1", got)
	}
	for i := 1; i < n; i++ {
		s.Pop()
		want := i + 1
		if got := s.GetMin(); got != want {
			t.Fatalf("after %d pops, GetMin() = %d, want %d", i, got, want)
		}
	}
}
