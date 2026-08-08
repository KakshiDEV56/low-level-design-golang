package jobscheduler

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

var start = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func at(seconds int) time.Time {
	return start.Add(time.Duration(seconds) * time.Second)
}

func TestScheduleOnce_DoesNotRunBeforeDue(t *testing.T) {
	s := NewScheduler()
	ran := 0
	s.ScheduleOnce(at(10), func() { ran++ })

	s.Tick(at(5))
	if ran != 0 {
		t.Fatalf("ran = %d before the due time, want 0", ran)
	}
}

func TestScheduleOnce_RunsWhenDue(t *testing.T) {
	s := NewScheduler()
	ran := 0
	s.ScheduleOnce(at(10), func() { ran++ })

	n := s.Tick(at(10))
	if ran != 1 {
		t.Fatalf("ran = %d at the due time, want 1", ran)
	}
	if n != 1 {
		t.Fatalf("Tick() returned %d executions, want 1", n)
	}
}

func TestScheduleOnce_RunsExactlyOnceAcrossMultipleTicks(t *testing.T) {
	s := NewScheduler()
	ran := 0
	s.ScheduleOnce(at(10), func() { ran++ })

	s.Tick(at(10))
	s.Tick(at(20))
	s.Tick(at(30))

	if ran != 1 {
		t.Fatalf("ran = %d across multiple later Ticks, want exactly 1", ran)
	}
}

func TestScheduleRecurring_RunsRepeatedly(t *testing.T) {
	s := NewScheduler()
	ran := 0
	s.ScheduleRecurring(at(10), 10*time.Second, func() { ran++ })

	s.Tick(at(10))
	if ran != 1 {
		t.Fatalf("after Tick(10s): ran = %d, want 1", ran)
	}
	s.Tick(at(20))
	if ran != 2 {
		t.Fatalf("after Tick(20s): ran = %d, want 2", ran)
	}
	s.Tick(at(30))
	if ran != 3 {
		t.Fatalf("after Tick(30s): ran = %d, want 3", ran)
	}
}

func TestScheduleRecurring_DoesNotRunTwiceForSameTick(t *testing.T) {
	s := NewScheduler()
	ran := 0
	s.ScheduleRecurring(at(10), 10*time.Second, func() { ran++ })

	s.Tick(at(10))
	s.Tick(at(10)) // calling Tick again with the SAME now must not re-fire
	if ran != 1 {
		t.Fatalf("ran = %d after two Tick() calls at the same time, want 1", ran)
	}
}

func TestScheduleRecurring_CatchUpMultipleOccurrencesInOneTick(t *testing.T) {
	s := NewScheduler()
	ran := 0
	s.ScheduleRecurring(at(0), 10*time.Second, func() { ran++ })

	// Jumping straight to t=35s should fire at t=0,10,20,30 -> 4 occurrences.
	n := s.Tick(at(35))
	if ran != 4 {
		t.Fatalf("ran = %d after a single Tick() jumping past 4 occurrences, want 4", ran)
	}
	if n != 4 {
		t.Fatalf("Tick() returned %d executions, want 4", n)
	}

	// The next occurrence should be t=40s, not t=30s or t=35s.
	n2 := s.Tick(at(39))
	if ran != 4 || n2 != 0 {
		t.Fatalf("after Tick(39s): ran = %d (want 4), Tick returned %d (want 0) — next occurrence is t=40s", ran, n2)
	}
	s.Tick(at(40))
	if ran != 5 {
		t.Fatalf("after Tick(40s): ran = %d, want 5", ran)
	}
}

func TestCancel_OneShotPreventsExecution(t *testing.T) {
	s := NewScheduler()
	ran := 0
	id := s.ScheduleOnce(at(10), func() { ran++ })

	if !s.Cancel(id) {
		t.Fatal("Cancel() on an existing not-yet-fired job = false, want true")
	}
	s.Tick(at(10))
	if ran != 0 {
		t.Fatalf("ran = %d after canceling before it fired, want 0", ran)
	}
}

func TestCancel_RecurringStopsFutureRuns(t *testing.T) {
	s := NewScheduler()
	ran := 0
	id := s.ScheduleRecurring(at(10), 10*time.Second, func() { ran++ })

	s.Tick(at(10))
	if ran != 1 {
		t.Fatalf("ran = %d after first occurrence, want 1", ran)
	}
	if !s.Cancel(id) {
		t.Fatal("Cancel() on an active recurring job = false, want true")
	}
	s.Tick(at(20))
	s.Tick(at(30))
	if ran != 1 {
		t.Fatalf("ran = %d after canceling a recurring job, want 1 (no further runs)", ran)
	}
}

func TestCancel_UnknownJobIDReturnsFalse(t *testing.T) {
	s := NewScheduler()
	if s.Cancel(9999) {
		t.Fatal("Cancel() on an unknown job ID = true, want false")
	}
}

func TestCancel_TwiceReturnsFalseSecondTime(t *testing.T) {
	s := NewScheduler()
	id := s.ScheduleOnce(at(10), func() {})
	if !s.Cancel(id) {
		t.Fatal("first Cancel() = false, want true")
	}
	if s.Cancel(id) {
		t.Fatal("second Cancel() on an already-canceled job = true, want false")
	}
}

func TestTick_OrdersJobsBySameTimeThenJobID(t *testing.T) {
	s := NewScheduler()
	var order []int

	type scheduledJob struct {
		id    int
		label int
	}
	var jobs []scheduledJob
	for label := 0; label < 3; label++ {
		label := label
		id := s.ScheduleOnce(at(10), func() { order = append(order, label) })
		jobs = append(jobs, scheduledJob{id: id, label: label})
	}

	// Expected execution order is by ascending job ID, regardless of the order the
	// jobs were scheduled in or how IDs happen to be assigned.
	sort.Slice(jobs, func(i, j int) bool { return jobs[i].id < jobs[j].id })
	wantOrder := make([]int, len(jobs))
	for i, j := range jobs {
		wantOrder[i] = j.label
	}

	s.Tick(at(10))

	if !reflect.DeepEqual(order, wantOrder) {
		t.Fatalf("execution order (by label) = %v, want %v (ascending job ID order)", order, wantOrder)
	}
}

func TestTick_NothingDueReturnsZero(t *testing.T) {
	s := NewScheduler()
	s.ScheduleOnce(at(100), func() {})
	if n := s.Tick(at(10)); n != 0 {
		t.Fatalf("Tick() with nothing due = %d, want 0", n)
	}
}

func TestTick_MultipleIndependentOneShots(t *testing.T) {
	s := NewScheduler()
	var ranIDs []string
	s.ScheduleOnce(at(5), func() { ranIDs = append(ranIDs, "a") })
	s.ScheduleOnce(at(10), func() { ranIDs = append(ranIDs, "b") })
	s.ScheduleOnce(at(15), func() { ranIDs = append(ranIDs, "c") })

	s.Tick(at(12))
	if len(ranIDs) != 2 || ranIDs[0] != "a" || ranIDs[1] != "b" {
		t.Fatalf("ranIDs = %v, want [a b] after Tick(12s)", ranIDs)
	}
	s.Tick(at(15))
	if len(ranIDs) != 3 || ranIDs[2] != "c" {
		t.Fatalf("ranIDs = %v, want [a b c] after Tick(15s)", ranIDs)
	}
}
