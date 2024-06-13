package main

import (
	"fmt"
	"os"
	"time"
)

const ArenaWidth = 10
const ArenaHeight = 20

var arena [20][10]byte
var b = [3][4]byte{
	{0, 1, 0, 0},
	{0, 1, 0, 0},
	{0, 1, 1, 0},
}

type Pos struct {
	x int
	y int
}

var pos Pos
var lastTick int64

func gameLoop() {
	if time.Now().UnixMilli()-lastTick > 2000 {
		pos.y = min(pos.y+1, ArenaHeight-1)
		draw(arena, pos)
		lastTick = time.Now().UnixMilli()
	}
}

func handleKey(key byte) {
	switch key {
	case KeyUp:
		pos.y = max(pos.y-1, 0)
	case KeyDown:
		pos.y = min(pos.y+1, ArenaHeight-1)
	case KeyLeft:
		pos.x = max(pos.x-1, 0)
	case KeyRight:
		pos.x = min(pos.x+1, ArenaWidth-1)
	case KeyEscape:
		os.Exit(0)
	}

	draw(arena, pos)
}

func main() {
	//

	lastTick = time.Now().UnixMilli()

	cls()
	draw(arena, pos)

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
