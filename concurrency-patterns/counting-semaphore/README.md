# Counting Semaphore

*Based on The Little Book of Semaphores, Chapter 2 ("Semaphores") and §9.2
("Make your own semaphores").*

## The concept

From the book's definition (Chapter 2), a semaphore is like an integer, with three
differences:

1. When you create it, you initialize its value to any integer. After that, the only
   operations allowed are **increment** (`Signal`, aka `V`) and **decrement**
   (`Wait`, aka `P`). You cannot read the current value.
2. When a thread decrements the semaphore, if the result is **negative**, the thread
   blocks and cannot continue until another thread increments the semaphore.
3. When a thread increments the semaphore, if there are other threads waiting, **one**
   of them gets unblocked.

That's the entire contract. Everything else in this folder (and much of the book) is
built on top of just this.

Go doesn't ship a semaphore in its standard library (there's `golang.org/x/sync/semaphore`
in the extended ecosystem, but nothing in `sync`). This exercise is about building one
yourself, because understanding *how* a semaphore is implemented — not just how to call
`Wait`/`Signal` — is what lets you recognize the pattern when you need it and reason
about correctness when something goes wrong.

## Why this matters for a Go engineer

- It's a genuinely common interview question ("implement a semaphore using a mutex and
  condition variable") precisely because it tests whether you understand blocking,
  waking, and the invariant a synchronization primitive has to maintain under
  concurrent access.
- `sync.Cond` is one of the least-used but most foundational primitives in Go's `sync`
  package — most Go engineers have never had a reason to reach for it. Building a
  semaphore is the natural way to actually learn it (`Cond.Wait()`, `Cond.Signal()`,
  `Cond.Broadcast()`, and the classic "wait in a `for` loop, not an `if`" rule for
  spurious/lost wakeups).
- Every later problem in this folder (`multiplex`, `readers-writers`, ...) is describable
  as "take a semaphore and use it this way" — so this is the base case.

## The constraint to implement

- `NewSemaphore(n int)` creates a semaphore initialized to `n`.
- `Wait()` decrements the internal counter. If the counter goes negative, the calling
  goroutine blocks until some other goroutine calls `Signal()`.
- `Signal()` increments the internal counter. If any goroutine is currently blocked in
  `Wait()`, exactly one of them is unblocked.
- Multiple `Signal()` calls that arrive while several goroutines are waiting must wake
  them one at a time as each `Signal()` arrives — not wake everyone on every signal.

## Hint

Build it with `sync.Mutex` + `sync.Cond` around a plain `int` counter: `Wait()` locks the
mutex, decrements the counter, and while the counter is negative, calls
`cond.Wait()` (which atomically unlocks the mutex and blocks — re-locking it upon
waking). `Signal()` locks the mutex, increments the counter, calls `cond.Signal()`, and
unlocks. Remember: always re-check the wait condition in a loop after waking, since
`Cond.Wait()` can wake up without a matching `Signal()` in some implementations/patterns
— checking `for s.value < 0` instead of `if s.value < 0` is the standard defensive habit.
