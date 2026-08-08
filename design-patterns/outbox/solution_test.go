package outbox

import (
	"errors"
	"reflect"
	"sync"
	"testing"
)

type recordingPublisher struct {
	mu        sync.Mutex
	published []Event
	failIDs   map[string]error
}

func newRecordingPublisher(failIDs map[string]error) *recordingPublisher {
	return &recordingPublisher{failIDs: failIDs}
}

func (p *recordingPublisher) Publish(e Event) error {
	if err, ok := p.failIDs[e.ID]; ok {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.published = append(p.published, e)
	return nil
}

func TestOutbox_AddAndPendingEventsPreservesOrder(t *testing.T) {
	o := NewOutbox()
	o.Add(Event{ID: "1", Payload: "a"})
	o.Add(Event{ID: "2", Payload: "b"})
	o.Add(Event{ID: "3", Payload: "c"})

	got := o.PendingEvents()
	want := []Event{{ID: "1", Payload: "a"}, {ID: "2", Payload: "b"}, {ID: "3", Payload: "c"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PendingEvents() = %v, want %v", got, want)
	}
}

func TestOutbox_MarkPublishedRemovesFromPending(t *testing.T) {
	o := NewOutbox()
	o.Add(Event{ID: "1", Payload: "a"})
	o.Add(Event{ID: "2", Payload: "b"})

	if !o.MarkPublished("1") {
		t.Fatal("MarkPublished(1) = false, want true")
	}
	got := o.PendingEvents()
	want := []Event{{ID: "2", Payload: "b"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PendingEvents() after MarkPublished(1) = %v, want %v", got, want)
	}
}

func TestOutbox_MarkPublishedUnknownReturnsFalse(t *testing.T) {
	o := NewOutbox()
	if o.MarkPublished("missing") {
		t.Fatal("MarkPublished(missing) = true, want false")
	}
}

func TestOutbox_MarkPublishedTwiceReturnsFalseSecondTime(t *testing.T) {
	o := NewOutbox()
	o.Add(Event{ID: "1", Payload: "a"})
	if !o.MarkPublished("1") {
		t.Fatal("first MarkPublished(1) = false, want true")
	}
	if o.MarkPublished("1") {
		t.Fatal("second MarkPublished(1) = true, want false")
	}
}

func TestRelay_FlushPublishesAllPending(t *testing.T) {
	o := NewOutbox()
	o.Add(Event{ID: "1", Payload: "a"})
	o.Add(Event{ID: "2", Payload: "b"})
	pub := newRecordingPublisher(nil)
	relay := NewRelay(o, pub)

	count, err := relay.Flush()
	if err != nil {
		t.Fatalf("Flush() error = %v, want nil", err)
	}
	if count != 2 {
		t.Fatalf("Flush() published count = %d, want 2", count)
	}
	if len(pub.published) != 2 {
		t.Fatalf("publisher received %d events, want 2", len(pub.published))
	}
	if len(o.PendingEvents()) != 0 {
		t.Fatalf("PendingEvents() after full Flush = %v, want empty", o.PendingEvents())
	}
}

func TestRelay_FailedPublishLeavesEventPendingForRetry(t *testing.T) {
	o := NewOutbox()
	o.Add(Event{ID: "1", Payload: "a"})
	failErr := errors.New("broker unavailable")
	pub := newRecordingPublisher(map[string]error{"1": failErr})
	relay := NewRelay(o, pub)

	count, err := relay.Flush()
	if count != 0 {
		t.Fatalf("Flush() published count = %d, want 0", count)
	}
	if !errors.Is(err, failErr) {
		t.Fatalf("Flush() error = %v, want it to wrap %v", err, failErr)
	}

	pending := o.PendingEvents()
	if len(pending) != 1 || pending[0].ID != "1" {
		t.Fatalf("PendingEvents() after a failed publish = %v, want the event still pending", pending)
	}
}

func TestRelay_OneFailureDoesNotStopOtherPublishes(t *testing.T) {
	o := NewOutbox()
	o.Add(Event{ID: "1", Payload: "a"})
	o.Add(Event{ID: "2", Payload: "b"})
	o.Add(Event{ID: "3", Payload: "c"})
	failErr := errors.New("broker unavailable")
	pub := newRecordingPublisher(map[string]error{"2": failErr})
	relay := NewRelay(o, pub)

	count, err := relay.Flush()
	if count != 2 {
		t.Fatalf("Flush() published count = %d, want 2 (events 1 and 3)", count)
	}
	if !errors.Is(err, failErr) {
		t.Fatalf("Flush() error = %v, want it to wrap %v", err, failErr)
	}

	pending := o.PendingEvents()
	if len(pending) != 1 || pending[0].ID != "2" {
		t.Fatalf("PendingEvents() = %v, want only event 2 still pending", pending)
	}
}

func TestRelay_RetryAfterFailureEventuallySucceeds(t *testing.T) {
	o := NewOutbox()
	o.Add(Event{ID: "1", Payload: "a"})
	failErr := errors.New("broker unavailable")
	pub := newRecordingPublisher(map[string]error{"1": failErr})
	relay := NewRelay(o, pub)

	relay.Flush() // fails
	pub.failIDs = nil // simulate the broker recovering

	count, err := relay.Flush() // retry
	if err != nil {
		t.Fatalf("retry Flush() error = %v, want nil", err)
	}
	if count != 1 {
		t.Fatalf("retry Flush() published count = %d, want 1", count)
	}
	if len(o.PendingEvents()) != 0 {
		t.Fatalf("PendingEvents() after a successful retry = %v, want empty", o.PendingEvents())
	}
}

func TestRelay_SuccessfullyPublishedEventNotRepublished(t *testing.T) {
	o := NewOutbox()
	o.Add(Event{ID: "1", Payload: "a"})
	pub := newRecordingPublisher(nil)
	relay := NewRelay(o, pub)

	relay.Flush()
	relay.Flush() // second flush, nothing pending

	if len(pub.published) != 1 {
		t.Fatalf("publisher received %d publishes across two Flush() calls, want 1 (no republish)", len(pub.published))
	}
}

func TestRelay_FlushWithNoPendingEvents(t *testing.T) {
	o := NewOutbox()
	pub := newRecordingPublisher(nil)
	relay := NewRelay(o, pub)

	count, err := relay.Flush()
	if count != 0 || err != nil {
		t.Fatalf("Flush() with nothing pending = (%d, %v), want (0, nil)", count, err)
	}
}
