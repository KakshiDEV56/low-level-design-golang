package designcirculardeque

/*
LeetCode 641. Design Circular Deque

Design your implementation of the circular double-ended queue (deque).

Implement the MyCircularDeque class:
  - NewMyCircularDeque(k int) initializes the deque with a maximum size of k.
  - InsertFront(value int) bool adds an item at the front of the deque. Returns true if
    successful.
  - InsertLast(value int) bool adds an item at the rear of the deque. Returns true if
    successful.
  - DeleteFront() bool deletes an item from the front of the deque. Returns true if
    successful.
  - DeleteLast() bool deletes an item from the rear of the deque. Returns true if
    successful.
  - GetFront() int returns the front item, or -1 if the deque is empty.
  - GetRear() int returns the last item, or -1 if the deque is empty.
  - IsEmpty() bool returns true if the deque is empty.
  - IsFull() bool returns true if the deque is full.

Example 1:

	Input:
	["MyCircularDeque", "insertLast", "insertLast", "insertFront", "insertFront",
	 "getRear", "isFull", "deleteLast", "insertFront", "getFront"]
	[[3], [1], [2], [3], [4], [], [], [], [4], []]
	Output:
	[null, true, true, true, false, 2, true, true, true, 4]

	Explanation:
	MyCircularDeque dq = new MyCircularDeque(3);
	dq.insertLast(1);  // true, deque: [1]
	dq.insertLast(2);  // true, deque: [1,2]
	dq.insertFront(3); // true, deque: [3,1,2]
	dq.insertFront(4); // false, deque is full
	dq.getRear();      // 2
	dq.isFull();       // true
	dq.deleteLast();   // true, deque: [3,1]
	dq.insertFront(4); // true, deque: [4,3,1]
	dq.getFront();     // 4

Constraints:
  - 1 <= k <= 1000
  - 0 <= value <= 1000
  - At most 2000 calls will be made to InsertFront, InsertLast, DeleteFront, DeleteLast,
    GetFront, GetRear, IsEmpty, and IsFull.

TODO:
 1. Add the fields you need on MyCircularDeque below. Same core idea as a circular queue —
    a fixed-size slice as a ring buffer — but now you need to insert/remove from BOTH
    ends in O(1), so track both a `front` index and a `size` (or `front`/`rear` indices),
    using modulo arithmetic to wrap in either direction (careful: wrapping backwards from
    index 0 needs `(front - 1 + k) % k`, not `front - 1`, since Go's `%` can return
    negative results for negative operands).
 2. Implement NewMyCircularDeque(k): allocate a buffer of capacity k.
 3. Implement InsertFront/InsertLast: if full, return false; otherwise write the value at
    the appropriate end (wrapping as needed) and return true.
 4. Implement DeleteFront/DeleteLast: if empty, return false; otherwise advance/retreat
    the appropriate index (wrapping as needed) and return true.
 5. Implement GetFront/GetRear: return -1 if empty, otherwise the value at that end.
 6. Implement IsEmpty/IsFull based on your size/count tracking.
*/

type MyCircularDeque struct {
	// TODO: add your fields here
}

func NewMyCircularDeque(k int) *MyCircularDeque {
	// TODO: implement
	return &MyCircularDeque{}
}

func (d *MyCircularDeque) InsertFront(value int) bool {
	// TODO: implement
	return false
}

func (d *MyCircularDeque) InsertLast(value int) bool {
	// TODO: implement
	return false
}

func (d *MyCircularDeque) DeleteFront() bool {
	// TODO: implement
	return false
}

func (d *MyCircularDeque) DeleteLast() bool {
	// TODO: implement
	return false
}

func (d *MyCircularDeque) GetFront() int {
	// TODO: implement
	return -1
}

func (d *MyCircularDeque) GetRear() int {
	// TODO: implement
	return -1
}

func (d *MyCircularDeque) IsEmpty() bool {
	// TODO: implement
	return true
}

func (d *MyCircularDeque) IsFull() bool {
	// TODO: implement
	return false
}
