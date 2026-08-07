package designbrowserhistory

import "testing"

func TestBrowserHistory_Example1(t *testing.T) {
	b := NewBrowserHistory("leetcode.com")
	b.Visit("google.com")
	b.Visit("facebook.com")
	b.Visit("youtube.com")
	if got := b.Back(1); got != "facebook.com" {
		t.Fatalf("Back(1) = %q, want facebook.com", got)
	}
	if got := b.Back(1); got != "google.com" {
		t.Fatalf("Back(1) = %q, want google.com", got)
	}
	if got := b.Forward(1); got != "facebook.com" {
		t.Fatalf("Forward(1) = %q, want facebook.com", got)
	}
	b.Visit("linkedin.com")
	if got := b.Forward(2); got != "linkedin.com" {
		t.Fatalf("Forward(2) = %q, want linkedin.com", got)
	}
	if got := b.Back(2); got != "google.com" {
		t.Fatalf("Back(2) = %q, want google.com", got)
	}
	if got := b.Back(7); got != "leetcode.com" {
		t.Fatalf("Back(7) = %q, want leetcode.com", got)
	}
}

func TestBrowserHistory_InitialStateIsHomepage(t *testing.T) {
	b := NewBrowserHistory("home.com")
	if got := b.Back(1); got != "home.com" {
		t.Fatalf("Back(1) with no history = %q, want home.com", got)
	}
	if got := b.Forward(1); got != "home.com" {
		t.Fatalf("Forward(1) with no history = %q, want home.com", got)
	}
}

func TestBrowserHistory_BackClampedAtStart(t *testing.T) {
	b := NewBrowserHistory("a.com")
	b.Visit("b.com")
	b.Visit("c.com")
	if got := b.Back(100); got != "a.com" {
		t.Fatalf("Back(100) = %q, want a.com", got)
	}
}

func TestBrowserHistory_ForwardClampedAtEnd(t *testing.T) {
	b := NewBrowserHistory("a.com")
	b.Visit("b.com")
	b.Visit("c.com")
	if got := b.Forward(100); got != "c.com" {
		t.Fatalf("Forward(100) = %q, want c.com", got)
	}
}

func TestBrowserHistory_VisitClearsForwardHistory(t *testing.T) {
	b := NewBrowserHistory("a.com")
	b.Visit("b.com")
	b.Visit("c.com")
	b.Back(2) // now at a.com, with b.com and c.com in forward history
	b.Visit("d.com")
	if got := b.Forward(1); got != "d.com" {
		t.Fatalf("Forward(1) = %q, want d.com (b.com/c.com forward history should be gone)", got)
	}
	if got := b.Back(1); got != "a.com" {
		t.Fatalf("Back(1) = %q, want a.com", got)
	}
}

func TestBrowserHistory_ZeroStepsReturnsCurrentPage(t *testing.T) {
	// Note: steps >= 1 per constraints, but Back(1)/Forward(1) at a boundary is
	// effectively a no-op, which this exercises indirectly via repeated boundary calls.
	b := NewBrowserHistory("a.com")
	b.Visit("b.com")
	if got := b.Back(1); got != "a.com" {
		t.Fatalf("Back(1) = %q, want a.com", got)
	}
	if got := b.Back(1); got != "a.com" {
		t.Fatalf("Back(1) again at boundary = %q, want a.com", got)
	}
}

func TestBrowserHistory_MultipleVisitsThenFullBackAndForward(t *testing.T) {
	b := NewBrowserHistory("h.com")
	pages := []string{"p1.com", "p2.com", "p3.com", "p4.com", "p5.com"}
	for _, p := range pages {
		b.Visit(p)
	}
	if got := b.Back(5); got != "h.com" {
		t.Fatalf("Back(5) = %q, want h.com", got)
	}
	if got := b.Forward(5); got != "p5.com" {
		t.Fatalf("Forward(5) = %q, want p5.com", got)
	}
}

func TestBrowserHistory_PartialBackThenPartialForward(t *testing.T) {
	b := NewBrowserHistory("h.com")
	b.Visit("a.com")
	b.Visit("b.com")
	b.Visit("c.com")
	b.Visit("d.com")
	if got := b.Back(2); got != "b.com" {
		t.Fatalf("Back(2) = %q, want b.com", got)
	}
	if got := b.Forward(1); got != "c.com" {
		t.Fatalf("Forward(1) = %q, want c.com", got)
	}
	if got := b.Back(3); got != "h.com" {
		t.Fatalf("Back(3) = %q, want h.com (clamped)", got)
	}
}

func TestBrowserHistory_VisitAfterFullBackDropsEverything(t *testing.T) {
	b := NewBrowserHistory("h.com")
	b.Visit("a.com")
	b.Visit("b.com")
	b.Back(2) // back to h.com
	b.Visit("x.com")
	if got := b.Forward(10); got != "x.com" {
		t.Fatalf("Forward(10) = %q, want x.com (a.com/b.com history must be gone)", got)
	}
	if got := b.Back(1); got != "h.com" {
		t.Fatalf("Back(1) = %q, want h.com", got)
	}
}

func TestBrowserHistory_StressManyVisitsAndNavigations(t *testing.T) {
	b := NewBrowserHistory("start.com")
	const n = 200
	for i := 0; i < n; i++ {
		b.Visit(pageName(i))
	}
	if got := b.Back(n); got != "start.com" {
		t.Fatalf("Back(%d) = %q, want start.com", n, got)
	}
	if got := b.Forward(n); got != pageName(n-1) {
		t.Fatalf("Forward(%d) = %q, want %q", n, got, pageName(n-1))
	}
}

func pageName(i int) string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	return string(letters[i%26]) + string(letters[(i/26)%26]) + ".com"
}
