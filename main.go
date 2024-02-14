package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	displays "terminal-timer/display"
	"terminal-timer/util"

	"github.com/mattn/go-tty"
)


func init() {
    logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal("Error opening log file:", err)
    }

    log.SetOutput(logFile)
}

func main() {

	timerFlag, alarmFlag, reminderFlag := parseFlags()

	cleanupFunc := func() {
		cleanup()
	}
	setupSignalHandling(cleanupFunc)

	totalSeconds, err := util.CalculateTotalSeconds(timerFlag, alarmFlag)

    if err != nil {
        fmt.Println("Error parsing timer or alarm flag:", err)
        return
    }

	reminder := util.GetReminderMessage(reminderFlag)

	runTimerLoop(totalSeconds, reminder)

    defer cleanup()
}


func parseFlags() (timerFlag string, alarmFlag string, reminderFlag string) {
	timer := flag.String("t", "", "Duration in hh:mm format")
	alarm := flag.String("a", "", "Alarm time in 24-hour format hh:mm")
	reminder := flag.String("r", "Time is Up!", "Reminder message")
	flag.Parse()
	return *timer, *alarm, *reminder
}

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

        matrix := displays.NewDisplayMatrix(width, height)

        font := ""

        // Display message when timer ends
        
		startTimer(totalSeconds, font, matrix)
        
        bufferEndMessage(matrix, reminder, font)

        matrix.Print()

        util.EndOfTimer(soundPath, title, reminder)


        shouldRestart := waitForUserInput(matrix, reminder, font)
        
        if !shouldRestart {
            break
        }
    }
}


var userDecision = make(chan string, 1)

func waitForUserInput(matrix *displays.DisplayMatrix, reminder string, font string) bool {
    tty, err := tty.Open()
    if err != nil {
        log.Fatalf("failed to open tty: %v", err)
    }
    defer tty.Close()

    // Start a goroutine to read user input.
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
                return // Exit the goroutine after decision is made.
            }
        }
    }()

    // Continuously refresh the screen until a decision is made.
    for {
        select {
        case decision := <-userDecision:
            util.ShowCursor() // Show cursor before exiting.
            if decision == "q" {
                return false
            } else if decision == "r" {
                return true 
            }
        default:
            // No decision made yet, refresh.
            matrix.ResizeAndClear()
            bufferEndMessage(matrix, reminder, font)
            matrix.Print()

            time.Sleep(100 * time.Millisecond) // save the cpu
        }
    }
}


func bufferEndMessage(matrix *displays.DisplayMatrix, reminder string, font string) {
    matrix.ResizeAndClear()

    timeUpMessage := util.GetAsciiArt(reminder, font)
    matrix.AddCenteredAsciiArt(timeUpMessage, reminder)
    message := "Press 'q' to quit or 'r' to repeat."
    matrix.AddBottomLeftMessage(message)
}


func startTimer(totalSeconds int, font string, matrix *displays.DisplayMatrix) {
    endTime := time.Now().Add(time.Duration(totalSeconds) * time.Second)
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    updateTimerDisplay(endTime, font, matrix) 
    for range ticker.C {
        if time.Until(endTime) <= 0 {
            break
        }
        updateTimerDisplay(endTime, font, matrix)
    }
}

func updateTimerDisplay(endTime time.Time, font string, matrix *displays.DisplayMatrix) {
    remaining := time.Until(endTime)
    message := fmt.Sprintf("%02d:%02d:%02d", int(remaining.Hours()), int(remaining.Minutes())%60, int(remaining.Seconds())%60)
    asciiArt := util.GetAsciiArt(message, font)

    matrix.ResizeAndClear()
    matrix.AddCenteredAsciiArt(asciiArt, message)
    matrix.Print()
}

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


func cleanup() {
    util.ShowCursor()
	util.ClearTerminal()
    util.Render()
}


