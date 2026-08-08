package interfacesegregation

/*
This code compiles and works today — but it violates the Interface Segregation
Principle. See README.md for the full explanation.

Worker forces every implementer to provide BOTH Work() and Eat(), even though
RobotWorker has no meaningful way to "eat" — its Eat() method is a lie (a no-op that
exists purely to satisfy the interface).

TODO: refactor this file so that:
 1. The fat Worker interface is replaced by two small ones:
      type Workable interface { Work() string }
      type Eatable interface { Eat() string }
 2. HumanWorker keeps both Work() and Eat() (returning some descriptive string each,
    e.g. "human working" / "human eating" — exact wording is up to you, tests just
    check the methods are callable and distinct).
 3. RobotWorker keeps ONLY Work() — delete its Eat() method entirely. It must not
    implement Eatable at all.
 4. RunShift(w Workable) string replaces any code that previously required a full
    Worker — it should call Work() and return (or wrap) its result, and must accept
    BOTH HumanWorker and RobotWorker.

Once refactored, solution_test.go should compile and pass.
*/

// Worker is a "fat" interface: it forces every implementer to support both working
// AND eating, even when that doesn't make sense for the implementer.
// TODO: split this into Workable and Eatable (see above).
type Worker interface {
	Work() string
	Eat() string
}

type HumanWorker struct{}

func (HumanWorker) Work() string { return "human working" }
func (HumanWorker) Eat() string  { return "human eating" }

// RobotWorker is forced to implement Eat() even though robots don't eat. This method
// is a lie that only exists to satisfy the Worker interface.
type RobotWorker struct{}

func (RobotWorker) Work() string { return "robot working" }
func (RobotWorker) Eat() string  { return "" /* doesn't actually mean anything */ }
