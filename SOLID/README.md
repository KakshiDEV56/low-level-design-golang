# SOLID Principles

Five exercises, one per principle. No PDF backs this set — the SOLID principles (coined
by Robert C. Martin, "Uncle Bob") are well-established and don't need a specific text to
reference faithfully.

## The format is different from the other folders here

Everywhere else in this repo (`concurrency/`, `oops/`, `concurrency-patterns/`), a
problem gives you an **empty struct with `// TODO` method stubs** and a test suite that
fails until you implement it.

SOLID is about *recognizing and fixing a violation*, not implementing something from
scratch — so each folder instead gives you:

- **`solution.go`** — code that **compiles and "works"** today, but violates the
  principle in a specific, real way (with a comment explaining exactly what the
  violation is and why it matters).
- **`solution_test.go`** — tests written against the **refactored** shape the code
  should have. These tests will **not compile** against the starting code — that's
  intentional. Your job is to refactor `solution.go` (change the types, split
  responsibilities, introduce interfaces — whatever the principle calls for) until the
  test file compiles and passes.
- **`README.md`** — the principle, the specific smell in the starting code, and what a
  fixed version looks like structurally (without just handing you the diff).

This mirrors how SOLID violations actually show up in real code review: you're handed
something that runs fine today, and the question is "what breaks when requirements
change, and how do you restructure this so it doesn't?"

## The five folders

1. **`single-responsibility/`** — a type that bundles business logic and output
   formatting together, so a change to either one risks breaking the other.
2. **`open-closed/`** — a type-switch that has to be edited every time a new case is
   added, instead of being extensible without modification.
3. **`liskov-substitution/`** — a subtype that structurally satisfies an interface but
   violates the behavioral contract callers reasonably assume (the classic
   Rectangle/Square trap, in Go terms).
4. **`interface-segregation/`** — a "fat" interface that forces an implementer to
   provide a method that doesn't make sense for it.
5. **`dependency-inversion/`** — a high-level type that directly constructs and depends
   on a concrete low-level detail, instead of depending on an abstraction.

## How to work through each one

```
cd SOLID/<folder>
go build ./...   # currently succeeds — this is the "before" state
go test ./...    # currently fails to COMPILE — refactor solution.go until it passes
```
