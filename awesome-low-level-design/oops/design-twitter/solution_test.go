package designtwitter

import (
	"reflect"
	"testing"
)

func TestTwitter_Example1(t *testing.T) {
	tw := NewTwitter()
	tw.PostTweet(1, 5)
	if got := tw.GetNewsFeed(1); !reflect.DeepEqual(got, []int{5}) {
		t.Fatalf("GetNewsFeed(1) = %v, want [5]", got)
	}
	tw.Follow(1, 2)
	tw.PostTweet(2, 6)
	if got := tw.GetNewsFeed(1); !reflect.DeepEqual(got, []int{6, 5}) {
		t.Fatalf("GetNewsFeed(1) = %v, want [6 5]", got)
	}
	tw.Unfollow(1, 2)
	if got := tw.GetNewsFeed(1); !reflect.DeepEqual(got, []int{5}) {
		t.Fatalf("GetNewsFeed(1) = %v, want [5]", got)
	}
}

func TestTwitter_EmptyFeedForNewUser(t *testing.T) {
	tw := NewTwitter()
	if got := tw.GetNewsFeed(1); len(got) != 0 {
		t.Fatalf("GetNewsFeed(1) = %v, want empty", got)
	}
}

func TestTwitter_FeedLimitedToTenMostRecent(t *testing.T) {
	tw := NewTwitter()
	for i := 0; i < 15; i++ {
		tw.PostTweet(1, i)
	}
	got := tw.GetNewsFeed(1)
	if len(got) != 10 {
		t.Fatalf("GetNewsFeed(1) returned %d tweets, want 10", len(got))
	}
	want := []int{14, 13, 12, 11, 10, 9, 8, 7, 6, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetNewsFeed(1) = %v, want %v", got, want)
	}
}

func TestTwitter_FollowingSelfIsRedundant(t *testing.T) {
	tw := NewTwitter()
	tw.PostTweet(1, 100)
	tw.Follow(1, 1) // following yourself must not duplicate your own tweets in the feed
	got := tw.GetNewsFeed(1)
	if !reflect.DeepEqual(got, []int{100}) {
		t.Fatalf("GetNewsFeed(1) = %v, want [100]", got)
	}
}

func TestTwitter_UnfollowWithoutFollowIsNoOp(t *testing.T) {
	tw := NewTwitter()
	tw.PostTweet(2, 1)
	tw.Unfollow(1, 2) // never followed 2; must not panic
	if got := tw.GetNewsFeed(1); len(got) != 0 {
		t.Fatalf("GetNewsFeed(1) = %v, want empty", got)
	}
}

func TestTwitter_MultipleFolloweesInterleaved(t *testing.T) {
	tw := NewTwitter()
	tw.PostTweet(2, 10) // order 0
	tw.PostTweet(3, 20) // order 1
	tw.PostTweet(2, 11) // order 2
	tw.Follow(1, 2)
	tw.Follow(1, 3)
	tw.PostTweet(3, 21) // order 3
	got := tw.GetNewsFeed(1)
	want := []int{21, 11, 20, 10}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetNewsFeed(1) = %v, want %v", got, want)
	}
}

func TestTwitter_FeedDoesNotIncludeUnfollowedUsersTweets(t *testing.T) {
	tw := NewTwitter()
	tw.PostTweet(2, 1)
	got := tw.GetNewsFeed(1)
	if len(got) != 0 {
		t.Fatalf("GetNewsFeed(1) = %v, want empty (user 1 never followed user 2)", got)
	}
}

func TestTwitter_RefollowAfterUnfollowSeesNewTweetsOnly(t *testing.T) {
	tw := NewTwitter()
	tw.Follow(1, 2)
	tw.PostTweet(2, 1)
	tw.Unfollow(1, 2)
	tw.PostTweet(2, 2)
	tw.Follow(1, 2)
	// Feeds aren't "subscriptions with history gaps" — once followed again, ALL of 2's
	// tweets (both old and new) are visible again, since GetNewsFeed just looks at the
	// current follow graph, not follow/unfollow history.
	got := tw.GetNewsFeed(1)
	want := []int{2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetNewsFeed(1) = %v, want %v", got, want)
	}
}

func TestTwitter_IndependentUsersHaveIndependentFeeds(t *testing.T) {
	tw := NewTwitter()
	tw.PostTweet(1, 1)
	tw.PostTweet(2, 2)
	tw.Follow(1, 2)
	// User 2 does not follow user 1, so their feed should only have their own tweet.
	if got := tw.GetNewsFeed(2); !reflect.DeepEqual(got, []int{2}) {
		t.Fatalf("GetNewsFeed(2) = %v, want [2]", got)
	}
	if got := tw.GetNewsFeed(1); !reflect.DeepEqual(got, []int{2, 1}) {
		t.Fatalf("GetNewsFeed(1) = %v, want [2 1]", got)
	}
}

func TestTwitter_ManyUsersManyTweetsOrdering(t *testing.T) {
	tw := NewTwitter()
	// User 1 follows users 2..5.
	for u := 2; u <= 5; u++ {
		tw.Follow(1, u)
	}
	// Post tweets round-robin so global order is deterministic: tweet IDs 0..19 in
	// strictly increasing post order across users 2,3,4,5,2,3,4,5,...
	tweetID := 0
	for round := 0; round < 5; round++ {
		for u := 2; u <= 5; u++ {
			tw.PostTweet(u, tweetID)
			tweetID++
		}
	}
	got := tw.GetNewsFeed(1)
	want := []int{19, 18, 17, 16, 15, 14, 13, 12, 11, 10}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetNewsFeed(1) = %v, want %v", got, want)
	}
}
