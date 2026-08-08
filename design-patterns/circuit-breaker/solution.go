package circuitbreaker

import (
	"errors"
	"time"
)

/*
See README.md for the full problem statement, requirements, and rubric.

State, its constants, Clock, RealClock, and ErrCircuitOpen below are already complete
— do not modify them.
*/

var ErrCircuitOpen = errors.New("circuit breaker is open")

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

// Clock is a source of the current time. Do not modify.
type Clock interface {
	Now() time.Time
}

// RealClock is the production Clock, backed by the real wall clock. Do not modify.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// TODO: implement CircuitBreaker.
//   - NewCircuitBreaker(failureThreshold int, cooldown time.Duration, clock Clock) *CircuitBreaker
//   - Call(fn func() error) error
//   - State() State
//
// See README.md for the exact state-transition rules.
type CircuitBreaker struct {
	// TODO: add your fields here
}

func NewCircuitBreaker(failureThreshold int, cooldown time.Duration, clock Clock) *CircuitBreaker {
	// TODO: implement
	return &CircuitBreaker{}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	// TODO: implement
	return nil
}

func (cb *CircuitBreaker) State() State {
	// TODO: implement
	return Closed
}
