package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ArenaConfig struct {
	posX             int
	posy             int
	arenaScaleWidth  int
	arenaScaleHeight int
}

var arenaConfig = ArenaConfig{
	posX: 5,
	posy: 5,
	// Font blocks are usually higher than wider. So let's just scale x2
	arenaScaleWidth:  ArenaWidth * 2,
	arenaScaleHeight: ArenaHeight,
}

// https://en.wikipedia.org/wiki/ANSI_escape_code#3-bit_and_4-bit
var colors = []string{"40", "41", "42", "43", "44", "45", "46", "47"}

func csi(sequence string) {
	fmt.Print("\033[", sequence)
}

func cls() {
	csi("2J")
}

func cursorDown() {
	csi("1B")
}

func cursorYX(y int, x int) {
	csi(strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H")
}

func debug(msg string) {
	cursorYX(1, arenaConfig.posX)
	fmt.Print("[DEBUG] ", msg)
}

func drawGameOver() {
	y := arenaConfig.posy + arenaConfig.arenaScaleHeight/3
	x := arenaConfig.posX + 2
	cursorYX(y, x)
	fmt.Print(strings.Repeat("*", arenaConfig.arenaScaleWidth-4))
	cursorYX(y+1, x)
	fmt.Print("*  GAME OVER!  *")
	cursorYX(y+2, x)
	fmt.Print(strings.Repeat("*", arenaConfig.arenaScaleWidth-4))
}

func resetColor() {
	csi("0m")
}

func hideCursor() {
	csi("?25l")
}

func showCursor() {
	csi("?25h")
}

func drawUI() {
	y := 0
	for ; y < arenaConfig.arenaScaleHeight; y++ {
		cursorYX(arenaConfig.posy+y, arenaConfig.posX-1)
		fmt.Print("*", strings.Repeat(" ", arenaConfig.arenaScaleWidth), "*")
	}
	cursorYX(arenaConfig.posy+y, arenaConfig.posX-1)
	fmt.Print(strings.Repeat("*", arenaConfig.arenaScaleWidth+2))
}

func draw(arena [ArenaHeight][ArenaWidth]Cell, player Player) {
	y := arenaConfig.posy + arenaConfig.arenaScaleHeight/3
	x := arenaConfig.posX + arenaConfig.arenaScaleWidth + 2
	cursorYX(y, x)
	fmt.Printf("Level: %d", player.level)
	cursorYX(y+1, x)
	fmt.Printf("Score: %d", player.score)

	for y, row := range arena {
		rowBuffer := ""
		for x, cell := range row {
			// Block withing range
			if x >= player.pos.x && x < player.pos.x+BlockSize &&
				y >= player.pos.y && y < player.pos.y+BlockSize {
				if player.block[y-player.pos.y][x-player.pos.x] == 1 {
					cell.color = player.blockColor
				}
			}

			rowBuffer += "\033[30;" + colors[cell.color] + "m  \033[0m"
		}

		cursorYX(arenaConfig.posy+y, arenaConfig.posX)

		fmt.Print(rowBuffer)
	}
}
