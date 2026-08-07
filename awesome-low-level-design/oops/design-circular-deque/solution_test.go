package designcirculardeque

import "testing"

func TestMyCircularDeque_Example1(t *testing.T) {
	dq := NewMyCircularDeque(3)
	if got := dq.InsertLast(1); !got {
		t.Fatal("InsertLast(1) = false, want true")
	}
	if got := dq.InsertLast(2); !got {
		t.Fatal("InsertLast(2) = false, want true")
	}
	if got := dq.InsertFront(3); !got {
		t.Fatal("InsertFront(3) = false, want true")
	}
	if got := dq.InsertFront(4); got {
		t.Fatal("InsertFront(4) = true, want false (deque is full)")
	}
	if got := dq.GetRear(); got != 2 {
		t.Fatalf("GetRear() = %d, want 2", got)
	}
	if !dq.IsFull() {
		t.Fatal("IsFull() = false, want true")
	}
	if got := dq.DeleteLast(); !got {
		t.Fatal("DeleteLast() = false, want true")
	}
	if got := dq.InsertFront(4); !got {
		t.Fatal("InsertFront(4) = false, want true")
	}
	if got := dq.GetFront(); got != 4 {
		t.Fatalf("GetFront() = %d, want 4", got)
	}
}

func TestMyCircularDeque_EmptyDequeBehavior(t *testing.T) {
	dq := NewMyCircularDeque(2)
	if !dq.IsEmpty() {
		t.Fatal("IsEmpty() = false, want true on a fresh deque")
	}
	if got := dq.GetFront(); got != -1 {
		t.Fatalf("GetFront() on empty deque = %d, want -1", got)
	}
	if got := dq.GetRear(); got != -1 {
		t.Fatalf("GetRear() on empty deque = %d, want -1", got)
	}
	if got := dq.DeleteFront(); got {
		t.Fatal("DeleteFront() on empty deque = true, want false")
	}
	if got := dq.DeleteLast(); got {
		t.Fatal("DeleteLast() on empty deque = true, want false")
	}
}

func TestMyCircularDeque_CapacityOne(t *testing.T) {
	dq := NewMyCircularDeque(1)
	if got := dq.InsertFront(10); !got {
		t.Fatal("InsertFront(10) = false, want true")
	}
	if !dq.IsFull() {
		t.Fatal("IsFull() = false, want true")
	}
	if got := dq.InsertLast(20); got {
		t.Fatal("InsertLast(20) = true, want false (capacity 1 is full)")
	}
	if got := dq.GetFront(); got != 10 {
		t.Fatalf("GetFront() = %d, want 10", got)
	}
	if got := dq.GetRear(); got != 10 {
		t.Fatalf("GetRear() = %d, want 10", got)
	}
}

func TestMyCircularDeque_InsertFrontOnly(t *testing.T) {
	dq := NewMyCircularDeque(4)
	dq.InsertFront(1)
	dq.InsertFront(2)
	dq.InsertFront(3)
	// Order front->rear should be [3, 2, 1].
	if got := dq.GetFront(); got != 3 {
		t.Fatalf("GetFront() = %d, want 3", got)
	}
	if got := dq.GetRear(); got != 1 {
		t.Fatalf("GetRear() = %d, want 1", got)
	}
}

func TestMyCircularDeque_WrapAroundBothDirections(t *testing.T) {
	dq := NewMyCircularDeque(3)
	dq.InsertLast(1)
	dq.InsertLast(2)
	dq.InsertLast(3)
	dq.DeleteFront() // front->rear: [2,3]
	dq.InsertLast(4) // wraps rear around: [2,3,4]
	if got := dq.GetFront(); got != 2 {
		t.Fatalf("GetFront() = %d, want 2", got)
	}
	if got := dq.GetRear(); got != 4 {
		t.Fatalf("GetRear() = %d, want 4", got)
	}
	dq.DeleteLast()  // [2,3]
	dq.InsertFront(9) // wraps front backwards: [9,2,3]
	if got := dq.GetFront(); got != 9 {
		t.Fatalf("GetFront() = %d, want 9", got)
	}
	if got := dq.GetRear(); got != 3 {
		t.Fatalf("GetRear() = %d, want 3", got)
	}
}

func TestMyCircularDeque_FullyDrainAndRefill(t *testing.T) {
	dq := NewMyCircularDeque(2)
	dq.InsertLast(1)
	dq.InsertLast(2)
	dq.DeleteFront()
	dq.DeleteFront()
	if !dq.IsEmpty() {
		t.Fatal("IsEmpty() = false, want true after draining all elements")
	}
	dq.InsertFront(5)
	if got := dq.GetFront(); got != 5 {
		t.Fatalf("GetFront() = %d, want 5", got)
	}
	if got := dq.GetRear(); got != 5 {
		t.Fatalf("GetRear() = %d, want 5", got)
	}
}

func TestMyCircularDeque_MixedInsertDeleteBothEnds(t *testing.T) {
	dq := NewMyCircularDeque(5)
	dq.InsertLast(1)   // [1]
	dq.InsertFront(2)  // [2,1]
	dq.InsertLast(3)   // [2,1,3]
	dq.InsertFront(4)  // [4,2,1,3]
	dq.DeleteLast()    // [4,2,1]
	dq.InsertLast(5)   // [4,2,1,5]
	if got := dq.GetFront(); got != 4 {
		t.Fatalf("GetFront() = %d, want 4", got)
	}
	if got := dq.GetRear(); got != 5 {
		t.Fatalf("GetRear() = %d, want 5", got)
	}
}

func TestMyCircularDeque_RepeatedFullEmptyCycles(t *testing.T) {
	dq := NewMyCircularDeque(3)
	for cycle := 0; cycle < 5; cycle++ {
		for v := 1; v <= 3; v++ {
			if !dq.InsertLast(v) {
				t.Fatalf("cycle %d: InsertLast(%d) = false, want true", cycle, v)
			}
		}
		if !dq.IsFull() {
			t.Fatalf("cycle %d: IsFull() = false, want true", cycle)
		}
		for v := 1; v <= 3; v++ {
			if got := dq.GetFront(); got != v {
				t.Fatalf("cycle %d: GetFront() = %d, want %d", cycle, got, v)
			}
			if !dq.DeleteFront() {
				t.Fatalf("cycle %d: DeleteFront() = false, want true", cycle)
			}
		}
		if !dq.IsEmpty() {
			t.Fatalf("cycle %d: IsEmpty() = false, want true", cycle)
		}
	}
}

func TestMyCircularDeque_LargeCapacityStress(t *testing.T) {
	const k = 100
	dq := NewMyCircularDeque(k)

	for i := 0; i < k; i++ {
		dq.InsertLast(i)
	}
	for round := 0; round < 10; round++ {
		dq.DeleteFront()
		dq.InsertLast(1000 + round)
	}
	if !dq.IsFull() {
		t.Fatal("IsFull() = false, want true")
	}
	if got := dq.GetFront(); got != 10 {
		t.Fatalf("GetFront() = %d, want 10", got)
	}
	if got := dq.GetRear(); got != 1009 {
		t.Fatalf("GetRear() = %d, want 1009", got)
	}
}
