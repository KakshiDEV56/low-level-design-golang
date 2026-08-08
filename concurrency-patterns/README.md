# Concurrency Patterns — from *The Little Book of Semaphores*

This folder is a curated set of exercises based on **Allen B. Downey's *The Little Book
of Semaphores*** (`/home/ani/Documents/LLD/LittleBookOfSemaphores.pdf`), a free textbook
that builds up synchronization primitives and classic concurrency problems from first
principles using semaphores.

Go doesn't have a built-in semaphore type, so instead of translating the book's
pseudocode line-for-line, each problem here is reframed as **the idiomatic Go way to
solve the same underlying synchronization constraint** — using `sync.Mutex`,
`sync.Cond`, `sync.WaitGroup`, buffered channels as counting tokens, or `sync.Locker`,
depending on which primitive is the natural fit. The point isn't to memorize semaphore
syntax; it's to understand the *patterns* (mutual exclusion, bounded concurrency,
barriers, categorical exclusion, starvation-freedom, dedicated coordinators) well enough
to recognize and apply them in real Go code — worker pools, connection limiters, caches,
batch pipelines, and so on.

## What's NOT duplicated here

A few of the book's foundational patterns are already covered by problems in
[`../concurrency/`](../concurrency/), so they aren't repeated:

- **Signaling / Rendezvous** (book §3.1, §3.3) → `../concurrency/print-in-order`
- **Basic mutual exclusion** (book §3.4) → covered throughout `../concurrency/`
  (e.g. `Fizz-Buzz-Multithreaded`, `Print-Zero-Even-Odd`)
- **Producer-consumer with a finite buffer** (book §4.1.4) →
  `../concurrency/design-bounded-blocking-queue`
- **Dining philosophers** (book §4.4) → `../concurrency/The-Dining-Philosophers`
- **Building H2O / barrier-style rendezvous of mixed thread types** (book §5.6) →
  `../concurrency/Building-H2O`

## What's here, and the order to do them in

Each stage builds on intuition from the one before it.

1. **`counting-semaphore/`** — Build a real counting semaphore from scratch (book
   ch. 2 + §9.2). Everything else in synchronization theory is describable in terms of
   this one primitive, so understanding *how* to build one (not just how to use one)
   is the foundation for all the rest.
2. **`multiplex/`** — Bound the number of goroutines allowed into a critical section at
   once (book §3.5). This is the single most common real-world use of a semaphore in Go:
   worker pools, concurrent-request limiters, connection pool caps.
3. **`reusable-barrier/`** — A barrier that can be used across multiple rounds in a loop,
   without threads "lapping" each other (book §3.7, the two-phase turnstile solution).
   `sync.WaitGroup` is single-use; this teaches you what to reach for when you need a
   barrier that resets.
4. **`lightswitch/`** — The "first one in locks the room, last one out unlocks it"
   reusable component (book §4.2, page 70). This is the reference-counting pattern that
   readers-writers is built out of, and it's a genuinely reusable utility on its own.
5. **`readers-writers/`** — Full reader/writer categorical mutual exclusion, including
   the writer-starvation problem and its fix (book §4.2). This is what `sync.RWMutex`
   does internally — implementing it yourself is the best way to understand its
   fairness tradeoffs (and why Go's `RWMutex` docs warn about writer starvation under
   heavy read load).
6. **`no-starve-mutex/`** — A strictly FIFO-fair mutex, so no goroutine can be
   perpetually skipped by unlucky scheduling (book §4.3). Go's `sync.Mutex` has a
   "starvation mode" fallback for exactly this reason — this exercise builds the
   ticket-lock idea that motivates it.
7. **`cigarette-smokers/`** — The "dedicated coordinator" pattern: one goroutine (an
   "agent") mediates between resources that can't safely coordinate directly with each
   other (book §4.5). This generalizes to real coordination problems where several
   independent resources need to be assembled before work can proceed.

## How to work through each one

Every folder is its own Go module with:

- **`README.md`** — the concept, why it matters, and the synchronization constraint to
  satisfy, in the book's own terms translated to Go.
- **`solution.go`** — a `// TODO`-driven template: struct with no fields, method stubs,
  and hints toward the right primitive — no implementation given away.
- **`solution_test.go`** — a thorough test suite (correctness + concurrency stress,
  runnable with `-race`).

```
cd concurrency-patterns/<folder>
go test -v -race ./...
```
