# Architectural & Enterprise Patterns

This folder is deliberately **not** another pass at the GoF 23 (`awesome-low-level-design/design-patterns/`
already covers those in Go/C++/Java). These are patterns you reach for building a real
backend service — mostly from Martin Fowler's *Patterns of Enterprise Application
Architecture* territory, plus one purely idiomatic-Go pattern that doesn't come from
any enterprise textbook but shows up constantly in Go code review.

## The eight patterns

1. **`repository/`** — abstracts persistence behind a collection-like interface, so
   business logic never depends on a concrete storage mechanism.
2. **`unit-of-work/`** — coordinates changes across multiple entities as a single
   atomic commit-or-rollback, the way a real transaction would (motivating example:
   a bank transfer, where the debit and credit must both succeed or neither does).
3. **`cqrs/`** — separates the write model (commands) from the read model (queries),
   connected only through events — never by sharing mutable state directly.
4. **`circuit-breaker/`** — stops hammering a failing dependency: closed → open (fail
   fast) → half-open (probe) → closed again, driven by an injectable clock.
5. **`saga/`** — runs a sequence of steps where a later failure triggers compensating
   (undo) actions for every step that already succeeded, in reverse order.
6. **`outbox/`** — the transactional outbox: stage an event alongside a data change,
   then a separate relay publishes it later with retry-on-failure — the standard fix
   for "how do I not lose an event if the broker is down right when I need to publish."
7. **`specification/`** — encapsulates a business rule as a composable, combinable
   object (`And`/`Or`/`Not`), instead of scattering `if` conditions across the codebase.
8. **`functional-options/`** — the idiomatic Go constructor-configuration pattern
   (`New(required, WithX(...), WithY(...))`) — not an enterprise pattern, but the most
   commonly expected "do you know idiomatic Go" pattern in real interviews.

## Format

Same as `machine-coding/`: each folder has a `README.md` (the pattern, why it matters,
the exact contract to implement), a TODO-driven `solution.go`, and a `solution_test.go`
that also proves the pattern's actual payoff — not just "does it compile," but e.g. "does
swapping the repository implementation really require zero changes to the service," or
"does one failed compensation still let the others run."

```
cd design-patterns/<folder>
go test -v ./...
```
