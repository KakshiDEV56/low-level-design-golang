package designsearchautocompletesystem

/*
LeetCode 642. Design Search Autocomplete System

Design a search autocomplete system for a search engine. Users may input a sentence
(at least one word and end with a special character '#').

You are given a string array sentences and an integer array times both of length n,
where sentences[i] is a previously typed sentence and times[i] is the corresponding
number of times the sentence was typed. More sentences may be typed in during a session
of autocomplete.

Implement the AutocompleteSystem class:
  - NewAutocompleteSystem(sentences []string, times []int) initializes the object with
    the sentences and times arrays.
  - Input(c byte) []string this indicates that the user typed the character c.
      - If c == '#', the user finished typing the sentence: it should be saved as a
        historical sentence in your system (incrementing its hot degree by 1 if it
        already existed, or being added with hot degree 1 if it's new). The input buffer
        should be cleared, and an empty array should be returned.
      - Otherwise, c is appended to the input. Return the top 3 historical sentences that
        have the same prefix as the input (case-sensitive). These sentences should be
        sorted by hot degree (the number of times a sentence has been typed) in
        descending order. If several sentences have the same hot degree, sort them
        lexicographically in ascending order. If there are fewer than 3 matches, return
        as many as there are (possibly none).

Example 1:

	Input:
	["AutocompleteSystem", "input", "input", "input", "input"]
	[[["i love you", "island", "iroman", "i love leetcode"], [5, 3, 2, 2]],
	 ["i"], [" "], ["a"], ["#"]]
	Output:
	[null,
	 ["i love you", "island", "i love leetcode"],
	 ["i love you", "i love leetcode"],
	 [],
	 []]

	Explanation:
	AutocompleteSystem obj = new AutocompleteSystem(
	    ["i love you", "island", "iroman", "i love leetcode"], [5, 3, 2, 2]);
	obj.input('i');   // returns ["i love you", "island", "i love leetcode"], the 3
	                  // sentences that have prefix "i" sorted by hot degree (5, 3, 2),
	                  // ties broken lexicographically ("i love leetcode" < "iroman").
	obj.input(' ');   // returns ["i love you", "i love leetcode"], prefix is now "i ",
	                  // "island"/"iroman" no longer match.
	obj.input('a');   // returns [], prefix is now "i a", nothing matches.
	obj.input('#');   // the user finished typing "i a"; it's saved as a new sentence
	                  // with hot degree 1, input buffer clears, returns [].

Constraints:
  - n == sentences.length == times.length
  - 1 <= n <= 100
  - 1 <= sentences[i].length <= 100
  - 1 <= times[i] <= 50
  - c is a lowercase English letter, a space ' ', or a special character '#'.
  - Each tested sentence will be a sequence of characters c that end with the character
    '#'.
  - Each tested sentence will have a length in the range [1, 200].
  - The words in each input sentence are separated by single spaces.
  - At most 5000 calls will be made to Input.

TODO:
 1. Add the fields you need on AutocompleteSystem below. Two workable approaches:
      a) Simplest: a map from sentence -> hot degree (its typed count), plus a growing
         string buffer for "what the current user has typed so far". On each non-'#'
         Input call, scan every stored sentence, keep the ones with the buffer as a
         prefix (strings.HasPrefix), sort by (hot degree desc, sentence asc), return the
         top 3. This is O(n) per keystroke, which is plenty fast for n <= 100+5000 calls.
      b) More advanced (worth trying after (a) works): a trie where each node also
         tracks hot-degree info for the sentences passing through it, so each keystroke
         only has to look at the matching subtree instead of scanning everything.
 2. Implement NewAutocompleteSystem: seed your storage from sentences/times, buffer
    starts empty.
 3. Implement Input(c):
      - If c == '#': record the current buffer as a completed sentence (increment or
        insert with hot degree 1), clear the buffer, return an empty (non-nil, so JSON
        marshaling doesn't matter here, but keep it an empty slice not literally nil if
        your tests check that — check what your own tests expect) slice.
      - Otherwise: append c to the buffer, compute and return the top 3 matches as
        described above.
*/

type AutocompleteSystem struct {
	// TODO: add your fields here
}

func NewAutocompleteSystem(sentences []string, times []int) *AutocompleteSystem {
	// TODO: implement
	return &AutocompleteSystem{}
}

func (a *AutocompleteSystem) Input(c byte) []string {
	// TODO: implement
	return nil
}
