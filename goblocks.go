package main

import (
	"math/rand"
	"os"
	"time"
)

const ArenaWidth = 10
const ArenaHeight = 20
const BlockSize = 4

type Block [BlockSize][BlockSize]byte

var arena [ArenaHeight][ArenaWidth]byte
var blockL = Block{
	{1, 0, 0, 0},
	{1, 0, 0, 0},
	{1, 1, 0, 0},
	{0, 0, 0, 0},
}
var blockT = Block{
	{0, 1, 0, 0},
	{1, 1, 1, 0},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
}
var blockI = Block{
	{1, 1, 1, 1},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
	{0, 0, 0, 0},
}
var blockO = Block{
	{1, 1, 0, 0},
	{1, 1, 0, 0},
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

type Pos struct {
	x int
	y int
}

type Player struct {
	pos   Pos
	block Block
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

func newBlock() {
	player.block = allBlocks[rand.Intn(len(allBlocks))]
	player.pos.x = ArenaWidth/2 - 1
	player.pos.y = 0
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

func gameLoop() {
	if time.Now().UnixMilli()-lastTick > 1000 {
		if !tryMove(Pos{x: player.pos.x, y: player.pos.y + 1}) {
			landBlock()
			newBlock()
			draw(arena, player)
		}
		lastTick = time.Now().UnixMilli()
	}
}

func main() {
	lastTick = time.Now().UnixMilli()
	newBlock()

	cls()
	drawUI()
	draw(arena, player)

	ch := make(chan byte)
	go readKey(ch)

	for {
		select {
		case key, _ := <-ch:
			handleKey(key)
		default:
		}
		gameLoop()
		time.Sleep(time.Millisecond * 10)
	}
}
