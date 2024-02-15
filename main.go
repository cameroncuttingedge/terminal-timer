package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"terminal-timer/alert"
	"terminal-timer/art"
	"terminal-timer/config"
	"terminal-timer/display"
	"terminal-timer/random"
	"terminal-timer/util"

	"github.com/mattn/go-tty"
)

var userDecision = make(chan string, 1)

// main initializes the application, parses flags, and starts the timer loop.
func main() {
    
    util.ParseFlags()

    fmt.Printf(*util.PreviewSoundFlag)

    config.LoadOrCreateConfig()

    random.GenerateTempFileName()

    config.CheckIfconfigChangesRequested()

    if *util.EnableLogging {
        util.SetupLogger()
    }

	setupSignalHandling(util.Cleanup)

    var directInput string
	if len(flag.Args()) > 0 {
		directInput = flag.Arg(0)
	}

	totalSeconds, err := util.CalculateTotalSeconds(*util.TimerFlag, *util.AlarmFlag, directInput)
    if err != nil {
        fmt.Println("Error parsing timer or alarm flag:", err)
        return
    }

	reminder := util.GetReminderMessage(*util.ReminderFlag)
	runTimerLoop(totalSeconds, reminder)
    defer util.Cleanup()
}

// runTimerLoop runs the main timer loop, displaying time and handling user input.
func runTimerLoop(totalSeconds int, reminder string) {
    title := "Timer Completed"
    soundPath := config.Sound
    for {
        util.HideCursor()
        util.Render()

        width, height, err := util.GetSize()
        if err != nil {
            fmt.Println("Error getting terminal size:", err)
            return
        }

        matrix := display.NewDisplayMatrix(width, height)
        startTimer(totalSeconds, matrix)
        display.BufferEndMessage(matrix, reminder, config.Font)

        matrix.Print()
        alert.EndOfTimer(soundPath, title, reminder)

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
            display.BufferEndMessage(matrix, reminder, config.Font)
            matrix.Print()
            time.Sleep(100 * time.Millisecond)
        }
    }
}

// startTimer counts down the timer and updates the display.
func startTimer(totalSeconds int, matrix *display.DisplayMatrix) {
    endTime := time.Now().Add(time.Duration(totalSeconds) * time.Second)
    
    // Call the update function once before entering the loop to display the first tick
    updateTimerDisplay(endTime, matrix)

    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for range ticker.C {
        if time.Now().After(endTime) {
            break
        }
        updateTimerDisplay(endTime, matrix)
    }
}

// updateTimerDisplay handles updating and printing the timer to the display matrix
func updateTimerDisplay(endTime time.Time, matrix *display.DisplayMatrix) {
    remaining := time.Until(endTime)
    timerRemaining := fmt.Sprintf("%02d:%02d:%02d", int(remaining.Hours()), int(remaining.Minutes())%60, int(remaining.Seconds())%60)
    asciiArt := art.GetAsciiArt(timerRemaining, config.Font)
    matrix.AddCenteredAsciiArt(asciiArt, timerRemaining)
    matrix.Print()
    matrix.ResizeAndClear()
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