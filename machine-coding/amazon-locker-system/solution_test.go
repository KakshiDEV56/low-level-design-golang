package lockersystem

import "testing"

func newTestSystem() *LockerSystem {
	return NewLockerSystem(map[string]Size{
		"S1": Small, "S2": Small,
		"M1": Medium, "M2": Medium,
		"L1": Large,
	})
}

func TestAssignLocker_BestFitPicksSmallestSufficient(t *testing.T) {
	s := newTestSystem()
	lockerID, _, ok := s.AssignLocker("pkg1", Small)
	if !ok {
		t.Fatal("AssignLocker(Small) ok = false, want true")
	}
	if lockerID != "S1" && lockerID != "S2" {
		t.Fatalf("AssignLocker(Small) locker = %q, want a Small locker (S1 or S2)", lockerID)
	}
}

func TestAssignLocker_DeterministicTieBreak(t *testing.T) {
	s := newTestSystem()
	lockerID, _, ok := s.AssignLocker("pkg1", Small)
	if !ok {
		t.Fatal("AssignLocker ok = false, want true")
	}
	if lockerID != "S1" {
		t.Fatalf("AssignLocker() with both S1 and S2 free = %q, want %q (alphabetically first)", lockerID, "S1")
	}
}

func TestAssignLocker_FallsBackToLargerSizeWhenExactSizeFull(t *testing.T) {
	s := newTestSystem()
	s.AssignLocker("pkg1", Small)
	s.AssignLocker("pkg2", Small) // both Small lockers now taken

	lockerID, _, ok := s.AssignLocker("pkg3", Small)
	if !ok {
		t.Fatal("AssignLocker ok = false, want true (should fall back to Medium)")
	}
	if lockerID != "M1" && lockerID != "M2" {
		t.Fatalf("AssignLocker() fallback = %q, want a Medium locker", lockerID)
	}
}

func TestAssignLocker_NeverPicksOversizedLockerWhenSmallerFits(t *testing.T) {
	s := newTestSystem()
	// Only Large is free for Small packages after filling Small and Medium.
	s.AssignLocker("a", Small)
	s.AssignLocker("b", Small)
	lockerID, _, ok := s.AssignLocker("c", Small)
	if !ok || (lockerID != "M1" && lockerID != "M2") {
		t.Fatalf("AssignLocker() = (%q, ok=%v), want a Medium locker before touching Large", lockerID, ok)
	}
}

func TestAssignLocker_FailsWhenNoLockerBigEnough(t *testing.T) {
	s := NewLockerSystem(map[string]Size{"S1": Small})
	_, _, ok := s.AssignLocker("pkg1", Large)
	if ok {
		t.Fatal("AssignLocker(Large) with only a Small locker = true, want false")
	}
}

func TestAssignLocker_FailsWhenAllLockersOfSufficientSizeAreFull(t *testing.T) {
	s := NewLockerSystem(map[string]Size{"L1": Large})
	if _, _, ok := s.AssignLocker("pkg1", Large); !ok {
		t.Fatal("first AssignLocker(Large) failed, want true")
	}
	if _, _, ok := s.AssignLocker("pkg2", Large); ok {
		t.Fatal("second AssignLocker(Large) with no lockers free = true, want false")
	}
}

func TestAssignLocker_RejectsDoubleAssignmentOfSamePackage(t *testing.T) {
	s := newTestSystem()
	if _, _, ok := s.AssignLocker("pkg1", Small); !ok {
		t.Fatal("first AssignLocker(pkg1) failed, want true")
	}
	if _, _, ok := s.AssignLocker("pkg1", Small); ok {
		t.Fatal("second AssignLocker(pkg1) while still assigned = true, want false")
	}
}

func TestAssignLocker_CodesAreUniqueAcrossAssignments(t *testing.T) {
	s := newTestSystem()
	_, code1, _ := s.AssignLocker("pkg1", Small)
	_, code2, _ := s.AssignLocker("pkg2", Medium)
	if code1 == "" || code2 == "" {
		t.Fatal("pickup codes must not be empty")
	}
	if code1 == code2 {
		t.Fatalf("two different assignments got the same pickup code %q, want unique codes", code1)
	}
}

func TestPickup_WrongCodeFails(t *testing.T) {
	s := newTestSystem()
	_, code, _ := s.AssignLocker("pkg1", Small)
	if s.Pickup("pkg1", code+"-wrong") {
		t.Fatal("Pickup with a wrong code = true, want false")
	}
}

func TestPickup_UnknownPackageFails(t *testing.T) {
	s := newTestSystem()
	if s.Pickup("no-such-package", "anything") {
		t.Fatal("Pickup for an unknown package = true, want false")
	}
}

func TestPickup_CorrectCodeReleasesLockerForReuse(t *testing.T) {
	s := newTestSystem()
	before := s.AvailableCount(Small)

	lockerID, code, ok := s.AssignLocker("pkg1", Small)
	if !ok {
		t.Fatal("AssignLocker failed")
	}
	if got := s.AvailableCount(Small); got != before-1 {
		t.Fatalf("AvailableCount(Small) after assignment = %d, want %d", got, before-1)
	}

	if !s.Pickup("pkg1", code) {
		t.Fatal("Pickup with the correct code = false, want true")
	}
	if got := s.AvailableCount(Small); got != before {
		t.Fatalf("AvailableCount(Small) after pickup = %d, want %d (locker freed)", got, before)
	}

	// The freed locker must be genuinely reusable.
	newLockerID, _, ok := s.AssignLocker("pkg2", Small)
	if !ok {
		t.Fatal("AssignLocker after a pickup freed a locker = false, want true")
	}
	if newLockerID != lockerID {
		t.Fatalf("expected the freed locker %q to be reused, got %q", lockerID, newLockerID)
	}
}

func TestPickup_CannotPickupTwice(t *testing.T) {
	s := newTestSystem()
	_, code, _ := s.AssignLocker("pkg1", Small)
	if !s.Pickup("pkg1", code) {
		t.Fatal("first Pickup failed, want true")
	}
	if s.Pickup("pkg1", code) {
		t.Fatal("second Pickup with the same (now-consumed) code = true, want false")
	}
}

func TestAvailableCount_ExactSizeOnly(t *testing.T) {
	s := newTestSystem()
	if got := s.AvailableCount(Small); got != 2 {
		t.Fatalf("AvailableCount(Small) = %d, want 2", got)
	}
	if got := s.AvailableCount(Medium); got != 2 {
		t.Fatalf("AvailableCount(Medium) = %d, want 2", got)
	}
	if got := s.AvailableCount(Large); got != 1 {
		t.Fatalf("AvailableCount(Large) = %d, want 1", got)
	}

	s.AssignLocker("pkg1", Small)
	if got := s.AvailableCount(Small); got != 1 {
		t.Fatalf("AvailableCount(Small) after one assignment = %d, want 1", got)
	}
	// Assigning a Small package must not affect Medium/Large counts.
	if got := s.AvailableCount(Medium); got != 2 {
		t.Fatalf("AvailableCount(Medium) after an unrelated Small assignment = %d, want 2", got)
	}
}
