package display

import (
	"fmt"
	"math"
	"os"
	"strings"
	"terminal-timer/art"
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
    util.Clear()
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


func BufferEndMessage(matrix *DisplayMatrix, reminder string, font string) {
    matrix.ResizeAndClear()
    timeUpMessage := art.GetAsciiArt(reminder, font)
    matrix.AddCenteredAsciiArt(timeUpMessage, reminder)
    message := "Press 'q' to quit or 'r' to repeat."
    matrix.AddBottomLeftMessage(message)
}

func (dm *DisplayMatrix) PrintItemsInGrid(items []string, columns int) {
    // Calculate necessary dimensions
    rows := int(math.Ceil(float64(len(items)) / float64(columns)))
    maxItemLength := 0
    for _, item := range items {
        if len(item) > maxItemLength {
            maxItemLength = len(item)
        }
    }

    // Calculate starting positions
    totalWidth := maxItemLength*columns + (columns-1)*2 // Assuming 2 spaces as padding between columns
    startX := (dm.Width - totalWidth) / 2
    startY := (dm.Height - rows) / 2

    // Iterate over items and print them in the specified grid layout
    for i, item := range items {
        row := i / columns
        col := i % columns

        posX := startX + col*(maxItemLength+2) // +2 for padding between columns
        posY := startY + row

        if posY >= 0 && posY < dm.Height && posX >= 0 && posX+maxItemLength <= dm.Width {
            copy(dm.Matrix[posY][posX:posX+maxItemLength], []rune(fmt.Sprintf("%-*s", maxItemLength, item)))
        }
    }
}


func RenderPrettyData(fonts []string) {
    width, height, err := util.GetSize()
    if err != nil {
        fmt.Printf("Error getting terminal size: %v\n", err)
        return
    }
    matrix := NewDisplayMatrix(width, height)

    // Determine the number of columns based on the terminal width
    maxLen := 0
    for _, font := range fonts {
        if len(font) > maxLen {
            maxLen = len(font)
        }
    }
    columns := int(math.Floor(float64(width) / float64(maxLen + 2))) // 2 spaces padding

    // Use the DisplayMatrix to print the fonts in a grid
    matrix.PrintItemsInGrid(fonts, columns)
    matrix.Print()

    os.Exit(1)
}


func RenderFontExample (font string) {

    message := "420:69"

    width, height, err := util.GetSize()
    if err != nil {
        fmt.Printf("Error getting terminal size: %v\n", err)
        return
    }
    matrix := NewDisplayMatrix(width, height)

    art := art.GetAsciiArt(message, font)

    matrix.AddCenteredAsciiArt(art, message)

    matrix.Print()

    os.Exit(1)
}