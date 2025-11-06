package snake

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand/v2"
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
	foodPos              Position
	solving              bool
	running              bool
}

func NewSnake(writer bufio.Writer, height, width, leftOffset int, doubleCharWidth bool, snakeChar, spaceChar rune) Snake {
	var body []Position
	body = append(body, Position{X: 0, Y: 0})

	realWidth := width
	widthCharMulti := 1
	if doubleCharWidth {
		widthCharMulti = 2
		realWidth = width / widthCharMulti
	}

	s := Snake{
		writer:         writer,
		height:         height,
		width:          realWidth,
		snakeChar:      snakeChar,
		spaceChar:      spaceChar,
		snakeSize:      1,
		lOffset:        leftOffset,
		widthCharMulti: widthCharMulti,
		snakeBody:      body,
		direction:      Down,
		solving:        false,
		running:        false,
	}
	s.createFood()

	return s
}

func (s *Snake) SetDirection(dir Direction) error {
	if s.solving {
		return errors.New("Set direction is unavailable when solving")
	}

	// TODO: when moving too fast, I die
	// Ex.: going down and doing d+w
	// The set direction is chaging the direction bypassing this validation
	// Tick did not happen when the second hey is pressed
	if (s.direction == Up && dir == Down) ||
		(s.direction == Down && dir == Up) ||
		(s.direction == Left && dir == Right) ||
		(s.direction == Right && dir == Left) {
		return nil
	}
	s.direction = dir
	return nil
}

func (s *Snake) Start(solve bool) error {
	if s.running {
		return errors.New("Game is already running")
	}

	s.writer.WriteString("\033[2J\033[H")
	s.solving = solve
	s.running = true
	s.printSnake()
	return nil
}

func (s *Snake) Tick() error {
	if !s.running {
		return errors.New("Game has not been started")
	}

	oldTail := s.snakeBody[len(s.snakeBody)-1]
	s.move()
	newHead := s.snakeBody[0]

	s.updatePosition(oldTail, s.spaceChar)
	s.updatePosition(newHead, s.snakeChar)

	if newHead == s.foodPos {
		s.snakeBody = append(s.snakeBody, oldTail)
		s.snakeSize++
		s.createFood()
	}
	s.updatePosition(s.foodPos, s.snakeChar)

	if s.collidedWithBody() {
		return errors.New("Game Over")
	}

	e := s.writer.Flush()

	return e
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

		s.writer.WriteString(fmt.Sprintf("%s%s\r\n", leftPadding, strings.Join(row, "")))
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
		newPosition.Y = s.height - 1
	} else if newPosition.X == s.width {
		newPosition.X = 0
	} else if newPosition.X < 0 {
		newPosition.X = s.width - 1
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

func (s *Snake) createFood() {
	foodX := rand.IntN(s.width / s.widthCharMulti)
	foodY := rand.IntN(s.height)

	for _, v := range s.snakeBody {
		if v.X == foodX && v.Y == foodY {
			s.createFood()
			return
		}
	}

	s.foodPos.X = foodX
	s.foodPos.Y = foodY
}

func (s *Snake) collidedWithBody() bool {
	head := s.snakeBody[0]

	for i := 1; i < len(s.snakeBody); i++ {
		if head == s.snakeBody[i] {
			return true
		}
	}

	return false
}
