package designstackwithincrement

/*
LeetCode 1381. Design a Stack With Increment Operation

Design a stack that supports increment operations on its elements.

Implement the CustomStack class:
  - NewCustomStack(maxSize int) initializes the object with maxSize which is the maximum
    number of elements in the stack.
  - Push(x int) adds x to the top of the stack if the stack has not reached maxSize.
  - Pop() int pops and returns the top of stack, or -1 if the stack is empty.
  - Increment(k, val int) increments the bottom k elements of the stack by val. If there
    are fewer than k elements in the stack, increments all of them.

Example 1:

	Input:
	["CustomStack", "push", "push", "pop", "push", "push", "push", "increment",
	 "increment", "pop", "pop", "pop", "pop"]
	[[3], [1], [2], [], [2], [3], [4], [5, 100], [2, 100], [], [], [], []]
	Output:
	[null, null, null, 2, null, null, null, null, null, 103, 202, 201, -1]

	Explanation:
	CustomStack stk = new CustomStack(3);       // Stack is Empty []
	stk.push(1);                                // stack becomes [1]
	stk.push(2);                                // stack becomes [1, 2]
	stk.pop();                                  // return 2 -> [1]
	stk.push(2);                                // stack becomes [1, 2]
	stk.push(3);                                // stack becomes [1, 2, 3]
	stk.push(4);                                // stack still [1, 2, 3], maxSize reached
	stk.increment(5, 100);                      // stack becomes [101, 102, 103]
	stk.increment(2, 100);                      // stack becomes [201, 202, 103]
	stk.pop();                                  // return 103 -> [201, 202]
	stk.pop();                                  // return 202 -> [201]
	stk.pop();                                  // return 201 -> []
	stk.pop();                                  // return -1 -> [] stack is empty

Constraints:
  - 1 <= maxSize <= 1000
  - 1 <= x <= 1000
  - 1 <= k <= 1000
  - 0 <= val <= 100
  - At most 1000 calls will be made to each method of Increment, Push, and Pop each
    separately.

TODO:
 1. Add the fields you need on CustomStack below. A slice used as a stack (append/reslice)
    plus a `maxSize` field is enough — the "bottom k elements" phrasing means index 0 is
    the bottom, so Increment can be a simple loop over the first min(k, len) elements.
 2. Implement NewCustomStack(maxSize): remember the capacity limit.
 3. Implement Push: no-op if the stack is already at maxSize, otherwise append x.
 4. Implement Pop: return -1 if empty, otherwise remove and return the top (last) element.
 5. Implement Increment(k, val): add val to each of the bottom min(k, current size)
    elements. This should still work fine as an O(k) loop given the constraints — no
    need for a fancier lazy-propagation trick unless you want the extra challenge.
*/

type CustomStack struct {
	// TODO: add your fields here
}

func NewCustomStack(maxSize int) *CustomStack {
	// TODO: implement
	return &CustomStack{}
}

func (s *CustomStack) Push(x int) {
	// TODO: implement
}

func (s *CustomStack) Pop() int {
	// TODO: implement
	return -1
}

func (s *CustomStack) Increment(k int, val int) {
	// TODO: implement
}
