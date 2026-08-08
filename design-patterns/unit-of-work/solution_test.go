package unitofwork

import (
	"errors"
	"testing"
)

func TestUnitOfWork_GetFallsThroughToStore(t *testing.T) {
	store := NewAccountStore()
	store.Seed(Account{ID: "a", Balance: 100})

	uow := NewUnitOfWork(store)
	got, err := uow.Get("a")
	if err != nil || got.Balance != 100 {
		t.Fatalf("Get(a) = (%v, %v), want (Balance:100, nil)", got, err)
	}
}

func TestUnitOfWork_GetUnknownAccountReturnsError(t *testing.T) {
	store := NewAccountStore()
	uow := NewUnitOfWork(store)
	if _, err := uow.Get("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get(missing) = %v, want ErrNotFound", err)
	}
}

func TestUnitOfWork_SaveNotVisibleInStoreUntilCommit(t *testing.T) {
	store := NewAccountStore()
	store.Seed(Account{ID: "a", Balance: 100})

	uow := NewUnitOfWork(store)
	uow.Save(Account{ID: "a", Balance: 999})

	// The store itself must be untouched until Commit.
	got, _ := store.Get("a")
	if got.Balance != 100 {
		t.Fatalf("store.Get(a).Balance = %d before Commit, want unchanged 100", got.Balance)
	}
}

func TestUnitOfWork_GetSeesOwnStagedWrite(t *testing.T) {
	store := NewAccountStore()
	store.Seed(Account{ID: "a", Balance: 100})

	uow := NewUnitOfWork(store)
	uow.Save(Account{ID: "a", Balance: 999})

	got, err := uow.Get("a")
	if err != nil || got.Balance != 999 {
		t.Fatalf("uow.Get(a) after Save = (%v, %v), want (Balance:999, nil)", got, err)
	}
}

func TestUnitOfWork_CommitAppliesToStore(t *testing.T) {
	store := NewAccountStore()
	store.Seed(Account{ID: "a", Balance: 100})

	uow := NewUnitOfWork(store)
	uow.Save(Account{ID: "a", Balance: 999})
	if err := uow.Commit(); err != nil {
		t.Fatalf("Commit() failed: %v", err)
	}

	got, _ := store.Get("a")
	if got.Balance != 999 {
		t.Fatalf("store.Get(a).Balance after Commit = %d, want 999", got.Balance)
	}
}

func TestUnitOfWork_RollbackDiscardsStagedWrites(t *testing.T) {
	store := NewAccountStore()
	store.Seed(Account{ID: "a", Balance: 100})

	uow := NewUnitOfWork(store)
	uow.Save(Account{ID: "a", Balance: 999})
	uow.Rollback()

	got, _ := store.Get("a")
	if got.Balance != 100 {
		t.Fatalf("store.Get(a).Balance after Rollback = %d, want unchanged 100", got.Balance)
	}
}

func TestUnitOfWork_IndependentFromASeparateUnitOfWork(t *testing.T) {
	store := NewAccountStore()
	store.Seed(Account{ID: "a", Balance: 100})

	uow1 := NewUnitOfWork(store)
	uow1.Save(Account{ID: "a", Balance: 999})

	uow2 := NewUnitOfWork(store)
	got, _ := uow2.Get("a")
	if got.Balance != 100 {
		t.Fatalf("a separate UnitOfWork's Get(a) = %d, want 100 (must not see uow1's uncommitted staged write)", got.Balance)
	}
}

func TestTransfer_SuccessUpdatesBothAccounts(t *testing.T) {
	store := NewAccountStore()
	store.Seed(Account{ID: "alice", Balance: 100})
	store.Seed(Account{ID: "bob", Balance: 50})

	if err := Transfer(store, "alice", "bob", 30); err != nil {
		t.Fatalf("Transfer() failed: %v", err)
	}

	alice, _ := store.Get("alice")
	bob, _ := store.Get("bob")
	if alice.Balance != 70 {
		t.Fatalf("alice.Balance = %d, want 70", alice.Balance)
	}
	if bob.Balance != 80 {
		t.Fatalf("bob.Balance = %d, want 80", bob.Balance)
	}
}

func TestTransfer_InsufficientFundsLeavesBothAccountsUnchanged(t *testing.T) {
	store := NewAccountStore()
	store.Seed(Account{ID: "alice", Balance: 10})
	store.Seed(Account{ID: "bob", Balance: 50})

	err := Transfer(store, "alice", "bob", 100)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf("Transfer() with insufficient funds = %v, want ErrInsufficientFunds", err)
	}

	alice, _ := store.Get("alice")
	bob, _ := store.Get("bob")
	if alice.Balance != 10 {
		t.Fatalf("alice.Balance after failed Transfer = %d, want unchanged 10", alice.Balance)
	}
	if bob.Balance != 50 {
		t.Fatalf("bob.Balance after failed Transfer = %d, want unchanged 50", bob.Balance)
	}
}

func TestTransfer_UnknownFromAccountLeavesToAccountUnchanged(t *testing.T) {
	store := NewAccountStore()
	store.Seed(Account{ID: "bob", Balance: 50})

	err := Transfer(store, "missing", "bob", 10)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Transfer() with unknown source account = %v, want ErrNotFound", err)
	}
	bob, _ := store.Get("bob")
	if bob.Balance != 50 {
		t.Fatalf("bob.Balance after failed Transfer = %d, want unchanged 50", bob.Balance)
	}
}

func TestTransfer_UnknownToAccountLeavesFromAccountUnchanged(t *testing.T) {
	store := NewAccountStore()
	store.Seed(Account{ID: "alice", Balance: 100})

	err := Transfer(store, "alice", "missing", 10)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Transfer() with unknown destination account = %v, want ErrNotFound", err)
	}
	alice, _ := store.Get("alice")
	if alice.Balance != 100 {
		t.Fatalf("alice.Balance after failed Transfer = %d, want unchanged 100 (must not be debited if the credit can't happen)", alice.Balance)
	}
}
