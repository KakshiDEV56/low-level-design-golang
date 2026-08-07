package lfucache

/*
LeetCode 460. LFU Cache

Design and implement a data structure for a Least Frequently Used (LFU) cache.

Implement the LFUCache class:
  - NewLFUCache(capacity int) initializes the object with the capacity of the data structure.
  - Get(key int) int gets the value of the key if present in the cache, otherwise returns -1.
  - Put(key, value int) updates the value of key if present, or inserts key if not already
    present. When the cache reaches its capacity, it should invalidate and remove the least
    frequently used key before inserting a new item. For this problem, when there is a tie
    (i.e., two or more keys with the same frequency), the least recently used key is evicted.

When a key is first inserted (via Put), its use counter (frequency) is set to 1 (unless it
gets evicted, of course, in which case a new use counter starts). For a key, Get and Put
operations BOTH count as uses, incrementing the key's use counter.

Get and Put must each run in O(1) average time complexity.

Example 1:

	Input:
	["LFUCache", "put", "put", "get", "put", "get", "get", "put", "get", "get", "get"]
	[[2], [1, 1], [2, 2], [1], [3, 3], [2], [3], [4, 4], [1], [3], [4]]
	Output:
	[null, null, null, 1, null, -1, 3, null, -1, 3, 4]

	Explanation:
	LFUCache lfu = new LFUCache(2);
	lfu.put(1, 1);   // cache={1=1}
	lfu.put(2, 2);   // cache={1=1, 2=2}
	lfu.get(1);      // returns 1, cache={1=1, 2=2} freq: {1:2, 2:1}
	lfu.put(3, 3);   // 2 is LFU (freq=1), evicts 2. cache={1=1, 3=3} freq: {1:2, 3:1}
	lfu.get(2);      // returns -1 (not found)
	lfu.get(3);      // returns 3, cache={1=1, 3=3} freq: {1:2, 3:2}
	lfu.put(4, 4);   // 1 and 3 tie at freq=2; 1 is LRU among them, evicts 1.
	                 // cache={3=3, 4=4} freq: {3:2, 4:1}
	lfu.get(1);      // returns -1 (not found)
	lfu.get(3);      // returns 3
	lfu.get(4);      // returns 4

Constraints:
  - 0 <= capacity <= 10^4
  - 0 <= key <= 10^5
  - 0 <= value <= 10^9
  - At most 2 * 10^5 calls will be made to Get and Put.

TODO:
 1. Add the fields you need on LFUCache below. Hint: the classic O(1) solution keeps two
    layers — a hash map from key -> node (value + frequency), and a hash map from
    frequency -> an ordered list of keys at that frequency (so you can find/evict the LRU
    key within the lowest frequency bucket in O(1)). Track the minimum frequency currently
    in use so eviction doesn't have to search for it.
 2. Implement NewLFUCache with the given capacity (capacity may be 0 — nothing should ever
    be stored, and every Get should return -1).
 3. Implement Get: return the value, or -1 if missing. A successful Get bumps frequency by 1
    and moves the key to the back (most-recently-used position) of its new frequency bucket.
 4. Implement Put:
      - If capacity is 0, do nothing.
      - If the key already exists, update its value and bump frequency by 1 (same rule as Get).
      - If the key is new and the cache is full, evict the least frequently used key,
        breaking ties by least recently used, before inserting.
      - A newly inserted key starts at frequency 1.
*/

type LFUCache struct {
	// TODO: add your fields here
}

func NewLFUCache(capacity int) *LFUCache {
	// TODO: implement
	return &LFUCache{}
}

func (c *LFUCache) Get(key int) int {
	// TODO: implement
	return -1
}

func (c *LFUCache) Put(key int, value int) {
	// TODO: implement
}
