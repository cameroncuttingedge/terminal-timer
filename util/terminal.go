package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/term"
)

var screen = bufio.NewWriter(os.Stdout)


func Clear() {
	clearCode := "\033[H\033[2J"
	fmt.Fprint(os.Stdout, clearCode)
}

func CmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func HideCursor() {
	fmt.Fprint(screen, "\033[?25l")
}

func ShowCursor() {
	fmt.Fprint(screen, "\033[?25h")
}

func MoveCursor(pos [2]int) {
	fmt.Fprintf(screen, "\033[%d;%dH", pos[1], pos[0])
}

func Draw(str string) {
	fmt.Fprint(screen, str)
}

func Render() {
	screen.Flush()
}

func GetSize() (int, int, error) {
    width, height, err := term.GetSize(int(os.Stdout.Fd()))
    if err != nil {
        return 0, 0, err
    }

    return width, height, nil
}