package designleaderboard

import "testing"

func TestLeaderboard_Example1(t *testing.T) {
	l := NewLeaderboard()
	l.AddScore(1, 73)
	l.AddScore(2, 56)
	l.AddScore(3, 39)
	l.AddScore(4, 51)
	l.AddScore(5, 4)
	if got := l.Top(1); got != 73 {
		t.Fatalf("Top(1) = %d, want 73", got)
	}
	l.Reset(1)
	l.Reset(2)
	l.AddScore(2, 51)
	if got := l.Top(3); got != 141 {
		t.Fatalf("Top(3) = %d, want 141", got)
	}
}

func TestLeaderboard_SinglePlayer(t *testing.T) {
	l := NewLeaderboard()
	l.AddScore(1, 10)
	if got := l.Top(1); got != 10 {
		t.Fatalf("Top(1) = %d, want 10", got)
	}
}

func TestLeaderboard_AddScoreAccumulates(t *testing.T) {
	l := NewLeaderboard()
	l.AddScore(1, 10)
	l.AddScore(1, 20)
	l.AddScore(1, 5)
	if got := l.Top(1); got != 35 {
		t.Fatalf("Top(1) = %d, want 35", got)
	}
}

func TestLeaderboard_TopAllPlayers(t *testing.T) {
	l := NewLeaderboard()
	l.AddScore(1, 10)
	l.AddScore(2, 20)
	l.AddScore(3, 30)
	if got := l.Top(3); got != 60 {
		t.Fatalf("Top(3) = %d, want 60", got)
	}
}

func TestLeaderboard_ResetThenReaddStartsFresh(t *testing.T) {
	l := NewLeaderboard()
	l.AddScore(1, 100)
	l.Reset(1)
	l.AddScore(1, 5)
	if got := l.Top(1); got != 5 {
		t.Fatalf("Top(1) after reset+readd = %d, want 5", got)
	}
}

func TestLeaderboard_TopKSelectsHighestScores(t *testing.T) {
	l := NewLeaderboard()
	l.AddScore(1, 100)
	l.AddScore(2, 90)
	l.AddScore(3, 80)
	l.AddScore(4, 70)
	l.AddScore(5, 60)
	if got := l.Top(2); got != 190 {
		t.Fatalf("Top(2) = %d, want 190 (100+90)", got)
	}
	if got := l.Top(4); got != 340 {
		t.Fatalf("Top(4) = %d, want 340 (100+90+80+70)", got)
	}
}

func TestLeaderboard_ResetRemovesFromTopCalculation(t *testing.T) {
	l := NewLeaderboard()
	l.AddScore(1, 100)
	l.AddScore(2, 50)
	l.Reset(1)
	if got := l.Top(1); got != 50 {
		t.Fatalf("Top(1) after resetting the top scorer = %d, want 50", got)
	}
}

func TestLeaderboard_TiedScores(t *testing.T) {
	l := NewLeaderboard()
	l.AddScore(1, 50)
	l.AddScore(2, 50)
	l.AddScore(3, 50)
	if got := l.Top(2); got != 100 {
		t.Fatalf("Top(2) with all-tied scores = %d, want 100", got)
	}
	if got := l.Top(3); got != 150 {
		t.Fatalf("Top(3) with all-tied scores = %d, want 150", got)
	}
}

func TestLeaderboard_ManyPlayersStress(t *testing.T) {
	l := NewLeaderboard()
	total := 0
	for id := 1; id <= 100; id++ {
		l.AddScore(id, id) // player id gets score = id
		total += id
	}
	if got := l.Top(100); got != total {
		t.Fatalf("Top(100) = %d, want %d", got, total)
	}
	// Top 10 should be players 91..100, sum = 91+...+100
	want := 0
	for id := 91; id <= 100; id++ {
		want += id
	}
	if got := l.Top(10); got != want {
		t.Fatalf("Top(10) = %d, want %d", got, want)
	}
}

func TestLeaderboard_InterleavedAddResetTop(t *testing.T) {
	l := NewLeaderboard()
	l.AddScore(1, 20)
	l.AddScore(2, 30)
	if got := l.Top(2); got != 50 {
		t.Fatalf("Top(2) = %d, want 50", got)
	}
	l.AddScore(1, 100) // player 1 now has 120
	if got := l.Top(1); got != 120 {
		t.Fatalf("Top(1) = %d, want 120", got)
	}
	l.Reset(2)
	l.AddScore(3, 10)
	if got := l.Top(2); got != 130 {
		t.Fatalf("Top(2) = %d, want 130 (120+10)", got)
	}
}
