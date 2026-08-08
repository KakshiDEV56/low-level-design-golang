# Open-Closed Principle (OCP)

> Software entities should be open for extension, but closed for modification.

## The violation in `solution.go`

`DiscountCalculator.Discount` decides how much discount to apply based on a
`customerType` string, via a `switch`. This works for the three cases it knows about
today — but the moment the business adds a fourth customer tier (say, "employee"), you
have to **open up and edit `Discount` itself** to add a new case. Every future customer
type means editing code that's already shipped and already tested, which is exactly
what "closed for modification" is warning you away from. It also means
`DiscountCalculator` has to know about every customer type that will ever exist, which
doesn't scale and couples unrelated concerns (the calculator's job vs. each customer
tier's specific discount rule) into one function.

## What a fix looks like

Replace the type-switch with **polymorphism via an interface**: each customer type
becomes its own type that knows its own discount rule, and the calculator just calls an
interface method — it never needs to know the full list of customer types that exist.

Adding a new customer type then means *adding a new type that implements the
interface* — not touching the calculator at all. That's the "open for extension, closed
for modification" property, and it's directly testable: the test file defines a brand
new customer type that didn't exist when `solution.go` was written, and proves it works
with the unmodified calculator.

## Target shape

```go
type Discounter interface {
    Discount(amount float64) float64
}

type RegularCustomer struct{}
func (RegularCustomer) Discount(amount float64) float64 { ... }

type PremiumCustomer struct{}
func (PremiumCustomer) Discount(amount float64) float64 { ... }

type VIPCustomer struct{}
func (VIPCustomer) Discount(amount float64) float64 { ... }

// ApplyDiscount must not need to know which concrete Discounter it was given.
func ApplyDiscount(d Discounter, amount float64) float64 { ... }
```

Preserve the original discount rates: Regular = 0%, Premium = 10%, VIP = 20%.
