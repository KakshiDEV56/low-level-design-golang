package jobscheduler

import "time"

/*
See README.md for the full problem statement, requirements, and rubric.

TODO:
 1. Add the fields you need on Scheduler below. Hint: a map (or slice) of job records,
    each tracking { id int; nextRun time.Time; interval time.Duration (zero =
    one-shot); fn func(); canceled bool }, plus a counter for generating unique job
    IDs.
 2. Implement NewScheduler.
 3. Implement ScheduleOnce, ScheduleRecurring, Cancel, Tick per the README — pay
    special attention to the catch-up behavior for recurring jobs and the
    (scheduled-time, job-ID) ordering requirement for jobs due at the same instant.
*/

type Scheduler struct {
	// TODO: add your fields here
}

func NewScheduler() *Scheduler {
	// TODO: implement
	return &Scheduler{}
}

func (s *Scheduler) ScheduleOnce(at time.Time, fn func()) int {
	// TODO: implement
	return 0
}

func (s *Scheduler) ScheduleRecurring(at time.Time, interval time.Duration, fn func()) int {
	// TODO: implement
	return 0
}

func (s *Scheduler) Cancel(jobID int) bool {
	// TODO: implement
	return false
}

func (s *Scheduler) Tick(now time.Time) int {
	// TODO: implement
	return 0
}
