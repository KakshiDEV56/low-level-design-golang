package designtictactoe

/*
LeetCode 348. Design Tic-Tac-Toe

Design a Tic-Tac-Toe game that is played between two players on an n x n board.

Implement the TicTacToe class:
  - NewTicTacToe(n int) initializes the object the size of the board n.
  - Move(row, col, player int) int indicates that the player with id player plays at the
    cell (row, col) of the board. The move is guaranteed to be a valid move, and the two
    players alternate in making moves. Returns:
      - 0 if there is no winner after the move,
      - 1 if player 1 is the winner after the move, or
      - 2 if player 2 is the winner after the move.

Example 1:

	Input:
	["TicTacToe", "move", "move", "move", "move", "move", "move", "move"]
	[[3], [0, 0, 1], [0, 2, 2], [2, 2, 1], [1, 1, 2], [2, 0, 1], [1, 0, 2], [2, 1, 1]]
	Output:
	[null, 0, 0, 0, 0, 0, 0, 1]

	Explanation:
	TicTacToe ttt = new TicTacToe(3);
	Assume that player 1 is "X" and player 2 is "O" in the board.
	ttt.move(0, 0, 1); // "X  "  -- player 1 places at (0,0); no winner yet, return 0
	                   // "   "
	                   // "   "
	ttt.move(0, 2, 2); // "X  O"  -- player 2 places at (0,2); no winner, return 0
	ttt.move(2, 2, 1); // player 1 places at (2,2); return 0
	ttt.move(1, 1, 2); // player 2 places at (1,1); return 0
	ttt.move(2, 0, 1); // player 1 places at (2,0); return 0
	ttt.move(1, 0, 2); // player 2 places at (1,0); return 0
	ttt.move(2, 1, 1); // player 1 places at (2,1); this completes row 2 (2,0),(2,1),(2,2)
	                   // for player 1 -> return 1

Constraints:
  - 2 <= n <= 100
  - player is 1 or 2.
  - 0 <= row, col < n
  - (row, col) are unique for each different call to Move.
  - At most n^2 calls will be made to Move.

TODO:
 1. Add the fields you need on TicTacToe below. A naive O(n) win-check per move (re-scan
    the whole affected row/column/diagonal) is fine for n <= 100 and n^2 calls total, but
    the elegant O(1)-per-move trick is worth practicing:
      - Keep a per-row counter and per-column counter for EACH player (or use signed
        counts: +1 for player 1, -1 for player 2, so hitting +n or -n means a win), plus
        two counters for the two diagonals.
      - A cell is on the main diagonal when row == col, and on the anti-diagonal when
        row + col == n-1.
 2. Implement NewTicTacToe(n): allocate your counters sized for n rows/cols (e.g. two
    slices of length n, one per player, or one slice of signed ints), and diagonal
    counters.
 3. Implement Move(row, col, player): update the counters affected by this move (its row,
    its column, and its diagonal(s) if applicable). If any updated counter reaches n (or
    -n for player 2 under the signed scheme), that player just won — return their id.
    Otherwise return 0.
*/

type TicTacToe struct {
	// TODO: add your fields here
}

func NewTicTacToe(n int) *TicTacToe {
	// TODO: implement
	return &TicTacToe{}
}

func (g *TicTacToe) Move(row int, col int, player int) int {
	// TODO: implement
	return 0
}
