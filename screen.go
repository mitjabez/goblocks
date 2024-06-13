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

func draw() {
	cursorXY(10, 10)
	arena[3][5] = 1
	for y, row := range arena {
		var displayRow [10]byte

		for x, cell := range row {
			if cell == 0 {
				displayRow[x] = '.'
			} else {
				displayRow[x] = 'X'
			}
		}

		fmt.Println(string(displayRow[:]))
		cursorXY(10, 11+y)
	}
}
