package maxstack

/*
LeetCode 716. Max Stack

Design a max stack data structure that supports the stack operations and supports finding
the stack's maximum element.

Implement the MaxStack class:
  - NewMaxStack() initializes the stack object.
  - Push(x int) pushes element x onto the stack.
  - Pop() int removes the element on top of the stack and returns it.
  - Top() int gets the element on the top of the stack without removing it.
  - PeekMax() int retrieves the maximum element in the stack without removing it.
  - PopMax() int retrieves the maximum element in the stack and removes it. If there is more
    than one maximum element, only remove the top-most one.

Example 1:

	Input:
	["MaxStack", "push", "push", "push", "top", "popMax", "top", "peekMax", "pop", "top"]
	[[], [5], [1], [5], [], [], [], [], [], []]
	Output:
	[null, null, null, null, 5, 5, 1, 5, 1, 5]

	Explanation:
	MaxStack stk = new MaxStack();
	stk.push(5);   // [5]
	stk.push(1);   // [5, 1]
	stk.push(5);   // [5, 1, 5]
	stk.top();     // return 5
	stk.popMax();  // return 5, stack becomes [5, 1]
	stk.top();     // return 1
	stk.peekMax(); // return 5
	stk.pop();     // return 1, stack becomes [5]
	stk.top();     // return 5

Constraints:
  - -10^7 <= x <= 10^7
  - At most 10^4 calls will be made to Push, Pop, Top, PeekMax, and PopMax.
  - There will be at least one element in the stack when Pop, Top, PeekMax, or PopMax is
    called.

TODO:
 1. Add the fields you need on MaxStack below. A straightforward approach that's plenty
    fast enough for the constraints here: keep two stacks —
      - a main stack holding every pushed value in order, and
      - a max-tracking stack, where the top always holds "the maximum among everything
        currently in the main stack" (same trick as Min Stack).
    Push/Pop/Top/PeekMax are then all O(1).
 2. PopMax is the interesting one, since the max isn't necessarily on top of the main
    stack. A simple O(n) technique: pop elements off both stacks into a temporary buffer
    until the popped value equals the current max; that's the element to remove. Then
    push everything you buffered back onto both stacks (re-pushing also keeps the
    max-stack invariant correct). If you want to go further, look into implementing this
    with O(log n) PopMax using a balanced structure (e.g. a doubly linked list plus a
    sorted map of value -> list of node pointers) — but get the simple version working
    and tested first.
 3. Implement NewMaxStack: initialize your stack(s).
 4. Implement Push, Pop, Top, PeekMax, PopMax per the description above.
*/

type MaxStack struct {
	// TODO: add your fields here
}

func NewMaxStack() *MaxStack {
	// TODO: implement
	return &MaxStack{}
}

func (s *MaxStack) Push(x int) {
	// TODO: implement
}

func (s *MaxStack) Pop() int {
	// TODO: implement
	return 0
}

func (s *MaxStack) Top() int {
	// TODO: implement
	return 0
}

func (s *MaxStack) PeekMax() int {
	// TODO: implement
	return 0
}

func (s *MaxStack) PopMax() int {
	// TODO: implement
	return 0
}
