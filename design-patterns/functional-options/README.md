# Functional Options

## The concept

The only pattern in this folder that isn't from an enterprise-architecture textbook —
it's specific to Go, and it's the idiomatic answer to "how do I write a constructor
with a bunch of optional configuration?" without falling into either of the two usual
traps:

- A giant constructor with 8 positional parameters, most of which most callers don't
  care about and will get wrong the order of.
- A `Config` struct passed to the constructor, where you can't tell the difference
  between "the caller explicitly wants 0" and "the caller didn't set this field" — Go's
  zero values make that genuinely ambiguous.

The fix: required parameters stay positional; everything optional is set via variadic
`...Option` arguments, where each `Option` is just a function that mutates the object
being built. Callers only mention the options they actually want to override, in any
order, and it reads like self-documenting configuration:
`NewServer("localhost", WithPort(9090), WithTLS())`.

## Functional requirements

1. `Server` (given) has `Host string`, `Port int`, `Timeout time.Duration`, `Retries
   int`, `TLS bool`.
2. `Option` is a fixed type: `func(*Server)`.
3. `NewServer(host string, opts ...Option) *Server` — `host` is required and
   positional. Before applying any options, the server must start with these defaults:
   `Port: 8080`, `Timeout: 30 * time.Second`, `Retries: 3`, `TLS: false`. Then apply
   each option **in the order given**, so a later option for the same field wins over
   an earlier one.
4. Implement these options:
   - `WithPort(port int) Option`
   - `WithTimeout(d time.Duration) Option`
   - `WithRetries(n int) Option`
   - `WithTLS() Option` — takes no arguments (it's just "turn TLS on"); note this one
     is NOT itself a function returning an `Option` from a parameter — well, it is, it
     just doesn't need to close over anything.

## Rubric

- [ ] Calling `NewServer(host)` with no options at all produces exactly the documented
      defaults.
- [ ] Options apply in argument order — `NewServer(host, WithPort(80), WithPort(443))`
      ends with port 443, not 80.
- [ ] A brand-new `Option`, written entirely in the test file (e.g. one that sets
      multiple fields, or is built by composing two existing options together), works
      with `NewServer` unmodified — proving `Option` is a genuinely open type, not a
      closed enum of cases `NewServer` switches on internally.
- [ ] `NewServer` itself contains no `if` branching on which specific options were
      passed — it just applies each one as an opaque function.

## Hints

- `Option` being `func(*Server)` (mutating a pointer) rather than `func(Server)
  Server` (returning a new value) is what makes multiple options composable in a
  simple loop: `for _, opt := range opts { opt(server) }`.
