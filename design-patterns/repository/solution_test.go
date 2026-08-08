package repository

import (
	"errors"
	"reflect"
	"testing"
)

// Compile-time proof InMemoryUserRepository satisfies UserRepository.
var _ UserRepository = (*InMemoryUserRepository)(nil)

func TestInMemoryUserRepository_AddAndGetByID(t *testing.T) {
	r := NewInMemoryUserRepository()
	u := User{ID: "1", Email: "a@example.com", Name: "Alice"}
	if err := r.Add(u); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}
	got, err := r.GetByID("1")
	if err != nil || got != u {
		t.Fatalf("GetByID(1) = (%v, %v), want (%v, nil)", got, err, u)
	}
}

func TestInMemoryUserRepository_AddDuplicateID(t *testing.T) {
	r := NewInMemoryUserRepository()
	r.Add(User{ID: "1", Email: "a@example.com", Name: "Alice"})
	err := r.Add(User{ID: "1", Email: "b@example.com", Name: "Bob"})
	if !errors.Is(err, ErrDuplicateID) {
		t.Fatalf("Add() with duplicate ID = %v, want ErrDuplicateID", err)
	}
}

func TestInMemoryUserRepository_GetByIDMissing(t *testing.T) {
	r := NewInMemoryUserRepository()
	_, err := r.GetByID("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetByID(missing) = %v, want ErrNotFound", err)
	}
}

func TestInMemoryUserRepository_Update(t *testing.T) {
	r := NewInMemoryUserRepository()
	r.Add(User{ID: "1", Email: "a@example.com", Name: "Alice"})
	if err := r.Update(User{ID: "1", Email: "a@example.com", Name: "Alicia"}); err != nil {
		t.Fatalf("Update() failed: %v", err)
	}
	got, _ := r.GetByID("1")
	if got.Name != "Alicia" {
		t.Fatalf("Name after Update = %q, want %q", got.Name, "Alicia")
	}
}

func TestInMemoryUserRepository_UpdateMissing(t *testing.T) {
	r := NewInMemoryUserRepository()
	err := r.Update(User{ID: "missing"})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Update(missing) = %v, want ErrNotFound", err)
	}
}

func TestInMemoryUserRepository_Delete(t *testing.T) {
	r := NewInMemoryUserRepository()
	r.Add(User{ID: "1", Email: "a@example.com", Name: "Alice"})
	if err := r.Delete("1"); err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}
	if _, err := r.GetByID("1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetByID after Delete = %v, want ErrNotFound", err)
	}
}

func TestInMemoryUserRepository_DeleteMissing(t *testing.T) {
	r := NewInMemoryUserRepository()
	if err := r.Delete("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete(missing) = %v, want ErrNotFound", err)
	}
}

func TestInMemoryUserRepository_ListSortedByID(t *testing.T) {
	r := NewInMemoryUserRepository()
	r.Add(User{ID: "3", Email: "c@example.com", Name: "Carol"})
	r.Add(User{ID: "1", Email: "a@example.com", Name: "Alice"})
	r.Add(User{ID: "2", Email: "b@example.com", Name: "Bob"})

	got := r.List()
	want := []User{
		{ID: "1", Email: "a@example.com", Name: "Alice"},
		{ID: "2", Email: "b@example.com", Name: "Bob"},
		{ID: "3", Email: "c@example.com", Name: "Carol"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %v, want %v", got, want)
	}
}

func TestInMemoryUserRepository_FindByEmail(t *testing.T) {
	r := NewInMemoryUserRepository()
	r.Add(User{ID: "1", Email: "a@example.com", Name: "Alice"})

	got, err := r.FindByEmail("a@example.com")
	if err != nil || got.ID != "1" {
		t.Fatalf("FindByEmail() = (%v, %v), want ID=1, nil", got, err)
	}
	if _, err := r.FindByEmail("missing@example.com"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("FindByEmail(missing) = %v, want ErrNotFound", err)
	}
}

// ---- UserService ----

func TestUserService_RegisterSuccess(t *testing.T) {
	repo := NewInMemoryUserRepository()
	svc := NewUserService(repo)

	if err := svc.Register("1", "a@example.com", "Alice"); err != nil {
		t.Fatalf("Register() failed: %v", err)
	}
	got, err := repo.GetByID("1")
	if err != nil || got.Email != "a@example.com" {
		t.Fatalf("repo state after Register = (%v, %v), want a user with the registered email", got, err)
	}
}

func TestUserService_RegisterDuplicateEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()
	svc := NewUserService(repo)
	svc.Register("1", "a@example.com", "Alice")

	err := svc.Register("2", "a@example.com", "Alicia-Clone")
	if !errors.Is(err, ErrDuplicateEmail) {
		t.Fatalf("Register() with a taken email = %v, want ErrDuplicateEmail", err)
	}
	if _, err := repo.GetByID("2"); !errors.Is(err, ErrNotFound) {
		t.Fatal("a user was added to the repo despite Register() failing on duplicate email")
	}
}

