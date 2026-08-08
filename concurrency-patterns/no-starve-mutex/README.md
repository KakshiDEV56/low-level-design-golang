# No-Starve Mutex (Morris's Algorithm)

*Based on The Little Book of Semaphores, §4.3 ("No-starve mutex").*

## The concept

A plain mutex guarantees mutual exclusion, but it does **not** guarantee fairness. The
book demonstrates this with three threads running:

```
while True:
    mutex.wait()
    # critical section
    mutex.signal()
```

Thread A gets the mutex; B and C queue up. When A leaves, B enters — but nothing
stops A from looping back around and joining the queue again before B even leaves.
When B leaves, there's no guarantee C goes next: if A gets it again, and then B rejoins
the queue, we're back to the starting position — and we can repeat this cycle forever.
**C starves**, even though the mutex is working "correctly" in the narrow sense of never
letting two threads in at once.

This is only possible because a basic semaphore only guarantees **Property 3** ("if
threads are waiting when a thread signals, one of them gets woken") — it says nothing
about *which* one. A semaphore that also guarantees **Property 4** ("the number of
threads that get woken before a given waiting thread is bounded") is called a **strong
semaphore**; one with only Property 3 is a **weak semaphore**. Go's `sync.Mutex` is
explicitly documented as *not* providing strict FIFO ordering by default (though it has
an internal "starvation mode" fallback for pathological cases) — so this exercise is
about building a mutex with a hard fairness guarantee on top of primitives that don't
give you one for free.

### Morris's solution

Dijkstra conjectured this couldn't be done with only weak semaphores/mutexes. J.M.
Morris refuted that in 1979. His solution uses **two turnstiles, creating three
"rooms"** that a goroutine passes through in sequence on the way to the critical
section:

- **Room 1:** the initial holding area. A `room1` counter (protected by its own mutex)
  tracks how many goroutines are here.
- **Turnstile `t1`:** goroutines cascade through it one at a time into Room 2,
  incrementing `room1` down and `room2` up as they go.
- **Room 2:** goroutines wait here, gated by `t1`/`t2`.
- **Turnstile `t2`:** the final gate into the critical section (Room 3).

The key trick: responsibility for guarding the critical section is **split** between
`t1` and `t2` — a goroutine must hold exactly one of them (never both, never neither) to
be in the critical section. As long as `room1 > 0` (more goroutines are still queued
behind), a goroutine finishing Room 1 keeps `t2` locked and instead signals `t1`, letting
the *next* queued goroutine advance instead of letting anyone loop back around and cut
in. Only once `room1` hits zero does `t2` open up, and only once `room2` hits zero (i.e.
everyone who was in the "first wave" has gone through the critical section) does `t1`
reopen for a new wave. Since there's a finite number of goroutines, the current wave is
guaranteed to fully drain before any goroutine can start a second wave — so no goroutine
can ever overtake one that arrived before it.

## Why this matters for a Go engineer

- It's a real answer to "how would you make a lock fair?" — a legitimate question for
  systems where request order matters (e.g., a queue-processing service where you don't
  want one hot key to be perpetually re-serviced while others wait).
- It's the same two-turnstile idea as the reusable barrier (`../reusable-barrier`) and
  the writer-priority readers-writers fix (`../readers-writers`) — recognizing "I've seen
  this shape before" is exactly the kind of pattern-matching that makes someone fast at
  concurrent code review.
- It's a genuinely well-known "gotcha" interview follow-up to "implement a mutex":
  *"is your mutex fair? How do you know?"*

## The constraint to implement

- `NewNoStarveMutex()` creates the mutex, unlocked.
- `Lock()` / `Unlock()` provide standard mutual exclusion (never two goroutines in the
  critical section at once) — PLUS: if goroutine X calls `Lock()` before goroutine Y
  does, X is guaranteed to enter the critical section before Y, no matter how many times
  Y's goroutine loops around and calls `Lock()` again afterward.

## Hint

Translate Morris's algorithm directly: `room1`/`room2` (int counters) each protected the
way the book describes, and `t1`/`t2` as your two turnstiles (each can be a
1-buffered channel, exactly like the reusable barrier's turnstiles). Follow the pseudocode
in the book faithfully — this is one of the few problems in this set where the published
solution's structure is exactly what you should implement, translated line-by-line into
Go's mutex/channel idioms.
