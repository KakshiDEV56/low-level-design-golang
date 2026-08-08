package liskovsubstitution

import "testing"

// Compile-time proof that Rectangle satisfies both interfaces.
var (
	_ Shape     = (*Rectangle)(nil)
	_ Resizable = (*Rectangle)(nil)
)

// Compile-time proof that Square satisfies Shape...
var _ Shape = (*Square)(nil)

// ...but NOT Resizable. Uncommenting the line below should fail to compile once
// Square no longer has SetWidth/SetHeight methods — that's the whole point of the
// fix. (A compile error can't be asserted at runtime, so this is documented here
// instead of executed; try uncommenting it yourself after refactoring to see it fail.)
//
// var _ Resizable = (*Square)(nil)

func TestRectangle_IndependentResizing(t *testing.T) {
	r := NewRectangle(2, 3)
	if got := r.Area(); got != 6 {
		t.Fatalf("Area() = %d, want 6", got)
	}

	r.SetWidth(10) // setting width alone must NOT affect height
	if got := r.Area(); got != 30 {
		t.Fatalf("after SetWidth(10): Area() = %d, want 30 (10x3)", got)
	}

	r.SetHeight(4) // setting height alone must NOT affect width
	if got := r.Area(); got != 40 {
		t.Fatalf("after SetHeight(4): Area() = %d, want 40 (10x4)", got)
	}
}

func TestSquare_SetSizeKeepsInvariant(t *testing.T) {
	sq := NewSquare(5)
	if got := sq.Area(); got != 25 {
		t.Fatalf("Area() = %d, want 25", got)
	}

	sq.SetSize(3)
	if got := sq.Area(); got != 9 {
		t.Fatalf("after SetSize(3): Area() = %d, want 9", got)
	}

	sq.SetSize(0)
	if got := sq.Area(); got != 0 {
		t.Fatalf("after SetSize(0): Area() = %d, want 0", got)
	}
}

// resizableShape combines both interfaces — only a type that satisfies the FULL
// behavioral contract (independent resizing AND reporting its area) can be passed
// here.
type resizableShape interface {
	Resizable
	Shape
}

// resizeAndArea sets width and height INDEPENDENTLY (as any real Resizable must
// support) and returns the resulting area. This function has no idea what concrete
// type it's given — it just trusts the Resizable contract.
func resizeAndArea(rs resizableShape, w, h int) int {
	rs.SetWidth(w)
	rs.SetHeight(h)
	return rs.Area()
}

func TestResizeAndArea_WorksForAnyGenuineResizable(t *testing.T) {
	r := NewRectangle(1, 1)
	if got := resizeAndArea(r, 6, 7); got != 42 {
		t.Fatalf("resizeAndArea(Rectangle, 6, 7) = %d, want 42", got)
	}
}

func TestRectangle_MultipleIndependentResizesInSequence(t *testing.T) {
	r := NewRectangle(1, 1)
	widths := []int{2, 5, 3, 8}
	height := 4
	r.SetHeight(height)
	for _, w := range widths {
		r.SetWidth(w)
		if got, want := r.Area(), w*height; got != want {
			t.Fatalf("SetWidth(%d) then Area() = %d, want %d (height must stay %d)", w, got, want, height)
		}
	}
}
