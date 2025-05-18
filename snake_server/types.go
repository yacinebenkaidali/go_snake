package main

import (
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Direction int

const (
	boardWidth  = 70
	boardHeight = 20
	// Colors/characters for rendering
	snakeChar = 'O'
	headChar  = '@'
	foodChar  = '●'
	emptyChar = ' '
	// game options
	frameTIme = time.Millisecond * 200
	// initial snake size
	snakeSize = 5
)

const (
	UP = iota
	RIGHT
	LEFT
	DOWN
)

type Position struct {
	X, Y int
}

type Game struct {
	snake         []Position
	s             tcell.Screen
	quiCh         chan struct{}
	speed         time.Duration
	food          Position
	score         int
	currDirection Direction
	gameOver      bool
	mu            sync.Mutex
	paused        bool
	player        *Player
	scoreCh       chan int
}

var snakeStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorMediumVioletRed)
var endGameMessageStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)

var gameOverSound = "assets/game_over_beep.mp3"
var gameStartSound = "assets/game_start_beep.mp3"
