package repository

import "errors"

/*
See README.md for the full problem statement, requirements, and rubric.

UserRepository and User are fixed — don't change their shape.
*/

var (
	ErrNotFound       = errors.New("not found")
	ErrDuplicateID    = errors.New("duplicate id")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	ID    string
	Email string
	Name  string
}

// UserRepository is a collection-like abstraction over persistence. Business logic
// (UserService) depends only on this interface, never on a concrete implementation.
type UserRepository interface {
	Add(u User) error
	GetByID(id string) (User, error)
	Update(u User) error
	Delete(id string) error
	List() []User
	FindByEmail(email string) (User, error)
}

// TODO: implement InMemoryUserRepository, satisfying UserRepository. See README.md
// for exact error semantics and the List ordering requirement.
type InMemoryUserRepository struct {
	// TODO: add your fields here
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	// TODO: implement
	return &InMemoryUserRepository{}
}

func (r *InMemoryUserRepository) Add(u User) error {
	// TODO: implement
	return nil
}

func (r *InMemoryUserRepository) GetByID(id string) (User, error) {
	// TODO: implement
	return User{}, nil
}

func (r *InMemoryUserRepository) Update(u User) error {
	// TODO: implement
	return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
	// TODO: implement
	return nil
}

func (r *InMemoryUserRepository) List() []User {
	// TODO: implement
	return nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (User, error) {
	// TODO: implement
	return User{}, nil
}

// TODO: implement UserService. It must depend ONLY on the UserRepository interface
// (injected via NewUserService), never on InMemoryUserRepository directly.
type UserService struct {
	// TODO: add your fields here
}

func NewUserService(repo UserRepository) *UserService {
	// TODO: implement
	return &UserService{}
}

func (s *UserService) Register(id, email, name string) error {
	// TODO: implement
	return nil
}

func (s *UserService) Rename(id, newName string) error {
	// TODO: implement
	return nil
}
