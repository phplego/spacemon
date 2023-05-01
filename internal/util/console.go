package util

import "fmt"

func RemoveLines(n int) {
	// Remove last n lines
	for i := 0; i < n; i++ {
		fmt.Print("\033[1A") // Mover cursor up
		fmt.Print("\033[K")  // Clear current line
	}
}

func Backspace(n int) {
	fmt.Printf("\033[%dD\033[%dP", n, n)
}

func ClearScreen(gotoTop bool) {
	// clear
	fmt.Print("\033[2J")
	if gotoTop {
		// go to top left corner
		fmt.Print("\u001B[H")
	}
}

func ClearAndPrint(str string) {
	ClearScreen(true)
	fmt.Println(str)
}
