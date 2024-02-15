package util

import (
	"flag"
)

var (
    // Timer and alarm options
    TimerFlag     = flag.String("t", "", "Duration in hh:mm format")
    AlarmFlag     = flag.String("a", "", "Alarm time in 24-hour format hh:mm")
    ReminderFlag  = flag.String("r", "Time is Up!", "Reminder message")

    // Font options
    SetFontFlag      = flag.String("f", "", "Set a new font")
    PreviewFontFlag  = flag.String("pf", "", "Preview the font")
    ListValidFonts   = flag.Bool("lf", false, "List all valid fonts")

    // Sound options
    PreviewSoundFlag = flag.String("ps", "", "Preview the sound")
    ListValidSounds  = flag.Bool("ls", false, "List all default sounds")
    SetSoundFlag     = flag.String("s", "", "Set a new sound")

    // Logging options
    EnableLogging = flag.Bool("log", false, "Enable logging to a file")
    ShowCurrentConfig = flag.Bool("c", false, "Show current sound and font config")
)


func ParseFlags() {
    flag.Parse()
}
