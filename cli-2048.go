package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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

// Variables are exported to make reading from file possible.
type game struct {
	CurrentScore      uint32
	HighScore         uint32
	Grid              [4][4]uint32
	IsGameOver        bool
	NeedScreenRefresh bool
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var g = game{}
	g.loadSave()

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	var play = true
	for play {
		// TODO: Replace changed parts of string rather than re-render whole game.
		if g.NeedScreenRefresh {
			var printString = g.getGameDisplay()
			clearScreen()
			fmt.Print(printString)
			g.NeedScreenRefresh = false
		}
		// Error is ignored because combinations like shift + arrow key will throw.
		_, key, _ := keyboard.GetKey()
		switch key {
		case keyboard.KeyArrowUp:
			if g.canMove(UP) {
				g.move(UP)
				g.NeedScreenRefresh = true
			}
		case keyboard.KeyArrowDown:
			if g.canMove(DOWN) {
				g.move(DOWN)
				g.NeedScreenRefresh = true
			}
		case keyboard.KeyArrowLeft:
			if g.canMove(LEFT) {
				g.move(LEFT)
				g.NeedScreenRefresh = true
			}
		case keyboard.KeyArrowRight:
			if g.canMove(RIGHT) {
				g.move(RIGHT)
				g.NeedScreenRefresh = true
			}
		case keyboard.KeyCtrlN:
			g.newGame()
			g.NeedScreenRefresh = true
			g.IsGameOver = false
		case keyboard.KeyCtrlQ:
		case keyboard.KeyCtrlC:
			g.createSave()
			play = false
		}
		g.IsGameOver = g.checkIsGameOver()
	}
}

func (g *game) canMove(direction int) bool {
	for y := 0; y < g.gridHeight(); y++ {
		for x := 0; x < g.gridWidth(); x++ {
			if direction == UP && y > 0 &&
				(g.Grid[y][x] == g.Grid[y-1][x] || g.Grid[y-1][x] == 0) {
				return true
			}
			if direction == DOWN && y < g.gridHeight()-1 &&
				(g.Grid[y][x] == g.Grid[y+1][x] || g.Grid[y+1][x] == 0) {
				return true
			}
			if direction == LEFT && x > 0 &&
				(g.Grid[y][x] == g.Grid[y][x-1] || g.Grid[y][x-1] == 0) {
				return true
			}
			if direction == RIGHT && x < g.gridWidth()-1 &&
				(g.Grid[y][x] == g.Grid[y][x+1] || g.Grid[y][x+1] == 0) {
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
	if g.Grid[yFrom][xFrom] == 0 {
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
			if g.Grid[yTo][xTo] != 0 {
				if g.Grid[yFrom][xFrom] == g.Grid[yTo][xTo] {
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
			if g.Grid[yTo][xTo] != 0 {
				if g.Grid[yFrom][xFrom] == g.Grid[yTo][xTo] {
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
			if g.Grid[yTo][xTo] != 0 {
				if g.Grid[yFrom][xFrom] == g.Grid[yTo][xTo] {
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
			if g.Grid[yTo][xTo] != 0 {
				if g.Grid[yFrom][xFrom] == g.Grid[yTo][xTo] {
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
			g.Grid[yTo][xTo] = g.Grid[yFrom][xFrom] * 2
			g.Grid[yFrom][xFrom] = 0
			g.CurrentScore += g.Grid[yTo][xTo]
			if g.CurrentScore > g.HighScore {
				g.HighScore = g.CurrentScore
			}
		} else {
			g.Grid[yTo][xTo] = g.Grid[yFrom][xFrom]
			g.Grid[yFrom][xFrom] = 0
		}
	}
}

func (g *game) spawnNewTile() {
	var emptyCoords = []coord{}

	for y := 0; y < g.gridHeight(); y++ {
		for x := 0; x < g.gridWidth(); x++ {
			if g.Grid[y][x] == 0 {
				emptyCoords = append(emptyCoords, coord{y, x})
			}
		}
	}

	if len(emptyCoords) > 0 {
		var randomCoord = emptyCoords[int32(rand.Float64()*float64(len(emptyCoords)))]
		if rand.Float64() > 0.9 {
			g.Grid[randomCoord.y][randomCoord.x] = 4
		} else {
			g.Grid[randomCoord.y][randomCoord.x] = 2
		}
	}
}

func (g *game) gridHeight() int {
	return len(g.Grid)
}

func (g *game) gridWidth() int {
	return len(g.Grid[0])
}

func (g *game) newGame() {
	g.Grid = [4][4]uint32{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
	g.spawnNewTile()
	g.spawnNewTile()
	g.CurrentScore = 0
	g.NeedScreenRefresh = true
}

func (g *game) getGameDisplay() string {
	// Printing output all in one go removes appearance lag in the terminal.
	var output = ""
	output += fmt.Sprintf("\n            cli-2048          \n  ")
	for q := 0; q <= 8-len(strconv.Itoa(int(g.CurrentScore))); q++ {
		output += fmt.Sprintf(" ")
	}
	output += fmt.Sprintf("🎮 %d || %d 🏆\n\n  ", g.CurrentScore, g.HighScore)

	for _, row := range g.Grid {
		for i := range []int{0, 1, 2} {
			for _, val := range row {
				var tilePrinter = getTilePrinter(val)
				if i == 1 {
					if val == 0 {
						output += tilePrinter("   .   ")
					} else {
						var tileNumLength = len(strconv.Itoa(int(val)))
						var rightPadding = (5 - tileNumLength) / 2
						var isEven = 1 - tileNumLength%2
						var leftPadding = rightPadding + isEven
						for q := 0; q <= leftPadding; q++ {
							output += tilePrinter(" ")
						}
						output += tilePrinter("%d", val)
						for q := 0; q <= rightPadding; q++ {
							output += tilePrinter(" ")
						}
					}
				} else {
					output += tilePrinter("       ")
				}
			}
			output += fmt.Sprintf("\n  ")
		}
	}

	if g.IsGameOver {
		output += fmt.Sprintf("\n  ----------------------------")
		output += fmt.Sprintf("\n    >>> 💀 GAME OVER! 💀 <<<  ")
		output += fmt.Sprintf("\n  ----------------------------\n")
	}

	output += fmt.Sprintf("\n   ←,↑,→,↓  💾ctrl-c 🔄ctrl-n \n\n")

	return output
}

func (g *game) loadSave() {
	if _, err := os.Stat(getSavePath()); errors.Is(err, os.ErrNotExist) {
		g.createSaveFileIfNotExists()
		return
	}

	f, err := os.OpenFile(getSavePath(), os.O_RDONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	err = binary.Read(f, binary.LittleEndian, g)
	if err != nil {
		panic(err)
	}
	f.Close()
	g.NeedScreenRefresh = true
}

func (g *game) createSaveFileIfNotExists() {
	f, err := os.Create(getSavePath())
	if err != nil {
		panic(err)
	}
	f.Close()
	g.newGame()
}

func (g *game) createSave() {
	f, err := os.OpenFile(getSavePath(), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	err = binary.Write(f, binary.LittleEndian, g)
	if err != nil {
		panic(err)
	}
	f.Close()
}

func getSavePath() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, "Documents", ".2048-save")
}

func getTilePrinter(tile uint32) func(format string, a ...interface{}) string {
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
	return c.Sprintf
}

func clearScreen() {
	cmd := exec.Command("clear") // Only works for unix based systems.
	cmd.Stdout = os.Stdout
	cmd.Run()
}
