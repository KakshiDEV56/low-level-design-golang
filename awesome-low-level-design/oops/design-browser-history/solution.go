package designbrowserhistory

/*
LeetCode 1472. Design Browser History

You have a browser of one tab where you start on the homepage and you can visit another
url, get back in the history number of steps or move forward in the history number of
steps.

Implement the BrowserHistory class:
  - NewBrowserHistory(homepage string) initializes the object with the homepage of the
    browser.
  - Visit(url string) visits url from the current page. It clears up all the forward
    history.
  - Back(steps int) string moves steps back in history. If you can only return x steps in
    the history and steps > x, you will return only x steps. Returns the current url
    after moving back in history at most steps.
  - Forward(steps int) string moves steps forward in history. If you can only forward x
    steps in the history and steps > x, you will forward only x steps. Returns the
    current url after forwarding in history at most steps.

Example 1:

	Input:
	["BrowserHistory", "visit", "visit", "visit", "back", "back", "forward", "visit",
	 "forward", "back", "back"]
	[["leetcode.com"], ["google.com"], ["facebook.com"], ["youtube.com"], [1], [1], [1],
	 ["linkedin.com"], [2], [2], [7]]
	Output:
	[null, null, null, null, "facebook.com", "google.com", "facebook.com", null,
	 "linkedin.com", "google.com", "leetcode.com"]

	Explanation:
	BrowserHistory bh = new BrowserHistory("leetcode.com");
	bh.visit("google.com");       // history: leetcode.com | google.com
	bh.visit("facebook.com");     // history: leetcode.com | google.com | facebook.com
	bh.visit("youtube.com");      // history: leetcode.com | google.com | facebook.com | youtube.com
	bh.back(1);                   // returns "facebook.com"
	bh.back(1);                   // returns "google.com"
	bh.forward(1);                // returns "facebook.com"
	bh.visit("linkedin.com");     // clears forward history, history: ...| google.com |
	                              // facebook.com | linkedin.com
	bh.forward(2);                // returns "linkedin.com" (can't go forward, already at
	                              // the last page)
	bh.back(2);                   // returns "google.com"
	bh.back(7);                   // returns "leetcode.com"

Constraints:
  - 1 <= homepage.length, url.length <= 20
  - 1 <= steps <= 100
  - homepage and url consist of '.' or lowercase English letters and digits.
  - At most 5000 calls will be made to Visit, Back, and Forward.

TODO:
 1. Add the fields you need on BrowserHistory below. Hint: a slice of visited URLs acting
    as a simple history list, plus a `current` index into it. Visiting truncates
    everything after `current` before appending.
 2. Implement NewBrowserHistory: seed the history with a single entry, the homepage, and
    point `current` at it.
 3. Implement Visit: drop any entries after `current` (the "forward" history), append the
    new url, and move `current` to it.
 4. Implement Back(steps): move `current` back by steps, clamped to index 0 (can't go
    before the first page). Return the url at the new `current`.
 5. Implement Forward(steps): move `current` forward by steps, clamped to the last index
    (can't go past the most recently visited page). Return the url at the new `current`.
*/

type BrowserHistory struct {
	visitedURL       []string
	currentURL_index int
}

func NewBrowserHistory(homepage string) *BrowserHistory {

	return &BrowserHistory{
		visitedURL:       []string{homepage},
		currentURL_index: 0,
	}
}

func (b *BrowserHistory) Visit(url string) {
	// simply i am going to append the URL in the in the VisitedURL by creating the slice of the visitedURL from the 0 to the current index
	b.visitedURL = b.visitedURL[:b.currentURL_index+1]
	b.visitedURL = append(b.visitedURL, url)
	b.currentURL_index++
}

func (b *BrowserHistory) Back(steps int) string {
	if b.currentURL_index >= steps {
		// updating the current URL index also
		b.currentURL_index = b.currentURL_index - steps
		return b.visitedURL[b.currentURL_index]
	} else if b.currentURL_index < steps {
		b.currentURL_index = 0
		return b.visitedURL[b.currentURL_index]
	}
	return ""
}

func (b *BrowserHistory) Forward(steps int) string {
	if b.currentURL_index+steps < len(b.visitedURL) {
		b.currentURL_index = b.currentURL_index + steps
		return b.visitedURL[b.currentURL_index]
	} else if b.currentURL_index+steps >= len(b.visitedURL) {
		b.currentURL_index = len(b.visitedURL) - 1
		return b.visitedURL[b.currentURL_index]
	}
	return ""
}
