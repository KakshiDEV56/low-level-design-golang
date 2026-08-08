package openclosed

import "testing"

func TestApplyDiscount_RegularCustomer(t *testing.T) {
	got := ApplyDiscount(RegularCustomer{}, 100)
	if got != 100 {
		t.Fatalf("ApplyDiscount(RegularCustomer, 100) = %v, want 100 (0%% discount)", got)
	}
}

func TestApplyDiscount_PremiumCustomer(t *testing.T) {
	got := ApplyDiscount(PremiumCustomer{}, 100)
	if got != 90 {
		t.Fatalf("ApplyDiscount(PremiumCustomer, 100) = %v, want 90 (10%% discount)", got)
	}
}

func TestApplyDiscount_VIPCustomer(t *testing.T) {
	got := ApplyDiscount(VIPCustomer{}, 100)
	if got != 80 {
		t.Fatalf("ApplyDiscount(VIPCustomer, 100) = %v, want 80 (20%% discount)", got)
	}
}

func TestApplyDiscount_ZeroAmount(t *testing.T) {
	for _, d := range []Discounter{RegularCustomer{}, PremiumCustomer{}, VIPCustomer{}} {
		if got := ApplyDiscount(d, 0); got != 0 {
			t.Fatalf("ApplyDiscount(%T, 0) = %v, want 0", d, got)
		}
	}
}

// employeeCustomer is a BRAND NEW customer type, defined here in the test file — it
// did not exist when solution.go was written, and solution.go is never modified to
// know about it. If this compiles and works with the unmodified ApplyDiscount, that's
// the proof the calculator is genuinely open for extension without modification.
type employeeCustomer struct{}

func (employeeCustomer) Discount(amount float64) float64 {
	return amount * 0.5 // employees get 50% off
}

func TestApplyDiscount_NewCustomerTypeAddedWithoutModifyingSolution(t *testing.T) {
	got := ApplyDiscount(employeeCustomer{}, 100)
	if got != 50 {
		t.Fatalf("ApplyDiscount(employeeCustomer, 100) = %v, want 50 (50%% discount)", got)
	}
}

// A second brand-new type, to really hammer the point: as many customer types as the
// business wants can be added, forever, without ApplyDiscount or Discounter changing.
type loyaltyMemberCustomer struct{ yearsActive int }

func (c loyaltyMemberCustomer) Discount(amount float64) float64 {
	if c.yearsActive >= 5 {
		return amount * 0.3
	}
	return amount * 0.15
}

func TestApplyDiscount_AnotherNewCustomerType(t *testing.T) {
	newMember := loyaltyMemberCustomer{yearsActive: 1}
	if got := ApplyDiscount(newMember, 100); got != 85 {
		t.Fatalf("ApplyDiscount(new loyalty member, 100) = %v, want 85", got)
	}

	longTimeMember := loyaltyMemberCustomer{yearsActive: 10}
	if got := ApplyDiscount(longTimeMember, 100); got != 70 {
		t.Fatalf("ApplyDiscount(long-time loyalty member, 100) = %v, want 70", got)
	}
}
