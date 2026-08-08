# Meeting Room Scheduler

**Suggested time box: 30 minutes.**

## Problem statement

Design a scheduler for a fixed set of meeting rooms. Users can book a specific room for
a time slot, ask the system to find and book *any* free room for a slot, cancel a
booking, and query which rooms are free during a given slot.

## Functional requirements

1. `NewScheduler(roomIDs []string)` initializes the scheduler with a fixed set of room
   IDs (no rooms are added/removed later).
2. `Book(roomID string, slot TimeSlot) bool` reserves `roomID` for `slot`. Returns
   `false` (and books nothing) if `roomID` doesn't exist or if the slot conflicts with
   an existing booking for that room. Two slots conflict if they overlap at all — a
   meeting ending at 3:00 and one starting at 3:00 do **not** conflict (treat slots as
   half-open: `[Start, End)`).
3. `BookAnyAvailable(slot TimeSlot) string` finds any room with no conflicting booking
   for `slot`, books it, and returns its room ID — or `""` if every room is busy. If
   multiple rooms are free, which one gets picked doesn't matter, but the choice must be
   deterministic given the same scheduler state (don't use map iteration order or
   anything else non-deterministic — tests need a reproducible pick, e.g. by sorting
   candidate room IDs first).
4. `Cancel(roomID string, slot TimeSlot) bool` removes an exact existing booking.
   Returns `false` if no such booking exists.
5. `AvailableRooms(slot TimeSlot) []string` returns every room ID with no conflicting
   booking for `slot`, sorted alphabetically.
6. `BookingsForRoom(roomID string) []TimeSlot` returns all current bookings for a room,
   sorted by start time. Returns `nil`/empty for an unknown room ID or one with no
   bookings.

## Rubric

- [ ] Conflict detection correctly implements half-open interval overlap (adjacent
      slots don't conflict; a slot fully contained inside another does; partial overlap
      on either edge does).
- [ ] `BookAnyAvailable` is deterministic given the same state.
- [ ] `Cancel` only removes an exact match — it shouldn't accidentally free an
      overlapping-but-different slot.
- [ ] Unknown room IDs are handled gracefully everywhere (no panics).
- [ ] Booking/conflict-checking logic lives in one clearly-named place, not duplicated
      across `Book` and `BookAnyAvailable`.

## Hints

- Two half-open slots `[s1,e1)` and `[s2,e2)` overlap iff `s1 < e2 && s2 < e1`.
- Store each room's bookings as a slice (or sorted slice) of `TimeSlot`; a linear scan
  per room for conflict checking is completely fine at this scale — don't over-engineer
  an interval tree.
- Write one private helper, e.g. `hasConflict(roomID string, slot TimeSlot) bool`, and
  have both `Book` and `BookAnyAvailable` call it — this is the "don't duplicate the
  conflict logic" grading point above.
