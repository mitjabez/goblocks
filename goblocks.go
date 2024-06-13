package main

import (
	"fmt"
	"time"
)

var arena [20][10]byte
var b = [3][4]byte{
	{0, 1, 0, 0},
	{0, 1, 0, 0},
	{0, 1, 1, 0},
}

func draw() {
	fmt.Print("\033[10;10H")
	arena[3][5] = 1
	for y, row := range arena {
		var displayRow [10]byte

		for x, cell := range row {
			if cell == 0 {
				displayRow[x] = ' '
			} else {
				displayRow[x] = 'X'
			}
		}

		fmt.Println(string(displayRow[:]))
		fmt.Print("\033[1B")
		fmt.Printf("\033[%d;10H", 11+y)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// fmt.Print("\033[2J")

	// draw()

	ch := make(chan byte)
	go readKey(ch)

	fmt.Println("app start")
	for {
		select {
		case stdin, _ := <-ch:
			fmt.Println("Keys pressed:", stdin)
		default:
			fmt.Println("Working A ..")
		}
		time.Sleep(time.Millisecond * 10)
	}
}
