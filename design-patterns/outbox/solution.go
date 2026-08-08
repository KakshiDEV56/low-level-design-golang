package outbox

/*
See README.md for the full problem statement, requirements, and rubric.

Event and Publisher are fixed — don't change their shape.
*/

type Event struct {
	ID      string
	Payload string
}

// Publisher is anything that can publish an event to a real message broker.
type Publisher interface {
	Publish(e Event) error
}

// TODO: implement Outbox.
//   - NewOutbox() *Outbox
//   - Add(e Event)
//   - PendingEvents() []Event
//   - MarkPublished(eventID string) bool
type Outbox struct {
	// TODO: add your fields here
}

func NewOutbox() *Outbox {
	// TODO: implement
	return &Outbox{}
}

func (o *Outbox) Add(e Event) {
	// TODO: implement
}

func (o *Outbox) PendingEvents() []Event {
	// TODO: implement
	return nil
}

func (o *Outbox) MarkPublished(eventID string) bool {
	// TODO: implement
	return false
}

// TODO: implement Relay.
//   - NewRelay(outbox *Outbox, publisher Publisher) *Relay
//   - Flush() (published int, err error)
type Relay struct {
	// TODO: add your fields here
}

func NewRelay(outbox *Outbox, publisher Publisher) *Relay {
	// TODO: implement
	return &Relay{}
}

func (r *Relay) Flush() (int, error) {
	// TODO: implement
	return 0, nil
}
