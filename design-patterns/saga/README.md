# Saga (Orchestrated, with Compensation)

## The concept

A saga is a sequence of steps that together form one logical operation spanning
multiple services or resources — where you can't wrap the whole thing in a single
database transaction (the classic example: "book a flight, book a hotel, charge a
card," where each is a call to a different service). If a later step fails, you can't
roll back the earlier ones with a database `ROLLBACK` — instead, each step defines its
own **compensating action** (an explicit "undo"), and the saga runs those compensations,
in reverse order, for every step that already succeeded.

This is the orchestrated flavor: one central `Saga` runs the steps in order and decides
when to trigger compensation, as opposed to the choreographed flavor (each service
reacts to events from the previous one with no central coordinator) — orchestration is
simpler to reason about and is what this exercise implements.

## Functional requirements

1. `Step` (given) bundles a `Name`, an `Action func() error`, and a `Compensate func()
   error`.
2. `NewSaga(steps ...Step) *Saga` builds a saga from an ordered list of steps.
3. `Run() error`:
   - Executes each step's `Action`, in order.
   - If a step's `Action` fails, **stop** running further steps. Then call
     `Compensate` for every step that already succeeded, **in reverse order** (most
     recently succeeded first).
   - A compensation failing must NOT stop the rest of the compensations from running —
     attempt every one of them regardless.
   - Returns an error that satisfies `errors.Is` against BOTH the original failing
     step's error AND any compensation error(s) that occurred (use `errors.Join`).
     Returns `nil` if every step's `Action` succeeded (no compensation needed, nothing
     to report).

## Rubric

- [ ] A step failure never causes a LATER step's `Action` to run.
- [ ] Compensations run in exact reverse order of the successful actions.
- [ ] A step that never got to run (because an earlier step failed first) never has its
      `Compensate` called either — there's nothing to undo.
- [ ] One compensation failing doesn't prevent the others from being attempted.
- [ ] A fully successful saga (`Run() == nil`) never calls any `Compensate` function.

## Hints

- Track successfully-completed steps in a slice as you go; on failure, iterate that
  slice backwards calling `Compensate`.
- `errors.Join(err1, err2, ...)` (standard library, Go 1.20+) is exactly the tool for
  combining the step failure with zero or more compensation failures into one error
  that `errors.Is` still works against for each individual error.
