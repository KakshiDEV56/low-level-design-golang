package alloonedatastructure

/*
LeetCode 432. All O`one Data Structure

Design a data structure to store the strings' count with the ability to return the
strings with minimum and maximum counts.

Implement the AllOne class:
  - NewAllOne() initializes the object of the data structure.
  - Inc(key string) increments the count of the string key by 1. If key does not exist
    in the data structure, insert it with count 1.
  - Dec(key string) decrements the count of the string key by 1. If the count of key is 0
    after the decrement, remove it from the data structure. It is guaranteed that key
    exists in the data structure before the decrement.
  - GetMaxKey() string returns one of the keys with the maximal count. If no element
    exists, return an empty string "".
  - GetMinKey() string returns one of the keys with the minimum count. If no element
    exists, return an empty string "".

Note: each function must run in O(1) average time complexity.

Example 1:

	Input:
	["AllOne", "inc", "inc", "getMaxKey", "getMinKey", "inc", "getMaxKey", "getMinKey"]
	[[], ["hello"], ["hello"], [], [], ["leet"], [], []]
	Output:
	[null, null, null, "hello", "hello", null, "hello", "leet"]

	Explanation:
	AllOne ao = new AllOne();
	ao.inc("hello");         // counts: {hello: 1}
	ao.inc("hello");         // counts: {hello: 2}
	ao.getMaxKey();          // return "hello"
	ao.getMinKey();          // return "hello"
	ao.inc("leet");          // counts: {hello: 2, leet: 1}
	ao.getMaxKey();          // return "hello"
	ao.getMinKey();          // return "leet"

Constraints:
  - 1 <= key.length <= 10
  - key consists of lowercase English letters.
  - It is guaranteed that for each call to Dec, key is existing in the data structure.
  - At most 5*10^4 calls will be made to Inc, Dec, GetMaxKey, and GetMinKey.

TODO:
 1. Add the fields you need on AllOne below. The true O(1)-per-op solution is the classic
    "doubly linked list of count-buckets" trick:
      - A doubly linked list where each node represents a distinct count value (e.g. 1,
        2, 3, ...) IN SORTED ORDER, and holds the set of keys currently at that count.
      - A map from key -> the node it currently lives in, so Inc/Dec can jump straight to
        the right bucket instead of searching.
      - GetMaxKey/GetMinKey just look at the tail/head bucket's key set in O(1).
      - Inc(key): find key's current bucket (or none, meaning count 0->1), move it to the
        NEXT bucket (count+1), creating that bucket right after the current one if it
        doesn't exist yet; remove the key from the old bucket, deleting the old bucket if
        it's now empty.
      - Dec(key) is the mirror image: move to the PREVIOUS bucket (count-1), creating it
        before the current one if needed (unless count-1 == 0, in which case the key is
        just removed entirely); delete the old bucket if now empty.
    This is a meaningfully harder structure than most problems in this set — it's a
    great one to come back to after building confidence on the others. A simpler
    (non-O(1)) version using a plain map[string]int plus a linear scan for
    GetMaxKey/GetMinKey is a fine first pass to get the behavior correct before
    optimizing.
 2. Implement NewAllOne: initialize your map(s)/list.
 3. Implement Inc, Dec, GetMaxKey, GetMinKey per the description above.
*/

type AllOne struct {
	// TODO: add your fields here
}

func NewAllOne() *AllOne {
	// TODO: implement
	return &AllOne{}
}

func (a *AllOne) Inc(key string) {
	// TODO: implement
}

func (a *AllOne) Dec(key string) {
	// TODO: implement
}

func (a *AllOne) GetMaxKey() string {
	// TODO: implement
	return ""
}

func (a *AllOne) GetMinKey() string {
	// TODO: implement
	return ""
}
