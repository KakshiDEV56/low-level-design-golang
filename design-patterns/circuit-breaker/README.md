# Circuit Breaker

## The concept

When a downstream dependency (another service, a database, a third-party API) starts
failing, retrying every call against it makes things worse: you pile on latency, waste
resources, and can even prevent the dependency from recovering (the "retry storm"
problem). A circuit breaker wraps calls to that dependency and, after enough
consecutive failures, starts **failing fast** — rejecting calls immediately without
even attempting them — for a cooldown period, then cautiously lets a single probe call
through to check whether the dependency has recovered.

Three states:

- **Closed** — normal operation. Calls go through. Consecutive failures are counted;
  once the count reaches a threshold, the breaker **opens**.
- **Open** — calls fail immediately with `ErrCircuitOpen`, without ever touching the
  wrapped function, until a cooldown period has elapsed. Once it has, the breaker moves
  to **half-open**.
- **Half-open** — exactly one probe call is allowed through. If it succeeds, the
  breaker closes (fully reset). If it fails, the breaker reopens (and the cooldown
  starts over from now).

## Functional requirements

1. `State` is one of `Closed`, `Open`, `HalfOpen` (given as constants below).
2. `NewCircuitBreaker(failureThreshold int, cooldown time.Duration, clock Clock)
   *CircuitBreaker` — `Clock` (given, same shape as in `../../machine-coding/rate-limiter`)
   makes the cooldown testable without real sleeps.
3. `Call(fn func() error) error`:
   - If `Closed`: call `fn`. On success, reset the consecutive-failure count to 0. On
     failure, increment it; if it reaches `failureThreshold`, transition to `Open`
     (recording the time it opened) — either way, return `fn`'s error.
   - If `Open`: if less than `cooldown` has elapsed since opening, return
     `ErrCircuitOpen` immediately WITHOUT calling `fn`. If `cooldown` has elapsed,
     transition to `HalfOpen` and proceed to the half-open behavior below for this same
     call (don't just transition and reject this call too — the whole point is this
     call becomes the probe).
   - If `HalfOpen`: call `fn`. On success, transition to `Closed` (reset failure
     count). On failure, transition back to `Open` (reset the opened-at time to now).
4. `State() State` reports the current state.

## Rubric

- [ ] While `Open` and still within the cooldown, `fn` is **never invoked** — verify
      this directly (e.g. with a counter), not just that `Call` returned an error.
- [ ] Exactly one call is allowed through in `HalfOpen` per transition — if that probe
      fails, the VERY NEXT call must be rejected again (fresh cooldown), not treated as
      still half-open.
- [ ] A single success in `Closed` state resets the failure count — 2 failures, 1
      success, 2 more failures should NOT open a breaker with `failureThreshold=3`.
- [ ] All timing is via the injected `Clock`, never `time.Now()` directly.

## Hints

- Track: `state State`, `consecutiveFailures int`, `openedAt time.Time`.
- The `Open -> HalfOpen` transition check belongs at the TOP of `Call`, before deciding
  whether to invoke `fn` — so a single `Call` invocation can both perform that
  transition AND immediately act as the probe.
