package timebasedkeyvaluestore

/*
LeetCode 981. Time Based Key-Value Store

Design a time-based key-value data structure that can store multiple values for the same
key at different time stamps and retrieve the key's value at a certain timestamp.

Implement the TimeMap class:
  - NewTimeMap() initializes the object of the data structure.
  - Set(key, value string, timestamp int) stores the key with the value at the given time
    timestamp.
  - Get(key string, timestamp int) string returns a value such that Set was called
    previously, with timestamp_prev <= timestamp. If there are multiple such values, it
    returns the value associated with the largest timestamp_prev. If there are no values,
    it returns "".

Example 1:

	Input:
	["TimeMap", "set", "get", "get", "set", "get", "get"]
	[[], ["foo", "bar", 1], ["foo", 1], ["foo", 3], ["foo", "bar2", 4], ["foo", 4], ["foo", 5]]
	Output:
	[null, null, "bar", "bar", null, "bar2", "bar2"]

	Explanation:
	TimeMap tm = new TimeMap();
	tm.set("foo", "bar", 1);   // store the key "foo" and value "bar" along with timestamp = 1
	tm.get("foo", 1);          // return "bar"
	tm.get("foo", 3);          // return "bar", since there is no value at timestamp 3 and
	                           // timestamp 2, but there is a value at timestamp 1
	tm.set("foo", "bar2", 4);  // store the key "foo" and value "bar2" along with timestamp = 4
	tm.get("foo", 4);          // return "bar2"
	tm.get("foo", 5);          // return "bar2"

Constraints:
  - 1 <= key.length, value.length <= 100
  - key and value consist of lowercase English letters and digits.
  - 1 <= timestamp <= 10^7
  - All the timestamps timestamp of Set are strictly increasing for the same key.
  - At most 2 * 10^5 calls will be made to Set and Get.

TODO:
 1. Add the fields you need on TimeMap below. Hint: a map from key -> a slice of
    (timestamp, value) pairs. Since Set is guaranteed to be called with strictly
    increasing timestamps for a given key, each key's slice is naturally kept sorted by
    timestamp just by appending — no need to sort it yourself.
 2. Implement NewTimeMap: initialize your map.
 3. Implement Set: append (timestamp, value) to the slice for key (creating it if needed).
 4. Implement Get: look up key's slice (return "" if the key was never set). Binary-search
    the slice for the largest timestamp that is <= the queried timestamp (see the `sort`
    package's `sort.Search`), and return its value — or "" if every stored timestamp is
    greater than the queried one.
*/

type TimeMap struct {
	// TODO: add your fields here
}

func NewTimeMap() *TimeMap {
	// TODO: implement
	return &TimeMap{}
}

func (tm *TimeMap) Set(key string, value string, timestamp int) {
	// TODO: implement
}

func (tm *TimeMap) Get(key string, timestamp int) string {
	// TODO: implement
	return ""
}
