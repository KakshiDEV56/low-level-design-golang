# Multiplex

*Based on The Little Book of Semaphores, §3.5 ("Multiplex").*

## The concept

A mutex enforces "at most **one** goroutine in the critical section." A multiplex
generalizes that to "at most **n** goroutines in the critical section at once."

From the book:

> In real life, the multiplex problem occurs at busy nightclubs where there is a maximum
> number of people allowed in the building at a time... a bouncer usually enforces the
> synchronization constraint by keeping track of the number of people inside and barring
> arrivals when the room is at capacity. Then, whenever one person leaves, another is
> allowed to enter.

The book's solution is almost embarrassingly simple: initialize a semaphore to `n`
instead of `1`. Each arrival calls `Wait()` (decrementing the "tokens" available); each
departure calls `Signal()` (returning a token). If the value is zero, the next arrival
blocks until someone leaves.

## Why this matters for a Go engineer

This is, by a wide margin, **the most common real-world use of a semaphore-style
primitive in production Go code**:

- Limiting how many goroutines hit a downstream API or database concurrently, so you
  don't overwhelm it.
- Capping the size of a worker pool.
- Bounding concurrent file descriptors, DB connections, or in-flight HTTP requests.

The idiomatic Go implementation of this exact pattern is **a buffered channel used as a
token bucket** — `make(chan struct{}, n)`. Sending into the channel is "acquire a
token" (blocks once the buffer is full, i.e. once `n` tokens are checked out); receiving
is "release a token." This is worth knowing cold: it comes up constantly, both in
interviews and in real code review ("how would you limit concurrency here?").

## The constraint to implement

- `NewMultiplex(n int)` creates a limiter allowing at most `n` concurrent holders.
- `Acquire()` blocks until fewer than `n` goroutines currently hold a token, then takes
  one.
- `Release()` returns a token, potentially unblocking one goroutine waiting in
  `Acquire()`.
- At all times, the number of goroutines that have called `Acquire()` without a matching
  `Release()` must never exceed `n`.

## Hint

Use a buffered channel of capacity `n` as your token bucket. `Acquire()` sends a value
into the channel (which blocks once it's full); `Release()` receives a value from the
channel (freeing a slot). You don't need a `sync.Mutex` at all — the channel itself
provides all the synchronization you need.