func TestUserService_RenameSuccess(t *testing.T) {
	repo := NewInMemoryUserRepository()
	svc := NewUserService(repo)
	svc.Register("1", "a@example.com", "Alice")

	if err := svc.Rename("1", "Alicia"); err != nil {
		t.Fatalf("Rename() failed: %v", err)
	}
	got, _ := repo.GetByID("1")
	if got.Name != "Alicia" {
		t.Fatalf("Name after Rename = %q, want %q", got.Name, "Alicia")
	}
}

func TestUserService_RenameUnknownUser(t *testing.T) {
	repo := NewInMemoryUserRepository()
	svc := NewUserService(repo)
	if err := svc.Rename("missing", "X"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Rename(missing) = %v, want ErrNotFound", err)
	}
}

// fakeUserRepository is a completely independent UserRepository implementation,
// defined only in the test file, backed by a slice instead of a map — deliberately a
// different internal structure than InMemoryUserRepository. If UserService works with
// this unmodified, that's the concrete proof of the Repository abstraction.
type fakeUserRepository struct {
	users []User
}

func (f *fakeUserRepository) Add(u User) error {
	for _, existing := range f.users {
		if existing.ID == u.ID {
			return ErrDuplicateID
		}
	}
	f.users = append(f.users, u)
	return nil
}

func (f *fakeUserRepository) GetByID(id string) (User, error) {
	for _, u := range f.users {
		if u.ID == id {
			return u, nil
		}
	}
	return User{}, ErrNotFound
}

func (f *fakeUserRepository) Update(u User) error {
	for i, existing := range f.users {
		if existing.ID == u.ID {
			f.users[i] = u
			return nil
		}
	}
	return ErrNotFound
}

func (f *fakeUserRepository) Delete(id string) error {
	for i, u := range f.users {
		if u.ID == id {
			f.users = append(f.users[:i], f.users[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (f *fakeUserRepository) List() []User { return f.users }

func (f *fakeUserRepository) FindByEmail(email string) (User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return User{}, ErrNotFound
}

func TestUserService_WorksWithAnyRepositoryImplementation(t *testing.T) {
	repo := &fakeUserRepository{}
	svc := NewUserService(repo)

	if err := svc.Register("1", "a@example.com", "Alice"); err != nil {
		t.Fatalf("Register() with a custom repository implementation failed: %v", err)
	}
	if err := svc.Register("2", "a@example.com", "Bob"); !errors.Is(err, ErrDuplicateEmail) {
		t.Fatalf("Register() duplicate email with custom repository = %v, want ErrDuplicateEmail", err)
	}
	if err := svc.Rename("1", "Alicia"); err != nil {
		t.Fatalf("Rename() with a custom repository implementation failed: %v", err)
	}
	got, _ := repo.GetByID("1")
	if got.Name != "Alicia" {
		t.Fatalf("Name after Rename via custom repository = %q, want %q", got.Name, "Alicia")
	}
}
