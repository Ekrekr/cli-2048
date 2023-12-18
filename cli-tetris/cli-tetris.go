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
	PIECE_L     = 0
	PIECE_R     = 1
	PIECE_COUNT = 2 // Num of pieces in total.
)

const kGridWidth = 10
const kGridHeight = 19

type Coord struct {
	Y uint32
	X uint32
}

type Piece struct {
	ActiveCoords []Coord
	Color        color.Attribute
	Width        uint32
	Height       uint32
}

// Pieces here are stored with zeroed coords.
var pieces = map[uint32]Piece{
	PIECE_L: {ActiveCoords: []Coord{{0, 0}, {0, 1}, {0, 2}, {1, 0}}, Color: color.BgBlue, Width: 3, Height: 2},
	PIECE_R: {ActiveCoords: []Coord{{1, 0}, {1, 1}, {1, 2}, {0, 0}}, Color: color.BgGreen, Width: 3, Height: 2},
}

// Variables are exported to make reading from file possible.
type game struct {
	CurrentScore      uint32
	HighScore         uint32
	Grid              [kGridHeight][kGridWidth]color.Attribute
	IsGameOver        bool
	NeedScreenRefresh bool
	UpcomingPieces    []uint32
	ActivePiece       Piece
	DebugMode         bool
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var g = game{}

	// TODO: Load save instead.
	// g.loadSave()
	g.newGame()
	g.NeedScreenRefresh = true
	g.IsGameOver = false

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
		charRune, key, _ := keyboard.GetKey()
		switch key {
		case keyboard.KeyArrowLeft:
			// Move left.
		case keyboard.KeyArrowRight:
			// Move right.
		case keyboard.KeyCtrlN:
			g.newGame()
			g.NeedScreenRefresh = true
			g.IsGameOver = false
		case keyboard.KeyCtrlQ:
		case keyboard.KeyCtrlC:
			g.createSave()
			play = false
		}
		// TODO: Verify this works.
		if !g.IsGameOver {
			switch charRune {
			case []rune("w")[0]:
				// Drop piece.
				// Swap to new peice.
				// Refresh piece store.
			case []rune("s")[0]:
				// Move down.
				for _, oldCoord := range g.ActivePiece.ActiveCoords {
					// TODO: Place based on timer instead.
					if g.Grid[oldCoord.Y+1][oldCoord.X] != color.Reset {
						g.spawnNewPiece()
						continue
					}
				}
				for i, oldCoord := range g.ActivePiece.ActiveCoords {
					newCoord := Coord{Y: oldCoord.Y + 1, X: oldCoord.X}
					g.Grid[newCoord.Y][oldCoord.X] = g.Grid[oldCoord.Y][oldCoord.X]
					g.ActivePiece.ActiveCoords[i] = newCoord
				}
			case []rune("a")[0]:
				// // Move left.
				// for _, oldCoord := range g.ActivePiece.ActiveCoords {
				// 	newX := oldCoord.X - 1
				// 	// TODO: Place based on timer instead.
				// 	if newX < 0 || g.Grid[oldCoord.Y][newX] != color.Reset {
				// 		return
				// 	}
				// }
				// for _, oldCoord := range g.ActivePiece.ActiveCoords {
				// 	g.Grid[oldCoord.Y][oldCoord.X-1] = g.Grid[oldCoord.Y][oldCoord.X]
				// }
			case []rune("d")[0]:
				// Move right.
			case []rune("q")[0]:
				// Rotate left.
			case []rune("e")[0]:
				// Rotate right.
			}
		}
	}
}

func (g *game) newGame() {
	g.Grid = [kGridHeight][kGridWidth]color.Attribute{}
	for i := range [kGridWidth]int{} {
		g.Grid[i] = [kGridWidth]color.Attribute{}
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
	newPiece := pieces[g.UpcomingPieces[0]]
	newPiece.ActiveCoords = []Coord{}
	// Place the piece from the localised piece shape grid onto the main game grid.
	for _, tile := range newPiece.ActiveCoords {
		gridCoord := Coord{X: (kGridWidth-newPiece.Width)/2 + tile.X, Y: tile.Y}
		if (g.Grid[gridCoord.X][gridCoord.Y]) != color.Reset {
			// Reset indicates a populated tile; the new piece cannot fit on
			// the board so the game is over.
			g.IsGameOver = true
		} else {
			g.Grid[gridCoord.X][gridCoord.Y] = newPiece.Color
		}
		newPiece.ActiveCoords = append(newPiece.ActiveCoords, gridCoord)
	}
	g.ActivePiece = newPiece
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

	borderTile := printTile(color.BgHiBlue)
	emptyTile := "  "

	for range [kGridWidth + 2]int{} {
		output += borderTile
	}
	output += "\n"
	for _, row := range g.Grid {
		output += emptyTile + borderTile // Left padding and border.
		for _, colorCode := range row {
			output += printTile(colorCode)
		}
		output += borderTile + "\n"
	}
	output += emptyTile
	for range [kGridWidth + 2]int{} {
		output += borderTile
	}
	output += "\n"

	if g.IsGameOver {
		output += fmt.Sprintf("\n  ----------------------------")
		output += fmt.Sprintf("\n    >>> 💀 GAME OVER! 💀 <<<  ")
		output += fmt.Sprintf("\n  ----------------------------\n")
	}

	output += fmt.Sprintf("\n   ←,↑,→,↓  💾ctrl-c 🔄ctrl-n \n\n")

	if g.DebugMode {
		output += fmt.Sprintf("\nActive piece: %+v\n", g.ActivePiece)
		output += fmt.Sprintf("Upcoming pieces: %+v\n", g.UpcomingPieces)
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

// Note that tiles are two spaces wide, to make them appear square in a console.
func printTile(col color.Attribute) string {
	c := color.New()
	c.Add(col)
	return c.Sprintf("  ")
}

func clearScreen() {
	cmd := exec.Command("clear") // Only works for unix based systems.
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func startKeyboard() func() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	return func() {
		_ = keyboard.Close()
	}
}
