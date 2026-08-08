package openclosed

/*
This code compiles and works today — but it violates the Open-Closed Principle. See
README.md for the full explanation.

Every time a new customer type is added, this switch statement has to be edited —
DiscountCalculator is not "closed for modification."

TODO: refactor this file so that:
 1. A Discounter interface is introduced:
      type Discounter interface {
          Discount(amount float64) float64
      }
 2. RegularCustomer, PremiumCustomer, and VIPCustomer each become their own type
    implementing Discounter, preserving the original rates (0%, 10%, 20%).
 3. A free function ApplyDiscount(d Discounter, amount float64) float64 replaces
    DiscountCalculator.Discount — it must work with ANY Discounter without needing to
    know which concrete type it was given (no switch, no type assertion on customer
    type).
 4. Delete DiscountCalculator and the string-based switch entirely.

Once refactored, solution_test.go (which defines its OWN new customer type, entirely
outside this file, and passes it to ApplyDiscount unmodified) should compile and
pass — proving new customer types can be added without touching this file again.
*/

// DiscountCalculator decides the discount using a switch on a string customer type.
// TODO: replace this whole approach with an interface (see above).
type DiscountCalculator struct{}

func (DiscountCalculator) Discount(customerType string, amount float64) float64 {
	switch customerType {
	case "regular":
		return amount * 0.0
	case "premium":
		return amount * 0.1
	case "vip":
		return amount * 0.2
	default:
		return 0
	}
}
