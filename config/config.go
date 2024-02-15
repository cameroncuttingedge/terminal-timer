package config

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
)

//go:embed check/fonts.txt
var fontsFile embed.FS


var (
    Font  string
    Sound string
)

const (
    defaultFont  = ""          
    defaultSound = "Beeper.wav"
)



func LoadOrCreateConfig() error {
    configPath := getConfigFilePath()

    // Check if the config file exists
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        // File does not exist, create it with default values
        if err := SaveConfig(defaultFont, defaultSound, configPath); err != nil {
            return fmt.Errorf("failed to create default config file: %v", err)
        }
    }

    // Load the configuration from the file
    if err := LoadConfig(configPath); err != nil {
        return fmt.Errorf("failed to load config: %v", err)
    }

    return nil
}


func LoadConfig(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.SplitN(line, "=", 2)
        if len(parts) == 2 {
            key := parts[0]
            value := parts[1]
            switch key {
            case "font":
                Font = value
            case "sound":
                Sound = value
            }
        }
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    return nil
}

// SaveConfig writes the provided font and sound values to the specified file path.
func SaveConfig(font, sound, filePath string) error {
    content := fmt.Sprintf("font=%s\nsound=%s", font, sound)
    if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
        return err
    }
    return os.WriteFile(filePath, []byte(content), 0644)
}


func getConfigFilePath() string {
    homeDir, err := homedir.Dir()
    if err != nil {
        panic(err) 
    }

    filename := "timer-config.txt"
    var dbPath string

    appName := "timer"

    switch runtime.GOOS {
    case "windows":
        appDataDir := os.Getenv("LOCALAPPDATA")
        if appDataDir == "" {
            appDataDir = filepath.Join(homeDir, "AppData", "Local")
        }
        dbPath = filepath.Join(appDataDir, appName, filename)

    case "darwin":
        dbPath = filepath.Join(homeDir, "Library", "Application Support", appName, filename)

    default:
        xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
        if xdgConfigHome == "" {
            //if not xdgConfigHome the use /home/$USER/.confg 
            xdgConfigHome = filepath.Join(homeDir, ".config")
        }
        dbPath = filepath.Join(xdgConfigHome, appName, filename)
    }

    if err := os.MkdirAll(filepath.Dir(dbPath), 0700); err != nil {
        panic(err)
    }

    return dbPath
}

// Assumes that the sound or font has been validated
func UpdateConfig(key, value string) error {
    configPath := getConfigFilePath()

    err := LoadConfig(configPath)
    if err != nil {
        return err
    }

    switch key {
    case "font":
        Font = value
    case "sound":
        Sound = value
    default:
        return fmt.Errorf("invalid configuration key: %s", key)
    }

    // Save the updated configuration
    return SaveConfig(Font, Sound, configPath)
}


func PrintCurrentConfig() {
    fmt.Println("Current Configuration:")
    fmt.Printf("Font: %s\n", Font)
    fmt.Printf("Sound: %s\n", Sound)
}