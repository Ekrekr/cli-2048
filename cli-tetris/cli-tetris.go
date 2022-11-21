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
	// Types of pieces.
	PIECE_L         = 0
	PIECE_L_REVERSE = 1
	PIECE_COUNT     = 2 // Num of pieces in total.
)

const kGridWidth = 16
const kGridHeight = 48

type Coord struct {
	y uint32
	x uint32
}

// Pieces are stored in a smaller grid representation.
var pieces = map[uint32][]Coord{
	PIECE_L:         {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 0}},
	PIECE_L_REVERSE: {{1, 0}, {1, 1}, {1, 2}, {1, 3}, {0, 0}},
}

// Variables are exported to make reading from file possible.
type game struct {
	CurrentScore      uint32
	HighScore         uint32
	Grid              [kGridWidth][kGridHeight]uint32
	IsGameOver        bool
	NeedScreenRefresh bool
	UpcomingPieces    []uint32
	ActivePiece       []Coord
	DebugMode         bool
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var g = game{}
	// g.loadSave()

	defer startKeyboard()

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
			// Drop piece.
			// Swap to new peice.
			// Refresh piece store.
		case keyboard.KeyArrowDown:
			// Move piece down 1.
		case keyboard.KeyArrowLeft:
			// Rotate left.
		case keyboard.KeyArrowRight:
			// Rotate right.
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

func startKeyboard() func() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	return func() {
		_ = keyboard.Close()
	}
}

func (g *game) checkIsGameOver() bool {
	// Check if piece is above top.
	return false
}

func (g *game) newGame() {
	g.Grid = [kGridWidth][kGridHeight]uint32{}
	for i := range [kGridWidth]int{} {
		g.Grid[i] = [kGridHeight]uint32{}
	}
	for range [5]int{} {
		g.UpcomingPieces = append(g.UpcomingPieces, rand.Uint32()%PIECE_COUNT)
	}
	g.spawnNewPiece()
	g.CurrentScore = 0
	g.NeedScreenRefresh = true
	g.DebugMode = true
}

func (g *game) spawnNewPiece() {
	g.ActivePiece = pieces[g.UpcomingPieces[0]]
	g.UpcomingPieces = g.UpcomingPieces[1:]
	g.UpcomingPieces = append(g.UpcomingPieces, rand.Uint32()%PIECE_COUNT)
}

func (g *game) getGameDisplay() string {
	// Printing output all in one go removes appearance lag in the terminal.
	var output = ""
	output += fmt.Sprintf("\n            cli-tetris          \n  ")
	for q := 0; q <= 8-len(strconv.Itoa(int(g.CurrentScore))); q++ {
		output += fmt.Sprintf(" ")
	}
	output += fmt.Sprintf("🎮 %d || %d 🏆\n\n  ", g.CurrentScore, g.HighScore)

	for _, row := range g.Grid {
		output += "  " // Left padding.
		for _, colorCode := range row {
			printTile(colorCode)
		}
	}

	if g.IsGameOver {
		output += fmt.Sprintf("\n  ----------------------------")
		output += fmt.Sprintf("\n    >>> 💀 GAME OVER! 💀 <<<  ")
		output += fmt.Sprintf("\n  ----------------------------\n")
	}

	output += fmt.Sprintf("\n   ←,↑,→,↓  💾ctrl-c 🔄ctrl-n \n\n")

	if g.DebugMode {
		output += fmt.Sprintf("\nActive piece: %v\n", g.ActivePiece)
		output += fmt.Sprintf("Upcoming pieces: %v\n", g.UpcomingPieces)
	}

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

func printTile(colorCode uint32) string {
	c := color.New()
	c.Add([]color.Attribute{
		color.BgRed,
		color.BgGreen,
		color.BgYellow,
		color.BgBlue,
		color.BgMagenta,
		color.BgCyan,
		color.BgBlack,
		color.BgHiRed,
		color.BgHiGreen,
		color.BgHiYellow,
		color.BgHiBlue,
		color.BgHiMagenta,
		color.BgHiCyan,
		color.BgHiWhite,
		color.FgBlack,
	}[colorCode])
	return c.Sprintf(" ")
}

func clearScreen() {
	cmd := exec.Command("clear") // Only works for unix based systems.
	cmd.Stdout = os.Stdout
	cmd.Run()
}
