# Liskov Substitution Principle (LSP)

> Objects of a superclass (or, in Go, anything satisfying a shared interface) should be
> replaceable with objects of a subclass without breaking the correctness of the
> program.

## The violation in `solution.go`

This is the classic Rectangle/Square trap, in Go's terms. `Rectangle` has independent
`SetWidth`/`SetHeight` methods — the reasonable, implied contract is "setting one
dimension doesn't affect the other." `Square` embeds `*Rectangle` (Go's version of
"inheriting" its methods) but **overrides** `SetWidth` and `SetHeight` so that setting
either one forces both dimensions to match (since a square, by definition, must have
equal sides).

Both types structurally satisfy the same shape (`SetWidth`, `SetHeight`, `Area`) — so
you could write a `Resizable` interface and pass either one to code that expects that
contract. But any code written against `Resizable` that reasonably assumes "calling
`SetWidth` doesn't change `Area` in a way that also depends on the previous
`SetHeight` call" will get silently wrong answers when handed a `Square` — the
*behavioral* contract is violated even though the *structural* one (the interface) is
satisfied. This is exactly what LSP is about: interface satisfaction alone isn't
enough; substitutability requires behavioral compatibility too.

## What a fix looks like

The honest fix is to stop pretending `Square` is a kind of `Rectangle` that can be
resized on either axis independently — because it isn't one. Two changes:

1. `Rectangle` keeps `SetWidth`/`SetHeight`/`Area` and is the only type that implements
   a `Resizable` interface (`{ SetWidth(int); SetHeight(int) }`) — since it's the only
   type for which "set either dimension independently" is actually a valid operation.
2. `Square` becomes its own independent type (not embedding `Rectangle` at all), with
   its own `SetSize(int)` method that keeps its one invariant (all sides equal) — and
   critically, it does **not** implement `Resizable`, because doing so would always be a
   lie.

This isn't just a naming change — it means it becomes a **compile-time error**, not a
runtime surprise, to pass a `Square` anywhere a `Resizable` is expected. The test file
proves this by using `*Rectangle` as a `Resizable` successfully, and documents (via a
comment showing what would fail to compile) why `*Square` can't be.

## Target shape

```go
type Shape interface {
    Area() int
}

type Resizable interface {
    SetWidth(w int)
    SetHeight(h int)
}

type Rectangle struct{ width, height int }

func NewRectangle(w, h int) *Rectangle
func (r *Rectangle) SetWidth(w int)
func (r *Rectangle) SetHeight(h int)
func (r *Rectangle) Area() int

type Square struct{ size int }

func NewSquare(size int) *Square
func (s *Square) SetSize(size int)
func (s *Square) Area() int
```

`*Rectangle` implements both `Shape` and `Resizable`. `*Square` implements only `Shape`
— it has no `SetWidth`/`SetHeight` methods at all.
