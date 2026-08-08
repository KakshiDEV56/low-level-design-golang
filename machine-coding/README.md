# Machine Coding Round Practice

This folder is different in *kind* from `awesome-low-level-design/oops/` (single-class
LeetCode-style OOD problems) and `awesome-low-level-design/solutions/` (fully worked
reference systems). This is the format used in **live "machine coding round"
interviews** (common at companies like Flipkart, Swiggy, Ola, Amazon, and in platforms
like Scaler/InterviewBit): you're given a multi-class system to design and build, under
a time limit, and graded against a rubric — not just "does it work," but "would this
survive code review."

Every problem here is something genuinely **not already covered** by the other LLD
folders in this repo (checked against all ~65 systems in `awesome-low-level-design/`).

## How to actually practice with these

1. **Read only the problem statement and functional requirements** in the folder's
   `README.md` — skip the hints section until you're stuck or done.
2. **Set a timer** — the suggested time box is in each README (30–60 minutes is typical
   for a real round). Build your solution directly in `solution.go`, replacing the
   `// TODO` stubs.
3. **Stop when the timer goes off**, even if unfinished — that's realistic. Then run
   `go test ./...` to check correctness.
4. **Grade yourself against the rubric** in the README. A real interviewer is scoring
   more than "do the tests pass":
   - **Correctness** — does it handle the stated requirements, including edge cases?
   - **Code organization** — are responsibilities cleanly separated (SOLID, not just
     "it compiles")? Would a reviewer need to ask "wait, why does this type do that?"
   - **Extensibility** — if a plausible new requirement showed up 5 minutes after time
     was called (a new eviction policy, a new notification channel, a new constraint),
     how much of the existing code would you have to touch?
   - **Naming and readability** — could a stranger understand this without you
     narrating it?
   - **Tests** — `solution_test.go` is provided so you're not grading your own
     homework, but in a real round you'd often be expected to write a few of your own
     too; consider adding 2-3 more edge-case tests of your own once the given ones pass.

## The problems, roughly easiest to hardest

1. **`rate-limiter/`** (30 min) — Token Bucket and Sliding Window Log, swappable via a
   common interface.
2. **`meeting-room-scheduler/`** (30 min) — book rooms, detect interval conflicts, find
   availability.
3. **`pluggable-cache/`** (40 min) — a cache with a swappable eviction policy (Strategy
   pattern) instead of one hardcoded algorithm.
4. **`amazon-locker-system/`** (40 min) — assign packages to size-appropriate lockers,
   retrieve via pickup code, release capacity.
5. **`notification-system/`** (45 min) — dispatch a notification across multiple
   channels (email/SMS/push) without the core service knowing channel-specific detail.
6. **`job-scheduler/`** (45 min) — schedule one-shot and recurring jobs against a
   controllable virtual clock (no real `time.Sleep`-based flakiness).
7. **`in-memory-transactional-kv-store/`** (45–60 min) — a key-value store with nested
   BEGIN/COMMIT/ROLLBACK transactions.
8. **`text-editor-undo-redo/`** (45–60 min) — a text buffer where every edit is
   undoable and redoable (Command + Memento).

Each folder is its own Go module: `cd machine-coding/<folder> && go test -v ./...`
