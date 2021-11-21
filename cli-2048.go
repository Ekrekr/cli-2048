package main

import (
	"fmt"
	"math/rand"
	"time"
)

var testGrid = [4][4]int{
	{0, 4, 2, 0},
	{2, 0, 0, 0},
	{0, 128, 64, 0},
	{0, 128, 2048, 4194304}}

type game struct {
	currentScore int
	highScore    int
	grid         [4][4]int
	newGrid      [4][4]int
	gridWidth    int
	gridHeight   int
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// args := os.Args[1:]
	var g = game{currentScore: 2048, highScore: 4096, grid: testGrid, gridWidth: 4, gridHeight: 4}
	g.print()
	g.moveDown()
	g.print()
}

func (g *game) moveDown() {
	for y := g.gridHeight - 1; y > 0; y-- {
		for x := 0; x < g.gridWidth; x++ {
			if g.grid[y][x] == g.grid[y-1][x] {
				// Add tiles if they're equal.
				g.grid[y][x] += g.grid[y-1][x]
				g.grid[y-1][x] = 0
			}
			if g.grid[y][x] == 0 {
				// Replace the tile value with the tile above if the selected tile is empty.
				g.grid[y][x] = g.grid[y-1][x]
				g.grid[y-1][x] = 0
			}
		}
	}
	var emptySpaces []int
	for x := 0; x < g.gridWidth; x++ {
		if g.grid[0][x] == 0 {
			emptySpaces = append(emptySpaces, x)
		}
	}
	if len(emptySpaces) > 0 {
		var emptySpacesIndex = int32(rand.Float64() * float64(len(emptySpaces)))
		var newTilePlace = emptySpaces[emptySpacesIndex]
		if rand.Float64() > 0.9 {
			g.grid[0][newTilePlace] = 4
		} else {
			g.grid[0][newTilePlace] = 2
		}
	}
}

func (g *game) print() {
	fmt.Printf("  cli-2048                 %07d pts\n", g.currentScore)
	fmt.Printf("  High score:              %07d pts\n\n", g.highScore)

	for _, row := range g.grid {
		fmt.Print("\n\n  ")
		for _, val := range row {
			if val == 0 {
				fmt.Printf("    .    ")
			} else {
				fmt.Printf(" %07d ", val)
			}
		}
		fmt.Print("\n\n")
	}

	fmt.Printf("\n                ←,↑,→,↓              \n\n")
}
