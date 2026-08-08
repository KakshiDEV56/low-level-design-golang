package multiplex

/*
Implement a bounded-concurrency limiter: at most n goroutines may hold the "multiplex"
at once. See README.md for the full explanation — this is the semaphore-as-nightclub-
bouncer pattern from The Little Book of Semaphores, §3.5.

TODO:
 1. Add the fields you need on Multiplex below. Hint: the idiomatic Go implementation is
    a buffered channel of capacity n used as a token bucket — no sync.Mutex needed.
 2. Implement NewMultiplex(n): allocate a channel with capacity n.
 3. Implement Acquire(): take a token (block if none are available).
 4. Implement Release(): return a token (potentially unblocking a goroutine in
    Acquire()). Calling Release() without a matching prior Acquire() is not a case you
    need to handle — assume callers always pair them correctly.
*/

type Multiplex struct {
	// TODO: add your fields here
}

func NewMultiplex(n int) *Multiplex {
	// TODO: implement
	return &Multiplex{}
}

func (m *Multiplex) Acquire() {
	// TODO: implement
}

func (m *Multiplex) Release() {
	// TODO: implement
}
