package designhitcounter

import "testing"

func TestHitCounter_Example1(t *testing.T) {
	h := NewHitCounter()
	h.Hit(1)
	h.Hit(2)
	h.Hit(3)
	if got := h.GetHits(4); got != 3 {
		t.Fatalf("GetHits(4) = %d, want 3", got)
	}
	h.Hit(300)
	if got := h.GetHits(300); got != 4 {
		t.Fatalf("GetHits(300) = %d, want 4", got)
	}
	if got := h.GetHits(301); got != 3 {
		t.Fatalf("GetHits(301) = %d, want 3", got)
	}
}

func TestHitCounter_NoHitsYet(t *testing.T) {
	h := NewHitCounter()
	if got := h.GetHits(100); got != 0 {
		t.Fatalf("GetHits(100) on fresh counter = %d, want 0", got)
	}
}

func TestHitCounter_ExactWindowBoundary(t *testing.T) {
	h := NewHitCounter()
	h.Hit(1)
	// Window at timestamp 300 is [1, 300] inclusive (300 seconds wide) -> hit at 1 counts.
	if got := h.GetHits(300); got != 1 {
		t.Fatalf("GetHits(300) = %d, want 1 (hit at ts=1 is still within the 300s window)", got)
	}
	// Window at timestamp 301 is [2, 301] -> hit at 1 has just expired.
	if got := h.GetHits(301); got != 0 {
		t.Fatalf("GetHits(301) = %d, want 0 (hit at ts=1 should have expired)", got)
	}
}

func TestHitCounter_MultipleHitsAtSameTimestamp(t *testing.T) {
	h := NewHitCounter()
	h.Hit(5)
	h.Hit(5)
	h.Hit(5)
	if got := h.GetHits(5); got != 3 {
		t.Fatalf("GetHits(5) = %d, want 3", got)
	}
}

func TestHitCounter_HitsExpireGradually(t *testing.T) {
	h := NewHitCounter()
	for ts := 1; ts <= 5; ts++ {
		h.Hit(ts)
	}
	// A hit at timestamp k remains visible to GetHits(T) exactly while T <= k+299, i.e.
	// it expires the instant T reaches k+300.
	if got := h.GetHits(5); got != 5 {
		t.Fatalf("GetHits(5) = %d, want 5", got)
	}
	if got := h.GetHits(300); got != 5 {
		t.Fatalf("GetHits(300) = %d, want 5 (window [1,300] still includes ts=1..5)", got)
	}
	if got := h.GetHits(301); got != 4 {
		t.Fatalf("GetHits(301) = %d, want 4 (ts=1 just expired)", got)
	}
	if got := h.GetHits(304); got != 1 {
		t.Fatalf("GetHits(304) = %d, want 1 (only ts=5 remains)", got)
	}
	if got := h.GetHits(305); got != 0 {
		t.Fatalf("GetHits(305) = %d, want 0 (all expired)", got)
	}
}

func TestHitCounter_InterleavedHitsAndQueries(t *testing.T) {
	h := NewHitCounter()
	h.Hit(10)
	if got := h.GetHits(10); got != 1 {
		t.Fatalf("GetHits(10) = %d, want 1", got)
	}
	h.Hit(20)
	if got := h.GetHits(20); got != 2 {
		t.Fatalf("GetHits(20) = %d, want 2", got)
	}
	h.Hit(400) // window at 400 is [101,400]; both earlier hits (10, 20) are now expired
	if got := h.GetHits(400); got != 1 {
		t.Fatalf("GetHits(400) = %d, want 1 (only the ts=400 hit itself remains)", got)
	}
}

func TestHitCounter_ManyHitsWithinLargeWindow(t *testing.T) {
	h := NewHitCounter()
	for ts := 1; ts <= 300; ts++ {
		h.Hit(ts)
	}
	if got := h.GetHits(300); got != 300 {
		t.Fatalf("GetHits(300) = %d, want 300", got)
	}
}

func TestHitCounter_SlidingWindowStress(t *testing.T) {
	h := NewHitCounter()
	// One hit per second for 600 seconds; at each point, GetHits should report at most
	// min(ts, 300).
	for ts := 1; ts <= 600; ts++ {
		h.Hit(ts)
		want := ts
		if want > 300 {
			want = 300
		}
		if got := h.GetHits(ts); got != want {
			t.Fatalf("at ts=%d: GetHits(%d) = %d, want %d", ts, ts, got, want)
		}
	}
}
