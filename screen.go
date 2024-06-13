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

func draw(arena [20][10]byte, pos Pos) {
	fmt.Printf("%+v", pos)

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
			if x >= pos.x && x < pos.x+BlockSize &&
				y >= pos.y && y < pos.y+BlockSize {
				if b[y-pos.y][x-pos.x]|cell == 1 {
					displayRow[x] = 'X'
				}
			}
		}

		fmt.Println(string(displayRow[:]))
		cursorXY(10, 11+y)
	}
}
