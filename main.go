package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

var KeyEscape byte = 0
var KeySpace byte = 1
var KeyUp byte = 11
var KeyDown byte = 12
var KeyLeft byte = 13
var KeyRight byte = 14

var arena [20][10]byte
var b = [3][4]byte{
	{0, 1, 0, 0},
	{0, 1, 0, 0},
	{0, 1, 1, 0},
}

func readKey() {
	// https://stackoverflow.com/a/54423725
	ch := make(chan byte)
	go func(ch chan byte) {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b []byte = make([]byte, 5)
		for {
			nRead, err := os.Stdin.Read(b)
			if err != nil {
				panic(err)
			}

			// CSI codes
			if nRead == 3 && b[0] == '\033' && b[1] == '[' {
				switch b[2] {
				case 'A':
					ch <- KeyUp
				case 'B':
					ch <- KeyDown
				case 'C':
					ch <- KeyRight
				case 'D':
					ch <- KeyLeft
				}
			} else if nRead == 1 {
				switch b[0] {
				case 27:
					ch <- KeyEscape
				case 32:
					ch <- KeySpace
				}
			}
		}
	}(ch)

	for {
		select {
		case stdin, _ := <-ch:
			fmt.Println("Keys pressed:", stdin)
		default:
			fmt.Println("Working..")
		}
		time.Sleep(time.Millisecond * 10)
	}
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
	fmt.Print("\033[2J")

	// https://en.wikipedia.org/wiki/ANSI_escape_code#CSI_(Control_Sequence_Introducer)_sequences
	// for i := 0; i < 80; i++ {
	// 	fmt.Printf("\033[%d;%dHError is\n ", i, i)
	// 	time.Sleep(50 * time.Millisecond)
	// }

	// draw()

	fmt.Println("Reading keys ...")
	readKey()
}
