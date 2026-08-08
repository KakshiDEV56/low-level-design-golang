# Transactional Outbox

## The concept

You need to update your own data and publish an event about it — say, saving an order
and publishing `OrderPlaced` to a message broker. Doing both as one atomic operation is
hard: your database and your broker are different systems, and there's no distributed
transaction spanning both in most real setups. If you save the order and the broker
call fails, the event is lost forever. If you publish the event first and the database
write fails, you've told the world about an order that doesn't exist.

The **transactional outbox** pattern sidesteps this: instead of publishing directly,
you write the event into an **outbox** — a table/store in the SAME local transaction as
your actual data change, so it's atomic almost for free (it's just one more local
write). A separate, independent **relay** process then reads pending outbox entries and
publishes them to the real broker, marking each one as published only after a
successful publish — and leaving anything that failed in the outbox to retry later.
This gives you reliable **at-least-once** delivery without needing a distributed
transaction.

## Functional requirements

1. `Event{ID, Payload string}` (given).
2. `Outbox`:
   - `Add(e Event)` stages an event (simulating "written in the same local transaction
     as the business data change").
   - `PendingEvents() []Event` returns every event not yet marked published, **in the
     order they were added**.
   - `MarkPublished(eventID string) bool` marks one event as published (it will no
     longer appear in `PendingEvents`). Returns `false` if no such pending event
     exists.
3. `Publisher` is a fixed interface: `Publish(e Event) error`.
4. `Relay`:
   - `NewRelay(outbox *Outbox, publisher Publisher) *Relay`
   - `Flush() (published int, err error)` — attempts to publish every currently
     pending event, in order. For each one that publishes successfully, mark it
     published in the outbox. A publish failure for one event must NOT stop the relay
     from attempting the rest — leave failed events in the outbox (don't mark them
     published) so a later `Flush()` call retries them. Returns the count of events
     successfully published during this call, and a combined error (via
     `errors.Join`) if any failed — `nil` error if all pending events published
     successfully (or there were none).

## Rubric

- [ ] A failed publish does not stop the relay from attempting later pending events in
      the same `Flush()` call.
- [ ] A failed event remains in `PendingEvents()` after `Flush()` — i.e. is retried on
      the next `Flush()` call.
- [ ] A successfully published event does NOT get published again on a later
      `Flush()` call.
- [ ] `PendingEvents()` preserves insertion order.

## Hints

- `Outbox` needs an ordered collection (a slice is fine) plus a way to mark entries
  published without disturbing the order of the rest — e.g. a slice of `{Event; published
  bool}`, filtering on `published == false` for `PendingEvents`.
- `Relay.Flush` is a straightforward loop over `outbox.PendingEvents()`, calling
  `publisher.Publish` and `outbox.MarkPublished` on success, collecting errors with
  `errors.Join` on failure — very similar in shape to the `notification-system`
  machine-coding problem's partial-failure handling, if you want a reference for the
  "don't stop on one failure" loop structure.
