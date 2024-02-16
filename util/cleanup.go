package util

import (
	"log"
	"os"
	"terminal-timer/random"
)

// cleanup performs application cleanup tasks.
func Cleanup() {
	ShowCursor()
	Clear()
	Render()
	if random.TempFileName != "" {
		err := os.Remove(random.TempFileName)
		if err != nil {
			log.Printf("Failed to delete temporary file %s: %v\n", random.TempFileName, err)
		}
	}
}
