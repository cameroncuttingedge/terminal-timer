package util

import (
	"log"
	"os"

	"github.com/cameroncuttingedge/terminal-timer/random"
)

// cleanup performs application cleanup tasks.
func Cleanup() {
	ShowCursor()
	Clear()
	Render()
	if random.TempFileName != "" {
		err := os.Remove(random.TempFileName)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to delete temporary file %s: %v\n", random.TempFileName, err)
		}
	}
}
