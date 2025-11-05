package snake

import (
	"bufio"
	"fmt"
	"slices"
	"strings"

	"github.com/gabriGutiz/snake.go/pkg/utils"
)

type Direction int

const (
	Left Direction = iota
	Right
	Down
	Up
)

type Position struct {
	X, Y int
}

type Snake struct {
	writer               bufio.Writer
	height, width        int
	snakeChar, spaceChar rune
	snakeSize            int
	lOffset              int
	widthCharMulti       int
	snakeBody            []Position
	direction            Direction
}

func NewSnake(writer bufio.Writer, height, width, size, leftOffset int, doubleCharWidth bool, snakeChar, spaceChar rune) Snake {
	body := make([]Position, size)

	for i := range size {
		body[i] = Position{X: i, Y: 0}
	}

	realWidth := width
	widthCharMulti := 1
	if doubleCharWidth {
		widthCharMulti = 2
		realWidth = width / widthCharMulti
	}

	return Snake{
		writer:         writer,
		height:         height,
		width:          realWidth,
		snakeChar:      snakeChar,
		spaceChar:      spaceChar,
		snakeSize:      size,
		lOffset:        leftOffset,
		widthCharMulti: widthCharMulti,
		snakeBody:      body,
		direction:      Down,
	}
}

func (s *Snake) printSnake() {
	leftPadding := strings.Repeat(" ", s.lOffset)
	for y := range s.height {
		row := make([]string, s.width)
		for x := range s.width {
			if slices.Contains(s.snakeBody, Position{X: x, Y: y}) {
				row[x] = strings.Repeat(string(s.snakeChar), s.widthCharMulti)
			} else {
				row[x] = strings.Repeat(string(s.spaceChar), s.widthCharMulti)
			}
		}

		s.writer.WriteString(fmt.Sprintf("%s%s\n", leftPadding, strings.Join(row, "")))
	}
}

func (s *Snake) move() {
	var newBody []Position
	newPosition := s.snakeBody[0]

	if s.direction == Down {
		newPosition.Y++
	} else if s.direction == Up {
		newPosition.Y--
	} else if s.direction == Right {
		newPosition.X++
	} else if s.direction == Left {
		newPosition.X--
	}

	if newPosition.Y == s.height {
		newPosition.Y = 0
	} else if newPosition.Y < 0 {
		newPosition.Y = s.height
	} else if newPosition.X == s.width {
		newPosition.X = 0
	} else if newPosition.X < 0 {
		newPosition.X = s.width
	}

	utils.Assert(len(s.snakeBody) == s.snakeSize, "Size of the snake should be equals to size of body")

	newBody = append(newBody, newPosition)
	newBody = append(newBody, s.snakeBody[:s.snakeSize-1]...)

	s.snakeBody = newBody
}

func (s *Snake) smartMove() {
	var newBody []Position
	newPosition := s.snakeBody[0]

	if s.direction == Down {
		newPosition.Y++
	} else if s.direction == Up {
		newPosition.Y--
	} else if s.direction == Right {
		newPosition.X++
	} else if s.direction == Left {
		newPosition.X--
	}

	if newPosition.Y == s.height {
		newPosition.Y = 0
	} else if newPosition.Y < 0 {
		newPosition.Y = s.height
	} else if newPosition.X == s.width {
		newPosition.X = 0
	} else if newPosition.X < 0 {
		newPosition.X = s.width
	}

	if newPosition.Y == s.height-1 && s.direction == Down {
		s.direction = Right
	} else if newPosition.Y == 0 && s.direction == Up {
		s.direction = Left
	} else if newPosition.X == s.width-1 && s.direction == Right {
		s.direction = Up
	} else if newPosition.X == 0 && s.direction == Left {
		s.direction = Down
	}

	utils.Assert(len(s.snakeBody) == s.snakeSize, "Size of the snake should be equals to size of body")

	newBody = append(newBody, newPosition)
	newBody = append(newBody, s.snakeBody[:s.snakeSize-1]...)

	s.snakeBody = newBody
}

func (s *Snake) updatePosition(pos Position, char rune) {
	col := s.lOffset + (pos.X * s.widthCharMulti) + 1
	s.writer.WriteString(fmt.Sprintf("\033[%d;%dH%s\033[u", pos.Y+1, col, strings.Repeat(string(char), s.widthCharMulti)))
}

func (s *Snake) Start() {
	s.writer.WriteString("\033[2J\033[H")
	s.printSnake()
}

func (s *Snake) Tick() {
	oldTail := s.snakeBody[len(s.snakeBody)-1]
	s.smartMove()
	newHead := s.snakeBody[0]

	s.updatePosition(oldTail, s.spaceChar)
	s.updatePosition(newHead, s.snakeChar)

	s.writer.Flush()
}
