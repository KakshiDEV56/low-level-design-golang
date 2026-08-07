package designaddandsearchwords

/*
LeetCode 211. Design Add and Search Words Data Structure

Design a data structure that supports adding new words and finding if a string matches
any previously added string.

Implement the WordDictionary class:
  - NewWordDictionary() initializes the object.
  - AddWord(word string) adds word to the data structure, it can be matched later.
  - Search(word string) bool returns true if there is any string in the data structure
    that matches word, or false otherwise. word may contain dots '.' where dots can be
    matched with any letter.

Example 1:

	Input:
	["WordDictionary", "addWord", "addWord", "addWord", "search", "search", "search",
	 "search"]
	[[], ["bad"], ["dad"], ["mad"], ["pad"], ["bad"], [".ad"], ["b.."]]
	Output:
	[null, null, null, null, false, true, true, true]

	Explanation:
	WordDictionary wd = new WordDictionary();
	wd.addWord("bad");
	wd.addWord("dad");
	wd.addWord("mad");
	wd.search("pad"); // return false, "pad" was never added
	wd.search("bad"); // return true, "bad" was added
	wd.search(".ad"); // return true, matches "bad"/"dad"/"mad" via the wildcard
	wd.search("b.."); // return true, matches "bad"

Constraints:
  - 1 <= word.length <= 25
  - word in AddWord consists of lowercase English letters.
  - word in Search consists of '.' or lowercase English letters.
  - There will be at most 2 dots in word for each call to Search.
  - At most 10^4 calls will be made to AddWord and Search.

TODO:
 1. Add the fields you need on WordDictionary below. This is a classic Trie (prefix
    tree): a node has up to 26 children (one per lowercase letter) and a flag marking
    "a complete word ends here". A common Go representation:
      type trieNode struct {
          children [26]*trieNode
          isWord   bool
      }
 2. Implement NewWordDictionary: create an empty root trieNode.
 3. Implement AddWord: walk from the root, one character at a time, creating child nodes
    as needed (word consists only of lowercase letters here, no dots), and mark the
    final node's isWord = true.
 4. Implement Search: this is the interesting part because of '.'. A recursive helper
    that takes (node, word, index) works well:
      - If index == len(word): return node.isWord.
      - If word[index] is a normal letter: follow that single child if it exists (return
        false if it doesn't), recurse with index+1.
      - If word[index] is '.': try EVERY non-nil child, recursing into each; return true
        if any of them leads to a match.
    Kick off the recursion from the root with index 0.
*/

type WordDictionary struct {
	// TODO: add your fields here
}

func NewWordDictionary() *WordDictionary {
	// TODO: implement
	return &WordDictionary{}
}

func (w *WordDictionary) AddWord(word string) {
	// TODO: implement
}

func (w *WordDictionary) Search(word string) bool {
	// TODO: implement
	return false
}
