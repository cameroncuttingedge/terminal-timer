package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"terminal-timer/art"
	"terminal-timer/display"
	"terminal-timer/random"
	"terminal-timer/sound"
	"terminal-timer/util"

	"github.com/mattn/go-tty"
)

var userDecision = make(chan string, 1)

// main initializes the application, parses flags, and starts the timer loop.
func main() {
	timerFlag, alarmFlag, reminderFlag, enableLogging := util.ParseFlags()
    
    random.GenerateTempFileName()

    if enableLogging {
        util.SetupLogger()
    }

	setupSignalHandling(cleanup)

    var directInput string
	if len(flag.Args()) > 0 {
		directInput = flag.Arg(0)
	}

	totalSeconds, err := util.CalculateTotalSeconds(timerFlag, alarmFlag, directInput)
    if err != nil {
        fmt.Println("Error parsing timer or alarm flag:", err)
        return
    }

	reminder := util.GetReminderMessage(reminderFlag)
	runTimerLoop(totalSeconds, reminder)
    defer cleanup()
}

// runTimerLoop runs the main timer loop, displaying time and handling user input.
func runTimerLoop(totalSeconds int, reminder string) {
    title := "Timer Completed"
    soundPath := "Jinja.wav"
    for {
        util.HideCursor()
        util.Render()

        width, height, err := util.GetSize()
        if err != nil {
            fmt.Println("Error getting terminal size:", err)
            return
        }

        matrix := display.NewDisplayMatrix(width, height)
        startTimer(totalSeconds, "", matrix)
        display.BufferEndMessage(matrix, reminder, "")

        matrix.Print()
        sound.EndOfTimer(soundPath, title, reminder)

        if !waitForUserInput(matrix, reminder, "") {
            break
        }
    }
}

// waitForUserInput waits for user input to restart or quit the timer.
func waitForUserInput(matrix *display.DisplayMatrix, reminder, font string) bool {
    tty, err := tty.Open()
    if err != nil {
        log.Fatalf("failed to open tty: %v", err)
    }
    defer tty.Close()

    go func() {
        for {
            r, err := tty.ReadRune()
            if err != nil {
                fmt.Printf("Error reading rune: %v", err)
                continue
            }
            switch r {
            case 'q', 'r':
                userDecision <- string(r)
                return
            }
        }
    }()

    for {
        select {
        case decision := <-userDecision:
            util.ShowCursor()
            return decision == "r"
        default:
            matrix.ResizeAndClear()
            display.BufferEndMessage(matrix, reminder, font)
            matrix.Print()
            time.Sleep(100 * time.Millisecond)
        }
    }
}

// startTimer counts down the timer and updates the display.
func startTimer(totalSeconds int, font string, matrix *display.DisplayMatrix) {
    endTime := time.Now().Add(time.Duration(totalSeconds) * time.Second)
	for range time.Tick(time.Second) {
		remaining := time.Until(endTime)
		if remaining <= 0 {
			break
		}

		message := fmt.Sprintf("%02d:%02d:%02d", int(remaining.Hours()), int(remaining.Minutes())%60, int(remaining.Seconds())%60)
		asciiArt := art.GetAsciiArt(message, font)
		matrix.AddCenteredAsciiArt(asciiArt, message)
		matrix.Print()
		matrix.ResizeAndClear()
	}
}

// setupSignalHandling configures handling for SIGINT and SIGTERM.
func setupSignalHandling(cleanupFunc func()) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-c 
        cleanupFunc()
        fmt.Println("\nReceived Ctrl+C, exiting...")
        os.Exit(0)
    }()
}

// cleanup performs application cleanup tasks.
func cleanup() {
    util.ShowCursor()
	util.Clear()
    util.Render()
    if random.TempFileName != "" {
        err := os.Remove(random.TempFileName)
        if err != nil {
            log.Printf("Failed to delete temporary file %s: %v\n", random.TempFileName, err)
        } else {
            log.Printf("Temporary file %s deleted successfully\n", random.TempFileName)
        }
    }
}