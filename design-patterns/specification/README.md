# Specification Pattern

## The concept

Business rules ("is this product eligible for the discount," "is this order overdue
for shipping") tend to accumulate as scattered `if` conditions across a codebase —
duplicated in a filter here, a validator there, slightly out of sync each time someone
edits one copy and forgets the other. The Specification pattern encapsulates a rule as
its own object with a single `IsSatisfiedBy(candidate) bool` method, so it can be:

- written once, and reused everywhere it applies,
- unit tested completely independently of where it's used, and
- **combined** with other specifications via `And`, `Or`, `Not` to build up complex
  rules out of simple, individually-understandable pieces.

## Functional requirements

1. `Specification[T any]` is a fixed generic interface: `IsSatisfiedBy(candidate T)
   bool`.
2. `And[T](left, right Specification[T]) Specification[T]` — satisfied only if both
   are.
3. `Or[T](left, right Specification[T]) Specification[T]` — satisfied if either is.
4. `Not[T](spec Specification[T]) Specification[T]` — satisfied iff `spec` isn't.
5. `Filter[T](items []T, spec Specification[T]) []T` — returns the elements of `items`
   that satisfy `spec`, preserving their original relative order.

## Rubric

- [ ] `And`/`Or`/`Not` are themselves valid `Specification[T]` implementations — they
      must be composable with EACH OTHER (an `And` of an `Or` of a `Not`, arbitrarily
      nested), not just usable as one-level wrappers around "leaf" specifications.
- [ ] None of `And`/`Or`/`Not`/`Filter` know anything about any specific domain (no
      product/order-specific logic anywhere in `solution.go`) — domain rules are
      defined entirely in the test file, using only the generic combinators.
- [ ] `Filter` preserves order and doesn't mutate the input slice.

## Hints

- The cleanest way to implement `And`/`Or`/`Not` in Go without a full class hierarchy
  is with small adapter types wrapping a closure, e.g.:
  ```go
  type specFunc[T any] func(T) bool
  func (f specFunc[T]) IsSatisfiedBy(candidate T) bool { return f(candidate) }
  ```
  Then `And(left, right)` can just return
  `specFunc[T](func(c T) bool { return left.IsSatisfiedBy(c) && right.IsSatisfiedBy(c) })`
  — this one adapter type is all three combinators need.
