package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

func NewGame() (*Game, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	// Set default style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	s.Clear()

	g := &Game{
		s:             s,
		quiCh:         make(chan struct{}),
		speed:         frameTIme,
		currDirection: RIGHT,
		mu:            sync.Mutex{},
		scoreCh:       make(chan int),
	}
	g.GenerateFood()
	g.InitializeSnake()

	return g, nil
}

func (g *Game) resetGame() {
	g.paused = false
	g.score = 0
	g.gameOver = false
	g.currDirection = RIGHT
	g.InitializeSnake()
}

func (g *Game) GenerateFood() {
	// generate food away from snake
	var x, y int

	for {
		foodHasCollisionWithSnake := false
		x = rand.Intn(boardWidth-2) + 1
		y = rand.Intn(boardHeight-2) + 1

		for _, part := range g.snake {
			if part.X == x && part.Y == y {
				foodHasCollisionWithSnake = true
			}
		}
		if !foodHasCollisionWithSnake {
			g.food = Position{X: x, Y: y}
			break
		}
	}
}

func (g *Game) InitializeSnake() {
	x := boardWidth / 6
	y := boardHeight / 4

	snake := make([]Position, snakeSize)
	for i := 0; i < snakeSize; i++ {
		snake[i] = Position{X: x + i, Y: y}
	}

	g.snake = snake
}

func (s *ConnectionManager) DrawSnakeAndFood() {
	// compute next position
	g := s.g
	if !g.paused {
		newHead := Position{X: g.snake[len(g.snake)-1].X, Y: g.snake[len(g.snake)-1].Y}
		switch g.currDirection {
		case UP:
			newHead.Y = newHead.Y - 1
		case RIGHT:
			newHead.X = newHead.X + 1
		case DOWN:
			newHead.Y = newHead.Y + 1
		case LEFT:
			newHead.X = newHead.X - 1
		}

		// check collision with walls
		if newHead.X == 0 || newHead.X > boardWidth || newHead.Y == 0 || newHead.Y > boardHeight {

			g.gameOver = true

			g.player.PlaySound(gameOverSound, true)
			return
		}

		// check collision with self
		for _, part := range g.snake {
			if newHead.X == part.X && newHead.Y == part.Y {
				g.gameOver = true
				g.player.PlaySound(gameOverSound, true)
				return
			}
		}
		// check collision with food
		if newHead.X == g.food.X && newHead.Y == g.food.Y {
			g.snake = append(g.snake, newHead)
			g.score++
			// send score to client
			g.scoreCh <- g.score
			g.WriteAt(g.score)
			g.GenerateFood()
		} else {
			g.snake = append(g.snake[1:], newHead)
		}
	}
	// draw food
	g.s.SetContent(g.food.X, g.food.Y, foodChar, nil, tcell.StyleDefault)

	// draw the snake
	for i := 0; i < len(g.snake)-1; i++ {
		g.s.SetContent(g.snake[i].X, g.snake[i].Y, snakeChar, nil, snakeStyle)
	}
	lastPosition := len(g.snake) - 1
	g.s.SetContent(g.snake[lastPosition].X, g.snake[lastPosition].Y, headChar, nil, snakeStyle)
}

func (g *Game) DrawMessage() {
	var gameOverText = []string{
		" +-----------------+ ",
		" |    GAME OVER    | ",
		" +-----------------+ ",
	}
	// Calculate starting position to center the text
	startY := boardHeight/2 - len(gameOverText)/2

	for i, line := range gameOverText {
		startX := boardWidth/2 - len(line)/2
		for j, char := range line {
			g.s.SetContent(startX+j, startY+i, rune(char), nil, endGameMessageStyle)
		}
	}
}

func (g *Game) DrawBoard() {
	wallStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)

	// Draw top anD bottom lines
	for x := 1; x < boardWidth+1; x++ {
		g.s.SetContent(x, 0, '═', nil, wallStyle)
		g.s.SetContent(x, boardHeight+1, '═', nil, wallStyle)
	}

	// Draw left and right borders
	for y := 1; y < boardHeight+1; y++ {
		g.s.SetContent(0, y, '║', nil, wallStyle)
		g.s.SetContent(boardWidth+1, y, '║', nil, wallStyle)
	}
	// Draw corners
	g.s.SetContent(0, 0, '╔', nil, wallStyle)
	g.s.SetContent(boardWidth+1, boardHeight+1, '╝', nil, wallStyle)
	g.s.SetContent(boardWidth+1, 0, '╗', nil, wallStyle)
	g.s.SetContent(0, boardHeight+1, '╚', nil, wallStyle)
}

func (g *Game) DrawInfo() {
	head := fmt.Sprintf("Game paused %+v", g.paused)
	for i := range len(head) {
		g.s.SetContent(i, boardHeight+3, rune(head[i]), nil, tcell.StyleDefault)
	}

	score := fmt.Sprintf("Score %d, Game status %+v", g.score, g.gameOver)
	for i := range len(score) {
		g.s.SetContent(i, boardHeight+4, rune(score[i]), nil, tcell.StyleDefault)
	}
}
func (g *Game) WriteAt(y int) {
	msg := fmt.Sprintf("Write score %d", g.score)
	for i := range len(msg) {
		g.s.SetContent(i, boardHeight+y, rune(msg[i]), nil, tcell.StyleDefault)
	}
}

func (g *Game) HandleInput(ev *tcell.EventKey) {
	if g.gameOver {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	// Update direction based on key press
	switch ev.Key() {
	case tcell.KeyUp:
		if g.currDirection != DOWN {
			g.currDirection = UP
		}
	case tcell.KeyRight:
		if g.currDirection != LEFT {
			g.currDirection = RIGHT
		}
	case tcell.KeyDown:
		if g.currDirection != UP {
			g.currDirection = DOWN
		}
	case tcell.KeyLeft:
		if g.currDirection != RIGHT {
			g.currDirection = LEFT
		}
	}

	// Also support WASD controls
	switch ev.Rune() {
	case 'w', 'W':
		if g.currDirection != DOWN {
			g.currDirection = UP
		}
	case 'd', 'D':
		if g.currDirection != LEFT {
			g.currDirection = RIGHT
		}
	case 's', 'S':
		if g.currDirection != UP {
			g.currDirection = DOWN
		}
	case 'a', 'A':
		if g.currDirection != RIGHT {
			g.currDirection = LEFT
		}
	}
}

func (cm *ConnectionManager) RunGame() {
	g := cm.g
	player := NewPlayer([]string{gameOverSound, gameStartSound})
	g.player = player
	defer g.player.Close()

	// go g.gameControlEvents()

	ticker := time.NewTicker(g.speed)
	defer ticker.Stop()
	g.player.PlaySound(gameStartSound, false)
	for {
		select {
		case <-g.quiCh:
			{
				g.s.Fini()
				return
			}
		case <-ticker.C:
			{
				g.s.Clear()
				g.DrawBoard()
				if g.gameOver {
					// draw ascii game over
					g.DrawMessage()
				} else {
					cm.DrawSnakeAndFood()
				}
				g.DrawInfo()
				g.s.Show()
			}
		}
	}
}

var debugLogger *log.Logger

func main() {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		os.Exit(1)
	}

	// Initialize the logger
	debugLogger = log.New(logFile, "DEBUG: ", log.Ltime|log.Lshortfile)

	server, err := NewServer("", "3000")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	go server.sendScore()
	go server.StartServer()

	server.RunGame()
}
