package specification

/*
See README.md for the full problem statement, requirements, and rubric.

Specification is a fixed interface — don't change its shape.
*/

// Specification is a composable business rule over a candidate of type T.
type Specification[T any] interface {
	IsSatisfiedBy(candidate T) bool
}

// TODO: implement And, Or, Not, and Filter. See README.md hints for a suggested
// adapter-type approach that lets one small type back all three combinators.

func And[T any](left, right Specification[T]) Specification[T] {
	// TODO: implement
	return nil
}

func Or[T any](left, right Specification[T]) Specification[T] {
	// TODO: implement
	return nil
}

func Not[T any](spec Specification[T]) Specification[T] {
	// TODO: implement
	return nil
}

func Filter[T any](items []T, spec Specification[T]) []T {
	// TODO: implement
	return nil
}
