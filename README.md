<div align="center">

# LLD, Concurrency & OOD Interview Prep

**A hands-on, test-driven repository for Low-Level Design, Object-Oriented Design,
concurrency, and system design — in idiomatic Go.**

[![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Format](https://img.shields.io/badge/format-TDD%20%2F%20TODO%20templates-informational?style=flat-square)](#how-to-work-through-a-problem)
[![License](https://img.shields.io/badge/license-mixed%20%E2%80%94%20see%20below-lightgrey?style=flat-square)](#license-and-attribution)
[![Status](https://img.shields.io/badge/status-work%20in%20progress-yellow?style=flat-square)](#repository-layout)

</div>

---

Most folders in this repository share one format: a `README.md` explaining the
problem, a `// TODO`-driven `solution.go` template with no implementation given away,
and a `solution_test.go` that stays red until the problem is solved correctly. Each
problem is its own Go module, so any folder can be opened and run in isolation.

## Table of Contents

- [At a Glance](#at-a-glance)
- [Repository Layout](#repository-layout)
- [How to Work Through a Problem](#how-to-work-through-a-problem)
- [Suggested Order](#suggested-order)
- [License and Attribution](#license-and-attribution)

## At a Glance

| Area | Folder | Count |
| :-- | :-- | :-: |
| Concurrency puzzles | `concurrency/` | 7 |
| Semaphore-based synchronization patterns | `concurrency-patterns/` | 7 |
| SOLID principles | `SOLID/` | 5 |
| Machine coding problems | `machine-coding/` | 8 |
| Architectural / enterprise patterns | `design-patterns/` | 8 |
| OOD practice problems | `awesome-low-level-design/oops/` | 30 |
| Fully-worked LLD systems | `awesome-low-level-design/solutions/` | ~34 |
| Gang-of-Four design patterns | `awesome-low-level-design/design-patterns/` | 19 |

## Repository Layout

### `concurrency/`

Classic multithreading puzzles — signaling, mutual exclusion, producer/consumer-style
coordination — each its own Go module.

| Folder | Status |
| :-- | :-- |
| `print-in-order/` | implemented, tested |
| `Print-Zero-Even-Odd/` | implemented, tested |
| `Fizz -Buzz-Multithreaded/` | implemented, tested |
| `Building-H2O/` | implemented, tested |
| `The-Dining-Philosophers/` | implemented, tested |
| `print-foo-bar/` | solved, not yet wired up with `go.mod` / tests |
| `design-bounded-blocking-queue/` | open `// TODO` stub |

### `concurrency-patterns/`

Seven exercises building semaphore-based synchronization primitives and patterns from
first principles, translated into idiomatic Go (`sync.Mutex` / `sync.Cond`, buffered
channels as token buckets) rather than literal semaphore syntax. Based on *The Little
Book of Semaphores* — see [License and Attribution](#license-and-attribution).

`counting-semaphore` · `multiplex` · `reusable-barrier` · `lightswitch` ·
`readers-writers` · `no-starve-mutex` · `cigarette-smokers`

### `SOLID/`

The five SOLID principles, one per folder — a different format from the rest of the
repository: `solution.go` compiles and "works" today but visibly violates the
principle; the tests are written against the *refactored* shape and won't compile
until the violation is fixed.

`single-responsibility` · `open-closed` · `liskov-substitution` ·
`interface-segregation` · `dependency-inversion`

### `machine-coding/`

Multi-class systems designed and built under a suggested time limit, graded against a
rubric (correctness, code organization, extensibility, naming) — the format used in
live "machine coding round" interviews, distinct from the single-class problems in
`awesome-low-level-design/oops/` and the fully-worked systems in
`awesome-low-level-design/solutions/`.

`rate-limiter` · `meeting-room-scheduler` · `pluggable-cache` · `amazon-locker-system` ·
`notification-system` · `job-scheduler` · `in-memory-transactional-kv-store` ·
`text-editor-undo-redo`

### `design-patterns/` (repository root)

Architectural and enterprise patterns for real backend services — mostly Martin
Fowler's *Patterns of Enterprise Application Architecture* territory, plus one
idiomatic-Go-specific pattern. Deliberately distinct from the classic Gang-of-Four
patterns already covered in `awesome-low-level-design/design-patterns/`.

`repository` · `unit-of-work` · `cqrs` · `circuit-breaker` · `saga` · `outbox` ·
`specification` · `functional-options`

### `awesome-low-level-design/`

A larger, pre-existing collection this repository builds on top of:

| Subfolder | Contents |
| :-- | :-- |
| `oop/golang/` | OOP fundamentals — abstraction, inheritance, polymorphism, composition |
| `design-patterns/golang/` | the classic Gang-of-Four design patterns |
| `solutions/{cpp,golang}/` | ~34 fully-worked, multi-class LLD systems (parking lot, elevator system, ATM, chess, splitwise, ride sharing, and more) |
| `oops/` | 30 original single-class, LeetCode-style OOD interview questions (LRU/LFU cache, design Twitter, design a leaderboard, and more), each with a `// TODO` template and test suite — added as part of this repository's own work, not the upstream source |
| `problems/` | Markdown problem statements for the systems in `solutions/` |
| `class-diagrams/` | class diagram images for the systems in `solutions/` |

See [License and Attribution](#license-and-attribution) regarding the parts of this
folder not authored in this repository.

### `lld/` and `system-design/`

Reserved, not yet built out. `lld/` is planned to hold interview-*process* notes
(clarifying questions to ask, how to structure a 45-minute round, evaluation criteria)
rather than more solved problems. `system-design/` is not yet scoped.

### `LittleBookOfSemaphores.pdf`

The source text for `concurrency-patterns/` — see
[License and Attribution](#license-and-attribution).

## How to Work Through a Problem

Every problem folder is its own Go module:

```sh
cd <folder>
go test -v ./...          # see it fail against the TODO stub / unrefactored code

# implement solution.go

go test -v ./...          # confirm it passes
go test -race ./...       # for anything concurrent
```

## Suggested Order

| Step | Folder |
| :-: | :-- |
| 1 | `concurrency/` — finish `design-bounded-blocking-queue`, wire up `print-foo-bar` |
| 2 | `awesome-low-level-design/oops/` — see that folder's own tiered ordering |
| 3 | `concurrency-patterns/` |
| 4 | `SOLID/` |
| 5 | `machine-coding/` |
| 6 | `design-patterns/` |
| 7 | `awesome-low-level-design/solutions/` — full multi-class systems |

## License and Attribution

This repository mixes original material written for personal interview preparation
with material derived from, or building on top of, outside sources. The terms below
apply specifically to those sources; original content in this repository (test suites,
`// TODO` templates, and README explanations written for `concurrency/`,
`concurrency-patterns/`, `SOLID/`, `machine-coding/`, `design-patterns/`, and
`awesome-low-level-design/oops/`) is personal educational work and is not separately
licensed for redistribution.

<table>
<tr>
<td width="100">

[![CC BY-NC-SA 4.0](https://licensebuttons.net/l/by-nc-sa/4.0/88x31.png)](https://creativecommons.org/licenses/by-nc-sa/4.0/)

</td>
<td>

**The Little Book of Semaphores** (`LittleBookOfSemaphores.pdf`), by Allen B. Downey,
Version 2.2.1, Copyright 2016 Allen B. Downey, is licensed under
**[CC BY-NC-SA 4.0](https://creativecommons.org/licenses/by-nc-sa/4.0/)**
(Attribution-NonCommercial-ShareAlike 4.0 International). The LaTeX source is
distributed by the author at [greenteapress.com/semaphores](http://greenteapress.com/semaphores).

</td>
</tr>
</table>

The `concurrency-patterns/` folder in this repository is a derivative educational
work: the underlying synchronization problems and patterns are drawn from this book,
but reframed and re-implemented from scratch in idiomatic Go (`sync.Mutex` /
`sync.Cond`, channels, and Go-specific patterns) rather than reproducing the book's own
pseudocode or prose. If any part of this repository is redistributed publicly, the
`concurrency-patterns/` folder should retain this attribution and, per the license's
ShareAlike term, remain available under CC BY-NC-SA 4.0 (or a compatible license) and
be used for non-commercial purposes only.

> **`awesome-low-level-design/`** predates the work done in this session (only the
> `oops/` subfolder within it was added here) and does not currently include its own
> `LICENSE` or attribution file in this local copy. If this folder originated from an
> external open-source project, its original license terms should be confirmed and
> restored before any public redistribution of this repository — do not treat the
> absence of a license file here as license-free.
