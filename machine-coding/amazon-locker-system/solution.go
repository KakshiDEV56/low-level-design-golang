package lockersystem

/*
See README.md for the full problem statement, requirements, and rubric.

TODO:
 1. Add the fields you need on LockerSystem below. Hint: a map from locker ID to a
    small struct tracking { size Size; occupied bool; packageID string; code string },
    plus a way to map packageID -> lockerID for O(1) Pickup lookups.
 2. Implement NewLockerSystem(lockers): initialize every locker as empty.
 3. Implement AssignLocker, Pickup, AvailableCount per the README.
*/

// Size represents a locker (and package) size class, in increasing order.
type Size int

const (
	Small Size = iota
	Medium
	Large
)

type LockerSystem struct {
	// TODO: add your fields here
}

func NewLockerSystem(lockers map[string]Size) *LockerSystem {
	// TODO: implement
	return &LockerSystem{}
}

func (s *LockerSystem) AssignLocker(packageID string, packageSize Size) (lockerID string, code string, ok bool) {
	// TODO: implement
	return "", "", false
}

func (s *LockerSystem) Pickup(packageID string, code string) bool {
	// TODO: implement
	return false
}

func (s *LockerSystem) AvailableCount(size Size) int {
	// TODO: implement
	return 0
}
