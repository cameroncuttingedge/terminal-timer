package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cameroncuttingedge/terminal-timer/alert"
	"github.com/cameroncuttingedge/terminal-timer/display"
	"github.com/cameroncuttingedge/terminal-timer/util"
)

type ValidItemChecker func() ([]string, error)

func CheckIfConfigChangesRequested() {
    var shouldExit bool

    // Setting a new font
    if *util.SetFontFlag != "" {
        updateConfiguration("font", *util.SetFontFlag, listValidFonts)
        shouldExit = true
    }

    // Setting a new sound
    if *util.SetSoundFlag != "" && !shouldExit {
        updateConfiguration("sound", *util.SetSoundFlag, alert.ListValidSounds)
        shouldExit = true
    }

    if *util.ListValidFonts && !shouldExit {
        fmt.Println("Listing all valid fonts...")
        fonts, err := listValidFonts()
        if err != nil {
            log.Printf("Error getting fonts: %v", err)
        } else {
            display.RenderPrettyData(fonts)
        }
        shouldExit = true
    }

    if *util.ListValidSounds && !shouldExit {
        fmt.Println("Listing all valid sounds...")
        sounds, err := alert.ListValidSounds()
        if err != nil {
            log.Printf("Error getting sounds: %v", err)
        } else {
            display.RenderPrettyData(sounds)
        }
        shouldExit = true
    }

    if *util.PreviewFontFlag != "" && !shouldExit {
        display.RenderFontExample(*util.PreviewFontFlag)
        shouldExit = true
    }

    if *util.PreviewSoundFlag != "" && !shouldExit {
        tmpfile, err := alert.PrepareSoundFile(*util.PreviewSoundFlag)
        if err != nil {
            log.Printf("Error preparing sound file: %v", err)
        } else {
            alert.ExecuteSoundPlayback(tmpfile)
        }
        shouldExit = true
    }

    if *util.ShowCurrentConfig && !shouldExit {
        PrintCurrentConfig()
        shouldExit = true
    }

    if shouldExit {
        util.Cleanup(false)
        os.Exit(0) 
    }
}

func listValidFonts() ([]string, error) {
	file, err := fontsFile.Open("check/fonts.txt")
	if err != nil {
		fmt.Printf("Error opening embedded fonts.txt: %v\n", err)
		return nil, err
	}
	defer file.Close()

	var fonts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fonts = append(fonts, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading embedded fonts.txt: %v\n", err)
		return nil, err
	}

	return fonts, nil
}

func updateConfiguration(configType, newValue string, checker ValidItemChecker) {
	fmt.Printf("Setting new %s to: %s\n", configType, newValue)

	// Use IsValidItem to validate the new value
	if !isValidItem(newValue, checker) {
		return
	}

	err := UpdateConfig(configType, newValue)
	if err != nil {
		fmt.Printf("Error updating %s configuration: %v\n", configType, err)
		os.Exit(1)
	} else {
		fmt.Printf("%s configuration updated successfully.\n", configType)
		os.Exit(0)
	}
}

func isValidItem(newValue string, checker ValidItemChecker) bool {

	//log.Printf("Entering isValidItem function with newValue:", newValue)

	validItems, err := checker()
	if err != nil {
		fmt.Printf("Error fetching list of valid items: %v\n", err)
		os.Exit(2) // Exit with an error code indicating failure to fetch valid items
	}

	//log.Printf("Valid items fetched:", validItems)

	for _, item := range validItems {
		if item == newValue {
			return true
		}
	}

	fmt.Printf("Invalid item specified. Please choose a valid item. Use -h to see how you can list valid items.\n")
	os.Exit(2)
	return false
}
