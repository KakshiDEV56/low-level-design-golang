package thediningphilosophers

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const numPhilosophers = 5

// Fork numbering for the operation encoding below, matching the problem statement:
// {1: left, 2: right}.
const (
	forkLeft  = 1
	forkRight = 2
)

// Operation codes, matching the problem statement's own output format exactly:
// output[i] = [a, b, c] where a = philosopher, b = fork (1=left, 2=right), c = operation
// (1=pick, 2=put, 3=eat). Recording events in this exact encoding means that if a test
// fails, the dumped log can be compared directly against the style of log LeetCode itself
// shows in the problem's "Example 1" explanation.
const (
	opPick = 1
	opPut  = 2
	opEat  = 3
)

// Physical fork layout: fork i sits between philosopher i and philosopher (i+1)%5.
// We fix "left fork of philosopher p" = fork p and "right fork of philosopher p" =
// fork (p+4)%5 (i.e. p-1 mod 5). This is just a labeling convention for the test's own
// bookkeeping — since a fork is a single physical resource shared by exactly two
// neighboring philosophers, this convention is consistent regardless of how the
// solution itself tracks fork state internally.
func leftForkOf(p int) int  { return p }
func rightForkOf(p int) int { return (p + numPhilosophers - 1) % numPhilosophers }

// event is one recorded call into pickLeftFork/pickRightFork/eat/putLeftFork/putRightFork,
// in the same [philosopher, fork, op] shape the problem statement itself uses.
type event struct {
	philosopher int
	fork        int // 0 for eat (not fork-specific), else forkLeft/forkRight
	op          int // opPick / opPut / opEat
}

func opName(op int) string {
	switch op {
	case opPick:
		return "pick"
	case opPut:
		return "put"
	case opEat:
		return "eat"
	default:
		return "?"
	}
}

func forkName(fork int) string {
	switch fork {
	case forkLeft:
		return "left"
	case forkRight:
		return "right"
	default:
		return "-"
	}
}

func (e event) String() string {
	return fmt.Sprintf("[philosopher=%d fork=%s op=%s]", e.philosopher, forkName(e.fork), opName(e.op))
}

// eventLog is a concurrency-safe append-only log of every fork/eat call made during a
// run, so that on failure we can print exactly what led up to the problem instead of
// just a bare "expected X got Y".
type eventLog struct {
	mu     sync.Mutex
	events []event
}

func (l *eventLog) record(e event) {
	l.mu.Lock()
	l.events = append(l.events, e)
	l.mu.Unlock()
}

// tail returns the most recent n events (or all of them if there are fewer than n),
// formatted one per line, for inclusion in a failure message.
func (l *eventLog) tail(n int) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	start := 0
	if len(l.events) > n {
		start = len(l.events) - n
	}
	var b strings.Builder
	for i := start; i < len(l.events); i++ {
		fmt.Fprintf(&b, "  #%d %s\n", i, l.events[i])
	}
	return b.String()
}

