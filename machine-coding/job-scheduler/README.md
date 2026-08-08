# Job Scheduler

**Suggested time box: 45 minutes.**

## Problem statement

Design a scheduler that supports one-shot jobs ("run this once at time T") and
recurring jobs ("run this every interval, starting at time T"). Deliberately, this
scheduler is **not** driven by a background goroutine polling the real wall clock —
instead it's driven by an explicit `Tick(now time.Time)` call, the same "inject the
clock" idea as `../rate-limiter`. This is a real pattern (plenty of production
schedulers separate "what's due" logic from "how time actually advances" for exactly
this reason): it keeps the scheduling logic itself simple, synchronous, and testable
without real sleeps or goroutine races.

## Functional requirements

1. `ScheduleOnce(at time.Time, fn func()) int` schedules `fn` to run once, the first
   time `Tick` is called with `now >= at`. Returns a unique job ID.
2. `ScheduleRecurring(at time.Time, interval time.Duration, fn func()) int` schedules
   `fn` to run at `at`, then again at `at+interval`, `at+2*interval`, ... forever
   (until canceled), each time `Tick` reaches or passes that occurrence. Returns a
   unique job ID.
3. `Cancel(jobID int) bool` stops a job from running again (a one-shot job that hasn't
   fired yet won't fire; a recurring job stops recurring). Returns `false` if `jobID`
   doesn't exist.
4. `Tick(now time.Time) int` advances the scheduler's clock to `now` and **synchronously
   runs every due job**, in ascending `(scheduled time, job ID)` order (so ties are
   deterministic). Returns the number of job *executions* that happened during this
   call (a recurring job that fires 3 times counts as 3).
5. **Catch-up behavior**: if `Tick` is called with a `now` that skips over multiple
   occurrences of a recurring job (e.g. the interval is 1 minute but `Tick` jumps
   forward by 5 minutes), that job must run once for **each** occurrence that's now due
   — not just once. A one-shot job never fires more than once no matter how far `now`
   advances or how many times `Tick` is called afterward.

## Rubric

- [ ] `Tick` is synchronous and deterministic — no goroutines, no real `time.Sleep`,
      no reliance on wall-clock time anywhere in the scheduling logic.
- [ ] Recurring jobs correctly catch up (fire multiple times) when `Tick` skips ahead
      by more than one interval.
- [ ] A canceled job — whether one-shot or recurring — never fires again, even if it
      was already due at cancellation time.
- [ ] Jobs due at the exact same scheduled time run in job-ID order, not map-iteration
      order (Go maps don't guarantee iteration order — this will bite you if you're not
      careful).
- [ ] `Tick` called with a `now` earlier than a previous `Tick` call, or with nothing
      due, doesn't panic and just returns 0.

## Hints

- Store jobs in a slice or map keyed by job ID, each with: next-run time, an interval
  (zero/absent for one-shot), the function, and a canceled flag.
- On each `Tick`, collect all non-canceled jobs whose next-run time is `<= now`, sort
  them by `(nextRun, jobID)`, and run them in that order. For a recurring job, after
  running it once, advance its next-run time by `interval` and — if the NEW next-run
  time is still `<= now` — it's due again in this same `Tick` call (loop until it's no
  longer due, or check due jobs with a `for` loop that re-scans rather than a single
  pass).
- A one-shot job that fires should be marked as done (e.g., removed or a "fired"
  flag) so it can never run again on a later `Tick`.
