package designstackwithincrement

import "testing"

func TestCustomStack_Example1(t *testing.T) {
	s := NewCustomStack(3)
	s.Push(1)
	s.Push(2)
	if got := s.Pop(); got != 2 {
		t.Fatalf("Pop() = %d, want 2", got)
	}
	s.Push(2)
	s.Push(3)
	s.Push(4) // stack is full (maxSize=3), this must be a no-op
	s.Increment(5, 100)
	s.Increment(2, 100)
	if got := s.Pop(); got != 103 {
		t.Fatalf("Pop() = %d, want 103", got)
	}
	if got := s.Pop(); got != 202 {
		t.Fatalf("Pop() = %d, want 202", got)
	}
	if got := s.Pop(); got != 201 {
		t.Fatalf("Pop() = %d, want 201", got)
	}
	if got := s.Pop(); got != -1 {
		t.Fatalf("Pop() = %d, want -1", got)
	}
}

func TestCustomStack_PopOnEmptyStack(t *testing.T) {
	s := NewCustomStack(5)
	if got := s.Pop(); got != -1 {
		t.Fatalf("Pop() on empty stack = %d, want -1", got)
	}
}

func TestCustomStack_PushBeyondMaxSizeIsNoOp(t *testing.T) {
	s := NewCustomStack(2)
	s.Push(1)
	s.Push(2)
	s.Push(3) // dropped, stack is full
	if got := s.Pop(); got != 2 {
		t.Fatalf("Pop() = %d, want 2", got)
	}
	if got := s.Pop(); got != 1 {
		t.Fatalf("Pop() = %d, want 1", got)
	}
	if got := s.Pop(); got != -1 {
		t.Fatalf("Pop() = %d, want -1", got)
	}
}

func TestCustomStack_IncrementOnEmptyStackIsNoOp(t *testing.T) {
	s := NewCustomStack(3)
	s.Increment(5, 100) // must not panic
	s.Push(1)
	if got := s.Pop(); got != 1 {
		t.Fatalf("Pop() = %d, want 1 (increment on empty stack shouldn't affect later pushes)", got)
	}
}

func TestCustomStack_IncrementKLargerThanSize(t *testing.T) {
	s := NewCustomStack(5)
	s.Push(1)
	s.Push(2)
	s.Increment(100, 10) // k larger than current size -> increments everything
	if got := s.Pop(); got != 12 {
		t.Fatalf("Pop() = %d, want 12", got)
	}
	if got := s.Pop(); got != 11 {
		t.Fatalf("Pop() = %d, want 11", got)
	}
}

func TestCustomStack_IncrementKEqualsOne(t *testing.T) {
	s := NewCustomStack(5)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Increment(1, 50) // only the bottom element changes
	if got := s.Pop(); got != 3 {
		t.Fatalf("Pop() = %d, want 3", got)
	}
	if got := s.Pop(); got != 2 {
		t.Fatalf("Pop() = %d, want 2", got)
	}
	if got := s.Pop(); got != 51 {
		t.Fatalf("Pop() = %d, want 51", got)
	}
}

func TestCustomStack_IncrementValZero(t *testing.T) {
	s := NewCustomStack(3)
	s.Push(5)
	s.Push(6)
	s.Increment(2, 0)
	if got := s.Pop(); got != 6 {
		t.Fatalf("Pop() = %d, want 6", got)
	}
	if got := s.Pop(); got != 5 {
		t.Fatalf("Pop() = %d, want 5", got)
	}
}

func TestCustomStack_PushPopIncrementInterleaved(t *testing.T) {
	s := NewCustomStack(4)
	s.Push(1)
	s.Increment(1, 10) // [11]
	s.Push(2)           // [11, 2]
	s.Increment(2, 5)   // [16, 7]
	s.Push(3)           // [16, 7, 3]
	if got := s.Pop(); got != 3 {
		t.Fatalf("Pop() = %d, want 3", got)
	}
	if got := s.Pop(); got != 7 {
		t.Fatalf("Pop() = %d, want 7", got)
	}
	if got := s.Pop(); got != 16 {
		t.Fatalf("Pop() = %d, want 16", got)
	}
}

func TestCustomStack_RefillAfterFullDrain(t *testing.T) {
	s := NewCustomStack(2)
	s.Push(1)
	s.Push(2)
	s.Pop()
	s.Pop()
	s.Push(9)
	if got := s.Pop(); got != 9 {
		t.Fatalf("Pop() = %d, want 9", got)
	}
}

func TestCustomStack_StressManyIncrements(t *testing.T) {
	const maxSize = 100
	s := NewCustomStack(maxSize)
	for i := 1; i <= maxSize; i++ {
		s.Push(i)
	}
	for i := 0; i < 50; i++ {
		s.Increment(maxSize, 1)
	}
	// Every element should have gone up by 50.
	for i := maxSize; i >= 1; i-- {
		want := i + 50
		if got := s.Pop(); got != want {
			t.Fatalf("Pop() = %d, want %d", got, want)
		}
	}
}
