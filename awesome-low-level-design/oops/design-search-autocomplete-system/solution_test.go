package designsearchautocompletesystem

import (
	"reflect"
	"testing"
)

// typeChars feeds each character of s through Input in order, returning only the
// result of the FINAL character (mirroring how the example walks through one char
// at a time and inspects the return value at each step).
func typeChars(a *AutocompleteSystem, s string) []string {
	var last []string
	for i := 0; i < len(s); i++ {
		last = a.Input(s[i])
	}
	return last
}

func TestAutocompleteSystem_Example1(t *testing.T) {
	a := NewAutocompleteSystem(
		[]string{"i love you", "island", "iroman", "i love leetcode"},
		[]int{5, 3, 2, 2},
	)

	got := a.Input('i')
	want := []string{"i love you", "island", "i love leetcode"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Input('i') = %v, want %v", got, want)
	}

	got = a.Input(' ')
	want = []string{"i love you", "i love leetcode"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Input(' ') = %v, want %v", got, want)
	}

	got = a.Input('a')
	if len(got) != 0 {
		t.Fatalf("Input('a') = %v, want empty", got)
	}

	got = a.Input('#')
	if len(got) != 0 {
		t.Fatalf("Input('#') = %v, want empty", got)
	}
}

func TestAutocompleteSystem_SentenceCompletedByHashIsRemembered(t *testing.T) {
	a := NewAutocompleteSystem([]string{}, []int{})
	typeChars(a, "hello#")
	got := typeChars(a, "hel")
	want := []string{"hello"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`after typing "hel" = %v, want %v`, got, want)
	}
}

func TestAutocompleteSystem_RepeatedTypingIncreasesHotDegree(t *testing.T) {
	a := NewAutocompleteSystem([]string{"cat", "car"}, []int{1, 1})
	// "car" is tied with "cat" at hot=1; typing it again should push it ahead.
	typeChars(a, "car#")
	got := typeChars(a, "ca")
	want := []string{"car", "cat"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`Input for "ca" = %v, want %v`, got, want)
	}
}

func TestAutocompleteSystem_SinglePrefixMatch(t *testing.T) {
	a := NewAutocompleteSystem([]string{"bb", "aa", "cc"}, []int{1, 1, 1})
	got := a.Input('a')
	// prefix "a" -> only "aa" matches
	want := []string{"aa"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Input('a') = %v, want %v", got, want)
	}
}

func TestAutocompleteSystem_TieBrokenLexicographicallySamePrefix(t *testing.T) {
	a := NewAutocompleteSystem([]string{"app", "apt", "aps"}, []int{3, 3, 3})
	got := a.Input('a')
	want := []string{"app", "aps", "apt"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Input('a') = %v, want %v", got, want)
	}
}

func TestAutocompleteSystem_OnlyTopThreeReturned(t *testing.T) {
	a := NewAutocompleteSystem(
		[]string{"a1", "a2", "a3", "a4", "a5"},
		[]int{5, 4, 3, 2, 1},
	)
	got := a.Input('a')
	want := []string{"a1", "a2", "a3"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Input('a') = %v, want %v", got, want)
	}
}

func TestAutocompleteSystem_NoMatchesReturnsEmpty(t *testing.T) {
	a := NewAutocompleteSystem([]string{"hello"}, []int{1})
	got := a.Input('z')
	if len(got) != 0 {
		t.Fatalf("Input('z') = %v, want empty", got)
	}
}

func TestAutocompleteSystem_HashClearsBufferForNextSentence(t *testing.T) {
	a := NewAutocompleteSystem([]string{}, []int{})
	typeChars(a, "abc#")
	typeChars(a, "xyz#")
	got := typeChars(a, "a")
	want := []string{"abc"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`after new session typing "a" = %v, want %v`, got, want)
	}
}

func TestAutocompleteSystem_EmptyInitialHistory(t *testing.T) {
	a := NewAutocompleteSystem([]string{}, []int{})
	got := a.Input('h')
	if len(got) != 0 {
		t.Fatalf("Input('h') on empty history = %v, want empty", got)
	}
	got = typeChars(a, "i there#")
	if len(got) != 0 {
		t.Fatalf(`Input('#') = %v, want empty`, got)
	}
	got = typeChars(a, "hi")
	want := []string{"hi there"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`typing "hi" = %v, want %v`, got, want)
	}
}

func TestAutocompleteSystem_MultipleSessionsAccumulateHotDegree(t *testing.T) {
	a := NewAutocompleteSystem([]string{"go", "golang"}, []int{1, 1})
	for i := 0; i < 3; i++ {
		typeChars(a, "go#")
	}
	got := a.Input('g')
	want := []string{"go", "golang"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Input('g') = %v, want %v", got, want)
	}
}
