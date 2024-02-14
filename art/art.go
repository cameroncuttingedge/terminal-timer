package art

import (
	"strings"

	"github.com/common-nighthawk/go-figure"
)

func GetAsciiArt(message string, font string) []string {
	mainAsciiArt := figure.NewFigure(message, font, true)
	return strings.Split(mainAsciiArt.String(), "\n")
}
