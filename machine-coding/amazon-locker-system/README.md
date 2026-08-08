# Amazon Locker System

**Suggested time box: 40 minutes.**

## Problem statement

Design a package locker system (like Amazon Lockers). Lockers come in fixed sizes.
Delivering a package assigns it to the smallest available locker that can fit it, and
generates a pickup code the recipient uses to open it. Picking up a package with the
correct code frees the locker back up.

## Functional requirements

1. Locker sizes: `Small`, `Medium`, `Large` (in that size order — a `Large` locker can
   fit anything, a `Small` locker only fits `Small` packages).
2. `NewLockerSystem(lockers map[string]Size)` initializes the system with a fixed set
   of lockers (`lockerID -> size`), all initially empty.
3. `AssignLocker(packageID string, packageSize Size) (lockerID, code string, ok bool)`:
   - Finds the **smallest available locker** whose size is `>= packageSize` (best fit —
     don't waste a `Large` locker on a `Small` package if a `Small` or `Medium` one is
     free).
   - If multiple free lockers tie for smallest-sufficient-size, pick the
     alphabetically-first locker ID (deterministic, for testability).
   - Generates a **unique** pickup code for this assignment.
   - Returns `ok=false` (with no locker assigned) if: no locker is big enough and free,
     OR `packageID` already has an active, not-yet-picked-up assignment.
4. `Pickup(packageID string, code string) bool`: if `packageID` has an active
   assignment and `code` matches exactly, releases the locker (making it available
   again) and returns `true`. Wrong code, unknown package, or an already-picked-up
   package all return `false` and change nothing.
5. `AvailableCount(size Size) int`: number of currently-free lockers of exactly that
   size (not "at least that size").

## Rubric

- [ ] Best-fit assignment never picks an oversized locker when a smaller sufficient one
      is free.
- [ ] Pickup codes are genuinely unique across assignments (don't just return a
      constant or an easily-guessable sequential value like locker ID).
- [ ] A wrong code does NOT release the locker.
- [ ] Picking up a package correctly makes its locker available for a **new**
      assignment (verify reuse, not just that `Pickup` returns `true`).
- [ ] Double-assigning an already-assigned (not picked up) `packageID` is rejected.

## Hints

- Track lockers as `map[string]lockerState` where `lockerState` has the size and
  whether it's occupied (and, if occupied, by which package + code).
- For "unique code," something as simple as an incrementing counter formatted as a
  string, or `crypto/rand`/`math/rand` generating a few random digits with a
  collision check, is fine — the bar is "not predictable/constant," not
  cryptographic security.
- For best-fit with alphabetical tie-breaking: iterate candidate lockers of the
  target size in sorted order first; only fall back to the next larger size if none of
  that size are free.
