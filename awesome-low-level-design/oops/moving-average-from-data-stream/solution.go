package movingaveragefromdatastream

/*
LeetCode 346. Moving Average from Data Stream

Given a stream of integers and a window size, calculate the moving average of all
integers in the sliding window.

Implement the MovingAverage class:
  - NewMovingAverage(size int) initializes the object with the size of the window size.
  - Next(val int) float64 returns the moving average of the last size values of the
    stream.

Example 1:

	Input:
	["MovingAverage", "next", "next", "next", "next"]
	[[3], [1], [10], [3], [5]]
	Output:
	[null, 1.0, 5.5, 4.66667, 6.0]

	Explanation:
	MovingAverage ma = new MovingAverage(3);
	ma.next(1);  // return 1.0 = 1 / 1
	ma.next(10); // return 5.5 = (1 + 10) / 2
	ma.next(3);  // return 4.66667 = (1 + 10 + 3) / 3
	ma.next(5);  // return 6.0 = (10 + 3 + 5) / 3 (window has slid, 1 fell out)

Constraints:
  - 1 <= size <= 1000
  - -10^5 <= val <= 10^5
  - At most 10^4 calls will be made to Next.

TODO:
 1. Add the fields you need on MovingAverage below. Hint: a fixed-capacity queue (a slice
    used as a FIFO, or a true ring buffer sized to `size`) holding the values currently
    in the window, plus a running sum so you don't have to re-sum the window on every
    call.
 2. Implement NewMovingAverage(size): remember the window size and initialize your
    queue/sum.
 3. Implement Next(val): add val to the window and to the running sum. If the window now
    exceeds `size` elements, remove the oldest one from both the window and the running
    sum. Return runningSum / float64(currentWindowLength) — note the divisor is the
    CURRENT number of values seen so far, capped at `size`, not always `size` (see the
    first two calls in the example, which divide by 1 and 2 respectively before the
    window fills up).
*/

type MovingAverage struct {
	// TODO: add your fields here
}

func NewMovingAverage(size int) *MovingAverage {
	// TODO: implement
	return &MovingAverage{}
}

func (m *MovingAverage) Next(val int) float64 {
	// TODO: implement
	return 0
}
