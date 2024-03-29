package util

import (
	"fmt"
	"strings"
	"time"
)

func CalculateTotalSeconds(timerDuration string, alarmTime string, directInput string) (int, error) {
	// Parse timer duration
	if timerDuration != "" {
		seconds, err := ParseDuration(timerDuration)
		if err != nil {
			return 0, fmt.Errorf("invalid timer duration: %w", err)
		}
		return seconds, err
	}

	// Parse alarm time
	if alarmTime != "" {
		secondsUntilAlarm, err := ParseAlarm(alarmTime)
		if err != nil {
			return 0, fmt.Errorf("invalid alarm time: %w", err)
		}
		return secondsUntilAlarm, err
	}

	if directInput != "" {
		seconds, err := ParseDuration(directInput)
		if err == nil {
			return seconds, nil
		}
		return seconds, err
	}
	return 3, nil
}

func ParseDuration(durationStr string) (int, error) {
	parts := strings.Split(durationStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format")
	}
	hours, err := time.ParseDuration(parts[0] + "h")
	if err != nil {
		return 0, err
	}
	minutes, err := time.ParseDuration(parts[1] + "m")
	if err != nil {
		return 0, err
	}
	return int(hours.Seconds() + minutes.Seconds()), nil
}

func ParseAlarm(alarmStr string) (int, error) {
	now := time.Now()
	alarmTime, err := time.Parse("15:04", alarmStr)
	if err != nil {
		return 0, err
	}
	alarmTime = time.Date(now.Year(), now.Month(), now.Day(), alarmTime.Hour(), alarmTime.Minute(), 0, 0, now.Location())
	if alarmTime.Before(now) {
		alarmTime = alarmTime.Add(24 * time.Hour)
	}
	return int(alarmTime.Sub(now).Seconds()), nil
}

func GetReminderMessage(reminderFlag string) string {
	if reminderFlag == "" {
		return "Time is Up!"
	}
	return reminderFlag
}
