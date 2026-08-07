package timebasedkeyvaluestore

import (
	"strconv"
	"testing"
)

func TestTimeMap_Example1(t *testing.T) {
	tm := NewTimeMap()
	tm.Set("foo", "bar", 1)
	if got := tm.Get("foo", 1); got != "bar" {
		t.Fatalf(`Get("foo", 1) = %q, want "bar"`, got)
	}
	if got := tm.Get("foo", 3); got != "bar" {
		t.Fatalf(`Get("foo", 3) = %q, want "bar"`, got)
	}
	tm.Set("foo", "bar2", 4)
	if got := tm.Get("foo", 4); got != "bar2" {
		t.Fatalf(`Get("foo", 4) = %q, want "bar2"`, got)
	}
	if got := tm.Get("foo", 5); got != "bar2" {
		t.Fatalf(`Get("foo", 5) = %q, want "bar2"`, got)
	}
}

func TestTimeMap_GetBeforeAnySet(t *testing.T) {
	tm := NewTimeMap()
	if got := tm.Get("missing", 100); got != "" {
		t.Fatalf(`Get("missing", 100) = %q, want ""`, got)
	}
}

func TestTimeMap_GetTimestampBeforeFirstSet(t *testing.T) {
	tm := NewTimeMap()
	tm.Set("k", "v", 10)
	if got := tm.Get("k", 5); got != "" {
		t.Fatalf(`Get("k", 5) = %q, want "" (queried before any Set)`, got)
	}
}

func TestTimeMap_GetExactTimestampMatch(t *testing.T) {
	tm := NewTimeMap()
	tm.Set("k", "v1", 1)
	tm.Set("k", "v2", 2)
	tm.Set("k", "v3", 3)
	if got := tm.Get("k", 2); got != "v2" {
		t.Fatalf(`Get("k", 2) = %q, want "v2"`, got)
	}
}

func TestTimeMap_GetBetweenTimestampsReturnsLatestBefore(t *testing.T) {
	tm := NewTimeMap()
	tm.Set("k", "v1", 1)
	tm.Set("k", "v2", 5)
	tm.Set("k", "v3", 10)
	if got := tm.Get("k", 7); got != "v2" {
		t.Fatalf(`Get("k", 7) = %q, want "v2"`, got)
	}
}

func TestTimeMap_MultipleIndependentKeys(t *testing.T) {
	tm := NewTimeMap()
	tm.Set("a", "alpha", 1)
	tm.Set("b", "beta", 1)
	tm.Set("a", "alpha2", 3)
	if got := tm.Get("a", 2); got != "alpha" {
		t.Fatalf(`Get("a", 2) = %q, want "alpha"`, got)
	}
	if got := tm.Get("b", 2); got != "beta" {
		t.Fatalf(`Get("b", 2) = %q, want "beta"`, got)
	}
	if got := tm.Get("a", 3); got != "alpha2" {
		t.Fatalf(`Get("a", 3) = %q, want "alpha2"`, got)
	}
}

func TestTimeMap_GetFarAfterLastSet(t *testing.T) {
	tm := NewTimeMap()
	tm.Set("k", "old", 1)
	tm.Set("k", "new", 100)
	if got := tm.Get("k", 10000000); got != "new" {
		t.Fatalf(`Get("k", 10000000) = %q, want "new"`, got)
	}
}

func TestTimeMap_SingleTimestampBoundary(t *testing.T) {
	tm := NewTimeMap()
	tm.Set("k", "onlyval", 1)
	if got := tm.Get("k", 1); got != "onlyval" {
		t.Fatalf(`Get("k", 1) = %q, want "onlyval"`, got)
	}
}

func TestTimeMap_ManyTimestampsBinarySearchCorrectness(t *testing.T) {
	tm := NewTimeMap()
	// Set values at timestamps 10, 20, 30, ..., 1000.
	for ts := 10; ts <= 1000; ts += 10 {
		tm.Set("k", stringValueFor(ts), ts)
	}

	cases := []struct {
		query int
		want  string
	}{
		{5, ""},              // before the first timestamp
		{10, "v10"},          // exact match on the first
		{15, "v10"},          // between 10 and 20
		{999, "v990"},        // just before the last
		{1000, "v1000"},      // exact match on the last
		{5000, "v1000"},      // far beyond the last
		{500, "v500"},        // exact match in the middle
		{505, "v500"},        // between two middle entries
	}
	for _, c := range cases {
		if got := tm.Get("k", c.query); got != c.want {
			t.Errorf("Get(%q, %d) = %q, want %q", "k", c.query, got, c.want)
		}
	}
}

func stringValueFor(ts int) string {
	return "v" + strconv.Itoa(ts)
}
