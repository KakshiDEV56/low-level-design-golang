package peekingiterator

/*
LeetCode 284. Peeking Iterator

Design an iterator that supports the peek operation on an existing iterator in addition
to the hasNext and the next operations.

Implement the PeekingIterator class:
  - NewPeekingIterator(iter *Iterator) initializes the object with the given underlying
    iterator.
  - Peek() int returns the next element in the iteration without advancing the iterator.
  - Next() int returns the next element in the iteration and advances the iterator.
  - HasNext() bool returns true if there are still elements in the iteration.

Note: each language's implementation may be slightly different from the original
LeetCode problem (which normally extends a language-specific Iterator interface). Here,
the underlying Iterator type below is provided fully implemented so you can focus on
PeekingIterator itself — do not modify Iterator.

Example 1:

	Input: ["PeekingIterator", "next", "peek", "next", "next", "hasNext"]
	      underlying data: [1, 2, 3]
	Output: [null, 1, 2, 2, 3, false]

	Explanation:
	PeekingIterator pi = new PeekingIterator(Iterator([1, 2, 3])); // [1,2,3]
	pi.next();    // return 1, the pointer moves to the next element [1,2,3].
	pi.peek();    // return 2, the pointer does not move [1,2,3].
	pi.next();    // return 2, the pointer moves to the next element [1,2,3]
	pi.next();    // return 3, the pointer moves to the next element [1,2,3]
	pi.hasNext(); // return false

Constraints:
  - 1 <= nums.length <= 1000
  - 1 <= nums[i] <= 1000
  - All the calls to Next and Peek are valid (i.e. it's guaranteed that there is a next
    element when Next or Peek is called).
  - At most 1000 calls will be made to Next, HasNext, and Peek.

TODO (PeekingIterator only — Iterator below is already complete):
 1. Add the fields you need on PeekingIterator below. Hint: wrap the underlying Iterator,
    plus a small cache for "the value we've already pulled out of the underlying
    iterator via Next(), that Peek() reported but the caller hasn't consumed yet" and a
    flag/marker for whether that cache is currently valid.
 2. Implement NewPeekingIterator(iter): store the underlying iterator, cache empty.
 3. Implement Peek(): if nothing is cached yet, pull the next value from the underlying
    iterator into the cache. Return the cached value WITHOUT clearing it (repeated Peek
    calls must return the same value until Next is called).
 4. Implement Next(): if something is cached (from a prior Peek), clear and return it.
    Otherwise pull directly from the underlying iterator.
 5. Implement HasNext(): true if there's a cached value, OR the underlying iterator still
    has more.
*/

// Iterator is the underlying, already-complete iterator that PeekingIterator wraps.
// Do not modify this type.
type Iterator struct {
	data []int
	idx  int
}

func NewIterator(data []int) *Iterator {
	cp := make([]int, len(data))
	copy(cp, data)
	return &Iterator{data: cp}
}

func (it *Iterator) HasNext() bool {
	return it.idx < len(it.data)
}

func (it *Iterator) Next() int {
	v := it.data[it.idx]
	it.idx++
	return v
}

type PeekingIterator struct {
	// TODO: add your fields here
}

func NewPeekingIterator(iter *Iterator) *PeekingIterator {
	// TODO: implement
	return &PeekingIterator{}
}

func (p *PeekingIterator) Peek() int {
	// TODO: implement
	return 0
}

func (p *PeekingIterator) Next() int {
	// TODO: implement
	return 0
}

func (p *PeekingIterator) HasNext() bool {
	// TODO: implement
	return false
}
