package nostarvemutex

/*
Implement a strictly fair (FIFO, no-starve) mutex: if goroutine X calls Lock() before
goroutine Y does, X is guaranteed to enter the critical section before Y, however many
times Y loops around and calls Lock() again. See README.md for the full explanation —
this is Morris's algorithm from The Little Book of Semaphores, §4.3.

TODO:
 1. Add the fields you need on NoStarveMutex below. Following Morris's algorithm
    directly:
      - room1, room2 int        // counts of goroutines in each waiting "room"
      - mutex  a sync.Mutex     // protects room1 (and briefly room2)
      - t1, t2 turnstiles       // each can be a channel of capacity 1; t1 starts
                                //   "open" (a value present), t2 starts "closed" (empty)
 2. Implement NewNoStarveMutex: room1 = room2 = 0, t1 pre-loaded with one value (open),
    t2 empty (closed).
 3. Implement Lock(), following the book's pseudocode structure exactly:
      mutex.Lock(); room1++; mutex.Unlock()

      <t1.wait()>                 // pull a value from t1
          room2++
          mutex.Lock(); room1--
          if room1 == 0:
              mutex.Unlock()
              <t2.signal()>       // push a value into t2 (opens it for this "wave")
          else:
              mutex.Unlock()
              <t1.signal()>       // push a value back into t1 (let the next one through)

      <t2.wait()>                 // pull a value from t2 — this is what actually grants
                                   // entry to the critical section
          room2--
          # (critical section begins; Unlock() does the rest below)
 4. Implement Unlock(), continuing the pseudocode:
      if room2 == 0:
          <t1.signal()>           // reopen t1 for the next wave
      else:
          <t2.signal()>           // let the next goroutine in THIS wave proceed
    (room2's decrement happens at the end of Lock(), right before entering the critical
    section, per the book's structure — track it however is cleanest for your Go
    implementation, as long as the net effect matches: exactly one of t1/t2 is signaled
    when a goroutine leaves, based on whether any goroutines from the current wave are
    still waiting in Room 2.)
*/

type NoStarveMutex struct {
	// TODO: add your fields here
}

func NewNoStarveMutex() *NoStarveMutex {
	// TODO: implement
	return &NoStarveMutex{}
}

func (m *NoStarveMutex) Lock() {
	// TODO: implement
}

func (m *NoStarveMutex) Unlock() {
	// TODO: implement
}
