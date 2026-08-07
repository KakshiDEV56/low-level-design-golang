package snakegame

/*
LeetCode 353. Design Snake Game

Design a Snake game that is played on a device with screen size height x width. The
snake is initially positioned at the top left corner (0, 0) with a length of 1 unit.

You are given an array food where food[i] = (ri, ci) is the row and column position of a
piece of food that the snake can eat. When a snake eats a piece of food, its length and
the game's score both increase by 1. Each piece of food appears one by one on the screen,
meaning the (i+1)th piece of food will not appear until the snake eats the ith piece of
food.

When food does not appear on the screen, the position where it will appear is out of
bounds. The game is over if the snake goes out of bounds (moves off the screen) or if its
head occupies a space that its body occupies after moving (i.e. a snake of length 4 cannot
run into itself). Note that eating food does not count as "running into itself" for the
tail — the tail moves out of that cell in the same step the head enters it, UNLESS the
snake just grew (ate food), in which case the tail stays put for that move.

Implement the SnakeGame class:
  - NewSnakeGame(width, height int, food [][]int) initializes the object with a screen of
    size height x width and the positions of the food.
  - Move(direction string) int returns the game's score after applying one direction move
    by the snake. If the game is over, return -1. direction is one of "U", "D", "L", "R"
    (up, down, left, right).

Example 1:

	Input:
	["SnakeGame", "move", "move", "move", "move", "move", "move"]
	[[3, 2, [[1, 2], [0, 1]]], ["R"], ["D"], ["R"], ["U"], ["L"], ["U"]]
	Output:
	[null, 0, 0, 1, 1, 2, -1]

	Explanation:
	SnakeGame sg = new SnakeGame(3, 2, [[1, 2], [0, 1]]);
	// board size: 2 rows x 3 cols, food at (1,2) then (0,1)
	sg.move("R");  // snake: [(0,1)], head moves right, no food there yet -> score 0
	sg.move("D");  // snake: [(1,1)], score 0
	sg.move("R");  // snake: [(1,2),(1,1)], head eats food at (1,2) -> score 1
	sg.move("U");  // snake: [(0,2),(1,2)], score 1
	sg.move("L");  // snake: [(0,1),(0,2)], head eats food at (0,1) -> score 2
	sg.move("U");  // head would move to (-1,1), out of bounds -> game over, -1

Constraints:
  - 1 <= width, height <= 25
  - 1 <= food.length <= 50
  - food[i].length == 2
  - 0 <= ri < height
  - 0 <= ci < width
  - direction.length == 1
  - direction is 'U', 'D', 'L', or 'R'.
  - At most 10^4 calls will be made to Move.

TODO:
 1. Add the fields you need on SnakeGame below. Hint: you'll want
      - width, height, and the food slice plus an index into it (which food is "next"),
      - the snake's body as an ordered sequence of (row, col) cells, head first — a slice
        used as a deque works well since you push a new head and (usually) pop the tail
        each move, and
      - a set (map) of the cells currently occupied by the body, for O(1) self-collision
        checks — keyed by e.g. row*width+col or a small struct.
      - a "game over" flag, since once the game ends it should keep returning -1.
 2. Implement NewSnakeGame: snake starts as a single segment at (0,0); board bounds and
    food list are stored as given.
 3. Implement Move(direction):
      - If the game is already over, return -1 immediately.
      - Compute the new head position from the current head + direction.
      - If the new head is out of bounds, mark game over and return -1.
      - Determine if this move eats the next food item (new head == next food position).
        If so, don't remove the tail this step (the snake grows) and advance the food
        index; otherwise remove the tail from both the body slice and the occupied set
        BEFORE checking self-collision (the tail cell vacates as the head moves in).
      - If the new head position collides with the (possibly just-updated) occupied set,
        mark game over and return -1.
      - Otherwise push the new head onto the body and occupied set, and return the
        current score (snake length - 1, or a running counter you increment on each
        food eaten — either works, just be consistent).
*/

type SnakeGame struct {
	// TODO: add your fields here
}

func NewSnakeGame(width int, height int, food [][]int) *SnakeGame {
	// TODO: implement
	return &SnakeGame{}
}

func (g *SnakeGame) Move(direction string) int {
	// TODO: implement
	return -1
}
