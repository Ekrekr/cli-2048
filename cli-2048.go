package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

const (
	TOP    = 1
	UP     = 1
	RIGHT  = 2
	BOTTOM = 3
	DOWN   = 3
	LEFT   = 4
)

var testGrid = [][]int{
	{0, 4, 2, 0},
	{2, 0, 0, 0},
	{0, 128, 64, 0},
	{0, 128, 2048, 32768}}

type game struct {
	currentScore      int
	highScore         int
	grid              [][]int
	newGrid           [][]int
	isGameOver        bool
	needScreenRefresh bool
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var g = game{
		currentScore:      2048,
		highScore:         84096,
		grid:              testGrid,
		isGameOver:        false,
		needScreenRefresh: true}

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	var play = true
	for play {
		if g.needScreenRefresh {
			clearScreen()
			g.print()
			g.needScreenRefresh = false
		}
		_, key, err := keyboard.GetKey()
		if err != nil {
			// TODO: Shift arrow keys will panic here. Ignore, or not?
			// panic(err)
		}
		if g.isGameOver == true {
			play = false
		}
		switch key {
		case keyboard.KeyArrowUp:
			if g.canMove(UP) {
				g.move(UP)
				g.needScreenRefresh = true
			}
		case keyboard.KeyArrowDown:
			if g.canMove(DOWN) {
				g.move(DOWN)
				g.needScreenRefresh = true
			}
		case keyboard.KeyArrowLeft:
			if g.canMove(LEFT) {
				g.move(LEFT)
				g.needScreenRefresh = true
			}
		case keyboard.KeyArrowRight:
			if g.canMove(RIGHT) {
				g.move(RIGHT)
				g.needScreenRefresh = true
			}
		case keyboard.KeyCtrlN:
			g.newGame()
			g.needScreenRefresh = true
		case keyboard.KeyCtrlQ:
		case keyboard.KeyCtrlC:
			play = false
		}
		g.isGameOver = g.checkIsGameOver()
	}
}

func (g *game) canMove(direction int) bool {
	for y := 0; y < g.gridHeight(); y++ {
		for x := 0; x < g.gridWidth(); x++ {
			if direction == UP && y > 0 &&
				(g.grid[y][x] == g.grid[y-1][x] || g.grid[y-1][x] == 0) {
				return true
			}
			if direction == DOWN && y < g.gridHeight()-1 &&
				(g.grid[y][x] == g.grid[y+1][x] || g.grid[y+1][x] == 0) {
				return true
			}
			if direction == LEFT && x > 0 &&
				(g.grid[y][x] == g.grid[y][x-1] || g.grid[y][x-1] == 0) {
				return true
			}
			if direction == RIGHT && x < g.gridWidth()-1 &&
				(g.grid[y][x] == g.grid[y][x+1] || g.grid[y][x+1] == 0) {
				return true
			}
		}
	}
	return false
}

func (g *game) checkIsGameOver() bool {
	return !g.canMove(UP) && !g.canMove(DOWN) && !g.canMove(LEFT) && !g.canMove(RIGHT)
}

func (g *game) move(direction int) {
	switch direction {
	case DOWN:
		for y := g.gridHeight() - 1; y > 0; y-- {
			for x := 0; x < g.gridWidth(); x++ {
				g.combineTilesIfEqual(y, x, y-1, x)
				g.moveTileIfSpaceEmpty(y, x, y-1, x)
			}
		}
		g.spawnNewTileIfNeeded(TOP)
	case UP:
		for y := 0; y < g.gridHeight()-1; y++ {
			for x := 0; x < g.gridWidth(); x++ {
				g.combineTilesIfEqual(y, x, y+1, x)
				g.moveTileIfSpaceEmpty(y, x, y+1, x)
			}
		}
		g.spawnNewTileIfNeeded(BOTTOM)
	case RIGHT:
		for y := 0; y < g.gridHeight(); y++ {
			for x := g.gridWidth() - 1; x > 0; x-- {
				g.combineTilesIfEqual(y, x, y, x-1)
				g.moveTileIfSpaceEmpty(y, x, y, x-1)
			}
		}
		g.spawnNewTileIfNeeded(LEFT)
	case LEFT:
		for y := 0; y < g.gridHeight(); y++ {
			for x := 0; x < g.gridWidth()-1; x++ {
				g.combineTilesIfEqual(y, x, y, x+1)
				g.moveTileIfSpaceEmpty(y, x, y, x+1)
			}
		}
		g.spawnNewTileIfNeeded(RIGHT)
	}
}

func (g *game) combineTilesIfEqual(yTo int, xTo int, yFrom int, xFrom int) {
	if g.grid[yTo][xTo] == g.grid[yFrom][xFrom] {
		g.grid[yTo][xTo] *= 2
		g.grid[yFrom][xFrom] = 0
	}
}