// runDiningPhilosophers has each of the 5 philosophers call wantsToEat n times
// sequentially (from n independent goroutines, one per philosopher — this is the
// "5 threads simultaneously use one object" setup the problem describes), and verifies:
//  1. No two philosophers ever hold the same fork at once (mutual exclusion).
//  2. eat() is only ever called while both forks are actually held.
//  3. Every philosopher eats exactly n times (no missed or duplicated meals).
//  4. Every fork ends up released (no fork left "stuck" held).
//  5. The whole run finishes within a timeout (no deadlock).
//
// On ANY failure, the full (or tail of the) event log is printed via t.Logf so you can
// see the exact sequence of pick/put/eat calls — in the same [philosopher, fork, op]
// shape as the problem statement's own example — that led up to the violation. Go's
// testing package only shows Logf output for a test that fails (or when run with -v),
// so this adds no noise to a passing run.
func runDiningPhilosophers(t *testing.T, n int) {
	t.Helper()

	dp := &DiningPhilosophers{}

	var forkState [numPhilosophers]int32 // 0 = free, 1 = held
	var eatCount [numPhilosophers]int32
	log := &eventLog{}

	var wg sync.WaitGroup
	wg.Add(numPhilosophers)

	for p := 0; p < numPhilosophers; p++ {
		p := p
		go func() {
			defer wg.Done()
			lf, rf := leftForkOf(p), rightForkOf(p)
			for i := 0; i < n; i++ {
				dp.wantsToEat(p,
					func() { // pickLeftFork
						log.record(event{p, forkLeft, opPick})
						if !atomic.CompareAndSwapInt32(&forkState[lf], 0, 1) {
							t.Errorf(
								"philosopher %d tried to PICK its LEFT fork (fork #%d), but that fork was already held by someone else.\n"+
									"This means two philosophers held the same fork at the same time — a mutual-exclusion violation.\n"+
									"Last events before the failure:\n%s",
								p, lf, log.tail(15))
						}
					},
					func() { // pickRightFork
						log.record(event{p, forkRight, opPick})
						if !atomic.CompareAndSwapInt32(&forkState[rf], 0, 1) {
							t.Errorf(
								"philosopher %d tried to PICK its RIGHT fork (fork #%d), but that fork was already held by someone else.\n"+
									"This means two philosophers held the same fork at the same time — a mutual-exclusion violation.\n"+
									"Last events before the failure:\n%s",
								p, rf, log.tail(15))
						}
					},
					func() { // eat
						log.record(event{p, 0, opEat})
						leftHeld := atomic.LoadInt32(&forkState[lf]) == 1
						rightHeld := atomic.LoadInt32(&forkState[rf]) == 1
						if !leftHeld || !rightHeld {
							t.Errorf(
								"philosopher %d called eat() without holding both forks (left fork #%d held=%v, right fork #%d held=%v).\n"+
									"A philosopher must never eat before acquiring BOTH forks.\n"+
									"Last events before the failure:\n%s",
								p, lf, leftHeld, rf, rightHeld, log.tail(15))
						}
						atomic.AddInt32(&eatCount[p], 1)
					},
					func() { // putLeftFork
						log.record(event{p, forkLeft, opPut})
						if !atomic.CompareAndSwapInt32(&forkState[lf], 1, 0) {
							t.Errorf(
								"philosopher %d tried to PUT DOWN its LEFT fork (fork #%d), but that fork wasn't marked as held.\n"+
									"This means putLeftFork was called without a matching pickLeftFork, or it was called twice.\n"+
									"Last events before the failure:\n%s",
								p, lf, log.tail(15))
						}
					},
					func() { // putRightFork
						log.record(event{p, forkRight, opPut})
						if !atomic.CompareAndSwapInt32(&forkState[rf], 1, 0) {
							t.Errorf(
								"philosopher %d tried to PUT DOWN its RIGHT fork (fork #%d), but that fork wasn't marked as held.\n"+
									"This means putRightFork was called without a matching pickRightFork, or it was called twice.\n"+
									"Last events before the failure:\n%s",
								p, rf, log.tail(15))
						}
					},
				)
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf(
			"wantsToEat(n=%d) timed out after 5s waiting for all %d philosophers to finish — this almost always means a deadlock "+
				"(e.g. every philosopher picked up one fork and is now stuck waiting forever for the second one).\n"+
				"Last events before the timeout:\n%s",
			n, numPhilosophers, log.tail(30))
	}

	// Every philosopher should have eaten exactly n times — not fewer (a call got lost
	// or stuck) and not more (a call got double-counted).
	for p := 0; p < numPhilosophers; p++ {
		if got := atomic.LoadInt32(&eatCount[p]); got != int32(n) {
			t.Errorf(
				"philosopher %d ate %d time(s), want exactly %d.\n"+
					"Full event log for this run:\n%s",
				p, got, n, log.tail(200))
		}
	}
	// Every fork should be back in the "free" state once every philosopher is done —
	// anything else means a pick was never matched by a put.
	for f := 0; f < numPhilosophers; f++ {
		if got := atomic.LoadInt32(&forkState[f]); got != 0 {
			t.Errorf(
				"fork #%d was left in the HELD state after all philosophers finished (want it free/released).\n"+
					"This means some pick*Fork() call was never matched by a corresponding put*Fork() call.\n"+
					"Full event log for this run:\n%s",
				f, log.tail(200))
		}
	}
}

// TestDiningPhilosophers runs the simulation across a range of meal counts n, including
// the problem statement's own Example 1 (n=1). Each subtest is named after its n so a
// failure immediately tells you which meal count broke — e.g. a failure only at large n
// often points to a race that needs enough iterations to surface, whereas a failure at
// n=1 points to a more fundamental logic bug.
func TestDiningPhilosophers(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"n1_matchesExample", 1}, // problem's Example 1 uses n=1
		{"n2", 2},
		{"n3", 3},
		{"n5", 5},
		{"n8", 8},
		{"n10", 10},
		{"n15", 15},
		{"n20", 20},
		{"n25", 25},
		{"n30", 30},
		{"n50", 50},
		{"n100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runDiningPhilosophers(t, tt.n)
		})
	}
}

// TestDiningPhilosophers_Repeated stress-tests a fixed n=10 across 10 independent runs.
// Starvation/race bugs in this problem are scheduler-dependent — a buggy solution can
// pass once and fail on the very next run — so each iteration is its own named subtest
// (iter0..iter9) to make it obvious which specific run(s) failed rather than just "some
// run somewhere failed". Run with `go test -race` to also catch data races on the fork
// state itself (a race that wouldn't necessarily show up as a wrong eat count, but is
// still a real bug).
func TestDiningPhilosophers_Repeated(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("iter%d", i), func(t *testing.T) {
			runDiningPhilosophers(t, 10)
		})
	}
}
