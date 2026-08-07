package designlogstoragesystem

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

// sortedInts returns a sorted copy so tests can compare "set of IDs" without depending
// on the implementation's internal iteration order (the real problem doesn't require a
// specific output order either).
func sortedInts(vals []int) []int {
	out := make([]int, len(vals))
	copy(out, vals)
	sort.Ints(out)
	return out
}

func TestLogSystem_Example1(t *testing.T) {
	l := NewLogSystem()
	l.Put(1, "2017:01:01:23:59:59")
	l.Put(2, "2017:01:01:22:59:59")
	l.Put(3, "2016:01:01:00:00:00")

	got := sortedInts(l.Retrieve("2016:01:01:01:01:01", "2017:01:01:23:00:00", "Year"))
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`Retrieve(..., "Year") = %v, want %v`, got, want)
	}

	got = sortedInts(l.Retrieve("2016:01:01:01:01:01", "2017:01:01:23:00:00", "Hour"))
	want = []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`Retrieve(..., "Hour") = %v, want %v`, got, want)
	}
}

func TestLogSystem_EmptyResultWhenNoneMatch(t *testing.T) {
	l := NewLogSystem()
	l.Put(1, "2015:01:01:00:00:00")
	got := l.Retrieve("2016:01:01:00:00:00", "2017:01:01:00:00:00", "Year")
	if len(got) != 0 {
		t.Fatalf("Retrieve() = %v, want empty", got)
	}
}

func TestLogSystem_MonthGranularity(t *testing.T) {
	l := NewLogSystem()
	l.Put(1, "2017:03:15:10:00:00")
	l.Put(2, "2017:04:01:00:00:00")
	l.Put(3, "2017:05:01:00:00:00")
	got := sortedInts(l.Retrieve("2017:03:01:00:00:00", "2017:04:30:23:59:59", "Month"))
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`Retrieve(..., "Month") = %v, want %v`, got, want)
	}
}

func TestLogSystem_DayGranularity(t *testing.T) {
	l := NewLogSystem()
	l.Put(1, "2017:01:01:00:00:01")
	l.Put(2, "2017:01:02:12:00:00")
	l.Put(3, "2017:01:03:23:59:59")
	got := sortedInts(l.Retrieve("2017:01:01:23:59:59", "2017:01:02:23:59:59", "Day"))
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`Retrieve(..., "Day") = %v, want %v`, got, want)
	}
}

func TestLogSystem_MinuteGranularity(t *testing.T) {
	l := NewLogSystem()
	l.Put(1, "2017:01:01:10:15:00")
	l.Put(2, "2017:01:01:10:15:59")
	l.Put(3, "2017:01:01:10:16:00")
	got := sortedInts(l.Retrieve("2017:01:01:10:15:00", "2017:01:01:10:15:59", "Minute"))
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`Retrieve(..., "Minute") = %v, want %v`, got, want)
	}
}

func TestLogSystem_SecondGranularityExactMatch(t *testing.T) {
	l := NewLogSystem()
	l.Put(1, "2017:01:01:10:15:30")
	l.Put(2, "2017:01:01:10:15:31")
	got := sortedInts(l.Retrieve("2017:01:01:10:15:30", "2017:01:01:10:15:30", "Second"))
	want := []int{1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`Retrieve(..., "Second") = %v, want %v`, got, want)
	}
}

func TestLogSystem_YearBoundaryInclusive(t *testing.T) {
	l := NewLogSystem()
	l.Put(1, "2000:01:01:00:00:00")
	l.Put(2, "2017:12:31:23:59:59")
	got := sortedInts(l.Retrieve("2000:01:01:00:00:00", "2017:12:31:23:59:59", "Year"))
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`Retrieve(..., "Year") = %v, want %v`, got, want)
	}
}

func TestLogSystem_HourGranularityExcludesAdjacentHours(t *testing.T) {
	l := NewLogSystem()
	l.Put(1, "2017:06:15:08:59:59")
	l.Put(2, "2017:06:15:09:00:00")
	l.Put(3, "2017:06:15:09:59:59")
	l.Put(4, "2017:06:15:10:00:00")
	got := sortedInts(l.Retrieve("2017:06:15:09:00:00", "2017:06:15:09:59:59", "Hour"))
	want := []int{2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`Retrieve(..., "Hour") = %v, want %v`, got, want)
	}
}

func TestLogSystem_StressManyLogsYearFilter(t *testing.T) {
	l := NewLogSystem()
	id := 1
	for year := 2010; year <= 2017; year++ {
		for month := 1; month <= 12; month++ {
			l.Put(id, monthTimestamp(year, month))
			id++
		}
	}
	got := l.Retrieve("2012:01:01:00:00:00", "2013:12:31:23:59:59", "Year")
	if len(got) != 24 {
		t.Fatalf("Retrieve(2012..2013, Year) returned %d ids, want 24", len(got))
	}
}

func monthTimestamp(year, month int) string {
	return fmt.Sprintf("%04d:%02d:01:00:00:00", year, month)
}
