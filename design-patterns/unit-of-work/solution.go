package unitofwork

import (
	"errors"
	"sync"
)

/*
See README.md for the full problem statement, requirements, and rubric.

Account and AccountStore below are already complete — do not modify them.
*/

var (
	ErrNotFound          = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type Account struct {
	ID      string
	Balance int
}

// AccountStore is the committed backing store (the "database"). Do not modify.
type AccountStore struct {
	mu       sync.Mutex
	accounts map[string]Account
}

func NewAccountStore() *AccountStore {
	return &AccountStore{accounts: make(map[string]Account)}
}

// Seed inserts an account directly into committed storage (test/setup helper — real
// code would normally only ever write via a UnitOfWork's Commit).
func (s *AccountStore) Seed(a Account) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accounts[a.ID] = a
}

func (s *AccountStore) Get(id string) (Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.accounts[id]
	if !ok {
		return Account{}, ErrNotFound
	}
	return a, nil
}

func (s *AccountStore) Save(a Account) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accounts[a.ID] = a
}

// TODO: implement UnitOfWork.
//   - NewUnitOfWork(store *AccountStore) *UnitOfWork
//   - Get(id string) (Account, error)
//   - Save(a Account)
//   - Commit() error
//   - Rollback()
//
// See README.md for exact semantics of each method.
type UnitOfWork struct {
	// TODO: add your fields here
}

func NewUnitOfWork(store *AccountStore) *UnitOfWork {
	// TODO: implement
	return &UnitOfWork{}
}

func (u *UnitOfWork) Get(id string) (Account, error) {
	// TODO: implement
	return Account{}, nil
}

func (u *UnitOfWork) Save(a Account) {
	// TODO: implement
}

func (u *UnitOfWork) Commit() error {
	// TODO: implement
	return nil
}

func (u *UnitOfWork) Rollback() {
	// TODO: implement
}

// TODO: implement Transfer using a UnitOfWork so the debit and credit either both
// happen (Commit) or neither does (Rollback). See README.md for the exact error
// conditions (ErrInsufficientFunds, ErrNotFound propagated from Get).
func Transfer(store *AccountStore, fromID, toID string, amount int) error {
	// TODO: implement
	return nil
}
