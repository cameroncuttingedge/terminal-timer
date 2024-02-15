package util

import (
	"log"
	"os"
)

func SetupLogger() {
    logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal("Error opening log file:", err)
    }

    log.SetOutput(logFile)
    log.SetFlags(0)
}
