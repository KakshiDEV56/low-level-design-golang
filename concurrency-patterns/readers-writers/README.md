# Readers-Writers (with the starvation fix)

*Based on The Little Book of Semaphores, §4.2 ("Readers-writers problem") including
§4.2.3 ("Starvation") and §4.2.4 ("No-starve readers-writers").*

## The concept

This is the problem behind any data structure, database, or file that's read and
written by concurrent threads:

1. Any number of readers can be in the critical section simultaneously.
2. Writers must have exclusive access — while a writer is in the critical section, no
   reader and no other writer may be there too.

The book calls this **categorical mutual exclusion**: a thread in the critical section
doesn't necessarily exclude other threads of the *same* category, but the presence of
one category excludes the other.

### The naive solution starves writers

The straightforward solution (readers use a `Lightswitch` — see `../lightswitch` — around
a `roomEmpty` semaphore that writers also wait on) is deadlock-free, but has a real
problem: if a writer arrives while readers are in the room, it queues on `roomEmpty` —
and if a new reader keeps arriving before the *last* current reader leaves, there is
always at least one reader in the room, so the waiting writer can starve forever. This
isn't a hang (the system is still making progress — readers keep getting served), but
it's a genuine fairness bug, and it's exactly the failure mode Go's own docs warn about:

> If a goroutine holds a RWMutex for reading and another goroutine might call Lock, no
> goroutine should expect to be able to acquire a read lock until the initial read lock
> is released. [...] a blocked Lock call excludes new readers from acquiring the lock.

Go's `sync.RWMutex` actually **does** solve this (a pending `Lock()` call blocks new
`RLock()` calls from jumping the queue) — this exercise asks you to build the mechanism
that makes that guarantee possible, so you understand *why* `RWMutex` behaves that way.

## Why this matters for a Go engineer

- `sync.RWMutex` is one of the most reached-for primitives in real Go code (any
  read-heavy cache, config store, or shared map). Understanding what it does internally
  — and specifically its writer-starvation protection — is the difference between using
  it correctly under load and being surprised when write latency spikes under heavy
  read traffic.
- The **turnstile** technique used to fix the starvation (a gate that both readers and
  writers pass through, which a waiting writer can lock to block new readers from
  queuing behind it) is a reusable idea beyond this one problem — it's the same
  mechanism used in the reusable barrier.

## The constraint to implement

Implement two variants (or one type with a mode flag — your choice, just cover both
behaviors with methods):

1. **Basic `RWLock`:** `RLock()`/`RUnlock()` allow concurrent readers; `Lock()`/`Unlock()`
   give a writer exclusive access; a writer never runs concurrently with anyone else.
2. **Starvation-free `RWLock`:** same guarantees, PLUS: once a writer is waiting, no
   *new* reader may enter ahead of it — existing readers already inside are allowed to
   finish, but newcomers must wait behind the pending writer.

## Hint

For the basic version: use a `Lightswitch` (see `../lightswitch` — literally import and
reuse it, or reimplement the same idea inline) wrapping a `roomEmpty` mutex/semaphore
that writers also acquire directly.

For the starvation-free version: add a `turnstile` that BOTH readers and writers must
pass through on the way in. A reader passes it and immediately releases it (so the next
reader isn't blocked). A writer holds the turnstile locked for the ENTIRE time it's
waiting for `roomEmpty` — which has the side effect of blocking any new readers from
even getting past the turnstile until the writer has gone through.
