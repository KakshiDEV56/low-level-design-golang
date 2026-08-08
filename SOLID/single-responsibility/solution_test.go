package singleresponsibility

import (
	"fmt"
	"strings"
	"testing"
)

func sampleInvoice() Invoice {
	return Invoice{
		Items: []LineItem{
			{Name: "Widget", Price: 9.99, Quantity: 3},
			{Name: "Gadget", Price: 19.99, Quantity: 1},
		},
	}
}

func TestInvoice_TotalIsIndependentOfFormatting(t *testing.T) {
	inv := sampleInvoice()
	want := 9.99*3 + 19.99*1
	if got := inv.Total(); absDiff(got, want) > 0.001 {
		t.Fatalf("Total() = %.4f, want %.4f", got, want)
	}
}

func TestInvoice_EmptyInvoiceTotalsZero(t *testing.T) {
	inv := Invoice{}
	if got := inv.Total(); got != 0 {
		t.Fatalf("Total() on empty invoice = %v, want 0", got)
	}
}

func TestPlainTextFormatter_ContainsAllItemsAndTotal(t *testing.T) {
	inv := sampleInvoice()
	var f Formatter = PlainTextFormatter{}
	out := f.Format(inv)

	for _, item := range inv.Items {
		if !strings.Contains(out, item.Name) {
			t.Errorf("formatted output %q does not mention item %q", out, item.Name)
		}
	}
	wantTotal := fmt.Sprintf("%.2f", inv.Total())
	if !strings.Contains(out, wantTotal) {
		t.Errorf("formatted output %q does not contain the total %q", out, wantTotal)
	}
}

// csvFormatter is a SECOND, entirely independent Formatter implementation, defined
// here in the test file — outside solution.go, with zero changes to Invoice needed to
// support it. If this compiles and works, it proves formatting was genuinely separated
// from Invoice's pricing logic, which is the whole point of the SRP fix.
type csvFormatter struct{}

func (csvFormatter) Format(inv Invoice) string {
	var b strings.Builder
	b.WriteString("name,quantity,price\n")
	for _, item := range inv.Items {
		fmt.Fprintf(&b, "%s,%d,%.2f\n", item.Name, item.Quantity, item.Price)
	}
	fmt.Fprintf(&b, "TOTAL,,%.2f\n", inv.Total())
	return b.String()
}

func TestCustomFormatter_WorksWithoutChangesToInvoice(t *testing.T) {
	inv := sampleInvoice()
	var f Formatter = csvFormatter{}
	out := f.Format(inv)

	if !strings.HasPrefix(out, "name,quantity,price\n") {
		t.Fatalf("csv output missing header, got %q", out)
	}
	for _, item := range inv.Items {
		if !strings.Contains(out, item.Name) {
			t.Errorf("csv output %q does not mention item %q", out, item.Name)
		}
	}
	wantTotal := fmt.Sprintf("%.2f", inv.Total())
	if !strings.Contains(out, wantTotal) {
		t.Errorf("csv output %q does not contain the total %q", out, wantTotal)
	}
}

func TestFormatters_ProduceDifferentOutputForSameInvoice(t *testing.T) {
	inv := sampleInvoice()
	plain := PlainTextFormatter{}.Format(inv)
	csv := csvFormatter{}.Format(inv)

	if plain == csv {
		t.Fatal("two different Formatter implementations produced identical output — they should differ in shape")
	}
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
