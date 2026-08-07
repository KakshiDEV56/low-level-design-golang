package designlogstoragesystem

/*
LeetCode 635. Design Log Storage System

You are given several logs, where each log contains a unique ID and timestamp. Timestamp
is a string that has the following format: Year:Month:Day:Hour:Minute:Second, for
example, 2017:01:01:23:59:59. All domains are zero-padded decimal numbers.

Implement the LogSystem class:
  - NewLogSystem() initializes the LogSystem object.
  - Put(id int, timestamp string) stores the given log (id, timestamp) in your storage
    system.
  - Retrieve(start, end, granularity string) []int returns the IDs of the logs whose
    timestamps are within the range from start to end, inclusive, after being
    granularized. The granularity is one of "Year", "Month", "Day", "Hour", "Minute", or
    "Second". It means the comparison should only care about the specified granularity
    (and anything coarser than it) of the timestamp; anything finer should be ignored
    (effectively treated as its minimum possible value, e.g. Month=01, Day=01, etc.).

	For example, say the given start = "2017:01:01:23:59:59" and end =
	"2017:01:02:23:59:59" with granularity "Day". This means the range of timestamps
	within the DAY 2017:01:01 to 2017:01:02 are considered, and the Hour, Minute, and
	Second fields are ignored for the comparison (a log at 2017:01:01:00:00:00 would
	still be included, and so would one at 2017:01:02:23:59:59).

Example 1:

	Input:
	["LogSystem", "put", "put", "put", "retrieve", "retrieve"]
	[[], [1, "2017:01:01:23:59:59"], [2, "2017:01:01:22:59:59"],
	 [3, "2016:01:01:00:00:00"],
	 ["2016:01:01:01:01:01", "2017:01:01:23:00:00", "Year"],
	 ["2016:01:01:01:01:01", "2017:01:01:23:00:00", "Hour"]]
	Output:
	[null, null, null, null, [3, 2, 1], [2, 1]]

	Explanation:
	LogSystem ls = new LogSystem();
	ls.put(1, "2017:01:01:23:59:59");
	ls.put(2, "2017:01:01:22:59:59");
	ls.put(3, "2016:01:01:00:00:00");
	// granularity = "Year", start/end are truncated to "2016" / "2017" -> all 3 logs
	// fall within years [2016, 2017], returned in any order (here, insertion order).
	ls.retrieve("2016:01:01:01:01:01", "2017:01:01:23:00:00", "Year"); // return [3,2,1]
	// granularity = "Hour", start truncates to "2016:01:01:01" and end truncates to
	// "2017:01:01:23" -> log 3 (year 2016) is now excluded, only 1 and 2 remain
	// (2017:01:01:23 and 2017:01:01:22 are both within range).
	ls.retrieve("2016:01:01:01:01:01", "2017:01:01:23:00:00", "Hour"); // return [2,1]

Constraints:
  - 1 <= id <= 500
  - 2000 <= Year <= 2017
  - 1 <= Month <= 12
  - 1 <= Day <= 28-31
  - 0 <= Hour <= 23
  - 0 <= Minute, Second <= 59
  - granularity is one of "Year", "Month", "Day", "Hour", "Minute", or "Second".
  - At most 500 calls will be made to Put and Retrieve.

TODO:
 1. Add the fields you need on LogSystem below. Hint: since timestamps are fixed-format
    strings ("YYYY:MM:DD:HH:MM:SS"), you don't need to parse them into numbers at all —
    Go string comparison (`<`, `<=`, `>=`) already works correctly lexicographically for
    this exact zero-padded format. A slice (or map) of (id, timestamp) pairs is enough
    storage.
 2. Write a helper that maps a granularity name ("Year".."Second") to how many characters
    of the "YYYY:MM:DD:HH:MM:SS" string to keep (Year -> 4, Month -> 7, Day -> 10,
    Hour -> 13, Minute -> 16, Second -> 19), so you can truncate any timestamp to a given
    granularity by simple slicing.
 3. Implement NewLogSystem: initialize your storage.
 4. Implement Put: append (id, timestamp) to storage.
 5. Implement Retrieve(start, end, granularity): truncate start and end to the given
    granularity's length. For every stored log, truncate its timestamp the same way and
    check whether it falls in [truncatedStart, truncatedEnd] (plain string comparison
    works here). Collect and return the matching IDs.
*/

type LogSystem struct {
	// TODO: add your fields here
}

func NewLogSystem() *LogSystem {
	// TODO: implement
	return &LogSystem{}
}

func (l *LogSystem) Put(id int, timestamp string) {
	// TODO: implement
}

func (l *LogSystem) Retrieve(start string, end string, granularity string) []int {
	// TODO: implement
	return nil
}
