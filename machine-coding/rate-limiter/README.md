# Rate Limiter

**Suggested time box: 30 minutes.**

## Problem statement

Design a rate limiter that can be applied per-key (e.g., per user ID or per IP address)
to control how many requests each key is allowed to make. Implement it as a common
`Limiter` interface with **two different algorithms** behind it, so callers can swap the
strategy without changing how they call it.

## Functional requirements

1. `Limiter` interface with `Allow(key string) bool` — returns whether the request for
   `key` is allowed right now (and, if allowed, counts against that key's quota).
2. **Token Bucket**: each key's bucket **starts full** (at capacity `C` tokens) the
   first time that key is seen, then refills at `R` tokens/second (fractional refill is
   fine — don't require whole-token ticks). Each `Allow` call that succeeds consumes one
   token. If the bucket is empty, `Allow` returns `false` without blocking.
3. **Sliding Window Log**: each key may make at most `N` requests in any trailing
   window of duration `W`. `Allow` returns `true` only if fewer than `N` requests are
   recorded for that key within `(now-W, now]`.
4. Both implementations must be safe for concurrent use by multiple goroutines calling
   `Allow` for the same or different keys.
5. Time must be **injectable**, not hardcoded to `time.Now()` — so tests can control it
   precisely instead of relying on real sleeps. A `Clock` interface is provided for you
   (see below); don't call `time.Now()` directly anywhere in your implementation.

## Rubric

- [ ] Both algorithms implement the same `Limiter` interface and are interchangeable.
- [ ] Per-key state doesn't leak between unrelated keys (a hammered key doesn't affect
      another key's quota).
- [ ] Concurrent `Allow` calls for the same key never over-admit (race-free counting).
- [ ] No direct `time.Now()` calls inside your limiter logic — everything goes through
      the injected `Clock`.
- [ ] Token bucket refill math is continuous (based on elapsed time), not "reset every
      tick."

## Hints

- `Clock` (below, already provided — don't modify) is a one-method interface so tests
  can inject a fake clock that only advances when told to.
- Token bucket: store, per key, `tokens float64` and `lastRefill time.Time`. On each
  `Allow`, compute `elapsed := clock.Now().Sub(lastRefill)`, refill
  `tokens = min(capacity, tokens + elapsed.Seconds()*refillRate)`, update
  `lastRefill = clock.Now()`, then check `tokens >= 1`.
- Sliding window log: store, per key, a slice (or ring buffer) of request timestamps.
  On each `Allow`, drop timestamps older than `now - window`, then check
  `len(remaining) < maxRequests`.
- A `sync.Mutex` per limiter (or a `sync.Map` keyed by request key, each with its own
  mutex) is enough — no need for anything fancier at this scale.
