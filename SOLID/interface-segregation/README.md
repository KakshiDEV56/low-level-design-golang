# Interface Segregation Principle (ISP)

> No client should be forced to depend on methods it does not use. Prefer several small,
> specific interfaces over one large, general-purpose one.

## The violation in `solution.go`

`Worker` bundles two unrelated capabilities into one "fat" interface: `Work()` and
`Eat()`. That's a fine fit for `HumanWorker` — but `RobotWorker` doesn't eat, and yet it
is **forced to implement `Eat()` anyway** just to satisfy the interface, resulting in a
method that either panics or silently does nothing. Either way, it's a lie: the type
signature promises `RobotWorker` can do something it fundamentally cannot do.

This is a smell any time you see a method implemented as `panic("not supported")`,
`return nil // no-op`, or a comment like `// N/A for this type` — it usually means the
interface was drawn too broadly for what its implementers actually need to provide.

## What a fix looks like

Split the fat interface into two small, focused ones: `Workable` (just `Work()`) and
`Eatable` (just `Eat()`). `HumanWorker` implements both. `RobotWorker` implements only
`Workable` — it simply has no `Eat()` method at all, which is honest, rather than having
one that lies about what it does.

Anything that only needs to manage *work* (like a shift scheduler) should depend on the
narrow `Workable` interface, not the fat one — so it can accept `RobotWorker` without
that type ever needing to pretend to eat. This is checkable at compile time in Go, which
is what makes ISP violations so directly testable here: once the split is done, it
becomes a compile error (not a runtime `panic`) to expect a `RobotWorker` to be
`Eatable`.

## Target shape

```go
type Workable interface {
    Work() string
}

type Eatable interface {
    Eat() string
}

type HumanWorker struct{}
func (HumanWorker) Work() string
func (HumanWorker) Eat() string

type RobotWorker struct{}
func (RobotWorker) Work() string
// RobotWorker has NO Eat() method.

// RunShift only needs to know how to make something work — it must accept
// RobotWorker just as easily as HumanWorker.
func RunShift(w Workable) string
```
