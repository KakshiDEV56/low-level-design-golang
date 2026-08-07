package designboundedblockingqueue

/*
Problem Statement
Implement a thread-safe bounded blocking queue that has the following methods :

BoundedBlockingQueue(int capacity): The constructor initializes the queue with a maximum capacity .

void enqueue(int element): Adds an element to the front of the queue. If the queue is full, the calling thread is blocked until the queue is no longer full .

int dequeue(): Returns the element at the rear of the queue and removes it. If the queue is empty, the calling thread is blocked until the queue is no longer empty .

int size(): Returns the number of elements currently in the queue .

Key Requirements:

Your implementation will be tested using multiple threads at the same time. Each thread will either be a producer thread that only makes calls to the enqueue method or a consumer thread that only makes calls to the dequeue method .

The size method will be called after every test case .

Please do not use built-in implementations of bounded blocking queue as this will not be accepted in an interview .

Example 1
Input:

text
1
1
["BoundedBlockingQueue","enqueue","dequeue","dequeue","enqueue","enqueue","enqueue","enqueue","dequeue"]
[[2],[1],[],[],[0],[2],[3],[4],[]]
Output:

text
[1,0,2,2]
Explanation:

Number of producer threads = 1, Number of consumer threads = 1 .

BoundedBlockingQueue queue = new BoundedBlockingQueue(2); // initialize queue with capacity = 2.

queue.enqueue(1); // Producer thread enqueues 1.

queue.dequeue(); // Consumer thread returns 1.

queue.dequeue(); // Consumer thread blocks because queue is empty.

queue.enqueue(0); // Producer enqueues 0. Consumer unblocks and returns 0.

queue.enqueue(2); // Producer enqueues 2.

queue.enqueue(3); // Producer enqueues 3.

queue.enqueue(4); // Producer blocks because capacity (2) is reached.

queue.dequeue(); // Consumer returns 2. Producer unblocks and enqueues 4.

queue.size(); // 2 elements remaining in the queue .

Example 2
Input:

text
3
4
["BoundedBlockingQueue","enqueue","enqueue","enqueue","dequeue","dequeue","dequeue","enqueue"]
[[3],[1],[0],[2],[],[],[],[3]]
Output:

text
[1,0,2,1]
Explanation:

Number of producer threads = 3, Number of consumer threads = 4 .

BoundedBlockingQueue queue = new BoundedBlockingQueue(3); // initialize queue with capacity = 3.

queue.enqueue(1); // Producer P1 enqueues 1.

queue.enqueue(0); // Producer P2 enqueues 0.

queue.enqueue(2); // Producer P3 enqueues 2.

queue.dequeue(); // Consumer C1 calls dequeue.

queue.dequeue(); // Consumer C2 calls dequeue.

queue.dequeue(); // Consumer C3 calls dequeue.

queue.enqueue(3); // One of the producers enqueues 3.

queue.size(); // 1 element remaining in the queue .

Important Note on Output Ordering: Since the number of threads for producer/consumer is greater than 1, we do not know how the threads will be scheduled in the operating system. Therefore, any order of the dequeued values like [1,0,2], [1,2,0], [0,1,2], etc., will be accepted .

Constraints
1 <= Number of Producers <= 8

1 <= Number of Consumers <= 8

1 <= size <= 30

0 <= element <= 20

The number of calls to enqueue is greater than or equal to the number of calls to dequeue.

At most 40 calls will be made to enqueue, dequeue, and size .


*/

type BoundedBlockingQueue struct {
	// Define your fields here
}

func NewBoundedBlockingQueue(capacity int) *BoundedBlockingQueue {
	// Initialize your queue here
	return &BoundedBlockingQueue{}
}

func (q *BoundedBlockingQueue) Enqueue(element int) {
	// Add element to the front of the queue. Block if full.
}

func (q *BoundedBlockingQueue) Dequeue() int {
	// Return element at the rear of the queue and remove it. Block if empty.
	return 0
}

func (q *BoundedBlockingQueue) Size() int {
	// Return the current number of elements in the queue.
	return 0
}
