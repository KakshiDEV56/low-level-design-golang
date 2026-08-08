package transactionalkvstore

import "testing"

func TestStore_BasicSetGetWithoutTransaction(t *testing.T) {
	s := NewStore()
	s.Set("x", "1")
	if got, ok := s.Get("x"); !ok || got != "1" {
		t.Fatalf("Get(x) = (%q, %v), want (\"1\", true)", got, ok)
	}
}

func TestStore_GetMissingKey(t *testing.T) {
	s := NewStore()
	if _, ok := s.Get("missing"); ok {
		t.Fatal("Get(missing) ok = true, want false")
	}
}

func TestStore_DeleteWithoutTransaction(t *testing.T) {
	s := NewStore()
	s.Set("x", "1")
	s.Delete("x")
	if _, ok := s.Get("x"); ok {
		t.Fatal("Get(x) after Delete ok = true, want false")
	}
}

func TestStore_InTransactionReflectsNestingDepth(t *testing.T) {
	s := NewStore()
	if s.InTransaction() {
		t.Fatal("InTransaction() at start = true, want false")
	}
	s.Begin()
	if !s.InTransaction() {
		t.Fatal("InTransaction() after Begin = false, want true")
	}
	s.Begin()
	if !s.InTransaction() {
		t.Fatal("InTransaction() after nested Begin = false, want true")
	}
	s.Rollback()
	if !s.InTransaction() {
		t.Fatal("InTransaction() after rolling back the inner transaction = false, want true (outer still open)")
	}
	s.Commit()
	if s.InTransaction() {
		t.Fatal("InTransaction() after closing the outermost transaction = true, want false")
	}
}

func TestStore_SetInsideTransactionVisibleImmediately(t *testing.T) {
	s := NewStore()
	s.Begin()
	s.Set("x", "1")
	if got, ok := s.Get("x"); !ok || got != "1" {
		t.Fatalf("Get(x) inside the transaction = (%q, %v), want (\"1\", true)", got, ok)
	}
}

func TestStore_RollbackDiscardsChanges(t *testing.T) {
	s := NewStore()
	s.Set("x", "1")
	s.Begin()
	s.Set("x", "2")
	s.Set("y", "new")
	if !s.Rollback() {
		t.Fatal("Rollback() on an open transaction = false, want true")
	}

	if got, ok := s.Get("x"); !ok || got != "1" {
		t.Fatalf("Get(x) after Rollback = (%q, %v), want (\"1\", true) — pre-transaction value restored", got, ok)
	}
	if _, ok := s.Get("y"); ok {
		t.Fatal("Get(y) after Rollback ok = true, want false — y was never committed")
	}
}

func TestStore_CommitPersistsChanges(t *testing.T) {
	s := NewStore()
	s.Begin()
	s.Set("x", "1")
	if !s.Commit() {
		t.Fatal("Commit() on an open transaction = false, want true")
	}

	if s.InTransaction() {
		t.Fatal("InTransaction() after Commit = true, want false")
	}
	if got, ok := s.Get("x"); !ok || got != "1" {
		t.Fatalf("Get(x) after Commit (no transaction open) = (%q, %v), want (\"1\", true)", got, ok)
	}
}

func TestStore_DeleteInsideTransactionHidesOuterValue(t *testing.T) {
	s := NewStore()
	s.Set("x", "1")
	s.Begin()
	s.Delete("x")
	if _, ok := s.Get("x"); ok {
		t.Fatal("Get(x) after Delete inside a transaction ok = true, want false")
	}
	s.Rollback()
	if got, ok := s.Get("x"); !ok || got != "1" {
		t.Fatalf("Get(x) after rolling back the delete = (%q, %v), want (\"1\", true)", got, ok)
	}
}

func TestStore_DeleteCommittedBecomesPermanent(t *testing.T) {
	s := NewStore()
	s.Set("x", "1")
	s.Begin()
	s.Delete("x")
	s.Commit()

	if _, ok := s.Get("x"); ok {
		t.Fatal("Get(x) after committing a delete ok = true, want false")
	}
}

func TestStore_CommitOrRollbackWithNoTransactionReturnsFalse(t *testing.T) {
	s := NewStore()
	if s.Commit() {
		t.Fatal("Commit() with no open transaction = true, want false")
	}
	if s.Rollback() {
		t.Fatal("Rollback() with no open transaction = true, want false")
	}
}

func TestStore_NestedTransactionsInnerWinsUntilRolledBack(t *testing.T) {
	s := NewStore()
	s.Set("x", "0")
	s.Begin() // T1
	s.Set("x", "1")
	s.Begin() // T2
	s.Set("x", "2")

	if got, _ := s.Get("x"); got != "2" {
		t.Fatalf("Get(x) with T2 open = %q, want \"2\"", got)
	}

	s.Rollback() // discard T2
	if got, _ := s.Get("x"); got != "1" {
		t.Fatalf("Get(x) after rolling back T2 = %q, want \"1\" (T1's value)", got)
	}

	s.Commit() // fold T1 into permanent storage
	if s.InTransaction() {
		t.Fatal("InTransaction() after committing the outermost transaction = true, want false")
	}
	if got, _ := s.Get("x"); got != "1" {
		t.Fatalf("Get(x) after committing T1 = %q, want \"1\"", got)
	}
}

func TestStore_NestedTransactionsRollbackOuterDiscardsBoth(t *testing.T) {
	s := NewStore()
	s.Set("x", "0")
	s.Begin() // T1
	s.Set("x", "1")
	s.Begin() // T2
	s.Set("x", "2")
	s.Commit() // fold T2 into T1: x=2 inside T1

	if got, _ := s.Get("x"); got != "2" {
		t.Fatalf("Get(x) after committing T2 into T1 = %q, want \"2\"", got)
	}

	s.Rollback() // discard T1 (and everything folded into it) entirely
	if got, _ := s.Get("x"); got != "0" {
		t.Fatalf("Get(x) after rolling back T1 = %q, want \"0\" (original value, untouched)", got)
	}
}

func TestStore_UntouchedKeyFallsThroughAllLayers(t *testing.T) {
	s := NewStore()
	s.Set("untouched", "base-value")
	s.Begin()
	s.Set("other", "1")
	s.Begin()
	s.Set("another", "2")

	if got, ok := s.Get("untouched"); !ok || got != "base-value" {
		t.Fatalf("Get(untouched) through two open transactions = (%q, %v), want (\"base-value\", true)", got, ok)
	}
}

func TestStore_MultipleKeysIndependentAcrossTransaction(t *testing.T) {
	s := NewStore()
	s.Set("a", "1")
	s.Set("b", "2")
	s.Begin()
	s.Set("a", "100")
	s.Delete("b")
	s.Set("c", "3")

	aVal, _ := s.Get("a")
	_, bOk := s.Get("b")
	cVal, _ := s.Get("c")
	if aVal != "100" || bOk || cVal != "3" {
		t.Fatalf("mid-transaction state = a:%q b-present:%v c:%q, want a:100 b-present:false c:3", aVal, bOk, cVal)
	}

	s.Rollback()
	aVal, _ = s.Get("a")
	bVal, bOk := s.Get("b")
	_, cOk := s.Get("c")
	if aVal != "1" || !bOk || bVal != "2" || cOk {
		t.Fatalf("post-rollback state = a:%q b:(%q,%v) c-present:%v, want a:1 b:(2,true) c-present:false", aVal, bVal, bOk, cOk)
	}
}
