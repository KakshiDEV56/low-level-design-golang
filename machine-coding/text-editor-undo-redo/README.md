# Text Editor with Undo/Redo (Command + Memento)

**Suggested time box: 45–60 minutes.**

## Problem statement

Design a plain-text buffer that supports inserting and deleting text at arbitrary
positions, where **every edit is undoable and redoable** — the classic Command pattern
(each edit is an object that knows how to apply itself and how to reverse itself)
combined with a Memento-style capture of "what was there before" so a delete can be
perfectly restored on undo.

## Functional requirements

1. `Text() string` returns the current full contents of the buffer (starts empty).
2. `Insert(pos int, text string) error` inserts `text` at byte offset `pos` (`0 <= pos
   <= len(current text)`; inserting at `len(text)` appends to the end). Returns an
   error (and changes nothing) if `pos` is out of that range.
3. `Delete(pos int, length int) error` removes `length` bytes starting at `pos`.
   Returns an error (and changes nothing) if `pos < 0`, `length < 0`, or
   `pos+length > len(current text)`.
4. `Undo() bool` reverses the most recent not-yet-undone edit (`Insert` or `Delete`).
   Returns `false` (no-op) if there's nothing to undo. A **failed** `Insert`/`Delete`
   call (one that returned an error) must NOT be something `Undo` can act on — only
   successful edits go on the undo history.
5. `Redo() bool` re-applies the most recently undone edit. Returns `false` (no-op) if
   there's nothing to redo.
6. Calling `Insert` or `Delete` (successfully) after one or more `Undo` calls discards
   the redo history — this is standard editor behavior (once you make a new edit, the
   "future" you undid from is gone).

## Rubric

- [ ] `Undo` on a `Delete` restores the **exact** text that was removed — not just "some
      placeholder" — meaning your `Delete` command needs to capture the deleted
      substring at delete-time (before it's gone) to hand back on undo.
- [ ] `Undo`/`Redo` are genuinely inverse operations — running `Undo()` then `Redo()`
      any number of times must always return the buffer to the exact same state it
      would have been in without the Undo/Redo pair.
- [ ] A failed edit call (invalid `pos`/`length`) never appears in the undo history and
      never changes buffer contents.
- [ ] Making a new edit after undoing clears the redo history (verified, not just
      assumed).
- [ ] The undo/redo mechanism is generic over the *kind* of edit (Insert vs. Delete) —
      look for a shared abstraction (an interface with `do`/`undo`, or similar) rather
      than parallel `insertUndoStack`/`deleteUndoStack` special-cased by type.

## Hints

- Define a small internal interface, e.g. `type command interface { do(e *TextEditor);
  undo(e *TextEditor) }`, with two implementations: one for insert, one for delete.
  This is the Command pattern — `TextEditor` just keeps two stacks of `command`
  (undo stack, redo stack) and never needs to know which concrete kind it's holding.
- The delete command must capture the substring being removed **before** mutating the
  buffer, and store it on the command itself — that stored copy is what makes `undo`
  possible (this is the Memento half of the pattern: the command carries a snapshot of
  what it's about to destroy).
- `Undo`: pop from the undo stack, call `.undo(e)`, push onto the redo stack.
  `Redo`: pop from the redo stack, call `.do(e)`, push onto the undo stack. A
  successful new `Insert`/`Delete` call: push the new command onto the undo stack,
  and clear the redo stack entirely.
