package main

import (
	"fmt"
	"os"
	"time"
)

const ArenaWidth = 10
const ArenaHeight = 20
const BlockSize = 4

type Block [4][4]byte

var arena [20][10]byte
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

func gameLoop() {
	if time.Now().UnixMilli()-lastTick > 2000 {
		player.pos.y = min(player.pos.y+1, ArenaHeight-1)
		draw(arena, player)
		lastTick = time.Now().UnixMilli()
	}
}

func handleKey(key byte) {
	switch key {
	case KeyUp:
		player.pos.y = max(player.pos.y-1, 0)
	case KeyDown:
		player.pos.y = min(player.pos.y+1, ArenaHeight-1)
	case KeyLeft:
		player.pos.x = max(player.pos.x-1, 0)
	case KeyRight:
		player.pos.x = min(player.pos.x+1, ArenaWidth-1)
	case KeyEscape:
		os.Exit(0)
	}

	draw(arena, player)
}

func main() {
	//

	lastTick = time.Now().UnixMilli()
	player.block = blockL

	cls()
	draw(arena, player)

	ch := make(chan byte)
	go readKey(ch)

	fmt.Println("app start")
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
