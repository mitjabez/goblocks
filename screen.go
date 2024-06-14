package main

import (
	"fmt"
	"strconv"
)

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
	cursorYX(7, 10)
	fmt.Print("[DEBUG] ", msg)
}

func drawUI() {
	for y := 0; y < ArenaHeight; y++ {
		cursorYX(y+10, 9)
		fmt.Print("|          |")
	}
	cursorYX(ArenaHeight+10, 9)
	fmt.Print("------------")
}

func drawGameOver() {
	cursorYX(17, 9)
	fmt.Print("************")
	cursorYX(18, 9)
	fmt.Print("*GAME OVER!*")
	cursorYX(19, 9)
	fmt.Print("************")
}

func draw(arena [ArenaHeight][ArenaWidth]byte, player Player) {
	cursorYX(10, 10)
	for y, row := range arena {
		var displayRow [10]byte

		for x, cell := range row {
			if cell == 0 {
				displayRow[x] = ' '
			} else {
				displayRow[x] = 'X'
			}

			// Block withing range
			if x >= player.pos.x && x < player.pos.x+BlockSize &&
				y >= player.pos.y && y < player.pos.y+BlockSize {
				if player.block[y-player.pos.y][x-player.pos.x]|cell == 1 {
					displayRow[x] = 'X'
				}
			}
		}

		fmt.Println(string(displayRow[:]))
		cursorYX(11+y, 10)
	}
	cursorYX(0, 0)
}
