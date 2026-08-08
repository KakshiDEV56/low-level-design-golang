package singleresponsibility

import (
	"fmt"
	"strings"
)

/*
This code compiles and works today — but it violates the Single Responsibility
Principle. See README.md for the full explanation.

Invoice bundles two unrelated responsibilities into one type:
 1. Pricing logic (Total()).
 2. Presentation logic (Print()) — which also happens to call Total() internally,
    coupling the two together.

TODO: refactor this file so that:
 1. Invoice (and LineItem) contain ONLY data plus the Total() pricing method — nothing
    about output formatting.
 2. A new Formatter interface is introduced:
      type Formatter interface {
          Format(inv Invoice) string
      }
 3. A PlainTextFormatter type implements Formatter, producing the same kind of output
    Print() used to (one line per item, then a "Total: $X.XX" line — exact formatting
    is up to you, the test only checks that it contains the right numbers/names).
 4. Delete Print() from Invoice entirely — formatting no longer lives there.

Once refactored, solution_test.go (which defines its OWN second Formatter
implementation, entirely outside this file) should compile and pass — proving Invoice
never needed to know anything about how it gets formatted.
*/

type LineItem struct {
	Name     string
	Price    float64
	Quantity int
}

// Invoice currently does pricing AND presentation. TODO: split these apart.
type Invoice struct {
	Items []LineItem
}

func (inv *Invoice) Total() float64 {
	total := 0.0
	for _, item := range inv.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// Print renders the invoice as text. This is a presentation concern that has nothing
// to do with pricing, yet it lives on the same type that owns Total() — and it calls
// Total() directly, coupling the two responsibilities together.
func (inv *Invoice) Print() string {
	var b strings.Builder
	for _, item := range inv.Items {
		fmt.Fprintf(&b, "%s x%d: $%.2f\n", item.Name, item.Quantity, item.Price*float64(item.Quantity))
	}
	fmt.Fprintf(&b, "Total: $%.2f\n", inv.Total())
	return b.String()
}
