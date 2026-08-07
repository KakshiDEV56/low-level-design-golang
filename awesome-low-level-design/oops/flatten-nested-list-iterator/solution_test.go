package flattennestedlistiterator

import (
	"reflect"
	"testing"
)

func drain(it *NestedIterator) []int {
	var out []int
	for it.HasNext() {
		out = append(out, it.Next())
	}
	return out
}

func TestNestedIterator_Example1(t *testing.T) {
	// [[1,1],2,[1,1]]
	nested := []*NestedInteger{
		NewList(NewInt(1), NewInt(1)),
		NewInt(2),
		NewList(NewInt(1), NewInt(1)),
	}
	it := NewNestedIterator(nested)
	got := drain(it)
	want := []int{1, 1, 2, 1, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("drain() = %v, want %v", got, want)
	}
}

func TestNestedIterator_Example2(t *testing.T) {
	// [1,[4,[6]]]
	nested := []*NestedInteger{
		NewInt(1),
		NewList(NewInt(4), NewList(NewInt(6))),
	}
	it := NewNestedIterator(nested)
	got := drain(it)
	want := []int{1, 4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("drain() = %v, want %v", got, want)
	}
}

func TestNestedIterator_FlatListNoNesting(t *testing.T) {
	nested := []*NestedInteger{NewInt(1), NewInt(2), NewInt(3)}
	it := NewNestedIterator(nested)
	got := drain(it)
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("drain() = %v, want %v", got, want)
	}
}

func TestNestedIterator_EmptyNestedLists(t *testing.T) {
	// [[],1,[],2,[]] -- empty sublists must be skipped entirely.
	nested := []*NestedInteger{
		NewList(),
		NewInt(1),
		NewList(),
		NewInt(2),
		NewList(),
	}
	it := NewNestedIterator(nested)
	got := drain(it)
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("drain() = %v, want %v", got, want)
	}
}

func TestNestedIterator_AllEmptyHasNoElements(t *testing.T) {
	nested := []*NestedInteger{NewList(), NewList(NewList(), NewList())}
	it := NewNestedIterator(nested)
	if it.HasNext() {
		t.Fatal("HasNext() = true, want false (no integers anywhere)")
	}
}

func TestNestedIterator_DeeplyNested(t *testing.T) {
	// [1,[2,[3,[4,[5]]]]]
	nested := []*NestedInteger{
		NewInt(1),
		NewList(
			NewInt(2),
			NewList(
				NewInt(3),
				NewList(
					NewInt(4),
					NewList(NewInt(5)),
				),
			),
		),
	}
	it := NewNestedIterator(nested)
	got := drain(it)
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("drain() = %v, want %v", got, want)
	}
}

func TestNestedIterator_NegativeNumbers(t *testing.T) {
	nested := []*NestedInteger{NewInt(-1), NewList(NewInt(-2), NewInt(3))}
	it := NewNestedIterator(nested)
	got := drain(it)
	want := []int{-1, -2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("drain() = %v, want %v", got, want)
	}
}

func TestNestedIterator_HasNextIsIdempotent(t *testing.T) {
	it := NewNestedIterator([]*NestedInteger{NewInt(7)})
	if !it.HasNext() {
		t.Fatal("HasNext() = false, want true")
	}
	if !it.HasNext() {
		t.Fatal("HasNext() second call = false, want true (must not consume)")
	}
	if got := it.Next(); got != 7 {
		t.Fatalf("Next() = %d, want 7", got)
	}
	if it.HasNext() {
		t.Fatal("HasNext() after consuming only element = true, want false")
	}
}

func TestNestedIterator_MixedTopLevelManyGroups(t *testing.T) {
	// [[1],[],[2,3],[[4]],5]
	nested := []*NestedInteger{
		NewList(NewInt(1)),
		NewList(),
		NewList(NewInt(2), NewInt(3)),
		NewList(NewList(NewInt(4))),
		NewInt(5),
	}
	it := NewNestedIterator(nested)
	got := drain(it)
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("drain() = %v, want %v", got, want)
	}
}
