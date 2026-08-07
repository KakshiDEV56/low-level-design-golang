package flattennestedlistiterator

/*
LeetCode 341. Flatten Nested List Iterator

You are given a nested list of integers nestedList. Each element is either an integer or
a list whose elements may also be integers or other lists. Implement an iterator to
flatten it.

Implement the NestedIterator class:
  - NewNestedIterator(nestedList []*NestedInteger) initializes the iterator with the
    nested list nestedList.
  - Next() int returns the next integer in the nested list.
  - HasNext() bool returns true if there are still some integers in the nested list and
    false otherwise.

Note: the NestedInteger type below is provided fully implemented so you can focus on
NestedIterator itself — do not modify NestedInteger.

Example 1:

	Input: nestedList = [[1,1],2,[1,1]]
	Output: [1,1,2,1,1]
	Explanation: By calling next repeatedly until hasNext returns false, the order of
	elements returned by next should be: 1, 1, 2, 1, 1.

Example 2:

	Input: nestedList = [1,[4,[6]]]
	Output: [1,4,6]
	Explanation: By calling next repeatedly until hasNext returns false, the order of
	elements returned by next should be: 1, 4, 6.

Constraints:
  - 1 <= nestedList.length <= 500
  - The values of the integers in the nested list is in the range [-10^6, 10^6].

TODO (NestedIterator only — NestedInteger below is already complete):
 1. Add the fields you need on NestedIterator below. Two classic approaches:
      a) Eager flatten: in NewNestedIterator, recursively walk nestedList and collect
         every integer (in order) into a flat slice up front, plus an index cursor.
         Simple, and fine given the constraints (<=500 top-level entries).
      b) Lazy flatten with a stack: push nestedList's elements onto a stack in REVERSE
         order, then in HasNext(), repeatedly peek the top: if it's an integer, you're
         done (return true); if it's a list, pop it and push its elements in reverse
         order instead, then re-check. This does the flattening work incrementally
         instead of all at once — worth implementing after the simple version works.
 2. Implement NewNestedIterator: build your chosen internal representation from
    nestedList.
 3. Implement HasNext(): true if there's a next integer available (see approach above).
 4. Implement Next(): return the next integer and advance past it. It's guaranteed to
    only be called when HasNext() is true.
*/

// NestedInteger is the underlying, already-complete type NestedIterator wraps. Do not
// modify this type. It mirrors LeetCode's NestedInteger interface: each value is either
// a single integer OR a list of more NestedIntegers, never both.
type NestedInteger struct {
	isInt bool
	val   int
	list  []*NestedInteger
}

// NewInt creates a NestedInteger holding a single integer.
func NewInt(v int) *NestedInteger {
	return &NestedInteger{isInt: true, val: v}
}

// NewList creates a NestedInteger holding a list of other NestedIntegers.
func NewList(items ...*NestedInteger) *NestedInteger {
	return &NestedInteger{isInt: false, list: items}
}

func (n *NestedInteger) IsInteger() bool {
	return n.isInt
}

// GetInteger returns the single integer value. Only valid when IsInteger() is true.
func (n *NestedInteger) GetInteger() int {
	return n.val
}

// GetList returns the nested list. Only valid when IsInteger() is false.
func (n *NestedInteger) GetList() []*NestedInteger {
	return n.list
}

type NestedIterator struct {
	// TODO: add your fields here
}

func NewNestedIterator(nestedList []*NestedInteger) *NestedIterator {
	// TODO: implement
	return &NestedIterator{}
}

func (it *NestedIterator) Next() int {
	// TODO: implement
	return 0
}

func (it *NestedIterator) HasNext() bool {
	// TODO: implement
	return false
}
