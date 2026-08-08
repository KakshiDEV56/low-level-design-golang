package saga

import (
	"errors"
	"reflect"
	"testing"
)

func TestSaga_EmptySagaSucceeds(t *testing.T) {
	s := NewSaga()
	if err := s.Run(); err != nil {
		t.Fatalf("Run() on an empty saga = %v, want nil", err)
	}
}

func TestSaga_AllStepsSucceedNoCompensation(t *testing.T) {
	var compensated []string
	s := NewSaga(
		Step{Name: "a", Action: func() error { return nil }, Compensate: func() error { compensated = append(compensated, "a"); return nil }},
		Step{Name: "b", Action: func() error { return nil }, Compensate: func() error { compensated = append(compensated, "b"); return nil }},
	)
	if err := s.Run(); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
	if len(compensated) != 0 {
		t.Fatalf("compensated = %v, want none (nothing failed)", compensated)
	}
}

func TestSaga_FirstStepFailsNoCompensationCalled(t *testing.T) {
	errBoom := errors.New("boom")
	var compensated []string
	s := NewSaga(
		Step{Name: "a", Action: func() error { return errBoom }, Compensate: func() error { compensated = append(compensated, "a"); return nil }},
	)
	err := s.Run()
	if !errors.Is(err, errBoom) {
		t.Fatalf("Run() = %v, want an error wrapping errBoom", err)
	}
	if len(compensated) != 0 {
		t.Fatalf("compensated = %v, want none (the failing step's own Compensate must not run — nothing preceded it)", compensated)
	}
}

func TestSaga_MiddleStepFailsCompensatesEarlierStepsInReverseOrder(t *testing.T) {
	errBoom := errors.New("boom")
	var actions, compensations []string

	s := NewSaga(
		Step{
			Name:       "book-flight",
			Action:     func() error { actions = append(actions, "book-flight"); return nil },
			Compensate: func() error { compensations = append(compensations, "cancel-flight"); return nil },
		},
		Step{
			Name:       "book-hotel",
			Action:     func() error { actions = append(actions, "book-hotel"); return nil },
			Compensate: func() error { compensations = append(compensations, "cancel-hotel"); return nil },
		},
		Step{
			Name:       "charge-card",
			Action:     func() error { actions = append(actions, "charge-card"); return errBoom },
			Compensate: func() error { compensations = append(compensations, "refund-card"); return nil },
		},
	)

	err := s.Run()
	if !errors.Is(err, errBoom) {
		t.Fatalf("Run() = %v, want an error wrapping errBoom", err)
	}

	wantActions := []string{"book-flight", "book-hotel", "charge-card"}
	if !reflect.DeepEqual(actions, wantActions) {
		t.Fatalf("actions = %v, want %v", actions, wantActions)
	}

	// Compensations must run for book-hotel and book-flight (the two that succeeded),
	// in REVERSE order — hotel was booked after flight, so it's undone first. The
	// failed step (charge-card) never succeeded, so its own compensation must NOT run.
	wantCompensations := []string{"cancel-hotel", "cancel-flight"}
	if !reflect.DeepEqual(compensations, wantCompensations) {
		t.Fatalf("compensations = %v, want %v", compensations, wantCompensations)
	}
}

func TestSaga_LaterStepNeverRunsAfterEarlierFailure(t *testing.T) {
	errBoom := errors.New("boom")
	stepCRan := false

	s := NewSaga(
		Step{Name: "a", Action: func() error { return nil }, Compensate: func() error { return nil }},
		Step{Name: "b", Action: func() error { return errBoom }, Compensate: func() error { return nil }},
		Step{Name: "c", Action: func() error { stepCRan = true; return nil }, Compensate: func() error { return nil }},
	)
	s.Run()
	if stepCRan {
		t.Fatal("step c's Action ran even though step b (before it) failed")
	}
}

func TestSaga_CompensationFailureDoesNotStopOthers(t *testing.T) {
	errBoom := errors.New("boom")
	errCompensateFailed := errors.New("compensate failed")
	var compensated []string

	s := NewSaga(
		Step{
			Name:       "a",
			Action:     func() error { return nil },
			Compensate: func() error { compensated = append(compensated, "a"); return nil },
		},
		Step{
			Name:       "b",
			Action:     func() error { return nil },
			Compensate: func() error { compensated = append(compensated, "b"); return errCompensateFailed },
		},
		Step{
			Name:       "c",
			Action:     func() error { return errBoom },
			Compensate: func() error { return nil },
		},
	)

	err := s.Run()
	if !errors.Is(err, errBoom) {
		t.Fatalf("Run() error does not wrap the original step failure: %v", err)
	}
	if !errors.Is(err, errCompensateFailed) {
		t.Fatalf("Run() error does not wrap the compensation failure: %v", err)
	}

	// Both "a" and "b" must have been compensated despite "b"'s compensation failing.
	want := []string{"b", "a"}
	if !reflect.DeepEqual(compensated, want) {
		t.Fatalf("compensated = %v, want %v (both attempted, in reverse order, despite b's failure)", compensated, want)
	}
}

func TestSaga_SingleStepFailure(t *testing.T) {
	errBoom := errors.New("boom")
	compensateCalled := false
	s := NewSaga(
		Step{
			Name:       "only",
			Action:     func() error { return errBoom },
			Compensate: func() error { compensateCalled = true; return nil },
		},
	)
	err := s.Run()
	if !errors.Is(err, errBoom) {
		t.Fatalf("Run() = %v, want an error wrapping errBoom", err)
	}
	if compensateCalled {
		t.Fatal("Compensate was called for the single step that itself failed")
	}
}
