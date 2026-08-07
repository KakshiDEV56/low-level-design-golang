package minstack

/*
LeetCode 155. Min Stack

Design a stack that supports push, pop, top, and retrieving the minimum element in
constant time.

Implement the MinStack class:
  - NewMinStack() initializes the stack object.
  - Push(val int) pushes the element val onto the stack.
  - Pop() removes the element on the top of the stack.
  - Top() int gets the top element of the stack.
  - GetMin() int retrieves the minimum element in the stack.

You must implement a solution with O(1) time complexity for each function.

Example 1:

	Input:
	["MinStack", "push", "push", "push", "getMin", "pop", "top", "getMin"]
	[[], [-2], [0], [-3], [], [], [], []]
	Output:
	[null, null, null, null, -3, null, 0, -2]

	Explanation:
	MinStack minStack = new MinStack();
	minStack.push(-2);
	minStack.push(0);
	minStack.push(-3);
	minStack.getMin(); // return -3
	minStack.pop();
	minStack.top();    // return 0
	minStack.getMin(); // return -2

Constraints:
  - -2^31 <= val <= 2^31 - 1
  - Methods Pop, Top and GetMin will always be called on non-empty stacks.
  - At most 3 * 10^4 calls will be made to Push, Pop, Top, and GetMin.

TODO:
 1. Add the fields you need on MinStack below. Hint: a plain stack alone can't answer
    GetMin in O(1) after pops — the classic trick is to keep a SECOND stack alongside the
    main one, where the top of the min-stack always holds "the minimum among everything
    currently in the main stack".
 2. Implement NewMinStack: initialize your stack(s).
 3. Implement Push: push val onto the main stack. Also push onto the min-stack: either
    val itself (if it's <= the current min-stack top) or a repeat of the current minimum,
    so the min-stack always has one entry per main-stack entry.
 4. Implement Pop: pop from both stacks together, keeping them in sync.
 5. Implement Top: return the top of the main stack.
 6. Implement GetMin: return the top of the min-stack.
*/

type MinStack struct {
	// TODO: add your fields here
}

func NewMinStack() *MinStack {
	// TODO: implement
	return &MinStack{}
}

func (s *MinStack) Push(val int) {
	// TODO: implement
}

func (s *MinStack) Pop() {
	// TODO: implement
}

func (s *MinStack) Top() int {
	// TODO: implement
	return 0
}

func (s *MinStack) GetMin() int {
	// TODO: implement
	return 0
}
