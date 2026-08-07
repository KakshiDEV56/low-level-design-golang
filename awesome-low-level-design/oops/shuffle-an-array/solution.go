package shufflearray

/*
LeetCode 384. Shuffle an Array

Given an integer array nums, design an algorithm to randomly shuffle the array. All
permutations of the array should be equally likely as a result of the shuffling.

Implement the Solution class:
  - NewSolution(nums []int) initializes the object with the integer array nums.
  - Reset() []int resets the array to its original configuration (the one passed to the
    constructor) and returns it.
  - Shuffle() []int returns a random shuffling of the array, with every permutation
    equally likely.

Example 1:

	Input:
	["Solution", "shuffle", "reset", "shuffle"]
	[[[1, 2, 3]], [], [], []]
	Output:
	[null, [3, 1, 2], [1, 2, 3], [1, 3, 2]]

	Explanation:
	Solution s = new Solution([1, 2, 3]);
	s.shuffle();  // shuffle the array [1,2,3] and return its result. Any permutation of
	              // array [1,2,3] must be equally likely to be returned. Example:
	              // return [3, 1, 2]
	s.reset();    // resets the array back to its original configuration [1,2,3].
	              // return [1, 2, 3]
	s.shuffle();  // returns the random shuffling of array [1,2,3]. Example: return
	              // [1, 3, 2]

Constraints:
  - 1 <= nums.length <= 200
  - -10^6 <= nums[i] <= 10^6
  - All the elements of nums are unique.
  - At most 5 * 10^4 calls in total will be made to Reset and Shuffle.

TODO:
 1. Add the fields you need on Solution below. Hint: you need to remember the ORIGINAL
    array (for Reset) separately from the CURRENT working array (which Shuffle mutates
    in place). Keep a copy of the original slice, and a separate slice for "current
    state" that starts as a copy of it.
 2. Implement NewSolution: store both the original and a working copy.
 3. Implement Reset: copy the original back into the working array (or just return a
    fresh copy of the original directly) and return it.
 4. Implement Shuffle: apply the Fisher-Yates shuffle to the working array — for i from
    len-1 down to 1, pick a uniformly random j in [0, i] (see math/rand) and swap
    working[i] and working[j]. This is the classic algorithm that guarantees every
    permutation is equally likely; a naive "swap with any random index" approach does
    NOT give a uniform distribution, so it's worth understanding why Fisher-Yates does.
    Return the shuffled working array.
*/

type Solution struct {
	// TODO: add your fields here
}

func NewSolution(nums []int) *Solution {
	// TODO: implement
	return &Solution{}
}

func (s *Solution) Reset() []int {
	// TODO: implement
	return nil
}

func (s *Solution) Shuffle() []int {
	// TODO: implement
	return nil
}
