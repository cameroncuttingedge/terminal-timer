// package util

// import (
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/ebitengine/oto/v3"
// 	"github.com/hajimehoshi/go-mp3"
// )

// func EndOfTimer() {
// 	filePath := "/home/cameron/personal/terminal-timer/util/public_sounds_boop.mp3"
//     // Open the MP3 file
//     file, err := os.Open(filePath)
//     if err != nil {
//         log.Fatalf("Failed to open MP3 file: %v", err)
//     }
//     defer file.Close()

//     // Decode the MP3 file
//     decoded, err := mp3.NewDecoder(file)
//     if err != nil {
//         log.Fatalf("Failed to decode MP3: %v", err)
//     }

//     // Initialize oto context with the decoded MP3 format
//     ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
//         SampleRate:       decoded.SampleRate(),
//         ChannelCount:     2, // Use 2 for stereo sound
//         Format:           oto.FormatSignedInt16LE, // go-mp3's format is signed 16bit
//     })
//     if err != nil {
//         log.Fatalf("Failed to initialize audio context: %v", err)
//     }
//     <-ready // Wait for the audio hardware to be ready

//     // Create a new player from the oto context
//     player := ctx.NewPlayer(decoded)

//     // Start playback
//     player.Play()

//     // Wait for the sound to finish playing
//     for player.IsPlaying() {
//         time.Sleep(100 * time.Millisecond)
//     }

//     // Clean up resources
//     if err := player.Close(); err != nil {
//         log.Fatalf("Failed to close player: %v", err)
//     }
// 	// Notify the user that the timer is done
// 	err = beeep.Notify("Timer Completed", reminder, "") // Path to an icon can be added as the third parameter
// 	if err != nil {
// 		log.Fatalf("Failed to send notification: %v", err)
// 		}
// }

package util

import (
	"embed"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/gen2brain/beeep"
)

//go:embed WAV/*
var wavFS embed.FS

//go:embed WAV/clock.png
var clockPNG embed.FS

func EndOfTimer(soundFilePath, title, message string) {
    // Play the end of timer sound in a non-blocking way
    
    go func() {
        tmpFileName, err := PrepareSoundFile(soundFilePath)
        if err != nil {
            log.Printf("Error preparing sound file: %v", err)
            return
        }
        err = ExecuteSoundPlayback(tmpFileName)
        if err != nil {
            log.Printf("Error playing sound: %v", err)
        }
    }()

    go func() {
        err := ShowNotification(title, message)
        if err != nil {
            log.Printf("Error showing notification: %v", err)
        }
    }()
}


func ExecuteSoundPlayback(tmpFileName string) error {
    var cmd *exec.Cmd

    switch runtime.GOOS {
    case "darwin":
        if CmdExists("afplay") {
            cmd = exec.Command("afplay", tmpFileName)
        } else {
            return errors.New("no compatible media player found")
        }
    case "linux":
        if CmdExists("ffplay") {
            cmd = exec.Command("ffplay", "-nodisp", "-autoexit", tmpFileName)
        } else if CmdExists("mpg123") {
            cmd = exec.Command("mpg123", tmpFileName)
        } else if CmdExists("paplay") {
            cmd = exec.Command("paplay", tmpFileName)
        } else if CmdExists("aplay") {
            cmd = exec.Command("aplay", tmpFileName)
        } else {
            return errors.New("no compatible media player found")
        }
    case "windows":
        if CmdExists("powershell") {
            cmdStr := `$player = New-Object System.Media.SoundPlayer;` +
                `$player.SoundLocation = '` + tmpFileName + `';` +
                `$player.PlaySync();`
            cmd = exec.Command("powershell", "-Command", cmdStr)
        } else {
            return errors.New("no compatible media player found")
        }
    default:
        return errors.New("unsupported platform")
    }

    return cmd.Run()
}


func PrepareSoundFile(filePath string) (string, error) {
    soundFilePath := "WAV/" + filePath

    soundFile, err := wavFS.Open(soundFilePath)
    if err != nil {
        log.Printf("Error opening embedded sound file '%s': %v", soundFilePath, err)
        return "", errors.New("failed to open embedded sound file")
    }
    defer soundFile.Close()

    tmpFile, err := os.CreateTemp("", "sound-*.wav")
    if err != nil {
        log.Printf("Error creating temporary file for sound: %v", err)
        return "", errors.New("failed to create temporary file for sound")
    }
    tmpFileName := tmpFile.Name()

    if _, err = io.Copy(tmpFile, soundFile); err != nil {
        tmpFile.Close()
        os.Remove(tmpFileName) // Clean up even in case of error
        log.Printf("Error copying sound file to temporary file '%s': %v", tmpFileName, err)
        return "", errors.New("failed to copy sound file to temporary file")
    }

    if err := tmpFile.Close(); err != nil {
        os.Remove(tmpFileName) // Clean up even in case of error
        log.Printf("Error closing temporary sound file '%s': %v", tmpFileName, err)
        return "", errors.New("failed to close temporary sound file")
    }

    return tmpFileName, nil
}



func ShowNotification(title, message string) error {
    // Open the embedded clock.png
    clockFile, err := clockPNG.Open("WAV/clock.png")
    if err != nil {
        log.Printf("Error opening embedded image 'WAV/clock.png': %v", err)
        return err
    }
    defer clockFile.Close()

    // Create a temporary file for the clock image
    tmpFile, err := os.CreateTemp("", "clock-*.png")
    if err != nil {
        log.Printf("Error creating temporary file for image: %v", err)
        return err
    }
    defer tmpFile.Close()
    defer os.Remove(tmpFile.Name()) // Clean up the temp file after use

    // Copy the embedded clock image content to the temporary file
    _, err = io.Copy(tmpFile, clockFile)
    if err != nil {
        log.Printf("Error copying image to temporary file '%s': %v", tmpFile.Name(), err)
        return err
    }

    // Use the path of the temp file for the icon in beeep.Notify
    iconPath := tmpFile.Name()
    err = beeep.Notify(title, message, iconPath)
    if err != nil {
        log.Printf("Error showing notification: %v", err)
        return err
    }
    return nil
}

