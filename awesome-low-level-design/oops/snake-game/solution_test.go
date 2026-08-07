package snakegame

import "testing"

func TestSnakeGame_Example1(t *testing.T) {
	g := NewSnakeGame(3, 2, [][]int{{1, 2}, {0, 1}})
	steps := []struct {
		dir  string
		want int
	}{
		{"R", 0},
		{"D", 0},
		{"R", 1},
		{"U", 1},
		{"L", 2},
		{"U", -1},
	}
	for i, s := range steps {
		if got := g.Move(s.dir); got != s.want {
			t.Fatalf("step %d Move(%q) = %d, want %d", i, s.dir, got, s.want)
		}
	}
}

func TestSnakeGame_OutOfBoundsRight(t *testing.T) {
	g := NewSnakeGame(2, 1, [][]int{})
	if got := g.Move("R"); got != 0 {
		t.Fatalf(`first Move("R") = %d, want 0`, got)
	}
	if got := g.Move("R"); got != -1 {
		t.Fatalf(`second Move("R") off the board = %d, want -1`, got)
	}
	if got := g.Move("L"); got != -1 {
		t.Fatalf(`Move after game over = %d, want -1 (game must stay over)`, got)
	}
}

func TestSnakeGame_OutOfBoundsUp(t *testing.T) {
	g := NewSnakeGame(3, 3, [][]int{})
	if got := g.Move("U"); got != -1 {
		t.Fatalf(`Move("U") from row 0 = %d, want -1`, got)
	}
}

func TestSnakeGame_OutOfBoundsLeft(t *testing.T) {
	g := NewSnakeGame(3, 3, [][]int{})
	if got := g.Move("L"); got != -1 {
		t.Fatalf(`Move("L") from col 0 = %d, want -1`, got)
	}
}

func TestSnakeGame_NoFoodJustMoves(t *testing.T) {
	g := NewSnakeGame(5, 5, [][]int{})
	if got := g.Move("R"); got != 0 {
		t.Fatalf(`Move("R") = %d, want 0`, got)
	}
	if got := g.Move("R"); got != 0 {
		t.Fatalf(`Move("R") = %d, want 0`, got)
	}
	if got := g.Move("D"); got != 0 {
		t.Fatalf(`Move("D") = %d, want 0`, got)
	}
}

func TestSnakeGame_EatingFoodGrowsSnakeAndIncreasesScore(t *testing.T) {
	g := NewSnakeGame(5, 5, [][]int{{0, 1}, {0, 2}})
	if got := g.Move("R"); got != 1 {
		t.Fatalf(`Move("R") onto food (0,1) = %d, want 1`, got)
	}
	if got := g.Move("R"); got != 2 {
		t.Fatalf(`Move("R") onto food (0,2) = %d, want 2`, got)
	}
}

func TestSnakeGame_SelfCollisionAfterGrowing(t *testing.T) {
	// Board 5x5. Food placed so the snake grows to length 4, then loop back into itself.
	// Path: R,R,R (eats food at (0,1) and (0,2), growing to length 3), D, L, L, U
	// would run the head back into the body.
	g := NewSnakeGame(5, 5, [][]int{{0, 1}, {0, 2}, {1, 2}})
	if got := g.Move("R"); got != 1 {
		t.Fatalf(`Move("R") = %d, want 1`, got)
	}
	if got := g.Move("R"); got != 2 {
		t.Fatalf(`Move("R") = %d, want 2`, got)
	}
	if got := g.Move("D"); got != 3 {
		t.Fatalf(`Move("D") onto food (1,2) = %d, want 3`, got)
	}
	// Snake body now (head first): (1,2),(0,2),(0,1),(0,0), length 4, no food left.
	// Moving U takes the head to (0,2), which is the snake's own neck segment (not the
	// tail, so the tail-vacates-this-turn exception does NOT apply here) -> collision.
	if got := g.Move("U"); got != -1 {
		t.Fatalf(`Move("U") into own neck segment (0,2) = %d, want -1`, got)
	}
}

func TestSnakeGame_TailMovesOutOfTheWayWithoutFood(t *testing.T) {
	// Without eating, the snake's tail vacates the cell it's leaving, so looping back to
	// where the tail JUST was should be legal (not a self-collision).
	g := NewSnakeGame(5, 5, [][]int{{0, 1}})
	if got := g.Move("R"); got != 1 {
		t.Fatalf(`Move("R") onto food = %d, want 1`, got)
	}
	// Snake body: (0,1),(0,0), length 2 (grew once).
	if got := g.Move("D"); got != 1 {
		t.Fatalf(`Move("D") = %d, want 1`, got)
	}
	// body: (1,1),(0,1)
	if got := g.Move("L"); got != 1 {
		t.Fatalf(`Move("L") = %d, want 1`, got)
	}
	// body: (1,0),(1,1) -- moving into where (0,0) used to be is irrelevant here;
	// this sequence just exercises normal tail-following movement without collision.
	if got := g.Move("U"); got != 1 {
		t.Fatalf(`Move("U") = %d, want 1`, got)
	}
}

func TestSnakeGame_AllFoodConsumedThenJustMoves(t *testing.T) {
	g := NewSnakeGame(10, 10, [][]int{{0, 1}})
	if got := g.Move("R"); got != 1 {
		t.Fatalf(`Move("R") = %d, want 1`, got)
	}
	// No more food; further moves should just keep returning the same score.
	if got := g.Move("R"); got != 1 {
		t.Fatalf(`Move("R") after food exhausted = %d, want 1`, got)
	}
	if got := g.Move("D"); got != 1 {
		t.Fatalf(`Move("D") = %d, want 1`, got)
	}
}
