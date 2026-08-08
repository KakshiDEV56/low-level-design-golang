package cqrs

import "errors"

/*
See README.md for the full problem statement, requirements, and rubric.

StockEvent is fixed — don't change its shape.
*/

var ErrInsufficientStock = errors.New("insufficient stock")

type StockEvent struct {
	ProductID string
	Delta     int
}

// TODO: implement CommandHandler (the write side).
//   - NewCommandHandler() *CommandHandler
//   - AddStock(productID string, qty int) error
//   - RemoveStock(productID string, qty int) error
//   - Subscribe(handler func(StockEvent))
//
// CommandHandler must track its own stock state to enforce business rules (qty > 0,
// no removing more than is in stock) — it is the source of truth, independent of
// whatever QueryHandler(s) happen to be subscribed.
type CommandHandler struct {
	// TODO: add your fields here
}

func NewCommandHandler() *CommandHandler {
	// TODO: implement
	return &CommandHandler{}
}

func (c *CommandHandler) AddStock(productID string, qty int) error {
	// TODO: implement
	return nil
}

func (c *CommandHandler) RemoveStock(productID string, qty int) error {
	// TODO: implement
	return nil
}

func (c *CommandHandler) Subscribe(handler func(StockEvent)) {
	// TODO: implement
}

// TODO: implement QueryHandler (the read side). It must have NO direct access to
// CommandHandler's state — the only way it learns anything is via Apply.
//   - NewQueryHandler() *QueryHandler
//   - Apply(e StockEvent)
//   - GetStock(productID string) int
//   - TotalStock() int
//   - LowStockProducts(threshold int) []string
type QueryHandler struct {
	// TODO: add your fields here
}

func NewQueryHandler() *QueryHandler {
	// TODO: implement
	return &QueryHandler{}
}

func (q *QueryHandler) Apply(e StockEvent) {
	// TODO: implement
}

func (q *QueryHandler) GetStock(productID string) int {
	// TODO: implement
	return 0
}

func (q *QueryHandler) TotalStock() int {
	// TODO: implement
	return 0
}

func (q *QueryHandler) LowStockProducts(threshold int) []string {
	// TODO: implement
	return nil
}
