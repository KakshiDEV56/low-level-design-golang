# Cigarette Smokers Problem (The "Pusher" / Dedicated-Coordinator Pattern)

*Based on The Little Book of Semaphores, §4.5 ("Cigarette smokers problem").*

## The concept

Originally posed by Suhas Patil, who claimed (with some qualifications) that it can't be
solved with semaphores at all — which is exactly why it's worth doing once you've built
confidence on the earlier problems in this set.

Four goroutines: an **agent**, and three **smokers**. Each smoker has an infinite supply
of exactly one ingredient — one has tobacco, one has paper, one has matches — and needs
the *other two* to make and smoke a cigarette. The agent has an infinite supply of all
three, and makes ingredients available one at a time, in no particular order.

The constraint that makes this hard: **the agent must not need to know which smoker is
waiting for what.** (The book calls a version where the agent explicitly signals "the
smoker who should go next" the "trivial version," and explicitly rejects it as
uninteresting — it defeats the whole point.) So something else has to watch what
ingredients are currently on the table and figure out, whenever a *pair* becomes
complete, which single smoker can now proceed.

### Why the obvious smoker-side code doesn't work

It's tempting to write, for the smoker who has matches (and thus needs tobacco and
paper):

```
tobacco.wait()
paper.wait()
# make and smoke a cigarette
```

This deadlocks. Imagine the agent puts out tobacco and paper. The smoker-with-matches
unblocks on `tobacco.wait()` — but so might the smoker-with-tobacco (who's waiting on
`paper.wait()` as their first step, since the roles are symmetric). Now the first thread
blocks on `paper` and the second blocks on `match`. Nobody can proceed. Two smokers each
grabbed one ingredient of a pair that only one of them can actually use, and neither will
ever release what they grabbed.

### The fix: a dedicated coordinator ("pusher")

The book's solution moves ALL the decision-making out of the smokers and into a small
set of helper threads (it calls them **"pushers"**) that:

- track, with three boolean flags protected by a mutex, which ingredients are currently
  sitting on the table,
- and — critically — whenever a *new* ingredient arrives, check whether it completes a
  pair with something already there. If so, they wake the *one* smoker who can use that
  exact pair (and clear the flags). If not, they just record the new ingredient and wait.

Neither the agent nor the smokers ever need to reason about each other. All the
cross-cutting logic lives in one place.

## Why this matters for a Go engineer

This is the general shape of **"multiple independent resources must be assembled before
any one consumer can proceed, and no single resource's arrival is enough information on
its own to know who to wake."** That's a real coordination pattern, well beyond
cigarettes:

- A build/CI system where a job needs *N* different artifacts to be ready (built by
  independent producers) before it can start, and you don't want producers to know
  about consumers or each other.
- Any "wait for K of N conditions, then dispatch to whichever specific handler can now
  proceed" scheduler.

The core lesson is architectural, not just syntactic: when a naive "just signal the
matching channel/semaphore directly" design creates a deadlock or requires resources to
know too much about their consumers, introducing **one small piece of centralized state
plus a decision function** (the pusher) is often the clean fix — much cleaner than
trying to make the resources or consumers smarter about each other.

## The constraint to implement

- `ProvideTobacco()`, `ProvidePaper()`, `ProvideMatch()` — the agent calls these (one
  ingredient at a time, any order) to put an ingredient on the table.
- `SmokeWithMatches(smoke func())`, `SmokeWithTobacco(smoke func())`,
  `SmokeWithPaper(smoke func())` — each blocks until its own smoker can make a
  cigarette (i.e. until the *other two* ingredients are both on the table), then calls
  `smoke()` exactly once and returns. (Method names describe what each smoker
  *possesses* — `SmokeWithMatches` is the smoker who has matches and needs tobacco+paper,
  matching the book's own naming.)
- At any given moment, at most one ingredient should be sitting unconsumed on the table
  waiting for its pair (callers are expected not to call `Provide*` for an ingredient
  that's already down and unconsumed — that mirrors the book's assumption that the
  agent doesn't put out a new round until the current one has been picked up).
- No smoker may ever be woken for a pair it can't actually use.

## Hint

Three boolean flags (`tobaccoOnTable`, `paperOnTable`, `matchOnTable`), one `sync.Mutex`
protecting them, and one "trigger" channel per smoker (buffered size 1, or a
`sync.Cond`) that the corresponding `Provide*` call sends into when it discovers the
complementary ingredient is already down. Each `Provide*` method is doing exactly the
job of one of the book's three "pushers": check whether the OTHER TWO flags include the
one ingredient this smoker needs alongside the one that just arrived... concretely, for
`ProvideTobacco()`: if `paperOnTable`, clear it and wake the matches-smoker (who needs
tobacco+paper); else if `matchOnTable`, clear it and wake the paper-smoker (who needs
tobacco+match); else set `tobaccoOnTable = true` and wait for the next ingredient.
