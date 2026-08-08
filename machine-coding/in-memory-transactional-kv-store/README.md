# In-Memory Transactional Key-Value Store

**Suggested time box: 45–60 minutes.**

## Problem statement

Design a key-value store that supports **nested transactions**: `Begin()` starts one,
`Commit()` folds its changes into the enclosing scope (or into permanent storage if it
was the outermost transaction), and `Rollback()` discards everything done since the
matching `Begin()` — without touching anything from an enclosing transaction.

## Functional requirements

1. `Get(key string) (string, bool)`, `Set(key, value string)`, `Delete(key string)` —
   basic key-value operations. Outside of any transaction, these act directly on
   permanent storage.
2. `Begin()` opens a new transaction, nested inside whatever transaction (if any) is
   currently open.
3. `Commit() bool` closes the **innermost** open transaction and merges its changes
   into the next level out (the enclosing transaction, or permanent storage if this was
   the outermost one). Returns `false` (no-op) if there's no open transaction.
4. `Rollback() bool` closes the innermost open transaction and **discards** every
   change made since its `Begin()` — the enclosing scope is completely unaffected.
   Returns `false` (no-op) if there's no open transaction.
5. `InTransaction() bool` reports whether at least one transaction is currently open.
6. **Visibility**: while inside a transaction, `Get` must see that transaction's own
   uncommitted writes (including deletes) first, then fall through to whatever the
   enclosing transaction (or permanent storage) has, for keys the current transaction
   hasn't touched.

## Example

```
Set("x", "1")           // x=1, no transaction
Begin()                 // T1
Set("x", "2")           // x=2 inside T1; permanent storage still has x=1
Begin()                 // T2, nested inside T1
Set("x", "3")           // x=3 inside T2
Get("x")                // -> "3", true  (innermost wins)
Rollback()               // discard T2 entirely
Get("x")                // -> "2", true  (back to T1's value)
Commit()                 // fold T1 into permanent storage
Get("x")                // -> "2", true  (now permanent, no transaction open)
```

## Rubric

- [ ] Nested transactions genuinely nest — rolling back an inner transaction never
      affects the outer one's uncommitted changes.
- [ ] `Delete` inside a transaction correctly hides an outer-scope/permanent value from
      `Get` (a deleted-in-this-transaction key must not fall through to an older
      value), but a `Rollback` correctly un-hides it.
- [ ] `Commit`/`Rollback` on a store with no open transaction return `false` and don't
      panic or corrupt state.
- [ ] `Get` on a key nobody has ever touched, at any nesting depth, returns `("",
      false)` — not a panic, not a false positive.
- [ ] No dependency on any external database — this is genuinely in-memory.

## Hints

- Model this as a **stack of layers**: layer 0 is permanent storage; each `Begin()`
  pushes a new, empty layer. Each layer maps `key -> *string`, where a `nil` pointer
  means "explicitly deleted in this layer" (distinct from "not touched in this layer,
  fall through to the layer below").
- `Get`: walk the stack from the top down; the first layer where the key has an entry
  (whether a value or a delete marker) determines the answer.
- `Set`/`Delete`: always write to the **top** layer (the innermost open transaction, or
  layer 0 if none is open).
- `Commit`: pop the top layer, then write each of its entries (values AND delete
  markers) into what is now the new top layer — this correctly propagates deletes, not
  just values.
- `Rollback`: pop the top layer and throw it away.
