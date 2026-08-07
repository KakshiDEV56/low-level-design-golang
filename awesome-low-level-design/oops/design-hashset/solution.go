package designhashset

/*
LeetCode 705. Design HashSet

Design a HashSet without using any built-in hash table libraries.

Implement MyHashSet class:
  - NewMyHashSet() initializes the object with an empty set.
  - Add(key int) inserts the value key into the HashSet.
  - Remove(key int) removes the value key if it exists in the HashSet.
  - Contains(key int) bool returns whether key exists in the HashSet or not.

Example 1:

	Input:
	["MyHashSet", "add", "add", "contains", "contains", "add", "contains", "remove", "contains"]
	[[], [1], [2], [1], [3], [2], [2], [2], [2]]
	Output:
	[null, null, null, true, false, null, true, null, false]

	Explanation:
	MyHashSet s = new MyHashSet();
	s.add(1);          // set = [1]
	s.add(2);          // set = [1, 2]
	s.contains(1);      // returns true
	s.contains(3);      // returns false (not found)
	s.add(2);           // set = [1, 2] (no change, already present)
	s.contains(2);      // returns true
	s.remove(2);         // set = [1]
	s.contains(2);       // returns false

Constraints:
  - 0 <= key <= 10^6
  - At most 10^4 calls will be made to Add, Remove, and Contains.

TODO:
 1. Like Design HashMap, the point is to build this yourself — do NOT back it with a Go
    `map[int]bool` or `map[int]struct{}`. Instead:
      - Pick a fixed number of buckets.
      - Give each bucket a chain (slice or linked list) of keys, to handle collisions.
 2. Write a hash function mapping a key to a bucket index (e.g. key % numBuckets).
 3. Implement NewMyHashSet: allocate the bucket array.
 4. Implement Add: hash the key; if it's not already present in that bucket's chain,
    append it. Adding an already-present key must be a no-op (no duplicates).
 5. Implement Remove: hash the key; remove it from its bucket's chain if present
    (no-op if the key isn't there).
 6. Implement Contains: hash the key, scan its bucket's chain for a match.
*/

type MyHashSet struct {
	// TODO: add your fields here (e.g. a slice of buckets)
}

func NewMyHashSet() *MyHashSet {
	// TODO: implement
	return &MyHashSet{}
}

func (s *MyHashSet) Add(key int) {
	// TODO: implement
}

func (s *MyHashSet) Remove(key int) {
	// TODO: implement
}

func (s *MyHashSet) Contains(key int) bool {
	// TODO: implement
	return false
}
