package designhitcounter

/*
LeetCode 362. Design Hit Counter

Design a hit counter which counts the number of hits received in the past 5 minutes
(i.e., the past 300 seconds).

Your system should accept a timestamp parameter (in seconds granularity), and you may
assume that calls are being made to the system in chronological order (i.e., timestamp
is monotonically increasing). Several hits may arrive roughly at the same time.

Implement the HitCounter class:
  - NewHitCounter() initializes the object of the hit counter system.
  - Hit(timestamp int) records a hit that happened at timestamp (in seconds). Several
    hits may happen at the same timestamp.
  - GetHits(timestamp int) int returns the number of hits in the past 5 minutes from
    timestamp (i.e., the past 300 seconds, in the inclusive range
    [timestamp-300+1, timestamp]).

Example 1:

	Input:
	["HitCounter", "hit", "hit", "hit", "getHits", "hit", "getHits", "getHits"]
	[[], [1], [2], [3], [4], [300], [300], [301]]
	Output:
	[null, null, null, null, 3, null, 4, 3]

	Explanation:
	HitCounter hc = new HitCounter();
	hc.hit(1);          // hit at timestamp 1
	hc.hit(2);          // hit at timestamp 2
	hc.hit(3);          // hit at timestamp 3
	hc.getHits(4);      // get hits at timestamp 4; window is [-296,4], all 3 hits count -> 3
	hc.hit(300);        // hit at timestamp 300
	hc.getHits(300);    // get hits at timestamp 300; window is [1,300], all 4 hits count -> 4
	hc.getHits(301);    // get hits at timestamp 301; window is [2,301], hit at ts=1 has
	                    // expired -> 3

Constraints:
  - 1 <= timestamp <= 2 * 10^9
  - All the calls are being made to the system in chronological order (i.e., timestamp is
    monotonically increasing).
  - At most 300 calls will be made to Hit and GetHits.

TODO:
 1. Add the fields you need on HitCounter below. Since calls arrive in chronological
    order, the simplest correct approach is a queue (a slice used as a FIFO) of hit
    timestamps: append on Hit, and on GetHits, pop expired timestamps off the front
    before reporting the remaining count. (A fancier O(1)-space bucket-per-second-slot
    ring buffer also works and is worth trying once the queue version passes, since it
    caps memory at 300 slots regardless of hit volume.)
 2. Implement NewHitCounter: initialize your queue.
 3. Implement Hit(timestamp): append timestamp to the queue (duplicates are fine — the
    problem allows several hits at the same timestamp).
 4. Implement GetHits(timestamp): remove all entries from the front of the queue that are
    <= timestamp-300 (i.e. older than the 300-second window ending at timestamp), then
    return the queue's remaining length.
*/

type HitCounter struct {
	// TODO: add your fields here
}

func NewHitCounter() *HitCounter {
	// TODO: implement
	return &HitCounter{}
}

func (h *HitCounter) Hit(timestamp int) {
	// TODO: implement
}

func (h *HitCounter) GetHits(timestamp int) int {
	// TODO: implement
	return 0
}
