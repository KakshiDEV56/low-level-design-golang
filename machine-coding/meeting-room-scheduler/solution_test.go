package meetingroomscheduler

import (
	"reflect"
	"testing"
	"time"
)

var day = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func slot(startHour, endHour int) TimeSlot {
	return TimeSlot{
		Start: day.Add(time.Duration(startHour) * time.Hour),
		End:   day.Add(time.Duration(endHour) * time.Hour),
	}
}

func TestScheduler_BookSucceedsWithNoConflict(t *testing.T) {
	s := NewScheduler([]string{"A", "B"})
	if !s.Book("A", slot(9, 10)) {
		t.Fatal("Book(A, 9-10) = false, want true")
	}
}

func TestScheduler_BookFailsOnUnknownRoom(t *testing.T) {
	s := NewScheduler([]string{"A"})
	if s.Book("Z", slot(9, 10)) {
		t.Fatal("Book on unknown room = true, want false")
	}
}

func TestScheduler_BookFailsOnOverlap(t *testing.T) {
	s := NewScheduler([]string{"A"})
	s.Book("A", slot(9, 11))

	cases := []struct {
		name string
		s    TimeSlot
	}{
		{"identical", slot(9, 11)},
		{"fully contained", slot(9, 10)},
		{"overlaps start", slot(8, 10)},
		{"overlaps end", slot(10, 12)},
		{"contains existing", slot(8, 12)},
	}
	for _, c := range cases {
		if s.Book("A", c.s) {
			t.Errorf("%s: Book() = true, want false (should conflict with 9-11)", c.name)
		}
	}
}

func TestScheduler_BookSucceedsOnAdjacentSlots(t *testing.T) {
	s := NewScheduler([]string{"A"})
	s.Book("A", slot(9, 10))

	if !s.Book("A", slot(10, 11)) {
		t.Fatal("Book(10-11) after existing 9-10 = false, want true (adjacent, half-open, no overlap)")
	}
	if !s.Book("A", slot(8, 9)) {
		t.Fatal("Book(8-9) before existing 9-10 = false, want true (adjacent, half-open, no overlap)")
	}
}

func TestScheduler_BookAnyAvailable_PicksAFreeRoom(t *testing.T) {
	s := NewScheduler([]string{"A", "B", "C"})
	s.Book("A", slot(9, 10))
	s.Book("B", slot(9, 10))

	got := s.BookAnyAvailable(slot(9, 10))
	if got != "C" {
		t.Fatalf("BookAnyAvailable() = %q, want %q (only C is free)", got, "C")
	}
	// The room must now actually be booked.
	if s.Book("C", slot(9, 10)) {
		t.Fatal("room C should already be booked by BookAnyAvailable, but Book() succeeded again")
	}
}

func TestScheduler_BookAnyAvailable_ReturnsEmptyWhenAllBusy(t *testing.T) {
	s := NewScheduler([]string{"A", "B"})
	s.Book("A", slot(9, 10))
	s.Book("B", slot(9, 10))

	if got := s.BookAnyAvailable(slot(9, 10)); got != "" {
		t.Fatalf("BookAnyAvailable() with all rooms busy = %q, want \"\"", got)
	}
}

func TestScheduler_BookAnyAvailable_IsDeterministic(t *testing.T) {
	s := NewScheduler([]string{"C", "A", "B"})
	first := s.BookAnyAvailable(slot(9, 10))
	// Cancel it and ask again on a fresh identical scheduler — must pick the same room.
	s.Cancel(first, slot(9, 10))

	second := s.BookAnyAvailable(slot(9, 10))
	if first != second {
		t.Fatalf("BookAnyAvailable() picked %q then %q on equivalent state, want the same room both times", first, second)
	}
}

func TestScheduler_CancelExactMatchOnly(t *testing.T) {
	s := NewScheduler([]string{"A"})
	s.Book("A", slot(9, 11))

	if s.Cancel("A", slot(9, 10)) {
		t.Fatal("Cancel() with a slot that doesn't exactly match an existing booking = true, want false")
	}
	if !s.Cancel("A", slot(9, 11)) {
		t.Fatal("Cancel() with the exact existing booking = false, want true")
	}
	if !s.Book("A", slot(9, 11)) {
		t.Fatal("room A should be free again after Cancel(), but Book() failed")
	}
}

func TestScheduler_CancelUnknownRoomOrBooking(t *testing.T) {
	s := NewScheduler([]string{"A"})
	if s.Cancel("Z", slot(9, 10)) {
		t.Fatal("Cancel() on unknown room = true, want false")
	}
	if s.Cancel("A", slot(9, 10)) {
		t.Fatal("Cancel() on a room with no such booking = true, want false")
	}
}

func TestScheduler_AvailableRoomsSortedAlphabetically(t *testing.T) {
	s := NewScheduler([]string{"C", "A", "B"})
	s.Book("A", slot(9, 10))

	got := s.AvailableRooms(slot(9, 10))
	want := []string{"B", "C"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AvailableRooms() = %v, want %v", got, want)
	}
}

func TestScheduler_AvailableRoomsWithNoBookings(t *testing.T) {
	s := NewScheduler([]string{"B", "A"})
	got := s.AvailableRooms(slot(9, 10))
	want := []string{"A", "B"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AvailableRooms() = %v, want %v", got, want)
	}
}

func TestScheduler_BookingsForRoomSortedByStart(t *testing.T) {
	s := NewScheduler([]string{"A"})
	s.Book("A", slot(14, 15))
	s.Book("A", slot(9, 10))
	s.Book("A", slot(11, 12))

	got := s.BookingsForRoom("A")
	want := []TimeSlot{slot(9, 10), slot(11, 12), slot(14, 15)}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BookingsForRoom() = %v, want %v", got, want)
	}
}

func TestScheduler_BookingsForRoomUnknownOrEmpty(t *testing.T) {
	s := NewScheduler([]string{"A"})
	if got := s.BookingsForRoom("Z"); len(got) != 0 {
		t.Fatalf("BookingsForRoom(unknown) = %v, want empty", got)
	}
	if got := s.BookingsForRoom("A"); len(got) != 0 {
		t.Fatalf("BookingsForRoom(no bookings) = %v, want empty", got)
	}
}

func TestScheduler_MultipleRoomsIndependentBookings(t *testing.T) {
	s := NewScheduler([]string{"A", "B"})
	if !s.Book("A", slot(9, 10)) {
		t.Fatal("Book(A, 9-10) = false, want true")
	}
	if !s.Book("B", slot(9, 10)) {
		t.Fatal("Book(B, 9-10) = false, want true (independent room, same slot)")
	}
}
