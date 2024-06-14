package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const ArenaWidth = 10
const ArenaHeight = 20
const BlockSize = 4
const (
	PlayingState  = iota
	GameOverState = iota
)

type Block [BlockSize][BlockSize]byte

var arena [ArenaHeight][ArenaWidth]byte
var blockL = Block{
	{0, 0, 0, 0},
	{0, 1, 0, 0},
	{0, 1, 0, 0},
	{0, 1, 1, 0},
}
var blockT = Block{
	{0, 0, 0, 0},
	{0, 1, 0, 0},
	{1, 1, 1, 0},
	{0, 0, 0, 0},
}
var blockI = Block{
	{0, 0, 0, 0},
	{1, 1, 1, 1},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
}
var blockO = Block{
	{0, 0, 0, 0},
	{0, 1, 1, 0},
	{0, 1, 1, 0},
	{0, 0, 0, 0},
}
var blockS = Block{
	{0, 0, 0, 0},
	{0, 1, 1, 0},
	{1, 1, 0, 0},
	{0, 0, 0, 0},
}
var allBlocks = []Block{blockL, blockT, blockI, blockO, blockS}
var state = PlayingState

type Pos struct {
	x int
	y int
}

type Player struct {
	pos   Pos
	block Block
	score uint
	level uint
}

var lastTick int64
var player Player

func canMove(block Block, newPos Pos) bool {
	for y, row := range block {
		for x, cell := range row {
			if cell == 0 {
				continue
			}

			arenaY := newPos.y + y
			arenaX := newPos.x + x

			// Wall detection
			if arenaX < 0 || arenaX >= ArenaWidth || arenaY < 0 || arenaY >= ArenaHeight {
				return false
			}

			if newPos.y+y < ArenaHeight && newPos.x+x < ArenaWidth {
				// Object hit detection
				if arena[arenaY][arenaX] == 1 {
					return false
				}
			}
		}
	}
	return true
}

func tryRotate() {
	var newBlock Block

	for y, row := range player.block {
		for x, cell := range row {
			newBlock[x][BlockSize-y-1] = cell
		}
	}

	if canMove(newBlock, player.pos) {
		player.block = newBlock
		draw(arena, player)
	}
}

func tryMove(newPos Pos) bool {
	if canMove(player.block, newPos) {
		player.pos = newPos
		draw(arena, player)
		return true
	}
	return false
}

func gameOver() {
	state = GameOverState
	drawGameOver()
}

func newBlock() bool {
	player.block = allBlocks[rand.Intn(len(allBlocks))]
	player.pos.x = ArenaWidth/2 - 1
	// First row is empty
	player.pos.y = -1

	return canMove(player.block, player.pos)
}

func removeRow(rowToRemove int) {
	for y := rowToRemove; y >= 0; y-- {
		for x := range arena[y] {
			if y > 0 {
				arena[y][x] = arena[y-1][x]
			} else {
				arena[y][x] = 0
			}
		}
	}

	player.score += ArenaWidth
	player.level = player.score/50 + 1
}

func landBlock() {
	for y, row := range arena {
		isFullRow := true
		for x := range row {
			if x >= player.pos.x && x < player.pos.x+BlockSize &&
				y >= player.pos.y && y < player.pos.y+BlockSize {
				blockCell := player.block[y-player.pos.y][x-player.pos.x]
				if blockCell == 1 {
					// Fill
					arena[y][x] = 1
				}
			}

			if arena[y][x] == 0 {
				isFullRow = false
			}
		}

		if isFullRow {
			removeRow(y)
		}
	}
}

func handleKey(key byte) {
	switch key {
	case KeyUp:
		tryRotate()
	case KeyDown:
		tryMove(Pos{x: player.pos.x, y: player.pos.y + 1})
	case KeyLeft:
		tryMove(Pos{x: player.pos.x - 1, y: player.pos.y})
	case KeyRight:
		tryMove(Pos{x: player.pos.x + 1, y: player.pos.y})
	case KeySpace:
		// TODO: pull down
	case KeyEscape:
		os.Exit(0)
	}
}

func handleGameOverKey(key byte) {
	switch key {
	case KeySpace:
		newGame()
	case KeyEscape:
		os.Exit(0)
	}
}

func gameLoop() {
	delay := int64(1000 - (player.level-1)*50)
	if time.Now().UnixMilli()-lastTick > delay {
		if !tryMove(Pos{x: player.pos.x, y: player.pos.y + 1}) {
			landBlock()
			if !newBlock() {
				gameOver()
				return
			}

			draw(arena, player)
		}
		lastTick = time.Now().UnixMilli()
	}
}

func newGame() {
	for y := range arena {
		for x := range arena[y] {
			arena[y][x] = 0
		}
	}

	newBlock()

	state = PlayingState
	player.score = 0
	player.level = 1

	cls()
	drawUI()
	draw(arena, player)

	lastTick = time.Now().UnixMilli()
}

func main() {
	newGame()

	ch := make(chan byte)
	go readKey(ch)

	for {
		select {
		case key, _ := <-ch:
			if state == PlayingState {
				handleKey(key)
			} else if state == GameOverState {
				handleGameOverKey(key)
			} else {
				panic(fmt.Sprint("Unknown state ", state))
			}
		default:
		}

		if state == PlayingState {
			gameLoop()
		}

		time.Sleep(time.Millisecond * 10)
	}
}
