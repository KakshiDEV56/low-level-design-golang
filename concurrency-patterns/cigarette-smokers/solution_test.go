package cigarettesmokers

import (
	"testing"
	"time"
)

// smokeSignals wires up all three smokers on a Table and gives each one a channel that
// closes the moment its smoke() callback fires, so tests can assert exactly which
// smoker (if any) was woken by a given sequence of Provide* calls.
type smokeSignals struct {
	matches chan struct{}
	tobacco chan struct{}
	paper   chan struct{}
}

func startSmokers(table *Table) *smokeSignals {
	s := &smokeSignals{
		matches: make(chan struct{}),
		tobacco: make(chan struct{}),
		paper:   make(chan struct{}),
	}
	go table.SmokeWithMatches(func() { close(s.matches) })
	go table.SmokeWithTobacco(func() { close(s.tobacco) })
	go table.SmokeWithPaper(func() { close(s.paper) })
	return s
}

func assertWoken(t *testing.T, ch chan struct{}, who string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s was never woken", who)
	}
}

func assertNotWoken(t *testing.T, ch chan struct{}, who string) {
	t.Helper()
	select {
	case <-ch:
		t.Fatalf("%s was woken, but shouldn't have been for this ingredient combination", who)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestTable_TobaccoThenPaper_WakesSmokerWithMatches(t *testing.T) {
	table := NewTable()
	s := startSmokers(table)

	table.ProvideTobacco()
	table.ProvidePaper()

	assertWoken(t, s.matches, "smoker with matches")
	assertNotWoken(t, s.tobacco, "smoker with tobacco")
	assertNotWoken(t, s.paper, "smoker with paper")
}

func TestTable_PaperThenTobacco_WakesSmokerWithMatches(t *testing.T) {
	// Order shouldn't matter — same pair, opposite arrival order.
	table := NewTable()
	s := startSmokers(table)

	table.ProvidePaper()
	table.ProvideTobacco()

	assertWoken(t, s.matches, "smoker with matches")
	assertNotWoken(t, s.tobacco, "smoker with tobacco")
	assertNotWoken(t, s.paper, "smoker with paper")
}

func TestTable_PaperThenMatch_WakesSmokerWithTobacco(t *testing.T) {
	table := NewTable()
	s := startSmokers(table)

	table.ProvidePaper()
	table.ProvideMatch()

	assertWoken(t, s.tobacco, "smoker with tobacco")
	assertNotWoken(t, s.matches, "smoker with matches")
	assertNotWoken(t, s.paper, "smoker with paper")
}

func TestTable_TobaccoThenMatch_WakesSmokerWithPaper(t *testing.T) {
	table := NewTable()
	s := startSmokers(table)

	table.ProvideTobacco()
	table.ProvideMatch()

	assertWoken(t, s.paper, "smoker with paper")
	assertNotWoken(t, s.matches, "smoker with matches")
	assertNotWoken(t, s.tobacco, "smoker with tobacco")
}

func TestTable_SingleIngredientWakesNoOne(t *testing.T) {
	table := NewTable()
	s := startSmokers(table)

	table.ProvideTobacco()

	assertNotWoken(t, s.matches, "smoker with matches")
	assertNotWoken(t, s.tobacco, "smoker with tobacco")
	assertNotWoken(t, s.paper, "smoker with paper")
}

// TestTable_SequentialRounds drives the agent through several full rounds, each one
// waiting for the previous smoker to finish before providing the next pair (mirroring
// the book's agentSem gating, done here at the test/driver level), and checks that the
// correct smoker fires every round regardless of which pair was provided.
func TestTable_SequentialRounds(t *testing.T) {
	table := NewTable()

	type round struct {
		provide func()
		who     string
	}
	rounds := []round{
		{func() { table.ProvideTobacco(); table.ProvidePaper() }, "matches"},
		{func() { table.ProvidePaper(); table.ProvideMatch() }, "tobacco"},
		{func() { table.ProvideTobacco(); table.ProvideMatch() }, "paper"},
		{func() { table.ProvideMatch(); table.ProvidePaper() }, "tobacco"},
		{func() { table.ProvideMatch(); table.ProvideTobacco() }, "paper"},
		{func() { table.ProvidePaper(); table.ProvideTobacco() }, "matches"},
	}

	for i, r := range rounds {
		matchesDone := make(chan struct{})
		tobaccoDone := make(chan struct{})
		paperDone := make(chan struct{})
		go table.SmokeWithMatches(func() { close(matchesDone) })
		go table.SmokeWithTobacco(func() { close(tobaccoDone) })
		go table.SmokeWithPaper(func() { close(paperDone) })

		r.provide()

		var got chan struct{}
		switch r.who {
		case "matches":
			got = matchesDone
		case "tobacco":
			got = tobaccoDone
		case "paper":
			got = paperDone
		}

		select {
		case <-got:
		case <-time.After(2 * time.Second):
			t.Fatalf("round %d: smoker with %s was never woken", i, r.who)
		}
	}
}

// TestTable_StressManyRounds runs many sequential rounds with varied ingredient
// ordering, both to check overall correctness at scale and to help surface races in
// the shared flag/mutex bookkeeping. Run with `go test -race`.
func TestTable_StressManyRounds(t *testing.T) {
	table := NewTable()

	pairs := []struct {
		provide func()
		who     string
	}{
		{func() { table.ProvideTobacco(); table.ProvidePaper() }, "matches"},
		{func() { table.ProvidePaper(); table.ProvideMatch() }, "tobacco"},
		{func() { table.ProvideTobacco(); table.ProvideMatch() }, "paper"},
	}

	const rounds = 100
	for i := 0; i < rounds; i++ {
		p := pairs[i%len(pairs)]

		done := make(chan struct{})
		switch p.who {
		case "matches":
			go table.SmokeWithMatches(func() { close(done) })
			go table.SmokeWithTobacco(func() {})
			go table.SmokeWithPaper(func() {})
		case "tobacco":
			go table.SmokeWithTobacco(func() { close(done) })
			go table.SmokeWithMatches(func() {})
			go table.SmokeWithPaper(func() {})
		case "paper":
			go table.SmokeWithPaper(func() { close(done) })
			go table.SmokeWithMatches(func() {})
			go table.SmokeWithTobacco(func() {})
		}

		p.provide()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatalf("round %d: smoker with %s was never woken", i, p.who)
		}
	}
}
