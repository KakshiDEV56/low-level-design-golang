package designtwitter

/*
LeetCode 355. Design Twitter

Design a simplified version of Twitter where users can post tweets, follow/unfollow
another user, and see the 10 most recent tweets in the user's news feed.

Implement the Twitter class:
  - NewTwitter() initializes your twitter object.
  - PostTweet(userID, tweetID int) composes a new tweet with ID tweetID by the user
    userID. Each call to this function will be made with a unique tweetID.
  - GetNewsFeed(userID int) []int retrieves the 10 most recent tweet IDs in the user's
    news feed. Each item in the news feed must be posted by users who the user followed
    or by the user themself. Tweets must be ordered from most recent to least recent.
  - Follow(followerID, followeeID int) the user with ID followerID started following the
    user with ID followeeID.
  - Unfollow(followerID, followeeID int) the user with ID followerID started unfollowing
    the user with ID followeeID.

Example 1:

	Input:
	["Twitter", "postTweet", "getNewsFeed", "follow", "postTweet", "getNewsFeed",
	 "unfollow", "getNewsFeed"]
	[[], [1, 5], [1], [1, 2], [2, 6], [1], [1, 2], [1]]
	Output:
	[null, null, [5], null, null, [6, 5], null, [5]]

	Explanation:
	Twitter tw = new Twitter();
	tw.postTweet(1, 5);      // user 1 posts a new tweet (id = 5)
	tw.getNewsFeed(1);       // user 1's news feed should return a list with 1 tweet id -> [5]
	tw.follow(1, 2);         // user 1 follows user 2
	tw.postTweet(2, 6);      // user 2 posts a new tweet (id = 6)
	tw.getNewsFeed(1);       // user 1's news feed should return a list with 2 tweet ids
	                         // -> [6, 5], tweet id 6 should precede tweet id 5 because it
	                         //    was posted after tweet id 5
	tw.unfollow(1, 2);       // user 1 unfollows user 2
	tw.getNewsFeed(1);       // user 1's news feed should return a list with 1 tweet id -> [5]

Constraints:
  - 1 <= userID, followerID, followeeID <= 500
  - 0 <= tweetID <= 10^4
  - All the tweets have unique IDs.
  - At most 3 * 10^4 calls will be made to PostTweet, GetNewsFeed, Follow, and Unfollow.

TODO:
 1. Add the fields you need on Twitter below. Hint: you'll want
      - a map from userID -> slice (or linked list) of that user's own tweets, each
        tagged with a global, monotonically increasing "post order" counter (so you can
        compare recency across different users without relying on wall-clock time), and
      - a map from userID -> set of followeeIDs they follow.
 2. Implement NewTwitter: initialize your maps and a counter starting at 0.
 3. Implement PostTweet: record (tweetID, counter) under userID's own tweet list, then
    increment the counter.
 4. Implement Follow/Unfollow: update followerID's followee set. A user does not need to
    explicitly follow themself — GetNewsFeed should always include the user's own tweets
    regardless of the follow set. Following yourself, or an already-followed/not-followed
    user, should be handled gracefully (no panics, no duplicate effects).
 5. Implement GetNewsFeed: gather the userID's own tweets plus every followee's tweets,
    then return the top 10 by post order, most recent first. A classic O(k log m)
    approach (k = 10, m = number of relevant users) is a k-way merge using a max-heap
    (container/heap) over "next tweet" pointers per user — but a simpler approach
    (collect everything, sort by counter descending, take the first 10) is perfectly fine
    to start with given the constraints.
*/

type Twitter struct {
	// TODO: add your fields here
}

func NewTwitter() *Twitter {
	// TODO: implement
	return &Twitter{}
}

func (tw *Twitter) PostTweet(userID int, tweetID int) {
	// TODO: implement
}

func (tw *Twitter) GetNewsFeed(userID int) []int {
	// TODO: implement
	return nil
}

func (tw *Twitter) Follow(followerID int, followeeID int) {
	// TODO: implement
}

func (tw *Twitter) Unfollow(followerID int, followeeID int) {
	// TODO: implement
}
