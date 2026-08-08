# Reusable Barrier (Cyclic Barrier)

*Based on The Little Book of Semaphores, §3.6–§3.7 ("Barrier" and "Reusable barrier").*

## The concept

A barrier makes `n` goroutines all wait at a point until every one of them has arrived,
then releases them all together:

```
rendezvous       # do some work
critical point   # nobody may proceed past here until ALL n have arrived
```

The book builds this up through several *broken* attempts before landing on the correct
one, and it's worth understanding why the naive versions fail, because the same mistakes
show up in real code:

- **Non-solution #1:** lock the barrier, increment a counter, and unlock the barrier
  when `count == n` — checking `count == n` outside the mutex that protects `count` is a
  race: two goroutines can both observe `count == n` and both call `signal()`.
- **Non-solution #2:** fixing that race still leaves a subtler bug — since a barrier is
  normally used *inside a loop* (each round does work, then hits the barrier again), a
  goroutine that reaches the end of round 1 fast enough can loop around and get to round
  2's barrier before the slowest goroutine has even left round 1's barrier. It "laps"
  the others. This is exactly the kind of bug that's nearly impossible to catch with
  testing and only shows up occasionally under specific scheduling — the book calls it
  Deadlock/bug #2 and #3.

The fix is the **two-phase turnstile**: one turnstile that opens once everyone has
arrived (letting them all through), and a *second* turnstile that opens once everyone
has *left* the first one (which is what makes it safe to loop back around). This is
sometimes called a "cyclic barrier" (Java has `java.util.concurrent.CyclicBarrier` built
on exactly this idea).

## Why this matters for a Go engineer

`sync.WaitGroup` is Go's built-in barrier-like primitive, but it's **single-use**:
calling `Add` after `Wait` has returned (in a new "round") without care is documented as
racy unless carefully sequenced, and `WaitGroup` provides no way to run through a fixed
`n` goroutines and reuse the same barrier every round in the natural "let each round hand
back control" way `CyclicBarrier` does. Round-based batch processing (map-reduce style
supersteps, simulation ticks, multi-stage pipelines where every worker must finish stage
`k` before anyone starts stage `k+1`) is a real, recurring shape of problem — knowing how
to build a barrier that's actually safe to reuse is what this teaches.

## The constraint to implement

- `NewCyclicBarrier(n int)` creates a barrier for exactly `n` participants.
- `Await()` blocks the calling goroutine until all `n` goroutines currently in this
  "round" have also called `Await()`. Once the `n`th one arrives, all `n` are released
  together.
- The barrier must be safe to call `Await()` on again immediately afterward, in a loop,
  potentially many rounds — and under NO circumstances should any goroutine be allowed
  to start round `k+1`'s work before every goroutine has finished round `k` (no lapping).

## Hint

Use two "turnstiles" (each a 1-buffered channel, or a semaphore-like gate) plus a mutex-
protected counter, per the book's two-phase solution: the first turnstile is initially
locked; when the `n`th goroutine arrives it locks the *second* turnstile and unlocks the
*first*. Each goroutine passes through the first turnstile (wait, then immediately
signal it again so the next goroutine can pass). On the way "back" (simulating loop
re-entry), each goroutine decrements the counter, and when it hits zero, the last one
relocks the first turnstile and unlocks the second — which is what everyone actually
waits on before returning from `Await()`.
