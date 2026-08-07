package designtictactoe

import "testing"

func TestTicTacToe_Example1(t *testing.T) {
	g := NewTicTacToe(3)
	moves := []struct {
		row, col, player, want int
	}{
		{0, 0, 1, 0},
		{0, 2, 2, 0},
		{2, 2, 1, 0},
		{1, 1, 2, 0},
		{2, 0, 1, 0},
		{1, 0, 2, 0},
		{2, 1, 1, 1},
	}
	for i, m := range moves {
		if got := g.Move(m.row, m.col, m.player); got != m.want {
			t.Fatalf("move %d Move(%d,%d,%d) = %d, want %d", i, m.row, m.col, m.player, got, m.want)
		}
	}
}

func TestTicTacToe_RowWin(t *testing.T) {
	g := NewTicTacToe(3)
	g.Move(0, 0, 1)
	g.Move(1, 0, 2)
	g.Move(0, 1, 1)
	g.Move(1, 1, 2)
	if got := g.Move(0, 2, 1); got != 1 {
		t.Fatalf("Move(0,2,1) completing row 0 = %d, want 1", got)
	}
}

func TestTicTacToe_ColumnWin(t *testing.T) {
	g := NewTicTacToe(3)
	g.Move(0, 0, 1)
	g.Move(0, 1, 2)
	g.Move(1, 0, 1)
	g.Move(1, 1, 2)
	if got := g.Move(2, 0, 1); got != 1 {
		t.Fatalf("Move(2,0,1) completing column 0 = %d, want 1", got)
	}
}

func TestTicTacToe_MainDiagonalWin(t *testing.T) {
	g := NewTicTacToe(3)
	g.Move(0, 0, 2)
	g.Move(0, 1, 1)
	g.Move(1, 1, 2)
	g.Move(0, 2, 1)
	if got := g.Move(2, 2, 2); got != 2 {
		t.Fatalf("Move(2,2,2) completing main diagonal = %d, want 2", got)
	}
}

func TestTicTacToe_AntiDiagonalWin(t *testing.T) {
	g := NewTicTacToe(3)
	g.Move(0, 2, 1)
	g.Move(0, 0, 2)
	g.Move(1, 1, 1)
	g.Move(0, 1, 2)
	if got := g.Move(2, 0, 1); got != 1 {
		t.Fatalf("Move(2,0,1) completing anti-diagonal = %d, want 1", got)
	}
}

func TestTicTacToe_NoWinnerYet(t *testing.T) {
	g := NewTicTacToe(3)
	if got := g.Move(0, 0, 1); got != 0 {
		t.Fatalf("Move(0,0,1) = %d, want 0", got)
	}
	if got := g.Move(1, 1, 2); got != 0 {
		t.Fatalf("Move(1,1,2) = %d, want 0", got)
	}
}

func TestTicTacToe_PartialRowNotFullDoesNotWin(t *testing.T) {
	g := NewTicTacToe(3)
	g.Move(0, 0, 1)
	g.Move(1, 0, 2)
	if got := g.Move(0, 1, 1); got != 0 {
		t.Fatalf("Move(0,1,1) = %d, want 0 (only 2 of 3 in the row)", got)
	}
}

func TestTicTacToe_LargerBoardN5RowWin(t *testing.T) {
	g := NewTicTacToe(5)
	for col := 0; col < 4; col++ {
		if got := g.Move(2, col, 1); got != 0 {
			t.Fatalf("Move(2,%d,1) = %d, want 0", col, got)
		}
		if got := g.Move(0, col, 2); got != 0 {
			t.Fatalf("Move(0,%d,2) = %d, want 0", col, got)
		}
	}
	if got := g.Move(2, 4, 1); got != 1 {
		t.Fatalf("Move(2,4,1) completing row 2 on a 5x5 board = %d, want 1", got)
	}
}

func TestTicTacToe_MinimumBoardSizeTwo(t *testing.T) {
	g := NewTicTacToe(2)
	g.Move(0, 0, 1)
	g.Move(1, 1, 2)
	if got := g.Move(0, 1, 1); got != 1 {
		t.Fatalf("Move(0,1,1) completing row 0 on a 2x2 board = %d, want 1", got)
	}
}

func TestTicTacToe_PlayerTwoWinsColumn(t *testing.T) {
	g := NewTicTacToe(4)
	// Player 1 fills unrelated cells; player 2 fills column 3.
	g.Move(0, 3, 2)
	g.Move(0, 0, 1)
	g.Move(1, 3, 2)
	g.Move(1, 0, 1)
	g.Move(2, 3, 2)
	g.Move(2, 0, 1)
	if got := g.Move(3, 3, 2); got != 2 {
		t.Fatalf("Move(3,3,2) completing column 3 = %d, want 2", got)
	}
}
