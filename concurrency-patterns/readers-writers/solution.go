package readerswriters

/*
Implement readers-writers mutual exclusion, in two versions. See README.md for the full
explanation — this is The Little Book of Semaphores §4.2, including the starvation
problem in §4.2.3 and its fix in §4.2.4.

TODO (RWLock — the basic version):
 1. Add the fields you need on RWLock below. Hint: a reader counter (protected by its
    own mutex — this is exactly the Lightswitch pattern from ../lightswitch, feel free
    to reimplement it inline here) plus a "roomEmpty" mutex that the first reader locks
    and the last reader unlocks, and that writers lock/unlock directly.
 2. Implement NewRWLock.
 3. Implement RLock/RUnlock: increment/decrement the reader count; the FIRST reader in
    locks roomEmpty, the LAST reader out unlocks it. Readers in between do neither.
 4. Implement Lock/Unlock for writers: just lock/unlock roomEmpty directly. Since
    readers only ever hold roomEmpty while readers > 0, and a writer holds it for its
    entire critical section, this guarantees writers never overlap with anyone.

TODO (FairRWLock — the starvation-free version):
 5. Add the fields you need on FairRWLock below. Same as RWLock, PLUS a "turnstile"
    mutex that BOTH readers and writers must pass through before doing anything else.
 6. Implement NewFairRWLock.
 7. Implement RLock: lock the turnstile, then IMMEDIATELY unlock it (so it's just a
    quick pass-through for readers, not something they hold), then do the same
    reader-count/roomEmpty dance as RWLock.
 8. Implement RUnlock: same as RWLock's RUnlock.
 9. Implement Lock: lock the turnstile and KEEP IT LOCKED while also locking roomEmpty
    (so no new reader can even get past the turnstile while this writer is waiting for
    or holding the room), then unlock the turnstile once roomEmpty is acquired.
 10. Implement Unlock: unlock roomEmpty.
*/

type RWLock struct {
	// TODO: add your fields here
}

func NewRWLock() *RWLock {
	// TODO: implement
	return &RWLock{}
}

func (l *RWLock) RLock() {
	// TODO: implement
}

func (l *RWLock) RUnlock() {
	// TODO: implement
}

func (l *RWLock) Lock() {
	// TODO: implement
}

func (l *RWLock) Unlock() {
	// TODO: implement
}

type FairRWLock struct {
	// TODO: add your fields here
}

func NewFairRWLock() *FairRWLock {
	// TODO: implement
	return &FairRWLock{}
}

func (l *FairRWLock) RLock() {
	// TODO: implement
}

func (l *FairRWLock) RUnlock() {
	// TODO: implement
}

func (l *FairRWLock) Lock() {
	// TODO: implement
}

func (l *FairRWLock) Unlock() {
	// TODO: implement
}
