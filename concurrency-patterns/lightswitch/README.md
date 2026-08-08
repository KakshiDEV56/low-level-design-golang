# Lightswitch

*Based on The Little Book of Semaphores, §4.2 (page 70, inside "Readers-writers
problem").*

## The concept

While solving readers-writers, the book notices a recurring shape: "the first thread
into a section locks a semaphore (or queues) and the last one out unlocks it." It's
common enough that the book names it and wraps it in a reusable component, called a
**Lightswitch** — "by analogy with the pattern where the first person into a room locks
the light switch on, and the last one out turns it off."

The book's own definition (translated from its pseudocode, unchanged in spirit):

```
class Lightswitch:
    def __init__(self):
        self.counter = 0
        self.mutex = Semaphore(1)

    def lock(self, semaphore):
        self.mutex.wait()
        self.counter += 1
        if self.counter == 1:
            semaphore.wait()      # first one in locks the resource
        self.mutex.signal()

    def unlock(self, semaphore):
        self.mutex.wait()
        self.counter -= 1
        if self.counter == 0:
            semaphore.signal()    # last one out unlocks the resource
        self.mutex.signal()
```

`lock` takes a semaphore (or, in Go, anything satisfying `sync.Locker`) that it will
check and possibly hold. If it's already held, the calling goroutine just blocks on the
Lightswitch's own internal mutex, not on the underlying resource directly. When the
counter reaches zero, the resource is genuinely released.

This is exactly the mechanism readers use in readers-writers: many readers can be
"inside," but the underlying room only actually gets locked once (by whichever reader
happens to be first), and only actually gets unlocked once (by whichever reader happens
to be last).

## Why this matters for a Go engineer

This is a **reference-counted gate**, and that shape shows up constantly:

- It's the mechanism underneath `sync.RWMutex` — many readers share one "read lock," but
  the underlying exclusive resource is only really contended by the first-in/last-out
  transition.
- It generalizes beyond readers-writers to anything with an "open while at least one
  user is active" resource: lazily opening a connection/file/subscription on first use,
  and closing it once every user is done — a pattern you'll recognize from connection
  pooling and resource lifecycle management.
- Go's `sync.Locker` interface (`Lock()` / `Unlock()`) is exactly the right shape to
  parametrize this over — it's a good exercise in writing to an interface instead of a
  concrete type.

## The constraint to implement

- `NewLightswitch()` creates a Lightswitch with an internal count of zero.
- `Lock(room sync.Locker)` — the calling goroutine "enters." If it is the *first*
  goroutine currently inside (count goes from 0 to 1), it must call `room.Lock()` on the
  underlying resource. Every other goroutine that enters while count > 0 does NOT touch
  `room` at all.
- `Unlock(room sync.Locker)` — the calling goroutine "leaves." If it is the *last*
  goroutine currently inside (count goes from 1 to 0), it must call `room.Unlock()`.
  Every other goroutine leaving while count stays > 0 does NOT touch `room`.
- The internal counter itself must be safe under concurrent `Lock`/`Unlock` calls from
  many goroutines (this is what the Lightswitch's own internal mutex protects).

## Hint

One `sync.Mutex` protecting a plain `int` counter is all the Lightswitch itself needs.
The tricky part isn't the counter — it's remembering that `room.Lock()`/`room.Unlock()`
should be called *exactly once* per "occupied" transition, not once per goroutine.
