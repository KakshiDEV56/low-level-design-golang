package reusablebarrier

/*
Implement a barrier that can be reused across many rounds by n goroutines running in a
loop, without any goroutine ever starting round k+1 before every goroutine has finished
round k ("lapping"). See README.md for the full explanation — this is the two-phase
turnstile barrier from The Little Book of Semaphores, §3.6-3.7.

TODO:
 1. Add the fields you need on CyclicBarrier below. Hint: you'll want
      - n (the number of participants),
      - a count of how many have arrived in the current phase, protected by a
        sync.Mutex,
      - TWO turnstiles (each can be a channel of capacity 1, or built on your own
        counting semaphore from ../counting-semaphore if you want to reuse it) —
        the first one gates entry, the second one gates the "everyone has definitely
        left the first turnstile" phase, which is what makes looping back around safe.
 2. Implement NewCyclicBarrier(n): initialize count to 0. The first turnstile should
    start "locked" (so early arrivals block); the second should start "open" (since
    nobody needs to wait on it until the first phase completes).
 3. Implement Await() using the two-phase pattern:
      - Lock the counter mutex, increment count. If count == n: lock the second
        turnstile (pull one value out of it, so it's now "empty"/blocking) and unlock
        the first turnstile (push a value in). Unlock the counter mutex.
      - Wait on the first turnstile (pull a value out), then immediately signal it
        again (push a value back in) — this is what lets exactly one goroutine at a
        time pass through, one after another, until all n have gone through.
      - This is the "critical point" / round boundary — participants are now released
        from phase 1.
      - Lock the counter mutex, decrement count. If count == 0: lock the first
        turnstile again (pull a value out, re-locking it for the NEXT round) and unlock
        the second turnstile (push a value in). Unlock the counter mutex.
      - Wait on the second turnstile (pull a value out), then immediately signal it
        again (push it back in).
    After Await() returns, it must be safe for the goroutine to immediately call
    Await() again for the next round.
*/

type CyclicBarrier struct {
	// TODO: add your fields here
}

func NewCyclicBarrier(n int) *CyclicBarrier {
	// TODO: implement
	return &CyclicBarrier{}
}

func (b *CyclicBarrier) Await() {
	// TODO: implement
}
