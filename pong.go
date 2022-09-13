package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

// const PaddleSymbol = '1'
const PaddleSymbol = 0x2588
const BallSymbol = 0x25CF

const PaddleHeight = 4
const VelocityRow = 1
const VelocityCol = 2

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen
var playerLeft *GameObject
var playerRight *GameObject
var ball *GameObject
var debugLog string

var gameObjects = []*GameObject{}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	initScreen()
	initGameState()
	inputChan := initUserInput()

	for {
		handleInput(readInput(inputChan))
		UpdateState()
		drawState()

		time.Sleep(75 * time.Millisecond)
	}
}

func UpdateState() {
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}

	if wallCollide(ball) {
		ball.velRow = -ball.velRow
	}

	if paddleCollide(ball, playerLeft) || paddleCollide(ball, playerRight) {
		ball.velCol = -ball.velCol
	}
}

func drawState() {
	screen.Clear()

	for _, obj := range gameObjects {
		drawObject(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}

	// drawObject(playerLeft.row, playerLeft.col, playerLeft.width, playerLeft.height, PaddleSymbol)
	// drawObject(playerRight.row, playerRight.col, playerRight.width, playerRight.height, PaddleSymbol)
	// drawObject(ball.row, ball.col, ball.width, ball.height, BallSymbol)
	screen.Show()
}

func wallCollide(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return !(obj.row+obj.velRow >= 0 && obj.row+obj.velRow < screenHeight)
}

func paddleCollide(ball, paddle *GameObject) bool {
	var collideCol bool
	if ball.col < paddle.col {
		collideCol = ball.col+ball.velCol >= paddle.col
	} else {
		collideCol = ball.col+ball.velCol <= paddle.col

	}

	return collideCol &&
		ball.row >= paddle.row &&
		ball.row < paddle.row+paddle.height
}

func gameEnd() bool {
	return getWinner() != ""
}

func getWinner() string {
	_, screenWidth := screen.Size()
	if ball.col < 0 {
		return "Player 1"
	} else if ball.col >= screenWidth {
		return "Player 2"
	} else {
		return ""
	}
}

func drawObject(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func initScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

}

func initGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	playerLeft = &GameObject{
		row: paddleStart, col: 0, width: 1, height: PaddleHeight,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}

	playerRight = &GameObject{
		row: paddleStart, col: width - 1, width: 1, height: PaddleHeight,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}

	ball = &GameObject{
		row: height / 2, col: width / 2, width: 1, height: 1,
		velRow: VelocityRow, velCol: VelocityCol,
		symbol: BallSymbol,
	}

	gameObjects = []*GameObject{
		playerLeft, playerRight, ball,
	}
}

func handleInput(key string) {
	_, screenHeight := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" && playerLeft.row > 0 {
		playerLeft.row--
	} else if key == "Rune[s]" && playerLeft.row+playerLeft.height < screenHeight {
		playerLeft.row++
	} else if key == "Up" && playerRight.row > 0 {
		playerRight.row--
	} else if key == "Down" && playerRight.row+playerRight.height < screenHeight {
		playerRight.row++
	}
}

func initUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func readInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}
	return key
}
