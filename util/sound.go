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
	"errors"
	"os/exec"
	"runtime"

	"github.com/gen2brain/beeep"
)

func EndOfTimer(soundFilePath, title, message string) error {
    done := make(chan error, 2)

    go func() {
        err := PlaySound(soundFilePath)
        done <- err
    }()

    go func() {
        err := ShowNotification(title, message)
        done <- err
    }()

    for i := 0; i < 2; i++ {
        err := <-done
        if err != nil {
            return err
        }
    }

    return nil
}


func PlaySound(filePath string) error {
	var cmd *exec.Cmd
    switch runtime.GOOS {
    case "darwin": // macOS
        if CmdExists("afplay") {
            cmd = exec.Command("afplay", filePath)
        } else {
            return errors.New("no compatible media player found")
        }
    case "linux":
        // Prioritize media players by quality and availability
        if CmdExists("ffplay") {
            cmd = exec.Command("ffplay", "-nodisp", "-autoexit", filePath)
        } else if CmdExists("mpg123") {
            cmd = exec.Command("mpg123", filePath)
        } else if CmdExists("paplay") {
            cmd = exec.Command("paplay", filePath)
        } else if CmdExists("aplay") {
            cmd = exec.Command("aplay", filePath)
        } else {
            return errors.New("no compatible media player found")
        }
    case "windows":
        // Windows: attempt to use PowerShell as a more reliable method
        if CmdExists("powershell") {
            cmdStr := `$player = New-Object System.Media.SoundPlayer;` +
                `$player.SoundLocation = '` + filePath + `';` +
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

func ShowNotification(title, message string) error {
    // Customize the icon path according to your application's needs
    iconPath := "path/to/icon.png" // Make sure to adjust this path
    return beeep.Notify(title, message, iconPath)
}
