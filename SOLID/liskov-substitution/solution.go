package liskovsubstitution

/*
This code compiles and works today — but it violates the Liskov Substitution
Principle. See README.md for the full explanation (the classic Rectangle/Square trap).

Square embeds *Rectangle and overrides SetWidth/SetHeight to force both dimensions to
match. This means Square structurally satisfies the same "shape" as Rectangle
(SetWidth, SetHeight, Area), but violates the behavioral contract that setting one
dimension shouldn't affect the other — silently breaking any code that assumes it can
resize a "Resizable" thing on either axis independently.

TODO: refactor this file so that:
 1. A Shape interface is introduced: { Area() int }.
 2. A Resizable interface is introduced: { SetWidth(w int); SetHeight(h int) }.
 3. Rectangle keeps NewRectangle/SetWidth/SetHeight/Area, and satisfies BOTH Shape and
    Resizable.
 4. Square becomes its OWN independent type — NOT embedding Rectangle at all — with
    NewSquare(size int) and a SetSize(size int) method (instead of SetWidth/SetHeight),
    plus Area(). Square satisfies Shape but must NOT satisfy Resizable (it should have
    no SetWidth/SetHeight methods whatsoever, so passing a *Square where a Resizable is
    expected is a compile error, not a runtime surprise).

Once refactored, solution_test.go should compile and pass.
*/

// Rectangle supports resizing on either axis independently.
type Rectangle struct {
	width, height int
}

func NewRectangle(w, h int) *Rectangle {
	return &Rectangle{width: w, height: h}
}

func (r *Rectangle) SetWidth(w int)  { r.width = w }
func (r *Rectangle) SetHeight(h int) { r.height = h }
func (r *Rectangle) Area() int       { return r.width * r.height }

// Square embeds *Rectangle and "inherits" its methods, but overrides SetWidth and
// SetHeight to keep both dimensions equal — which breaks any code written against
// Rectangle's implied contract that width and height can be set independently.
// TODO: this whole embedding relationship needs to go away (see above).
type Square struct {
	*Rectangle
}

func NewSquare(size int) *Square {
	return &Square{Rectangle: NewRectangle(size, size)}
}

func (s *Square) SetWidth(w int) {
	s.Rectangle.SetWidth(w)
	s.Rectangle.SetHeight(w)
}

func (s *Square) SetHeight(h int) {
	s.Rectangle.SetWidth(h)
	s.Rectangle.SetHeight(h)
}