func (g *game) moveTileIfSpaceEmpty(yTo int, xTo int, yFrom int, xFrom int) {
	if g.grid[yTo][xTo] == 0 {
		g.grid[yTo][xTo] = g.grid[yFrom][xFrom]
		g.grid[yFrom][xFrom] = 0
	}
}

func (g *game) spawnNewTileIfNeeded(side int) {
	var emptySpaces []int
	var spawnIndexLimit = g.gridWidth()
	if side == LEFT || side == RIGHT {
		spawnIndexLimit = g.gridHeight()
	}
	for i := 0; i < spawnIndexLimit; i++ {
		if side == TOP && g.grid[0][i] == 0 {
			emptySpaces = append(emptySpaces, i)
		}
		if side == BOTTOM && g.grid[g.gridHeight()-1][i] == 0 {
			emptySpaces = append(emptySpaces, i)
		}
		if side == LEFT && g.grid[i][0] == 0 {
			emptySpaces = append(emptySpaces, i)
		}
		if side == RIGHT && g.grid[i][g.gridWidth()-1] == 0 {
			emptySpaces = append(emptySpaces, i)
		}
	}

	if len(emptySpaces) > 0 {
		var newTileValue = 2
		if rand.Float64() > 0.9 {
			newTileValue = 4
		}

		var randomIndex = int32(rand.Float64() * float64(len(emptySpaces)))
		var newTilePlace = emptySpaces[randomIndex]
		switch side {
		case TOP:
			g.grid[0][newTilePlace] = newTileValue
		case BOTTOM:
			g.grid[g.gridHeight()-1][newTilePlace] = newTileValue
		case LEFT:
			g.grid[newTilePlace][0] = newTileValue
		case RIGHT:
			g.grid[newTilePlace][g.gridWidth()-1] = newTileValue
		}
	}
}

func (g *game) gridHeight() int {
	return len(g.grid)
}

func (g *game) gridWidth() int {
	return len(g.grid[0])
}

func (g *game) newGame() {
	g.grid = make([][]int, 4)
	for i := range g.grid {
		g.grid[i] = make([]int, 4)
	}
	// Spawn a tile on a random side.
	g.move(int(rand.Float64()*4 + 1))
	g.currentScore = 0
}

func (g *game) print() {
	// TODO: Save to string then print all in one go, to improve rendering.
	fmt.Printf("\n            cli-2048          \n  ")
	for q := 0; q <= 8-len(strconv.Itoa(g.currentScore)); q++ {
		fmt.Print(" ")
	}
	fmt.Printf("🎮 %d || %d 🏆\n\n  ", g.currentScore, g.highScore)

	for _, row := range g.grid {
		for i := range []int{0, 1, 2} {
			for _, val := range row {
				var tilePrinter = getTilePrinter(val)
				if i == 1 {
					if val == 0 {
						tilePrinter("   .   ")
					} else {
						var tileNumLength = len(strconv.Itoa(val))
						var rightPadding = (5 - tileNumLength) / 2
						var isEven = 1 - tileNumLength%2
						var leftPadding = rightPadding + isEven
						for q := 0; q <= leftPadding; q++ {
							tilePrinter(" ")
						}
						tilePrinter("%d", val)
						for q := 0; q <= rightPadding; q++ {
							tilePrinter(" ")
						}
					}
				} else {
					tilePrinter("       ")
				}
			}
			fmt.Print("\n  ")
		}
	}

	if g.isGameOver {
		fmt.Print("\n  ----------------------------")
		fmt.Print("\n    >>> 💀 GAME OVER! 💀 <<<  ")
		fmt.Print("\n  ----------------------------\n")
	}

	fmt.Printf("\n   ←,↑,→,↓  💾ctrl-c 🔄ctrl-n \n\n")
}

func getTilePrinter(tile int) func(format string, a ...interface{}) (n int, err error) {
	c := color.New()
	c.Add(color.FgWhite)
	switch tile {
	case 2:
		c.Add(color.BgRed)
	case 4:
		c.Add(color.BgGreen)
	case 8:
		c.Add(color.BgYellow)
	case 16:
		c.Add(color.BgBlue)
	case 32:
		c.Add(color.BgMagenta)
	case 64:
		c.Add(color.BgCyan)
	case 128:
		c.Add(color.BgWhite)
		c.Add(color.FgBlack)
	case 256:
		c.Add(color.BgBlack)
	case 512:
		c.Add(color.BgHiRed)
	case 1024:
		c.Add(color.BgHiGreen)
	case 2048:
		c.Add(color.BgHiYellow)
	case 4096:
		c.Add(color.BgHiBlue)
	case 8192:
		c.Add(color.BgHiMagenta)
	case 16384:
		c.Add(color.BgHiCyan)
	case 32768:
		c.Add(color.BgHiWhite)
		c.Add(color.FgBlack)
	}
	return c.Printf
}

func clearScreen() {
	cmd := exec.Command("clear") // Only works for unix based systems.
	cmd.Stdout = os.Stdout
	cmd.Run()
}
