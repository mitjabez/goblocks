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

func cursorXY(x int, y int) {
	csi(strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H")
}

func draw(arena [20][10]byte, player Player) {
	fmt.Printf("%+v", player.pos)

	cursorXY(10, 10)
	for y, row := range arena {
		var displayRow [10]byte

		for x, cell := range row {
			if cell == 0 {
				displayRow[x] = '.'
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
		cursorXY(10, 11+y)
	}
}
