package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gabriGutiz/snake.go/pkg/snake"
	"golang.org/x/term"
)

func runPlay(s snake.Snake, ticker time.Ticker, keyPressChan chan byte) error {
	err := s.Start(false)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			err := s.Tick()
			if err != nil {
				fmt.Print("\033[2J\033[H")
				return err
			}
		case key := <-keyPressChan:
			switch key {
			case 'k', 'w', 'W':
				s.SetDirection(snake.Up)
			case 'j', 's', 'S':
				s.SetDirection(snake.Down)
			case 'h', 'a', 'A':
				s.SetDirection(snake.Left)
			case 'l', 'd', 'D':
				s.SetDirection(snake.Right)
			case 'q', 'Q', 3: // 3 = Ctrl+C
				return nil
			}
		}
	}
}

func runSolve(s snake.Snake, ticker time.Ticker, keyPressChan chan byte) error {
	err := s.Start(true)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			err := s.Tick()
			if err != nil {
				fmt.Print("\033[2J\033[H")
				return err
			}
		case key := <-keyPressChan:
			switch key {
			case 'q', 'Q', 3: // 3 = Ctrl+C
				return nil
			}
		}
	}
}

func start(tick, height, width, leftOffset int, doubleCharWidth bool, snakeChar, spaceChar, foodChar rune, solve bool) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error setting raw mode:", err)
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	ticker := time.NewTicker(time.Duration(tick) * time.Millisecond)
	defer ticker.Stop()

	writer := bufio.NewWriter(os.Stdout)
	// TODO: understand if makes sense to hide cursor
	// and how to show again when program ends
	// writer.WriteString("\033[?25l")
	// defer writer.WriteString("\033[?25h")

	keyPressChan := make(chan byte)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			char, err := reader.ReadByte()
			if err != nil {
				panic("error reading key stroke")
			}
			keyPressChan <- char
		}
	}()

	s := snake.NewSnake(*writer, height, width, leftOffset,
		doubleCharWidth, snakeChar, spaceChar, foodChar)

	if solve {
		err = runSolve(s, *ticker, keyPressChan)
	} else {
		err = runPlay(s, *ticker, keyPressChan)
	}

	return err
}

var playCmd = flag.NewFlagSet("play", flag.ExitOnError)
var pHeightFlag = playCmd.Int("height", 10, "total height of the snake area")
var pWidthFlag = playCmd.Int("width", 10, "total width of the snake area")
var pTickFlag = playCmd.Int("tick", 100, "tick to update the snake (in ms)")
var pLeftOffsetFlag = playCmd.Int("l-off", 10, "an offset of whitespaces on the left")
var pSpaceCharFlag = playCmd.String("space-char", "░", "char used to make the spacing")
var pSnakeCharFlag = playCmd.String("snake-char", "█", "char used to make the snake")
var pFoodCharFlag = playCmd.String("food-char", "█", "char used to make the foods")
var pDCharWidthFlag = playCmd.Bool("double-char-w", true, "use two chars on horizontal")

var solveCmd = flag.NewFlagSet("solve", flag.ExitOnError)
var sHeightFlag = solveCmd.Int("height", 10, "total height of the snake area")
var sWidthFlag = solveCmd.Int("width", 10, "total width of the snake area")
var sTickFlag = solveCmd.Int("tick", 100, "tick to update the snake (in ms)")
var sLeftOffsetFlag = solveCmd.Int("l-off", 10, "an offset of whitespaces on the left")
var sSpaceCharFlag = solveCmd.String("space-char", "░", "char used to make the spacing")
var sSnakeCharFlag = solveCmd.String("snake-char", "█", "char used to make the snake")
var sFoodCharFlag = solveCmd.String("food-char", "█", "char used to make the foods")
var sDCharWidthFlag = solveCmd.Bool("double-char-w", true, "use two chars on horizontal")

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'play' or 'solve' subcommands")
		os.Exit(1)
	}

	var err error = nil

	// TODO: the program runs when giving play and solve commands
	switch os.Args[1] {
	case "play":
		playCmd.Parse(os.Args[2:])
		err = start(*pTickFlag, *pHeightFlag, *pWidthFlag, *pLeftOffsetFlag,
			*pDCharWidthFlag, []rune(*pSnakeCharFlag)[0], []rune(*pSpaceCharFlag)[0],
			[]rune(*pFoodCharFlag)[0], false)
	case "solve":
		solveCmd.Parse(os.Args[2:])
		err = start(*sTickFlag, *sHeightFlag, *sWidthFlag, *sLeftOffsetFlag,
			*sDCharWidthFlag, []rune(*sSnakeCharFlag)[0], []rune(*sSpaceCharFlag)[0],
			[]rune(*sFoodCharFlag)[0], true)
	default:
		fmt.Println("expected 'play' or 'solve' subcommands")
		return
	}

	if err != nil {
		fmt.Println("An error happened: ", err)
	}
}
