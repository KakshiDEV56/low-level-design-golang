package texteditor

import "errors"

/*
See README.md for the full problem statement, requirements, and rubric.

TODO:
 1. Add the fields you need on TextEditor below. Hint: the current text (string or
    []byte), plus an undo stack and a redo stack of some shared "command" abstraction
    — see the README hints for the suggested interface shape.
 2. Implement NewTextEditor: empty text, empty stacks.
 3. Implement Text, Insert, Delete, Undo, Redo per the README. Consider defining a
    small internal command interface with do/undo methods, and two implementations
    (insert, delete) — this is the Command pattern, and it's what lets Undo/Redo be
    generic over which kind of edit they're reversing.
*/

var ErrOutOfRange = errors.New("position out of range")

type TextEditor struct {
	// TODO: add your fields here
}

func NewTextEditor() *TextEditor {
	// TODO: implement
	return &TextEditor{}
}

func (e *TextEditor) Text() string {
	// TODO: implement
	return ""
}

func (e *TextEditor) Insert(pos int, text string) error {
	// TODO: implement
	return nil
}

func (e *TextEditor) Delete(pos int, length int) error {
	// TODO: implement
	return nil
}

func (e *TextEditor) Undo() bool {
	// TODO: implement
	return false
}

func (e *TextEditor) Redo() bool {
	// TODO: implement
	return false
}
