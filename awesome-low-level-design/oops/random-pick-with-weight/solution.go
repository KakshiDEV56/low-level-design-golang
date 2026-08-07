package randompickwithweight

/*
LeetCode 528. Random Pick with Weight

You are given an array of positive integers w where w[i] describes the weight of the ith
index. You need to implement the function PickIndex, which randomly picks an index in
the range [0, w.length - 1] (inclusive) and returns it. The probability of picking an
index i is w[i] / sum(w).

Implement the Solution class:
  - NewSolution(w []int) initializes the object with the array w.
  - PickIndex() int picks an index in the range [0, w.length - 1] and returns it. The
    probability of returning index i equals w[i] / sum(w).

Example 1:

	Input:
	["Solution", "pickIndex"]
	[[[1]], []]
	Output:
	[null, 0]

	Explanation:
	Solution s = new Solution([1]);
	s.pickIndex(); // there is only one index (0), so pickIndex always returns 0.

Example 2:

	Input:
	["Solution", "pickIndex", "pickIndex", "pickIndex", "pickIndex", "pickIndex"]
	[[[1, 3]], [], [], [], [], []]
	Output:
	[null, 1, 1, 1, 1, 0]

	Explanation:
	Solution s = new Solution([1, 3]);
	s.pickIndex(); // returns 0 with probability 1/4 = 0.25 (i.e., 25%), and returns 1
	              // with probability 3/4 = 0.75 (i.e., 75%). Any sequence of return
	              // values that follows this probability will be accepted.

Constraints:
  - 1 <= w.length <= 10^4
  - 1 <= w[i] <= 10^5
  - At most 10^4 calls will be made to PickIndex.

TODO:
 1. Add the fields you need on Solution below. Hint: precompute the prefix-sum array of
    weights (prefix[i] = w[0]+w[1]+...+w[i]) once, in NewSolution. This turns "pick index
    i with probability proportional to w[i]" into "pick a uniformly random integer in
    [1, prefix[len-1]], then binary-search the prefix-sum array for the first index whose
    prefix sum is >= that random integer".
 2. Implement NewSolution: build the prefix-sum slice from w.
 3. Implement PickIndex: generate a uniformly random integer target in
    [1, totalWeight] (see math/rand — e.g. rand.Intn(total) + 1), then binary-search
    (sort.Search) the prefix-sum slice for the smallest index whose prefix sum is >=
    target, and return that index.
*/

type Solution struct {
	// TODO: add your fields here
}

func NewSolution(w []int) *Solution {
	// TODO: implement
	return &Solution{}
}

func (s *Solution) PickIndex() int {
	// TODO: implement
	return 0
}
