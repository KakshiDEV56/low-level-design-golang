# Unit of Work

## The concept

A single repository call is easy to make atomic — the storage layer handles it. The
problem Unit of Work solves is **coordinating changes across multiple entities (or
multiple repositories) as one atomic operation**, so that either all of them take effect
or none do.

The motivating example here is a bank transfer: debiting one account and crediting
another must happen together. If you debit account A, then something fails before you
credit account B, the money has simply vanished. Unit of Work fixes this by **staging**
every change in-memory first, and only writing them all to the real store in one
`Commit()` — with `Rollback()` discarding the staged changes entirely if anything goes
wrong first.

## Functional requirements

1. `Account` has `ID string` and `Balance int`. `AccountStore` (given, complete — the
   "database") holds committed account state.
2. `NewUnitOfWork(store *AccountStore) *UnitOfWork` starts a new unit of work over that
   store.
3. `Get(id string) (Account, error)` reads an account — if it's already been staged for
   a write in this unit of work, return the **staged** value; otherwise fall through to
   the committed value in `AccountStore`. Returns `ErrNotFound` if the account exists
   nowhere.
4. `Save(a Account)` stages a write — this must NOT be visible in `AccountStore` until
   `Commit()` is called.
5. `Commit() error` atomically applies every staged write to `AccountStore`.
6. `Rollback()` discards every staged write; `AccountStore` is left completely
   untouched.
7. `Transfer(store *AccountStore, fromID, toID string, amount int) error` — the
   business operation that demonstrates why this matters: debits `fromID`, credits
   `toID`, and either both take effect (via `Commit`) or neither does (`Rollback`, e.g.
   on `ErrInsufficientFunds` if `fromID`'s balance is less than `amount`, or if either
   account doesn't exist).

## Rubric

- [ ] A failed `Transfer` (insufficient funds, unknown account) leaves **both**
      accounts' committed balances in `AccountStore` completely unchanged — verify this
      directly against the store, not just that `Transfer` returned an error.
- [ ] Staged writes within a unit of work are visible to `Get` calls on that SAME unit
      of work, but invisible to `AccountStore` (and to a separate, independent unit of
      work over the same store) until `Commit`.
- [ ] A successful `Transfer` updates both balances together — there's no window where
      one account is debited but the other isn't yet credited, from the store's
      perspective.

## Hints

- `AccountStore` (given) is just a `map[string]Account` behind `Get`/`Save`/thread
  safety — you don't need to modify it.
- `UnitOfWork` needs its own `map[string]Account` for staged writes, checked first in
  `Get` before falling through to the store.
- `Commit` just iterates the staged map and calls `store.Save` for each entry.
