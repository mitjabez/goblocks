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

const (
	Red     = iota
	Green   = iota
	Yellow  = iota
	Blue    = iota
	Magenta = iota
	Cyan    = iota
	White   = iota
)

type Block [BlockSize][BlockSize]byte

type Cell struct {
	isSet bool
	color uint
}

var arena [ArenaHeight][ArenaWidth]Cell
var blockL = Block{
	{0, 1, 0, 0},
	{0, 1, 0, 0},
	{0, 1, 1, 0},
	{0, 0, 0, 0},
}
var blockT = Block{
	{0, 1, 0, 0},
	{1, 1, 1, 0},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
}
var blockI = Block{
	{0, 0, 0, 0},
	{1, 1, 1, 1},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
}
var blockO = Block{
	{0, 1, 1, 0},
	{0, 1, 1, 0},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
}
var blockS = Block{
	{0, 1, 1, 0},
	{1, 1, 0, 0},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
}
var allBlocks = []Block{blockL, blockT, blockI, blockO, blockS}
var state = PlayingState

type Pos struct {
	x int
	y int
}

type Player struct {
	pos        Pos
	block      Block
	blockColor uint
	score      uint
	level      uint
}

var lastTick int64
var player Player

func exit() {
	showCursor()
	os.Exit(0)
}

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
				if arena[arenaY][arenaX].isSet {
					return false
				}
			}
		}
	}
	return true
}

func tryRotate() {
	var rotatedBlock Block

	for y, row := range player.block {
		for x, cell := range row {
			rotatedBlock[x][BlockSize-y-1] = cell
		}
	}

	if canMove(rotatedBlock, player.pos) {
		player.block = rotatedBlock
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

func generateBlock() Block {
	block := allBlocks[rand.Intn(len(allBlocks))]
	shouldFlip := lastTick%2 == 0

	if !shouldFlip {
		return block
	}

	var flippedBlock Block
	for y := range block {
		for x := 0; x < BlockSize/2; x++ {
			flippedBlock[y][x] = block[y][BlockSize-x-1]
			flippedBlock[y][BlockSize-x-1] = block[y][x]
		}
	}

	return flippedBlock

}

func newBlock() bool {
	player.block = generateBlock()

	player.pos.x = ArenaWidth/2 - 1
	// First row is empty
	player.pos.y = 0
	player.blockColor = newColor()

	return canMove(player.block, player.pos)
}

func removeRow(rowToRemove int) {
	for y := rowToRemove; y >= 0; y-- {
		for x := range arena[y] {
			if y > 0 {
				arena[y][x].isSet = arena[y-1][x].isSet
				arena[y][x].color = arena[y-1][x].color
			} else {
				arena[y][x].isSet = false
			}
		}
	}

	player.score += ArenaWidth
	player.level = player.score/30 + 1
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
					arena[y][x].isSet = true
					arena[y][x].color = player.blockColor
				}
			}

			if !arena[y][x].isSet {
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
		exit()
	}
}

func handleGameOverKey(key byte) {
	switch key {
	case KeySpace:
		newGame()
	case KeyEscape:
		exit()
	}
}

func gameLoop() {
	delay := int64(1000 - (player.level-1)*70)
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

func newColor() uint {
	return uint(rand.Intn(7-1) + 1)
}

func newGame() {
	for y := range arena {
		for x := range arena[y] {
			arena[y][x].isSet = false
			arena[y][x].color = 0
		}
	}

	newBlock()

	state = PlayingState
	player.score = 0
	player.level = 1
	player.blockColor = newColor()

	cls()
	hideCursor()
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
