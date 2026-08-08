package texteditor

import "testing"

func TestInsert_Basic(t *testing.T) {
	e := NewTextEditor()
	if err := e.Insert(0, "hello"); err != nil {
		t.Fatalf("Insert(0, hello) failed: %v", err)
	}
	if e.Text() != "hello" {
		t.Fatalf("Text() = %q, want %q", e.Text(), "hello")
	}
	if err := e.Insert(5, " world"); err != nil {
		t.Fatalf("Insert(5, ' world') failed: %v", err)
	}
	if e.Text() != "hello world" {
		t.Fatalf("Text() = %q, want %q", e.Text(), "hello world")
	}
	if err := e.Insert(5, ","); err != nil {
		t.Fatalf("Insert(5, ',') failed: %v", err)
	}
	if e.Text() != "hello, world" {
		t.Fatalf("Text() = %q, want %q", e.Text(), "hello, world")
	}
}

func TestInsert_OutOfRangeReturnsError(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "abc")
	if err := e.Insert(-1, "x"); err == nil {
		t.Fatal("Insert(-1, x) = nil error, want an error")
	}
	if err := e.Insert(4, "x"); err == nil {
		t.Fatal("Insert(4, x) on a 3-char buffer = nil error, want an error")
	}
	if e.Text() != "abc" {
		t.Fatalf("Text() after failed inserts = %q, want unchanged %q", e.Text(), "abc")
	}
}

func TestDelete_Basic(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "hello world")
	if err := e.Delete(5, 6); err != nil {
		t.Fatalf("Delete(5, 6) failed: %v", err)
	}
	if e.Text() != "hello" {
		t.Fatalf("Text() = %q, want %q", e.Text(), "hello")
	}
}

func TestDelete_InvalidRangeReturnsError(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "hello")
	if err := e.Delete(-1, 2); err == nil {
		t.Fatal("Delete(-1, 2) = nil error, want an error")
	}
	if err := e.Delete(0, -1); err == nil {
		t.Fatal("Delete(0, -1) = nil error, want an error")
	}
	if err := e.Delete(3, 10); err == nil {
		t.Fatal("Delete(3, 10) past the end = nil error, want an error")
	}
	if e.Text() != "hello" {
		t.Fatalf("Text() after failed deletes = %q, want unchanged %q", e.Text(), "hello")
	}
}

func TestUndo_ReversesInsert(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "hello")
	if !e.Undo() {
		t.Fatal("Undo() after an Insert = false, want true")
	}
	if e.Text() != "" {
		t.Fatalf("Text() after undoing the only insert = %q, want \"\"", e.Text())
	}
}

func TestUndo_ReversesDeleteRestoringExactText(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "hello world")
	e.Delete(5, 6) // removes " world"
	if !e.Undo() {
		t.Fatal("Undo() after a Delete = false, want true")
	}
	if e.Text() != "hello world" {
		t.Fatalf("Text() after undoing the delete = %q, want %q", e.Text(), "hello world")
	}
}

func TestRedo_ReappliesInsert(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "hello")
	e.Undo()
	if !e.Redo() {
		t.Fatal("Redo() after undoing an Insert = false, want true")
	}
	if e.Text() != "hello" {
		t.Fatalf("Text() after redo = %q, want %q", e.Text(), "hello")
	}
}

func TestRedo_ReappliesDelete(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "hello world")
	e.Delete(5, 6)
	e.Undo()
	if !e.Redo() {
		t.Fatal("Redo() after undoing a Delete = false, want true")
	}
	if e.Text() != "hello" {
		t.Fatalf("Text() after redo = %q, want %q", e.Text(), "hello")
	}
}

func TestUndoRedo_MultiStepSequenceIsExactlyReversible(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "abc")
	e.Insert(3, "def")
	e.Delete(1, 2) // removes "bc" -> "adef"
	final := e.Text()

	e.Undo()
	e.Undo()
	e.Undo()
	if e.Text() != "" {
		t.Fatalf("Text() after undoing all 3 edits = %q, want \"\"", e.Text())
	}

	e.Redo()
	e.Redo()
	e.Redo()
	if e.Text() != final {
		t.Fatalf("Text() after redoing all 3 edits = %q, want %q (the original final state)", e.Text(), final)
	}
}

func TestNewEditAfterUndoClearsRedoStack(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "hello")
	e.Undo()
	e.Insert(0, "goodbye") // new edit — should discard the redo history for "hello"

	if e.Redo() {
		t.Fatal("Redo() after a new edit following an Undo = true, want false (redo history should be cleared)")
	}
	if e.Text() != "goodbye" {
		t.Fatalf("Text() = %q, want %q", e.Text(), "goodbye")
	}
}

func TestUndo_NothingToUndoReturnsFalse(t *testing.T) {
	e := NewTextEditor()
	if e.Undo() {
		t.Fatal("Undo() on a fresh editor = true, want false")
	}
}

func TestRedo_NothingToRedoReturnsFalse(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "hello")
	if e.Redo() {
		t.Fatal("Redo() with nothing undone = true, want false")
	}
}

func TestFailedEditDoesNotAffectUndoHistory(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "hello")
	e.Insert(-1, "bad") // should fail and not be undoable
	e.Insert(999, "bad")

	if !e.Undo() {
		t.Fatal("Undo() = false, want true (should undo the successful 'hello' insert)")
	}
	if e.Text() != "" {
		t.Fatalf("Text() after undoing = %q, want \"\" (failed inserts must not be on the undo stack)", e.Text())
	}
	if e.Undo() {
		t.Fatal("second Undo() = true, want false (only one successful edit was ever made)")
	}
}

func TestUndo_InterleavedWithNewEdits(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "a")
	e.Insert(1, "b")
	e.Undo() // back to "a"
	e.Insert(1, "c") // "ac", discarding the "b" redo

	if e.Text() != "ac" {
		t.Fatalf("Text() = %q, want %q", e.Text(), "ac")
	}
	e.Undo() // undo the "c" insert
	if e.Text() != "a" {
		t.Fatalf("Text() after undoing 'c' = %q, want %q", e.Text(), "a")
	}
	e.Undo() // undo the "a" insert
	if e.Text() != "" {
		t.Fatalf("Text() after undoing 'a' = %q, want \"\"", e.Text())
	}
}

func TestDelete_UndoThenModifyDifferently(t *testing.T) {
	e := NewTextEditor()
	e.Insert(0, "abcdef")
	e.Delete(2, 2) // removes "cd" -> "abef"
	if e.Text() != "abef" {
		t.Fatalf("Text() = %q, want %q", e.Text(), "abef")
	}
	e.Undo() // restore "cd"
	if e.Text() != "abcdef" {
		t.Fatalf("Text() after undo = %q, want %q", e.Text(), "abcdef")
	}
	e.Delete(0, 3) // now delete "abc" instead -> "def"
	if e.Text() != "def" {
		t.Fatalf("Text() after a different delete post-undo = %q, want %q", e.Text(), "def")
	}
	if e.Redo() {
		t.Fatal("Redo() after a new edit replaced the undone one = true, want false")
	}
}
