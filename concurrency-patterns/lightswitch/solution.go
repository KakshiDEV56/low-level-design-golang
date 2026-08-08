package lightswitch

import "sync"

/*
Implement the Lightswitch pattern: the first goroutine to "enter" locks an underlying
resource, and the last goroutine to "leave" unlocks it — regardless of how many
goroutines pass through in between. See README.md for the full explanation and the
book's original pseudocode (The Little Book of Semaphores, §4.2, page 70).

TODO:
 1. Add the fields you need on Lightswitch below. Hint: just a counter (int) and a
    sync.Mutex to protect it — that's genuinely all the Lightswitch itself needs.
 2. Implement NewLightswitch: counter starts at 0.
 3. Implement Lock(room sync.Locker): increment the counter under the Lightswitch's own
    mutex; if the counter just became 1 (this goroutine is the first one in), call
    room.Lock().
 4. Implement Unlock(room sync.Locker): decrement the counter under the Lightswitch's
    own mutex; if the counter just became 0 (this goroutine is the last one out), call
    room.Unlock().
*/

type Lightswitch struct {
	// TODO: add your fields here
}

func NewLightswitch() *Lightswitch {
	// TODO: implement
	return &Lightswitch{}
}

func (l *Lightswitch) Lock(room sync.Locker) {
	// TODO: implement
}

func (l *Lightswitch) Unlock(room sync.Locker) {
	// TODO: implement
}
