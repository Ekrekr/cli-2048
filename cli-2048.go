package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

const (
	TOP    = 1
	RIGHT  = 2
	BOTTOM = 3
	LEFT   = 4
)

var testGrid = [][]int{
	{0, 4, 2, 0},
	{2, 0, 0, 0},
	{0, 128, 64, 0},
	{0, 128, 2048, 32768}}

type game struct {
	currentScore int
	highScore    int
	grid         [][]int
	newGrid      [][]int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// args := os.Args[1:]

	var g = game{currentScore: 2048, highScore: 4096, grid: testGrid}

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	var play = true
	for play {
		clearScreen()
		g.print()
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyArrowUp {
			g.moveUp()
		}
		if key == keyboard.KeyArrowDown {
			g.moveDown()
		}
		if key == keyboard.KeyArrowLeft {
			g.moveLeft()
		}
		if key == keyboard.KeyArrowRight {
			g.moveRight()
		}
		if key == keyboard.KeyCtrlQ || key == keyboard.KeyCtrlC {
			play = false
		}
	}

}

func (g *game) moveDown() {
	for y := g.gridHeight() - 1; y > 0; y-- {
		for x := 0; x < g.gridWidth(); x++ {
			g.combineTilesIfEqual(y, x, y-1, x)
			g.moveTileIfSpaceEmpty(y, x, y-1, x)
		}
	}
	g.spawnNewTileIfNeeded(TOP)
}

func (g *game) moveUp() {
	for y := 0; y < g.gridHeight()-1; y++ {
		for x := 0; x < g.gridWidth(); x++ {
			g.combineTilesIfEqual(y, x, y+1, x)
			g.moveTileIfSpaceEmpty(y, x, y+1, x)
		}
	}
	g.spawnNewTileIfNeeded(BOTTOM)
}

func (g *game) moveRight() {
	for y := 0; y < g.gridHeight(); y++ {
		for x := g.gridWidth() - 1; x > 0; x-- {
			g.combineTilesIfEqual(y, x, y, x-1)
			g.moveTileIfSpaceEmpty(y, x, y, x-1)
		}
	}
	g.spawnNewTileIfNeeded(LEFT)
}

func (g *game) moveLeft() {
	for y := 0; y < g.gridHeight(); y++ {
		for x := 0; x < g.gridWidth()-1; x++ {
			g.combineTilesIfEqual(y, x, y, x+1)
			g.moveTileIfSpaceEmpty(y, x, y, x+1)
		}
	}
	g.spawnNewTileIfNeeded(RIGHT)
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
		if side == TOP {
			g.grid[0][newTilePlace] = newTileValue
		}
		if side == BOTTOM {
			g.grid[g.gridHeight()-1][newTilePlace] = newTileValue
		}
		if side == LEFT {
			g.grid[newTilePlace][0] = newTileValue
		}
		if side == RIGHT {
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

func (g *game) print() {
	fmt.Printf("\ncli-2048              %07d pts\n", g.currentScore)
	fmt.Printf("High score:           %07d pts\n\n", g.highScore)

	for _, row := range g.grid {
		for i := range []int{0, 1, 2} {
			for _, val := range row {
				var tilePrinter = getTilePrinter(val)
				if i == 1 {
					if val == 0 {
						tilePrinter("   .   ")
					} else {
						tilePrinter(" %05d ", val)
					}
				} else {
					tilePrinter("       ")
				}
			}
			fmt.Print("\n")
		}
	}

	fmt.Printf("\n     ←,↑,→,↓ or ctrl-q   \n\n")
}

func getTilePrinter(tile int) func(format string, a ...interface{}) (n int, err error) {
	c := color.New()
	if tile == 2 {
		c.Add(color.BgRed)
	}
	if tile == 4 {
		c.Add(color.BgGreen)
	}
	if tile == 8 {
		c.Add(color.BgYellow)
	}
	if tile == 16 {
		c.Add(color.BgBlue)
	}
	if tile == 32 {
		c.Add(color.BgMagenta)
	}
	if tile == 64 {
		c.Add(color.BgCyan)
	}
	if tile == 128 {
		c.Add(color.BgWhite)
	}
	if tile == 256 {
		c.Add(color.BgHiBlack)
	}
	if tile == 512 {
		c.Add(color.BgHiRed)
	}
	if tile == 1024 {
		c.Add(color.BgHiGreen)
	}
	if tile == 2048 {
		c.Add(color.BgHiYellow)
	}
	if tile == 4096 {
		c.Add(color.BgHiBlue)
	}
	if tile == 8192 {
		c.Add(color.BgHiMagenta)
	}
	if tile == 16384 {
		c.Add(color.BgHiCyan)
	}
	if tile == 32768 {
		c.Add(color.BgHiWhite)
	}
	return c.Printf
}

func clearScreen() {
	cmd := exec.Command("clear") // Only works for unix based systems.
	cmd.Stdout = os.Stdout
	cmd.Run()
}
