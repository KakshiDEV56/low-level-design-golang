package designcircularqueue

/*
LeetCode 622. Design Circular Queue

Design your implementation of the circular queue. The circular queue is a linear data
structure in which the operations are performed based on FIFO (First In First Out)
principle, and the last position is connected back to the first position to make a
circle. It is also called "Ring Buffer".

Implement the MyCircularQueue class:
  - NewMyCircularQueue(k int) initializes the object with the size of the queue to be k.
  - EnQueue(value int) bool inserts an element into the circular queue. Returns true if
    the operation is successful.
  - DeQueue() bool deletes an element from the circular queue. Returns true if the
    operation is successful.
  - Front() int gets the front item from the queue. Returns -1 if the queue is empty.
  - Rear() int gets the last item from the queue. Returns -1 if the queue is empty.
  - IsEmpty() bool checks whether the circular queue is empty or not.
  - IsFull() bool checks whether the circular queue is full or not.

Example 1:

	Input:
	["MyCircularQueue", "enQueue", "enQueue", "enQueue", "enQueue", "Rear", "isFull",
	 "deQueue", "enQueue", "Rear"]
	[[3], [1], [2], [3], [4], [], [], [], [4], []]
	Output:
	[null, true, true, true, false, 3, true, true, true, 4]

	Explanation:
	MyCircularQueue q = new MyCircularQueue(3);
	q.enQueue(1); // true, queue: [1]
	q.enQueue(2); // true, queue: [1,2]
	q.enQueue(3); // true, queue: [1,2,3]
	q.enQueue(4); // false, queue is full
	q.Rear();     // 3
	q.isFull();   // true
	q.deQueue();  // true, queue: [2,3]
	q.enQueue(4); // true, queue: [2,3,4]
	q.Rear();     // 4

Constraints:
  - 1 <= k <= 1000
  - 0 <= value <= 1000
  - At most 3000 calls will be made to EnQueue, DeQueue, Front, Rear, IsEmpty, and IsFull.

TODO:
 1. Add the fields you need on MyCircularQueue below. Hint: a fixed-size slice used as a
    ring buffer, plus a `front` index and a `size` (count of elements currently stored) —
    or `front`/`rear` indices — is enough to support every operation in O(1). Avoid
    growing/shrinking a slice with append/re-slicing; the whole point is fixed capacity
    with index arithmetic (mod k).
 2. Implement NewMyCircularQueue(k): allocate a buffer of capacity k.
 3. Implement EnQueue: if full, return false. Otherwise write the value at the next free
    slot (wrapping with modulo) and return true.
 4. Implement DeQueue: if empty, return false. Otherwise advance the front index
    (wrapping with modulo) and return true.
 5. Implement Front/Rear: return -1 if empty, otherwise the value at the front/rear index.
 6. Implement IsEmpty/IsFull based on your size/count tracking.
*/

type MyCircularQueue struct {
	// TODO: add your fields here
}

func NewMyCircularQueue(k int) *MyCircularQueue {
	// TODO: implement
	return &MyCircularQueue{}
}

func (q *MyCircularQueue) EnQueue(value int) bool {
	// TODO: implement
	return false
}

func (q *MyCircularQueue) DeQueue() bool {
	// TODO: implement
	return false
}

func (q *MyCircularQueue) Front() int {
	// TODO: implement
	return -1
}

func (q *MyCircularQueue) Rear() int {
	// TODO: implement
	return -1
}

func (q *MyCircularQueue) IsEmpty() bool {
	// TODO: implement
	return true
}

func (q *MyCircularQueue) IsFull() bool {
	// TODO: implement
	return false
}
