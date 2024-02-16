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

func CheckIfconfigChangesRequested() {

	// Setting a new font
	if *util.SetFontFlag != "" {
		updateConfiguration("font", *util.SetFontFlag, listValidFonts)
		return
	}

	// Setting a new sound
	if *util.SetSoundFlag != "" {
		updateConfiguration("sound", *util.SetSoundFlag, alert.ListValidSounds)
	}

	if *util.ListValidFonts {
		fmt.Println("Listing all valid fonts...")
		sounds, err := listValidFonts()
		if err != nil {
			log.Printf("Error getting fonts")
			fmt.Println("Error getting gonts")
			return
		}
		display.RenderPrettyData(sounds)
		time.Sleep(1 * time.Minute)
		return
	}

	if *util.ListValidSounds {
		fmt.Println("Listing all valid sounds...")
		sounds, err := alert.ListValidSounds()
		if err != nil {
			log.Printf("Error getting sounds")
			fmt.Println("Error getting sounds")
			return
		}
		display.RenderPrettyData(sounds)
		time.Sleep(1 * time.Minute)
		return
	}

	// Handle previewing the font
	if *util.PreviewFontFlag != "" {
		if !isValidItem(*util.PreviewFontFlag, listValidFonts) {
			return
		}
		display.RenderFontExample(*util.PreviewFontFlag)
		os.Exit(1)
	}

	if *util.PreviewSoundFlag != "" {
		if !isValidItem(*util.PreviewSoundFlag, alert.ListValidSounds) {
			return
		}

		tmpfile, err := alert.PrepareSoundFile(*util.PreviewSoundFlag)

		if err != nil {
			log.Printf("Error getting sounds")
			fmt.Println("Error getting sounds")
			return
		}

		alert.ExecuteSoundPlayback(tmpfile)
		util.Cleanup()
		os.Exit(1)
	}

	if *util.ShowCurrentConfig {
		PrintCurrentConfig()
		os.Exit(1)
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
