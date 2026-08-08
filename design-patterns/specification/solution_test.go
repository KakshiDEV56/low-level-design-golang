package specification

import "testing"

// Product and the leaf specifications below are all domain-specific and defined ONLY
// in this test file — solution.go knows nothing about products, prices, or stock.
// This is the concrete proof that And/Or/Not/Filter are genuinely generic.
type Product struct {
	Name    string
	Price   int
	InStock bool
}

type minPriceSpec struct{ min int }

func (s minPriceSpec) IsSatisfiedBy(p Product) bool { return p.Price >= s.min }

type maxPriceSpec struct{ max int }

func (s maxPriceSpec) IsSatisfiedBy(p Product) bool { return p.Price <= s.max }

type inStockSpec struct{}

func (inStockSpec) IsSatisfiedBy(p Product) bool { return p.InStock }

func products() []Product {
	return []Product{
		{Name: "cheap-in-stock", Price: 5, InStock: true},
		{Name: "cheap-out-of-stock", Price: 5, InStock: false},
		{Name: "mid-in-stock", Price: 50, InStock: true},
		{Name: "expensive-in-stock", Price: 500, InStock: true},
		{Name: "expensive-out-of-stock", Price: 500, InStock: false},
	}
}

func names(ps []Product) []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = p.Name
	}
	return out
}

func TestAnd_BothMustBeSatisfied(t *testing.T) {
	spec := And[Product](minPriceSpec{min: 10}, inStockSpec{})
	got := names(Filter(products(), spec))
	want := []string{"mid-in-stock", "expensive-in-stock"}
	assertStringSliceEqual(t, got, want)
}

func TestOr_EitherSatisfied(t *testing.T) {
	spec := Or[Product](maxPriceSpec{max: 5}, minPriceSpec{min: 500})
	got := names(Filter(products(), spec))
	want := []string{"cheap-in-stock", "cheap-out-of-stock", "expensive-in-stock", "expensive-out-of-stock"}
	assertStringSliceEqual(t, got, want)
}

func TestNot_InvertsResult(t *testing.T) {
	spec := Not[Product](inStockSpec{})
	got := names(Filter(products(), spec))
	want := []string{"cheap-out-of-stock", "expensive-out-of-stock"}
	assertStringSliceEqual(t, got, want)
}

func TestFilter_PreservesOriginalOrder(t *testing.T) {
	spec := inStockSpec{}
	got := names(Filter(products(), spec))
	want := []string{"cheap-in-stock", "mid-in-stock", "expensive-in-stock"}
	assertStringSliceEqual(t, got, want)
}

func TestFilter_EmptyResultWhenNothingMatches(t *testing.T) {
	spec := minPriceSpec{min: 10000}
	got := Filter(products(), spec)
	if len(got) != 0 {
		t.Fatalf("Filter() = %v, want empty", got)
	}
}

func TestFilter_DoesNotMutateInput(t *testing.T) {
	original := products()
	originalCopy := append([]Product(nil), original...)

	Filter(original, inStockSpec{})

	for i := range original {
		if original[i] != originalCopy[i] {
			t.Fatalf("input slice was mutated: original[%d] = %v, want %v", i, original[i], originalCopy[i])
		}
	}
}

func TestNestedComposition_AndOfOrOfNot(t *testing.T) {
	// "in stock AND (cheap OR expensive)" — i.e. exclude the mid-priced item.
	cheapOrExpensive := Or[Product](maxPriceSpec{max: 5}, minPriceSpec{min: 500})
	spec := And[Product](inStockSpec{}, cheapOrExpensive)

	got := names(Filter(products(), spec))
	want := []string{"cheap-in-stock", "expensive-in-stock"}
	assertStringSliceEqual(t, got, want)
}

func TestNestedComposition_NotOfAnd(t *testing.T) {
	// NOT (in stock AND expensive) -- i.e. anything that's either out of stock or not
	// expensive.
	spec := Not[Product](And[Product](inStockSpec{}, minPriceSpec{min: 500}))
	got := names(Filter(products(), spec))
	want := []string{"cheap-in-stock", "cheap-out-of-stock", "mid-in-stock", "expensive-out-of-stock"}
	assertStringSliceEqual(t, got, want)
}

func TestSpecification_DoubleNegationEqualsOriginal(t *testing.T) {
	spec := inStockSpec{}
	doubleNegated := Not[Product](Not[Product](spec))
	if doubleNegated == nil {
		t.Fatal("Not(Not(spec)) = nil, want a valid Specification")
	}

	for _, p := range products() {
		if spec.IsSatisfiedBy(p) != doubleNegated.IsSatisfiedBy(p) {
			t.Fatalf("Not(Not(spec)).IsSatisfiedBy(%v) = %v, want %v (should equal spec.IsSatisfiedBy)", p, doubleNegated.IsSatisfiedBy(p), spec.IsSatisfiedBy(p))
		}
	}
}

func assertStringSliceEqual(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
