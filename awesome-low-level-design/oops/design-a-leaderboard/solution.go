package designleaderboard

import "container/heap"

/*
LeetCode 1244. Design A Leaderboard

Design a Leaderboard class, which has three functions:
  - AddScore(playerId, score int) updates the leaderboard by adding score to the given
    player's score. If there is no player with such id in the leaderboard, add him to the
    leaderboard with the given score.
  - Top(K int) int returns the score sum of the top K players.
  - Reset(playerId int) resets the score of the player with the given id to 0 (in other
    words, erase it from the leaderboard). It is guaranteed that the player was added to
    the leaderboard before calling this function.

Initially, the leaderboard is empty.

Example 1:

	Input:
	["Leaderboard", "addScore", "addScore", "addScore", "addScore", "addScore", "top",
	 "reset", "reset", "addScore", "top"]
	[[], [1, 73], [2, 56], [3, 39], [4, 51], [5, 4], [1], [1], [2], [2, 51], [3]]
	Output:
	[null, null, null, null, null, null, 73, null, null, null, 141]

	Explanation:
	Leaderboard lb = new Leaderboard();
	lb.addScore(1, 73);  // scores: {1:73}
	lb.addScore(2, 56);  // scores: {1:73, 2:56}
	lb.addScore(3, 39);  // scores: {1:73, 2:56, 3:39}
	lb.addScore(4, 51);  // scores: {1:73, 2:56, 3:39, 4:51}
	lb.addScore(5, 4);   // scores: {1:73, 2:56, 3:39, 4:51, 5:4}
	lb.top(1);           // returns 73
	lb.reset(1);         // scores: {2:56, 3:39, 4:51, 5:4}
	lb.reset(2);         // scores: {3:39, 4:51, 5:4}
	lb.addScore(2, 51);  // scores: {2:51, 3:39, 4:51, 5:4}
	lb.top(3);           // returns 51 + 51 + 39 = 141

Constraints:
  - 1 <= playerId, K <= 10^4
  - It's guaranteed that K is less than or equal to the current number of players.
  - 1 <= score <= 100
  - There will be at most 1000 function calls in total, including AddScore, Top, and
    Reset.
  - Reset will only be called on a player already existing in the leaderboard.

TODO:
 1. Add the fields you need on Leaderboard below. Hint: a map from playerId -> current
    score is enough state — you don't need to keep it sorted at all times.
 2. Implement NewLeaderboard: initialize your map.
 3. Implement AddScore: add `score` to the player's existing score (or start them at
    `score` if new).
 4. Implement Reset: remove the player from the map entirely (equivalent to erasing them,
    per the description) — a subsequent AddScore for that playerId should start fresh.
 5. Implement Top(K): collect all current scores, sort descending (see `sort.Slice`),
    sum the first K, and return the sum. Given the constraints (at most 1000 calls
    total), this straightforward "sort on every Top call" approach is plenty fast — no
    need for a fancier order-statistics structure unless you want the extra challenge.
*/

type MinHeap []int

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

/*How will i sort the overall */
type Leaderboard struct {
	player map[int]int
}

func NewLeaderboard() *Leaderboard {
	return &Leaderboard{
		player: make(map[int]int),
	}
}

func (l *Leaderboard) AddScore(playerId int, score int) {
	if _, exist := l.player[playerId]; !exist {
		l.player[playerId] = score
	} else {
		l.player[playerId] += score
	}

}

func (l *Leaderboard) Top(K int) int {
	total_sum := 0
	// using minHeap i can easily solve this
	MinHeap := &MinHeap{}
	heap.Init(MinHeap)
	for _, score := range l.player {
		if MinHeap.Len() < K {
			heap.Push(MinHeap, score)
		} else if score > (*MinHeap)[0] {
			heap.Pop(MinHeap)
			heap.Push(MinHeap, score)
		}
	}
	for MinHeap.Len() > 0 {
		total_sum += heap.Pop(MinHeap).(int)
	}
	return total_sum
}

func (l *Leaderboard) Reset(playerId int) {
	if _, exist := l.player[playerId]; exist {
		l.player[playerId] = 0
		delete(l.player, playerId)
	}

}
