package functionaloptions

import "time"

/*
See README.md for the full problem statement, requirements, and rubric.

Server and Option are fixed — don't change their shape.
*/

type Server struct {
	Host    string
	Port    int
	Timeout time.Duration
	Retries int
	TLS     bool
}

// Option mutates a Server being constructed.
type Option func(*Server)

// TODO: implement NewServer and the four Option constructors below. See README.md
// for the exact default values and application order.

func NewServer(host string, opts ...Option) *Server {
	// TODO: implement
	return &Server{}
}

func WithPort(port int) Option {
	// TODO: implement
	return nil
}

func WithTimeout(d time.Duration) Option {
	// TODO: implement
	return nil
}

func WithRetries(n int) Option {
	// TODO: implement
	return nil
}

func WithTLS() Option {
	// TODO: implement
	return nil
}
