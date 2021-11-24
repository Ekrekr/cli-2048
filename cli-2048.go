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

type coord struct {
	y int
	x int
}

type game struct {
	currentScore      int
	highScore         int
	grid              [][]int
	vectorsWithMerged []bool
	isGameOver        bool
	needScreenRefresh bool
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var g = game{
		currentScore:      0,
		highScore:         0,
		grid:              [][]int{},
		isGameOver:        false,
		needScreenRefresh: true}
	g.newGame()

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
			g.isGameOver = false
		case keyboard.KeyCtrlQ:
		case keyboard.KeyCtrlC:
			// TODO: Save game.
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
		for y := g.gridHeight() - 1; y >= 0; y-- {
			for x := 0; x < g.gridWidth(); x++ {
				g.moveTileIfAble(y, x, DOWN)
			}
		}
	case UP:
		for y := 0; y < g.gridHeight(); y++ {
			for x := 0; x < g.gridWidth(); x++ {
				g.moveTileIfAble(y, x, UP)
			}
		}
	case RIGHT:
		for x := g.gridWidth() - 1; x >= 0; x-- {
			for y := 0; y < g.gridHeight(); y++ {
				g.moveTileIfAble(y, x, RIGHT)
			}
		}
	case LEFT:
		for x := 0; x < g.gridWidth(); x++ {
			for y := 0; y < g.gridHeight(); y++ {
				g.moveTileIfAble(y, x, LEFT)
			}
		}
	}
	g.spawnNewTile()
}

func (g *game) moveTileIfAble(yFrom int, xFrom int, direction int) {
	if g.grid[yFrom][xFrom] == 0 {
		return
	}
	var yTo = yFrom
	var xTo = xFrom
	var merge = false
out:
	switch direction {
	case DOWN:
		for y := yFrom + 1; y < g.gridHeight(); y++ {
			yTo = y
			if g.grid[yTo][xTo] != 0 {
				if g.grid[yFrom][xFrom] == g.grid[yTo][xTo] {
					merge = true
				} else {
					yTo = y - 1
				}
				break out
			}
		}
	case UP:
		for y := yFrom - 1; y >= 0; y-- {
			yTo = y
			if g.grid[yTo][xTo] != 0 {
				if g.grid[yFrom][xFrom] == g.grid[yTo][xTo] {
					merge = true
				} else {
					yTo = y + 1
				}
				break out
			}
		}
	case LEFT:
		for x := xFrom - 1; x >= 0; x-- {
			xTo = x
			if g.grid[yTo][xTo] != 0 {
				if g.grid[yFrom][xFrom] == g.grid[yTo][xTo] {
					merge = true
				} else {
					xTo = x + 1
				}
				break out
			}
		}
	case RIGHT:
		for x := xFrom + 1; x < g.gridWidth(); x++ {
			xTo = x
			if g.grid[yTo][xTo] != 0 {
				if g.grid[yFrom][xFrom] == g.grid[yTo][xTo] {
					merge = true
				} else {
					xTo = x - 1
				}
				break out
			}
		}
	}

	if yTo != yFrom || xTo != xFrom {
		if merge {
			g.grid[yTo][xTo] = g.grid[yFrom][xFrom] * 2
			g.grid[yFrom][xFrom] = 0
			g.currentScore += g.grid[yTo][xTo]
			if g.currentScore > g.highScore {
				g.highScore = g.currentScore
			}
		} else {
			g.grid[yTo][xTo] = g.grid[yFrom][xFrom]
			g.grid[yFrom][xFrom] = 0
		}
	}
}

func (g *game) spawnNewTile() {
	var emptyCoords = []coord{}

	for y := 0; y < g.gridHeight(); y++ {
		for x := 0; x < g.gridWidth(); x++ {
			if g.grid[y][x] == 0 {
				emptyCoords = append(emptyCoords, coord{y, x})
			}
		}
	}

	if len(emptyCoords) > 0 {
		var randomCoord = emptyCoords[int32(rand.Float64()*float64(len(emptyCoords)))]
		if rand.Float64() > 0.9 {
			g.grid[randomCoord.y][randomCoord.x] = 4
		} else {
			g.grid[randomCoord.y][randomCoord.x] = 2
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
	g.spawnNewTile()
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
