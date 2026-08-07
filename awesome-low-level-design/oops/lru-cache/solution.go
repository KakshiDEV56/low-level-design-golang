package lrucache

/*
LeetCode 146. LRU Cache

Design a data structure that follows the constraints of a Least Recently Used (LRU) cache.

Implement the LRUCache class:
  - NewLRUCache(capacity int) initializes the LRU cache with positive size capacity.
  - Get(key int) int returns the value of the key if the key exists, otherwise returns -1.
  - Put(key, value int) updates the value of the key if the key exists. Otherwise, adds the
    key-value pair to the cache. If the number of keys exceeds the capacity from this
    operation, evicts the least recently used key.

Get and Put must each run in O(1) average time complexity.

Example 1:

	Input:
	["LRUCache", "put", "put", "get", "put", "get", "put", "get", "get", "get"]
	[[2], [1, 1], [2, 2], [1], [3, 3], [2], [4, 4], [1], [3], [4]]
	Output:
	[null, null, null, 1, null, -1, null, -1, 3, 4]

	Explanation:
	LRUCache lru = new LRUCache(2);
	lru.put(1, 1);  // cache is {1=1}
	lru.put(2, 2);  // cache is {1=1, 2=2}
	lru.get(1);     // returns 1
	lru.put(3, 3);  // evicts key 2, cache is {1=1, 3=3}
	lru.get(2);     // returns -1 (not found)
	lru.put(4, 4);  // evicts key 1, cache is {4=4, 3=3}
	lru.get(1);     // returns -1 (not found)
	lru.get(3);     // returns 3
	lru.get(4);     // returns 4

Constraints:
  - 1 <= capacity <= 3000
  - 0 <= key <= 10^4
  - 0 <= value <= 10^5
  - At most 2 * 10^5 calls will be made to Get and Put.

TODO:
 1. Add the fields you need on LRUCache below. Hint: O(1) lookup + O(1) "move to front /
    evict from back" usually means a hash map keyed by key, paired with a doubly linked
    list ordered by recency (most-recently-used at one end, least-recently-used at the
    other).
 2. Implement NewLRUCache to initialize your chosen data structures with the given capacity.
 3. Implement Get: return the value for key, or -1 if the key is absent. A successful Get
    also counts as "using" the key — it must become the most recently used entry.
 4. Implement Put: insert or update the key's value, and mark it as most recently used.
    If inserting a brand-new key would push the cache over capacity, evict the least
    recently used entry first.
*/

type LRUCache struct {
	// TODO: add your fields here
}

func NewLRUCache(capacity int) *LRUCache {
	// TODO: implement
	return &LRUCache{}
}

func (c *LRUCache) Get(key int) int {
	// TODO: implement
	return -1
}

func (c *LRUCache) Put(key int, value int) {
	// TODO: implement
}
