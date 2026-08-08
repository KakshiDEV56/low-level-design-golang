# Single Responsibility Principle (SRP)

> A class (or type) should have only one reason to change.

## The violation in `solution.go`

`Invoice` currently does two unrelated jobs:

1. **Pricing logic** — computing the total from line items (`Total()`).
2. **Presentation logic** — rendering itself as a human-readable string (`Print()`).

These are two independent reasons for this type to change: a change to pricing rules
(discounts, tax, currency rounding) and a change to output format (adding a JSON export,
changing column widths, localizing the currency symbol) both force you to touch the same
type. Worse, `Print()` calls `Total()` internally, so the two responsibilities are also
*coupled* — you can't easily reuse the pricing logic somewhere that doesn't want string
formatting (e.g., an API that returns structured JSON), or add a new output format
without risking a change to the pricing method it depends on.

## What a fix looks like

Separate "compute the total" from "render the invoice" into two independent things:

- `Invoice` (plus `LineItem`) stays pure data + pricing logic — `Total()` and nothing
  about output format.
- A separate `Formatter` **interface** with a `Format(Invoice) string` method, and at
  least one concrete implementation (e.g. `PlainTextFormatter`), takes an `Invoice` and
  renders it. Formatting logic never needs to know how `Total()` is computed internally
  — it just calls it.

The test in `solution_test.go` defines a *second*, different formatter right there in
the test file and checks it works against your `Invoice` with zero changes to `Invoice`
itself — that's the proof the separation actually happened.

## Target shape

```go
type LineItem struct {
    Name     string
    Price    float64
    Quantity int
}

type Invoice struct {
    Items []LineItem
}

func (inv Invoice) Total() float64 { ... }

type Formatter interface {
    Format(inv Invoice) string
}

type PlainTextFormatter struct{}

func (PlainTextFormatter) Format(inv Invoice) string { ... }
```
