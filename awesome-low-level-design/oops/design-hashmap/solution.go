package designhashmap

/*
LeetCode 706. Design HashMap

Design a HashMap without using any built-in hash table libraries.

Implement the MyHashMap class:
  - NewMyHashMap() initializes the object with an empty map.
  - Put(key, value int) inserts a (key, value) pair into the HashMap. If the key already
    exists in the map, update the corresponding value.
  - Get(key int) int returns the value to which the specified key is mapped, or -1 if this
    map contains no mapping for the key.
  - Remove(key int) removes the key and its corresponding value if the map contains the
    mapping for the key.

Example 1:

	Input:
	["MyHashMap", "put", "put", "get", "get", "put", "get", "remove", "get"]
	[[], [1, 1], [2, 2], [1], [3], [2, 1], [2], [2], [2]]
	Output:
	[null, null, null, 1, -1, null, 1, null, -1]

	Explanation:
	MyHashMap m = new MyHashMap();
	m.put(1, 1);    // cache: [[1,1]]
	m.put(2, 2);    // cache: [[1,1], [2,2]]
	m.get(1);       // return 1
	m.get(3);       // return -1 (not found)
	m.put(2, 1);    // cache: [[1,1], [2,1]] (update the existing value)
	m.get(2);       // return 1
	m.remove(2);    // remove the mapping for 2
	m.get(2);       // return -1 (not found)

Constraints:
  - 0 <= key, value <= 10^6
  - At most 10^4 calls will be made to Put, Get, and Remove.

TODO:
 1. The whole point of this exercise is to build the hash table yourself — do NOT store
    everything in a single Go `map[int]int`. Instead:
    - Pick a fixed number of buckets (e.g. 1000 or a prime like 769).
    - Add a bucket array field where each bucket holds a chain of (key, value) pairs
    (e.g. a slice or linked list per bucket) to handle collisions.
 2. Write a hash function mapping a key to a bucket index (e.g. key % numBuckets).
 3. Implement NewMyHashMap: allocate the bucket array.
 4. Implement Put: hash the key to find its bucket, then either update the existing
    (key, value) pair in that bucket's chain or append a new one.
 5. Implement Get: hash the key, scan its bucket's chain for a matching key, return the
    value or -1 if not found.
 6. Implement Remove: hash the key, remove the matching entry from its bucket's chain if
    present (no-op if the key isn't there).
*/

type Node struct {
	Key  int
	Val  int
	Next *Node
}
type MyHashMap struct {
	// TODO: add your fields here (e.g. a slice of buckets)

	bucketArr []*Node
}

const (
	bucketSize int = 769
)

func NewMyHashMap() *MyHashMap {
	// TODO: implement
	bucketarr := make([]*Node, bucketSize)

	return &MyHashMap{
		bucketArr: bucketarr,
	}
}

func hash(key int) int {
	hashval := key % bucketSize
	return hashval
}

func (m *MyHashMap) Put(key int, value int) {
	// TODO: implement
	hashIndex := hash(key)
	curr := m.bucketArr[hashIndex]
	if curr == nil {
		m.bucketArr[hashIndex] = &Node{
			Key: key,
			Val: value,
		}
		return
	}
	prev := (*Node)(nil)
	for curr != nil {
		if curr.Key == key {
			curr.Val = value
			return
		}
		prev = curr
		curr = curr.Next
	}

	prev.Next = &Node{
		Key: key, Val: value,
	}
}

func (m *MyHashMap) Get(key int) int {
	// TODO: implement
	hashIndex := hash(key)
	curr := m.bucketArr[hashIndex]
	for curr != nil {
		if curr.Key == key {
			return curr.Val
		}
		curr = curr.Next
	}
	return -1
}

func (m *MyHashMap) Remove(key int) {
	// TODO: implement
	hashIndex := hash(key)
	curr := m.bucketArr[hashIndex]
	prev := (*Node)(nil)
	for curr != nil {
		if curr.Key == key {
			if prev == nil {
				m.bucketArr[hashIndex] = curr.Next
			} else {
				prev.Next = curr.Next
			}
			return
		}
		prev = curr
		curr = curr.Next
	}
}
