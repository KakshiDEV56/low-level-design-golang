package saga

/*
See README.md for the full problem statement, requirements, and rubric.

Step is fixed — don't change its shape.
*/

// Step is one stage of a saga.
type Step struct {
	Name       string
	Action     func() error
	Compensate func() error
}

// TODO: implement Saga.
//   - NewSaga(steps ...Step) *Saga
//   - Run() error
//
// See README.md for the exact execution and compensation-ordering rules, and for how
// to combine errors with errors.Join.
type Saga struct {
	// TODO: add your fields here
}

func NewSaga(steps ...Step) *Saga {
	// TODO: implement
	return &Saga{}
}

func (s *Saga) Run() error {
	// TODO: implement
	return nil
}
