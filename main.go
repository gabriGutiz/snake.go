package main

import (
	"bufio"
	"flag"
	"os"
	"time"

	"github.com/gabriGutiz/snake.go/pkg/snake"
)

var heightFlag = flag.Int("height", 10, "total height of the snake area")
var widthFlag = flag.Int("width", 10, "total width of the snake area")
var sizeFlag = flag.Int("size", 10, "size of the snake")
var tickFlag = flag.Int("tick", 100, "tick to update the snake (in ms)")
var leftOffsetFlag = flag.Int("l-off", 10, "an offset of whitespaces on the left")
var spaceCharFlag = flag.String("space-char", "░", "char used to make the spacing")
var snakeCharFlag = flag.String("snake-char", "█", "char used to make the snake")
var dCharWidthFlag = flag.Bool("double-char-w", true, "use two chars on horizontal")

func main() {
	flag.Parse()

	ticker := time.NewTicker(time.Duration(*tickFlag) * time.Millisecond)
	defer ticker.Stop()

	writer := bufio.NewWriter(os.Stdout)
	// TODO: understand if makes sense to hide cursor
	// and how to show again when program ends
	// writer.WriteString("\033[?25l")
	// defer writer.WriteString("\033[?25h")

	snake := snake.NewSnake(*writer, *heightFlag, *widthFlag, *sizeFlag, *leftOffsetFlag,
		*dCharWidthFlag, []rune(*snakeCharFlag)[0], []rune(*spaceCharFlag)[0])

	snake.Start()

	for range ticker.C {
		snake.Tick()
	}
}
