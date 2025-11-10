
# snake.go

### Terminal Snake (Go) — Linux

A minimal, fast Snake game that runs directly in your terminal. Built with Go and ANSI escape codes for smooth rendering and responsive controls. It also includes a simple “solver” mode that demonstrates an automated pathing strategy.

#### Why this exists
- Learn-by-doing project for terminal rendering and input handling in raw mode.

---

### Features
- Playable in your terminal with WASD/HJKL controls
- Raw-mode keyboard input (no Enter required)
- Wrapping at edges (go off one side, appear on the other)
- Adjustable board size, tick speed, and characters
- Optional double-width rendering for better proportions
- Auto-solver mode with a “dumb” sweep pattern

---

### Requirements
- Linux
- Go 1.21+ (uses `math/rand/v2` and `golang.org/x/term`)
- A terminal that supports ANSI escape codes and UTF-8

---

### Install
```bash
# Clone
git clone https://github.com/gabriGutiz/snake.go.git
cd snake.go

# Build
go build -o snake

# Or run directly
go run .
```

---

### Usage

The program exposes two subcommands: `play` and `solve`.

General syntax:
```bash
./snake <subcommand> [options]
```

#### 1) Play mode
Interactive mode with keyboard controls.

```bash
./snake play [--height N] [--width N] [--tick MS] [--l-off N] \
  [--space-char CHAR] [--snake-char CHAR] [--food-char CHAR] [--double-char-w]
```

Options:
- `--height` (int, default 10): total height of the grid
- `--width` (int, default 10): total width of the grid
- `--tick` (ms, default 100): update interval in milliseconds
- `--l-off` (int, default 10): left padding (number of spaces before the grid)
- `--space-char` (string, default `░`): background cell character
- `--snake-char` (string, default `█`): snake character
- `--food-char` (string, default `█`): food character
- `--double-char-w` (flag): render cells as two characters wide (improves visual proportions)

Controls:
- Up: `W` or `K`
- Left: `A` or `H`
- Down: `S` or `J`
- Right: `D` or `L`
- Quit: `Q` or `Ctrl+C`

Examples:
```bash
# Default 10x10 at 100ms tick
./snake play

# Larger board with faster tick
./snake play --height 20 --width 30 --tick 60

# Use ASCII chars and no left offset
./snake play --space-char . --snake-char O --food-char X --l-off 0

# Wider cells for better aspect ratio
./snake play --double-char-w --height 18 --width 32
```

#### 2) Solve mode
Automated movement driven by a simple built-in strategy.

```bash
./snake solve [--height N] [--width N] [--tick MS] [--l-off N] \
  [--space-char CHAR] [--snake-char CHAR] [--food-char CHAR] \
  [--double-char-w] [--mode NAME]
```

Options:
- All options from `play` are supported.
- `--mode` (string, default `dumb`): select the solver. Currently supported:
  - `dumb`: a simple sweep-like strategy.

Controls:
- Quit: `Q` or `Ctrl+C`

Examples:
```bash
# Run solver on a 24x40 board
./snake solve --height 24 --width 40 --tick 50 --mode dumb

# Solver with ASCII and double width
./snake solve --space-char . --snake-char o --food-char * --double-char-w
```

---

### Known Behaviors and TODOs
- Fast opposite-direction inputs in the same tick can still feel “snappy”; intentional blocking keeps the snake from reversing into itself.
- Solver “dumb” mode works on even heights. There’s a TODO to refine odd-height behavior.
- Cursor hiding is TODO’d.
- Terminal is put into raw mode for the duration; it restores on exit.

---

### Development

Project structure:
- `main.go`: CLI, input handling (raw mode), ticker loop, and rendering bootstrap
- `pkg/snake/snake.go`: game state, logic, solver, and drawing functions
- `pkg/utils/utils.go`: small asserts helper
- `go.mod`, `go.sum`: dependencies (`golang.org/x/term`)

Run from source:
```bash
go run .
go run . play --height 20 --width 30
go run . solve --height 20 --width 30 --mode dumb
```

---

### Demo

TODO

---

### Troubleshooting
- Seeing garbled output or no colors/blocks? Make sure your terminal font supports UTF-8 block characters or switch to ASCII via `--space-char`, `--snake-char`, `--food-char`.
- Terminal not restoring properly after a crash? It should restore raw mode on exit; if it doesn’t, run `reset` or close and reopen the terminal.
- Input lag: try a larger tick (e.g., `--tick 120`) or a smaller grid.

---

### License
This repository includes a LICENSE file. See [LICENSE](./LICENSE) for details.

---

### Acknowledgments
- Uses `golang.org/x/term` for raw mode input
- ANSI escape sequences for terminal control and rendering

