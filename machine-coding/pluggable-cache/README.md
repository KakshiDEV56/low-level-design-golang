# Pluggable Cache (Strategy Pattern)

**Suggested time box: 40 minutes.**

## Problem statement

You've already implemented `lru-cache` and `lfu-cache` as standalone LeetCode-style
problems (`awesome-low-level-design/oops/`). This is a different, more "real system"
framing of the same idea: build a **generic cache whose eviction algorithm is a
swappable strategy**, not a hardcoded one — the way you'd actually design a caching
library meant to support more than one policy.

## Functional requirements

1. `EvictionPolicy[K comparable]` is a fixed interface (given below — don't change its
   shape) with three methods: `KeyAccessed(key K)`, `KeyRemoved(key K)`, and
   `Evict() (K, bool)`.
2. `Cache[K comparable, V any]` wraps a `map[K]V` plus an `EvictionPolicy[K]`:
   - `NewCache[K, V](capacity int, policy EvictionPolicy[K]) *Cache[K, V]`
   - `Get(key K) (V, bool)` — on a hit, notify the policy via `KeyAccessed` and return
     the value; on a miss, return the zero value and `false`.
   - `Put(key K, value V)` — if the key already exists, update its value and notify
     `KeyAccessed`. If it's new and the cache is at capacity, ask the policy for a
     victim via `Evict()` and remove it first. Then insert and notify `KeyAccessed`.
   - `Remove(key K)` — delete the key if present and notify the policy via
     `KeyRemoved`.
   - `Len() int` — current number of entries.
3. Implement **two** concrete policies against that same interface:
   - `NewLRUPolicy[K]() EvictionPolicy[K]` — evicts the least-recently-accessed key.
   - `NewFIFOPolicy[K]() EvictionPolicy[K]` — evicts the key that was inserted first,
     regardless of how often it's accessed afterward (FIFO does NOT move a key to the
     back on access — only `KeyAccessed` calls that correspond to a genuinely NEW key
     should affect eviction order; re-accessing an existing key does nothing to FIFO's
     ordering).
4. `Cache` itself must not know or care which policy it's using — swapping
   `NewLRUPolicy[string]()` for `NewFIFOPolicy[string]()` in the constructor call is the
   only change needed to change eviction behavior.

**Stretch goal (optional, not covered by the provided tests):** implement
`NewLFUPolicy[K]()` too, evicting the least-frequently-used key (ties broken by least
recently used among the tied keys) — this is the hardest of the three and worth doing
if you finish early.

## Rubric

- [ ] `Cache`'s `Get`/`Put`/`Remove` logic is identical regardless of which policy is
      injected — no policy-specific branching inside `Cache`.
- [ ] `LRUPolicy` and `FIFOPolicy` are genuinely interchangeable (both satisfy
      `EvictionPolicy[K]` and can be swapped with one constructor-call change).
- [ ] `KeyAccessed` on an EXISTING key correctly updates LRU order but must NOT affect
      FIFO order.
- [ ] `Evict()` on an empty policy returns `ok == false`, not a zero-value key mistaken
      for a real one.
- [ ] Generic type parameters are used correctly (no `interface{}`/`any` casting hacks
      to work around them).

## Hints

- `LRUPolicy`: a doubly linked list (or a `container/list.List` from the standard
  library — perfectly fine to use here) ordered by recency, plus a
  `map[K]*list.Element` for O(1) lookup. `KeyAccessed` moves the element to the front;
  `Evict` pops from the back.
- `FIFOPolicy`: much simpler — a plain slice or `container/list.List` of keys in
  insertion order. `KeyAccessed` should only append a NEW key (check a
  `map[K]struct{}` first to know if it's already tracked); it's a no-op for a key
  that's already there. `Evict` pops from the front.
- Don't forget: `Cache.Put` on an *existing* key must not trigger eviction — it's just
  an update.
