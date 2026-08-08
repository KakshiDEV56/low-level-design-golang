# Dependency Inversion Principle (DIP)

> High-level modules should not depend on low-level modules; both should depend on
> abstractions. Abstractions should not depend on details; details should depend on
> abstractions.

## The violation in `solution.go`

`OrderService` is a high-level policy (the *what* — "placing an order saves it
somewhere"). `MySQLDatabase` is a low-level implementation detail (the *how* —
specifically MySQL). Right now, `NewOrderService` **directly constructs** a
`MySQLDatabase` inside itself, which means:

- `OrderService` cannot be reused with any other storage mechanism (Postgres, an
  in-memory cache, S3, a message queue) without editing `OrderService` itself.
- `OrderService` cannot be **unit tested** without a real (or realistically faked)
  MySQL connection, because the concrete dependency is baked in — there's no seam to
  substitute a test double.

This is the most common, most consequential SOLID violation in real production code,
precisely because of that second point: **DIP violations are usually discovered as "why
can't I test this without spinning up a database."**

## What a fix looks like

Introduce a `Storage` interface that describes only what `OrderService` actually needs
(`Save`). `OrderService` depends on that abstraction, not on `MySQLDatabase` directly —
and the concrete storage implementation is **injected** (passed in as a constructor
argument) rather than constructed internally. `MySQLDatabase` still exists and still
implements `Storage` — it just no longer sits "above" `OrderService` in the dependency
graph; both now depend on the same abstraction.

The test file proves the inversion worked by injecting a completely fake `Storage`
implementation (no real database anywhere) and testing `OrderService`'s behavior in
isolation — including forcing a save failure, which would be awkward-to-impossible to
trigger reliably against a real MySQL instance in a unit test.

## Target shape

```go
type Storage interface {
    Save(order string) error
}

type MySQLDatabase struct{}
func (MySQLDatabase) Save(order string) error { ... }

type OrderService struct {
    storage Storage
}
func NewOrderService(storage Storage) *OrderService
func (s *OrderService) PlaceOrder(order string) error
```
