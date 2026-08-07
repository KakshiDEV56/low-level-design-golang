package designcircularqueue

import "testing"

func TestMyCircularQueue_Example1(t *testing.T) {
	q := NewMyCircularQueue(3)
	if got := q.EnQueue(1); !got {
		t.Fatal("EnQueue(1) = false, want true")
	}
	if got := q.EnQueue(2); !got {
		t.Fatal("EnQueue(2) = false, want true")
	}
	if got := q.EnQueue(3); !got {
		t.Fatal("EnQueue(3) = false, want true")
	}
	if got := q.EnQueue(4); got {
		t.Fatal("EnQueue(4) = true, want false (queue is full)")
	}
	if got := q.Rear(); got != 3 {
		t.Fatalf("Rear() = %d, want 3", got)
	}
	if got := q.IsFull(); !got {
		t.Fatal("IsFull() = false, want true")
	}
	if got := q.DeQueue(); !got {
		t.Fatal("DeQueue() = false, want true")
	}
	if got := q.EnQueue(4); !got {
		t.Fatal("EnQueue(4) = false, want true")
	}
	if got := q.Rear(); got != 4 {
		t.Fatalf("Rear() = %d, want 4", got)
	}
}

func TestMyCircularQueue_EmptyQueueBehavior(t *testing.T) {
	q := NewMyCircularQueue(2)
	if !q.IsEmpty() {
		t.Fatal("IsEmpty() = false, want true on a fresh queue")
	}
	if got := q.Front(); got != -1 {
		t.Fatalf("Front() on empty queue = %d, want -1", got)
	}
	if got := q.Rear(); got != -1 {
		t.Fatalf("Rear() on empty queue = %d, want -1", got)
	}
	if got := q.DeQueue(); got {
		t.Fatal("DeQueue() on empty queue = true, want false")
	}
}

func TestMyCircularQueue_CapacityOne(t *testing.T) {
	q := NewMyCircularQueue(1)
	if got := q.EnQueue(10); !got {
		t.Fatal("EnQueue(10) = false, want true")
	}
	if !q.IsFull() {
		t.Fatal("IsFull() = false, want true")
	}
	if got := q.EnQueue(20); got {
		t.Fatal("EnQueue(20) = true, want false (capacity 1 is full)")
	}
	if got := q.Front(); got != 10 {
		t.Fatalf("Front() = %d, want 10", got)
	}
	if got := q.Rear(); got != 10 {
		t.Fatalf("Rear() = %d, want 10", got)
	}
	if got := q.DeQueue(); !got {
		t.Fatal("DeQueue() = false, want true")
	}
	if !q.IsEmpty() {
		t.Fatal("IsEmpty() = false, want true after dequeuing the only element")
	}
}

func TestMyCircularQueue_WrapAroundIndices(t *testing.T) {
	q := NewMyCircularQueue(3)
	q.EnQueue(1)
	q.EnQueue(2)
	q.EnQueue(3)
	q.DeQueue() // front now logically at index 1
	q.DeQueue() // front now logically at index 2
	q.EnQueue(4)
	q.EnQueue(5) // this wraps the internal buffer back to index 0
	if got := q.Front(); got != 3 {
		t.Fatalf("Front() = %d, want 3", got)
	}
	if got := q.Rear(); got != 5 {
		t.Fatalf("Rear() = %d, want 5", got)
	}
	if !q.IsFull() {
		t.Fatal("IsFull() = false, want true")
	}
}

func TestMyCircularQueue_FullyDrainAndRefill(t *testing.T) {
	q := NewMyCircularQueue(2)
	q.EnQueue(1)
	q.EnQueue(2)
	q.DeQueue()
	q.DeQueue()
	if !q.IsEmpty() {
		t.Fatal("IsEmpty() = false, want true after draining all elements")
	}
	q.EnQueue(3)
	if got := q.Front(); got != 3 {
		t.Fatalf("Front() = %d, want 3", got)
	}
	if got := q.Rear(); got != 3 {
		t.Fatalf("Rear() = %d, want 3", got)
	}
}

func TestMyCircularQueue_RepeatedFullEmptyCycles(t *testing.T) {
	q := NewMyCircularQueue(3)
	for cycle := 0; cycle < 5; cycle++ {
		for v := 1; v <= 3; v++ {
			if !q.EnQueue(v) {
				t.Fatalf("cycle %d: EnQueue(%d) = false, want true", cycle, v)
			}
		}
		if !q.IsFull() {
			t.Fatalf("cycle %d: IsFull() = false, want true", cycle)
		}
		for v := 1; v <= 3; v++ {
			if got := q.Front(); got != v {
				t.Fatalf("cycle %d: Front() = %d, want %d", cycle, got, v)
			}
			if !q.DeQueue() {
				t.Fatalf("cycle %d: DeQueue() = false, want true", cycle)
			}
		}
		if !q.IsEmpty() {
			t.Fatalf("cycle %d: IsEmpty() = false, want true", cycle)
		}
	}
}

func TestMyCircularQueue_FrontAndRearWithSingleElement(t *testing.T) {
	q := NewMyCircularQueue(5)
	q.EnQueue(99)
	if got := q.Front(); got != 99 {
		t.Fatalf("Front() = %d, want 99", got)
	}
	if got := q.Rear(); got != 99 {
		t.Fatalf("Rear() = %d, want 99", got)
	}
}

func TestMyCircularQueue_LargeCapacityStress(t *testing.T) {
	const k = 100
	q := NewMyCircularQueue(k)

	// Push and pop repeatedly, far more than k times, to exercise wraparound heavily.
	for i := 0; i < k; i++ {
		q.EnQueue(i)
	}
	for round := 0; round < 10; round++ {
		q.DeQueue()
		q.EnQueue(1000 + round)
	}
	if !q.IsFull() {
		t.Fatal("IsFull() = false, want true after equal number of pushes and pops on a full queue")
	}
	if got := q.Front(); got != 10 {
		t.Fatalf("Front() = %d, want 10", got)
	}
	if got := q.Rear(); got != 1009 {
		t.Fatalf("Rear() = %d, want 1009", got)
	}
}
