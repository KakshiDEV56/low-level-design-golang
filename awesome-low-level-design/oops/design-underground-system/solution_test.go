package designundergroundsystem

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-5
}

func TestUndergroundSystem_Example1(t *testing.T) {
	u := NewUndergroundSystem()
	u.CheckIn(45, "Leyton", 3)
	u.CheckIn(32, "Paradise", 8)
	u.CheckIn(27, "Leyton", 10)
	u.CheckOut(45, "Waterloo", 15)
	u.CheckOut(27, "Waterloo", 20)
	u.CheckOut(32, "Cambridge", 22)

	if got := u.GetAverageTime("Paradise", "Cambridge"); !almostEqual(got, 14.0) {
		t.Fatalf(`GetAverageTime("Paradise","Cambridge") = %v, want 14.0`, got)
	}
	if got := u.GetAverageTime("Leyton", "Waterloo"); !almostEqual(got, 11.0) {
		t.Fatalf(`GetAverageTime("Leyton","Waterloo") = %v, want 11.0`, got)
	}
	u.CheckIn(10, "Leyton", 24)
	if got := u.GetAverageTime("Leyton", "Waterloo"); !almostEqual(got, 11.0) {
		t.Fatalf(`GetAverageTime("Leyton","Waterloo") = %v, want 11.0`, got)
	}
	u.CheckOut(10, "Waterloo", 38)
	if got := u.GetAverageTime("Leyton", "Waterloo"); !almostEqual(got, 12.0) {
		t.Fatalf(`GetAverageTime("Leyton","Waterloo") = %v, want 12.0`, got)
	}
}

func TestUndergroundSystem_SingleTrip(t *testing.T) {
	u := NewUndergroundSystem()
	u.CheckIn(1, "A", 0)
	u.CheckOut(1, "B", 10)
	if got := u.GetAverageTime("A", "B"); !almostEqual(got, 10.0) {
		t.Fatalf(`GetAverageTime("A","B") = %v, want 10.0`, got)
	}
}

func TestUndergroundSystem_DirectionMatters(t *testing.T) {
	u := NewUndergroundSystem()
	u.CheckIn(1, "A", 0)
	u.CheckOut(1, "B", 10) // A -> B took 10
	u.CheckIn(2, "B", 20)
	u.CheckOut(2, "A", 25) // B -> A took 5
	if got := u.GetAverageTime("A", "B"); !almostEqual(got, 10.0) {
		t.Fatalf(`GetAverageTime("A","B") = %v, want 10.0`, got)
	}
	if got := u.GetAverageTime("B", "A"); !almostEqual(got, 5.0) {
		t.Fatalf(`GetAverageTime("B","A") = %v, want 5.0`, got)
	}
}

func TestUndergroundSystem_SameCustomerReusedAfterCheckout(t *testing.T) {
	u := NewUndergroundSystem()
	u.CheckIn(1, "A", 0)
	u.CheckOut(1, "B", 10)
	u.CheckIn(1, "B", 15)
	u.CheckOut(1, "C", 20)
	if got := u.GetAverageTime("A", "B"); !almostEqual(got, 10.0) {
		t.Fatalf(`GetAverageTime("A","B") = %v, want 10.0`, got)
	}
	if got := u.GetAverageTime("B", "C"); !almostEqual(got, 5.0) {
		t.Fatalf(`GetAverageTime("B","C") = %v, want 5.0`, got)
	}
}

func TestUndergroundSystem_MultipleCustomersSameRouteInterleaved(t *testing.T) {
	u := NewUndergroundSystem()
	u.CheckIn(1, "X", 0)
	u.CheckIn(2, "X", 5)
	u.CheckOut(1, "Y", 10) // trip 1: 10
	u.CheckOut(2, "Y", 25) // trip 2: 20
	if got := u.GetAverageTime("X", "Y"); !almostEqual(got, 15.0) {
		t.Fatalf(`GetAverageTime("X","Y") = %v, want 15.0`, got)
	}
}

func TestUndergroundSystem_NonIntegerAverage(t *testing.T) {
	u := NewUndergroundSystem()
	u.CheckIn(1, "A", 0)
	u.CheckOut(1, "B", 10) // 10
	u.CheckIn(2, "A", 0)
	u.CheckOut(2, "B", 11) // 11
	u.CheckIn(3, "A", 0)
	u.CheckOut(3, "B", 12) // 12
	if got := u.GetAverageTime("A", "B"); !almostEqual(got, 11.0) {
		t.Fatalf(`GetAverageTime("A","B") = %v, want 11.0`, got)
	}
}

func TestUndergroundSystem_UnrelatedRoutesDoNotAffectEachOther(t *testing.T) {
	u := NewUndergroundSystem()
	u.CheckIn(1, "A", 0)
	u.CheckOut(1, "B", 5)
	u.CheckIn(2, "C", 0)
	u.CheckOut(2, "D", 100)
	if got := u.GetAverageTime("A", "B"); !almostEqual(got, 5.0) {
		t.Fatalf(`GetAverageTime("A","B") = %v, want 5.0`, got)
	}
	if got := u.GetAverageTime("C", "D"); !almostEqual(got, 100.0) {
		t.Fatalf(`GetAverageTime("C","D") = %v, want 100.0`, got)
	}
}

func TestUndergroundSystem_StressManyTripsAverage(t *testing.T) {
	u := NewUndergroundSystem()
	const trips = 100
	total := 0
	for i := 0; i < trips; i++ {
		id := i
		checkInT := i * 100
		checkOutT := checkInT + 10 + i%5 // varying travel time
		u.CheckIn(id, "Start", checkInT)
		u.CheckOut(id, "End", checkOutT)
		total += 10 + i%5
	}
	want := float64(total) / float64(trips)
	if got := u.GetAverageTime("Start", "End"); !almostEqual(got, want) {
		t.Fatalf(`GetAverageTime("Start","End") = %v, want %v`, got, want)
	}
}
