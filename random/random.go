package random

import (
	"os"
	"path/filepath"
	"time"
)

var TempFileName string

func GenerateTempSoundFileName() {
	timestamp := time.Now().Format("20060102-150405")
	tempFileName := "sound-" + timestamp + ".wav"
	TempFileName = filepath.Join(os.TempDir(), tempFileName)
}
