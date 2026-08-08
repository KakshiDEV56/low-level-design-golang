# Repository Pattern

## The concept

The Repository pattern puts a **collection-like interface** in front of persistence:
callers `Add`, `GetByID`, `Update`, `Delete`, and `List` domain objects as if they were
working with an in-memory collection, with no idea whether the actual storage is a SQL
database, a document store, or (as in this exercise) a plain map. Business logic depends
on the `UserRepository` **interface**, never on a concrete storage implementation — the
same inversion idea as the Dependency Inversion Principle (`../../SOLID/dependency-inversion`),
but Repository is a specific, named shape of that idea: a *collection* abstraction, not
just any interface.

## Functional requirements

1. `UserRepository` is a fixed interface (given below): `Add`, `GetByID`, `Update`,
   `Delete`, `List`, `FindByEmail`.
2. Implement `InMemoryUserRepository`, backed by a map, satisfying `UserRepository`:
   - `Add` fails with `ErrDuplicateID` if the ID is already taken.
   - `GetByID`/`Update`/`Delete` fail with `ErrNotFound` for an unknown ID.
   - `List` returns all users **sorted by ID** (deterministic — a real DB query would
     need an explicit `ORDER BY` too; don't rely on map iteration order).
   - `FindByEmail` fails with `ErrNotFound` if no user has that email.
3. Implement `UserService`, which depends on `UserRepository` (the interface, injected
   via constructor) — never on `InMemoryUserRepository` directly:
   - `Register(id, email, name string) error` fails with `ErrDuplicateEmail` if the
     email is already taken by a different user (checked via `FindByEmail`), otherwise
     adds the user.
   - `Rename(id, newName string) error` fails with `ErrNotFound` for an unknown ID,
     otherwise updates the name.

## Rubric

- [ ] `UserService` contains no storage-specific code anywhere — it only calls
      `UserRepository` interface methods.
- [ ] A brand-new `UserRepository` implementation, written entirely in the test file,
      works with `UserService` unmodified — the concrete proof of the abstraction.
- [ ] `List` is deterministically ordered.
- [ ] Every "not found" / "duplicate" case returns the specific sentinel error, not a
      generic `errors.New("error")` the caller can't distinguish.

## Hints

- `map[string]User` plus `sort.Slice` on `List` is all `InMemoryUserRepository` needs.
- `errors.Is` is how tests will check for the sentinel errors — use `errors.New` for
  each sentinel once at package scope, and return that exact value (or wrap it with
  `fmt.Errorf("...: %w", ErrNotFound)` if you want extra context — either is fine as
  long as `errors.Is` still matches).
