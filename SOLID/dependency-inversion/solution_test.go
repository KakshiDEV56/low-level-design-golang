package dependencyinversion

import (
	"errors"
	"testing"
)

// Compile-time proof MySQLDatabase still implements Storage after the refactor —
// production code isn't broken by the fix.
var _ Storage = MySQLDatabase{}

// fakeStorage is a test double with NO real database anywhere — it just records what
// was saved (or returns a configured error), entirely in memory. Being able to write
// this and have it work with OrderService is the direct, tangible payoff of DIP.
type fakeStorage struct {
	saved   []string
	failWith error
}

func (f *fakeStorage) Save(order string) error {
	if f.failWith != nil {
		return f.failWith
	}
	f.saved = append(f.saved, order)
	return nil
}

func TestOrderService_PlaceOrderSavesViaInjectedStorage(t *testing.T) {
	fake := &fakeStorage{}
	svc := NewOrderService(fake)

	if err := svc.PlaceOrder("order-1"); err != nil {
		t.Fatalf("PlaceOrder returned error: %v", err)
	}

	if len(fake.saved) != 1 || fake.saved[0] != "order-1" {
		t.Fatalf("fakeStorage.saved = %v, want [\"order-1\"]", fake.saved)
	}
}

func TestOrderService_MultipleOrdersAllRecorded(t *testing.T) {
	fake := &fakeStorage{}
	svc := NewOrderService(fake)

	orders := []string{"order-a", "order-b", "order-c"}
	for _, o := range orders {
		if err := svc.PlaceOrder(o); err != nil {
			t.Fatalf("PlaceOrder(%q) returned error: %v", o, err)
		}
	}

	if len(fake.saved) != len(orders) {
		t.Fatalf("fakeStorage.saved = %v, want %v", fake.saved, orders)
	}
	for i, o := range orders {
		if fake.saved[i] != o {
			t.Fatalf("fakeStorage.saved[%d] = %q, want %q", i, fake.saved[i], o)
		}
	}
}

func TestOrderService_PropagatesStorageError(t *testing.T) {
	wantErr := errors.New("disk full")
	fake := &fakeStorage{failWith: wantErr}
	svc := NewOrderService(fake)

	err := svc.PlaceOrder("order-1")
	if !errors.Is(err, wantErr) {
		t.Fatalf("PlaceOrder error = %v, want %v", err, wantErr)
	}
	if len(fake.saved) != 0 {
		t.Fatalf("fakeStorage.saved = %v, want empty (the save should have failed)", fake.saved)
	}
}

func TestOrderService_WorksWithRealMySQLDatabaseToo(t *testing.T) {
	// Proves the fix didn't break production usage: OrderService still works when
	// given the real (well, "real" in this exercise) MySQLDatabase.
	svc := NewOrderService(MySQLDatabase{})
	if err := svc.PlaceOrder("order-1"); err != nil {
		t.Fatalf("PlaceOrder with MySQLDatabase returned error: %v", err)
	}
}
