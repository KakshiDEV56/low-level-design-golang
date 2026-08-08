package meetingroomscheduler

import "time"

/*
See README.md for the full problem statement, requirements, and rubric.

TODO:
 1. Add the fields you need on Scheduler below. Hint: a map from room ID to that room's
    slice of booked TimeSlots is enough.
 2. Implement NewScheduler(roomIDs): initialize each room with an empty booking list.
 3. Implement Book, BookAnyAvailable, Cancel, AvailableRooms, BookingsForRoom per the
    README. Consider writing a private helper like hasConflict(roomID, slot) that both
    Book and BookAnyAvailable can share.
*/

// TimeSlot is a half-open interval [Start, End).
type TimeSlot struct {
	Start time.Time
	End   time.Time
}

type Scheduler struct {
	// TODO: add your fields here
}

func NewScheduler(roomIDs []string) *Scheduler {
	// TODO: implement
	return &Scheduler{}
}

func (s *Scheduler) Book(roomID string, slot TimeSlot) bool {
	// TODO: implement
	return false
}

func (s *Scheduler) BookAnyAvailable(slot TimeSlot) string {
	// TODO: implement
	return ""
}

func (s *Scheduler) Cancel(roomID string, slot TimeSlot) bool {
	// TODO: implement
	return false
}

func (s *Scheduler) AvailableRooms(slot TimeSlot) []string {
	// TODO: implement
	return nil
}

func (s *Scheduler) BookingsForRoom(roomID string) []TimeSlot {
	// TODO: implement
	return nil
}
