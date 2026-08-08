package dependencyinversion

/*
This code compiles and works today — but it violates the Dependency Inversion
Principle. See README.md for the full explanation.

OrderService (a high-level policy) directly constructs and depends on the concrete
MySQLDatabase (a low-level detail). It cannot be reused with any other storage
mechanism, and it cannot be unit tested without a real database.

TODO: refactor this file so that:
 1. A Storage interface is introduced:
      type Storage interface {
          Save(order string) error
      }
 2. MySQLDatabase implements Storage (its Save method can keep its current trivial
    body).
 3. OrderService holds a Storage (the interface, not MySQLDatabase directly) as a
    field.
 4. NewOrderService takes a Storage as a PARAMETER and stores it — it must NOT
    construct a MySQLDatabase internally anymore.
 5. PlaceOrder delegates to the injected Storage's Save method and returns whatever
    error it returns.

Once refactored, solution_test.go (which injects its OWN fake Storage implementation,
with no real database involved) should compile and pass.
*/

// MySQLDatabase is a concrete, low-level storage detail.
type MySQLDatabase struct{}

func (MySQLDatabase) Save(order string) error {
	// Pretend this writes to a real MySQL database.
	return nil
}

// OrderService is a high-level policy that directly constructs and depends on the
// concrete MySQLDatabase — it cannot be tested or reused with any other storage
// mechanism without modifying OrderService itself.
// TODO: depend on an abstraction instead (see above).
type OrderService struct {
	db MySQLDatabase
}

func NewOrderService() *OrderService {
	return &OrderService{db: MySQLDatabase{}}
}

func (s *OrderService) PlaceOrder(order string) error {
	return s.db.Save(order)
}
