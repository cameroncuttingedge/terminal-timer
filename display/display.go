package displays

import (
	"fmt"
	"strings"
	"terminal-timer/util"
)

type DisplayMatrix struct {
	Width  int
	Height int
	Matrix [][]rune
}

func NewDisplayMatrix(width, height int) *DisplayMatrix {
	matrix := make([][]rune, height)
	for i := range matrix {
		matrix[i] = make([]rune, width)
		for j := range matrix[i] {
			matrix[i][j] = ' ' // Initialize with spaces
		}
	}
	return &DisplayMatrix{Width: width, Height: height, Matrix: matrix}
}


func (dm *DisplayMatrix) AddCenteredMessage(message string) {
    lines := strings.Split(message, "\n")
    totalLines := len(lines)

    startY := (dm.Height - totalLines) / 2
    if startY < 0 {
        startY = 0
    }

    for i, line := range lines {
        maxWidth := len(line)
        startX := (dm.Width - maxWidth) / 2
        if startX < 0 {
            startX = 0
        }

        for x, char := range line {
            matrixY := startY + i
            matrixX := startX + x
            if matrixY < dm.Height && matrixX < dm.Width {
                dm.Matrix[matrixY][matrixX] = char
            }
        }
    }
}


func (dm *DisplayMatrix) AddCenteredAsciiArt(asciiArt []string, message string) {
    artHeight := len(asciiArt)
    artWidth := 0
    for _, line := range asciiArt {
        if len(line) > artWidth {
            artWidth = len(line)
        }
    }

    // Check if the art will fit vertically and horizontally
    if artHeight > dm.Height || artWidth > dm.Width {
        dm.AddCenteredMessage(message)
        return
    }

    startY := (dm.Height - artHeight) / 2
    startX := (dm.Width - artWidth) / 2

    for y, line := range asciiArt {
        for x, char := range line {
            matrixY := startY + y
            matrixX := startX + x
            if matrixY < dm.Height && matrixX < dm.Width {
                dm.Matrix[matrixY][matrixX] = char
            }
        }
    }
}

func (dm *DisplayMatrix) Print() {
    var output strings.Builder
    for _, row := range dm.Matrix {
        for _, char := range row {
            output.WriteRune(char)
        }
        output.WriteRune('\n')
    }
    util.HideCursor()
    util.ClearTerminal()
    fmt.Print(output.String())
}


func (dm *DisplayMatrix) ResizeAndClear() {
	width, height, err := util.GetSize()
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return
	}
    newMatrix := make([][]rune, height)
    for i := range newMatrix {
        newMatrix[i] = make([]rune, width)
        for j := range newMatrix[i] {
            newMatrix[i][j] = ' ' 
        }
    }
    dm.Width = width
    dm.Height = height
    dm.Matrix = newMatrix
}


func (dm *DisplayMatrix) AddBottomLeftMessage(message string) {
    lines := strings.Split(message, "\n")
    totalLines := len(lines)

    startY := dm.Height - totalLines
    startX := 0 

    for i, line := range lines {
        for x, char := range line {
            matrixY := startY + i
            matrixX := startX + x
            if matrixY >= 0 && matrixY < dm.Height && matrixX >= 0 && matrixX < dm.Width {
                dm.Matrix[matrixY][matrixX] = char
            }
            if matrixX >= dm.Width - 1 {
                break
            }
        }
        if startY+i >= dm.Height-1 {
            break
        }
    }
}