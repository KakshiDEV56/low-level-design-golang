# CQRS (Command Query Responsibility Segregation)

## The concept

CQRS splits a system into two sides that never share mutable state directly:

- The **command side** handles writes and enforces business rules (e.g. "stock can't
  go negative"). It doesn't know or care how reads are served.
- The **query side** handles reads, with a data shape optimized for the queries it
  needs to answer — which may look nothing like the command side's internal
  representation.

The two sides are connected **only through events**: every successful command emits an
event describing what happened, and the query side updates its own read model by
*applying* those events (a "projection"). This means the query side is always
technically a little behind the command side (eventual consistency) — but it also
means either side can evolve, scale, or even be rebuilt independently, since neither
one reaches into the other's internals.

## Functional requirements — a simple inventory system

1. `StockEvent{ProductID string; Delta int}` (given) — `Delta` is positive for stock
   added, negative for stock removed.
2. `CommandHandler` (the write side):
   - `AddStock(productID string, qty int) error` — `qty` must be `> 0`; emits a
     `StockEvent{ProductID: productID, Delta: qty}` to every subscriber.
   - `RemoveStock(productID string, qty int) error` — `qty` must be `> 0`; fails with
     `ErrInsufficientStock` if the command side's own tracked stock for that product is
     less than `qty` (the command side must track enough state itself to enforce this
     rule — it is the source of truth). On success, emits `StockEvent{ProductID:
     productID, Delta: -qty}`.
   - `Subscribe(handler func(StockEvent))` registers a function to be called with every
     event emitted by a successful command.
3. `QueryHandler` (the read side) — has **no direct access to `CommandHandler`'s
   state**, is built purely by having events applied to it:
   - `Apply(e StockEvent)` — the projector: updates the read model based on one event.
   - `GetStock(productID string) int` — current known stock for a product (0 if never
     seen).
   - `TotalStock() int` — sum of all products' stock.
   - `LowStockProducts(threshold int) []string` — product IDs with stock `<
     threshold`, sorted alphabetically.

## Rubric

- [ ] `QueryHandler` never reads `CommandHandler`'s fields directly — the only way
      data reaches it is through `Apply`.
- [ ] Wiring `cmdHandler.Subscribe(queryHandler.Apply)` is enough to keep them in
      sync — no other coupling between the two types.
- [ ] `RemoveStock` correctly rejects an over-removal using the command side's OWN
      state, not the query side's (which may not even be wired up).
- [ ] `QueryHandler` before any events are applied reports zero stock for everything,
      not a panic or a nil-map crash.

## Hints

- `CommandHandler` needs its own `map[string]int` tracking current stock (this is
  its private write-side state, separate from whatever `QueryHandler` ends up with),
  plus a slice of subscriber functions to call on every successful command.
- `QueryHandler` needs its own separate `map[string]int` — don't share a map between
  the two handlers, even though that would "work" for this in-memory exercise; sharing
  state defeats the entire point of the pattern (and the tests specifically check that
  swapping in a query handler with a completely different internal shape, e.g. built
  from a slice instead of a map, still works via `Apply`).
