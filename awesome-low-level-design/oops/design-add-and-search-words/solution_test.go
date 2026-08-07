package designaddandsearchwords

import "testing"

func TestWordDictionary_Example1(t *testing.T) {
	w := NewWordDictionary()
	w.AddWord("bad")
	w.AddWord("dad")
	w.AddWord("mad")

	if got := w.Search("pad"); got {
		t.Fatal(`Search("pad") = true, want false`)
	}
	if got := w.Search("bad"); !got {
		t.Fatal(`Search("bad") = false, want true`)
	}
	if got := w.Search(".ad"); !got {
		t.Fatal(`Search(".ad") = false, want true`)
	}
	if got := w.Search("b.."); !got {
		t.Fatal(`Search("b..") = false, want true`)
	}
}

func TestWordDictionary_SearchOnEmptyDictionary(t *testing.T) {
	w := NewWordDictionary()
	if got := w.Search("anything"); got {
		t.Fatal(`Search("anything") on empty dictionary = true, want false`)
	}
	if got := w.Search("."); got {
		t.Fatal(`Search(".") on empty dictionary = true, want false`)
	}
}

func TestWordDictionary_ExactWordNoDots(t *testing.T) {
	w := NewWordDictionary()
	w.AddWord("hello")
	if got := w.Search("hello"); !got {
		t.Fatal(`Search("hello") = false, want true`)
	}
	if got := w.Search("hell"); got {
		t.Fatal(`Search("hell") = true, want false (prefix is not itself a stored word)`)
	}
	if got := w.Search("helloo"); got {
		t.Fatal(`Search("helloo") = true, want false`)
	}
}

func TestWordDictionary_AllDots(t *testing.T) {
	w := NewWordDictionary()
	w.AddWord("cat")
	if got := w.Search("..."); !got {
		t.Fatal(`Search("...") = false, want true`)
	}
	if got := w.Search(".."); got {
		t.Fatal(`Search("..") = true, want false (wrong length)`)
	}
	if got := w.Search("...."); got {
		t.Fatal(`Search("....") = true, want false (wrong length)`)
	}
}

func TestWordDictionary_DotMatchesOnlyOneCharacter(t *testing.T) {
	w := NewWordDictionary()
	w.AddWord("a")
	if got := w.Search("."); !got {
		t.Fatal(`Search(".") = false, want true`)
	}
	if got := w.Search(".."); got {
		t.Fatal(`Search("..") = true, want false (dot matches exactly one char, not zero or many)`)
	}
}

func TestWordDictionary_MultipleWordsSamePrefix(t *testing.T) {
	w := NewWordDictionary()
	w.AddWord("car")
	w.AddWord("cart")
	w.AddWord("care")
	if got := w.Search("car"); !got {
		t.Fatal(`Search("car") = false, want true`)
	}
	if got := w.Search("care"); !got {
		t.Fatal(`Search("care") = false, want true`)
	}
	if got := w.Search("cart"); !got {
		t.Fatal(`Search("cart") = false, want true`)
	}
	if got := w.Search("ca"); got {
		t.Fatal(`Search("ca") = true, want false (not a complete added word)`)
	}
}

func TestWordDictionary_DotInMiddle(t *testing.T) {
	w := NewWordDictionary()
	w.AddWord("dog")
	w.AddWord("dig")
	w.AddWord("dug")
	if got := w.Search("d.g"); !got {
		t.Fatal(`Search("d.g") = false, want true`)
	}
	if got := w.Search("d.x"); got {
		t.Fatal(`Search("d.x") = true, want false`)
	}
}

func TestWordDictionary_TwoDotsAtOnce(t *testing.T) {
	w := NewWordDictionary()
	w.AddWord("abcd")
	if got := w.Search("a..d"); !got {
		t.Fatal(`Search("a..d") = false, want true`)
	}
	if got := w.Search(".b.d"); !got {
		t.Fatal(`Search(".b.d") = false, want true`)
	}
	if got := w.Search("a.c."); !got {
		t.Fatal(`Search("a.c.") = false, want true`)
	}
}

func TestWordDictionary_NoFalsePositiveAcrossUnrelatedBranches(t *testing.T) {
	w := NewWordDictionary()
	w.AddWord("apple")
	w.AddWord("banana")
	if got := w.Search(".anana"); !got {
		t.Fatal(`Search(".anana") = false, want true (matches "banana")`)
	}
	if got := w.Search(".pple"); !got {
		t.Fatal(`Search(".pple") = false, want true (matches "apple")`)
	}
	if got := w.Search(".range"); got {
		t.Fatal(`Search(".range") = true, want false (nothing added matches)`)
	}
}

func TestWordDictionary_DuplicateAddWordIsHarmless(t *testing.T) {
	w := NewWordDictionary()
	w.AddWord("dup")
	w.AddWord("dup")
	if got := w.Search("dup"); !got {
		t.Fatal(`Search("dup") = false, want true`)
	}
}

func TestWordDictionary_LargerVocabularyStress(t *testing.T) {
	w := NewWordDictionary()
	words := []string{
		"time", "timer", "timex", "tim", "tin", "tan", "ton", "tone",
		"gone", "gene", "genie", "green", "grep", "great",
	}
	for _, word := range words {
		w.AddWord(word)
	}
	for _, word := range words {
		if got := w.Search(word); !got {
			t.Fatalf("Search(%q) = false, want true", word)
		}
	}
	if got := w.Search("....."); !got {
		t.Fatal(`Search(".....") = false, want true (matches "genie" or "green")`)
	}
	if got := w.Search("t.m.r"); !got {
		t.Fatal(`Search("t.m.r") = false, want true (matches "timer")`)
	}
	if got := w.Search("zzzzz"); got {
		t.Fatal(`Search("zzzzz") = true, want false`)
	}
}
