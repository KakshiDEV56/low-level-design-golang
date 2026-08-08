package interfacesegregation

import "testing"

// Compile-time proof HumanWorker satisfies both narrow interfaces.
var (
	_ Workable = (HumanWorker{})
	_ Eatable  = (HumanWorker{})
)

// Compile-time proof RobotWorker satisfies Workable...
var _ Workable = (RobotWorker{})

// ...but NOT Eatable — RobotWorker should have no Eat() method at all. Uncommenting
// the line below should fail to compile once the fix is in place. (A compile error
// can't be asserted at runtime, so this is documented here instead of executed; try
// uncommenting it yourself after refactoring to see it fail.)
//
// var _ Eatable = (RobotWorker{})

func TestRunShift_AcceptsHumanWorker(t *testing.T) {
	got := RunShift(HumanWorker{})
	if got == "" {
		t.Fatal("RunShift(HumanWorker{}) returned an empty string")
	}
}

func TestRunShift_AcceptsRobotWorker(t *testing.T) {
	got := RunShift(RobotWorker{})
	if got == "" {
		t.Fatal("RunShift(RobotWorker{}) returned an empty string")
	}
}

func TestRunShift_DifferentWorkersProduceDifferentOutput(t *testing.T) {
	human := RunShift(HumanWorker{})
	robot := RunShift(RobotWorker{})
	if human == robot {
		t.Fatalf("RunShift produced identical output %q for both a human and a robot worker, want them distinguishable", human)
	}
}

func TestHumanWorker_WorkAndEatAreDistinctCapabilities(t *testing.T) {
	h := HumanWorker{}
	if h.Work() == h.Eat() {
		t.Fatalf("HumanWorker.Work() and Eat() returned the same value %q, want them distinct", h.Work())
	}
}

// customWorkable is a brand-new type defined only in the test file, implementing just
// Workable — proving RunShift genuinely only needs the narrow interface, not the fat
// original Worker.
type customWorkable struct{ label string }

func (c customWorkable) Work() string { return c.label }

func TestRunShift_AcceptsAnyWorkableImplementation(t *testing.T) {
	got := RunShift(customWorkable{label: "custom working"})
	if got == "" {
		t.Fatal("RunShift(customWorkable{...}) returned an empty string")
	}
}
