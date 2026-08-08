package countingsemaphore

/*
Build a counting semaphore from scratch (no golang.org/x/sync/semaphore, no buffered
channels — the point is to understand the primitive at the level of a mutex + condition
variable, the way it's actually implemented under the hood).

See README.md for the full explanation. In short, from The Little Book of Semaphores,
Chapter 2:

A semaphore is like an integer, with three differences:
 1. You initialize it to any integer, but afterward the only operations are increment
    (Signal) and decrement (Wait). You cannot read the current value.
 2. When a goroutine decrements the semaphore, if the result is negative, the goroutine
    blocks until another goroutine increments the semaphore.
 3. When a goroutine increments the semaphore, if there are other goroutines waiting,
    exactly one of them gets unblocked.

TODO:
 1. Add the fields you need on Semaphore below. Hint: a plain int counter, protected by
    a sync.Mutex, plus a sync.Cond built on that same mutex (sync.NewCond(&s.mu)) so
    that Wait() can block and Signal() can wake a blocked goroutine.
 2. Implement NewSemaphore(n): initialize the counter to n and construct the Cond.
 3. Implement Wait(): lock the mutex, decrement the counter, and while the counter is
    negative, call cond.Wait() (which atomically unlocks the mutex while blocked, and
    re-locks it before returning). Use a `for` loop around the condition check, not a
    single `if` — this is the standard defensive pattern for condition variables.
    Unlock before returning.
 4. Implement Signal(): lock the mutex, increment the counter, call cond.Signal() to
    wake (at most) one waiting goroutine, then unlock.
*/

type Semaphore struct {
	// TODO: add your fields here
}

func NewSemaphore(n int) *Semaphore {
	// TODO: implement
	return &Semaphore{}
}

func (s *Semaphore) Wait() {
	// TODO: implement
}

func (s *Semaphore) Signal() {
	// TODO: implement
}
