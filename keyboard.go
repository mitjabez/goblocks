package main

import (
	"os"
	"os/exec"
)

var KeyEscape byte = 0
var KeySpace byte = 1
var KeyUp byte = 11
var KeyDown byte = 12
var KeyLeft byte = 13
var KeyRight byte = 14

func readKey(ch chan byte) {
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
}
