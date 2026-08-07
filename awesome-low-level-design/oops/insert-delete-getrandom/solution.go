package insertdeletegetrandom

/*
LeetCode 380. Insert Delete GetRandom O(1)

Implement the RandomizedSet class:
  - NewRandomizedSet() initializes the RandomizedSet object.
  - Insert(val int) bool inserts an item val into the set if not present. Returns true if
    the item was not present, false otherwise.
  - Remove(val int) bool removes an item val from the set if present. Returns true if the
    item was present, false otherwise.
  - GetRandom() int returns a random element from the current set of elements (it's
    guaranteed that at least one element exists when this method is called). Each element
    must have the same probability of being returned.

You must implement the functions of the class such that each function works in average
O(1) time complexity.

Example 1:

	Input:
	["RandomizedSet", "insert", "remove", "insert", "getRandom", "remove", "insert", "getRandom"]
	[[], [1], [2], [2], [], [1], [2], []]
	Output:
	[null, true, false, true, 2, true, false, 2]

	Explanation:
	RandomizedSet rs = new RandomizedSet();
	rs.insert(1);    // true, set=[1]
	rs.remove(2);    // false, 2 not in set, set=[1]
	rs.insert(2);    // true, set=[1,2]
	rs.getRandom();  // returns 1 or 2 randomly
	rs.remove(1);    // true, set=[2]
	rs.insert(2);    // false, 2 already in set, set=[2]
	rs.getRandom();  // always returns 2

Constraints:
  - -2^31 <= val <= 2^31 - 1
  - At most 2 * 10^5 calls will be made to Insert, Remove, and GetRandom.
  - There will be at least one element in the data structure when GetRandom is called.

TODO:
 1. Add the fields you need on RandomizedSet below. Hint: O(1) Insert/Remove/GetRandom is
    the classic "array + hash map" trick:
      - a slice holding all current values (so GetRandom can pick a uniformly random
        index in O(1)), and
      - a map from value -> its current index in that slice (so Insert/Remove can find
        and update in O(1)).
    The tricky part is Remove in O(1): instead of shifting the slice, SWAP the element
    to remove with the LAST element in the slice, update the swapped element's index in
    the map, then shrink the slice by one (removing the true last element).
 2. Implement NewRandomizedSet: initialize your slice and map.
 3. Implement Insert: if val is already in the map, return false. Otherwise append to the
    slice, record its index in the map, return true.
 4. Implement Remove: if val isn't in the map, return false. Otherwise do the
    swap-with-last-element trick described above, delete val from the map, return true.
 5. Implement GetRandom: pick a uniformly random index in [0, len(slice)) — see the
    math/rand (or math/rand/v2) package — and return that element.
*/

type RandomizedSet struct {
	// TODO: add your fields here
}

func NewRandomizedSet() *RandomizedSet {
	// TODO: implement
	return &RandomizedSet{}
}

func (s *RandomizedSet) Insert(val int) bool {
	// TODO: implement
	return false
}

func (s *RandomizedSet) Remove(val int) bool {
	// TODO: implement
	return false
}

func (s *RandomizedSet) GetRandom() int {
	// TODO: implement
	return 0
}
