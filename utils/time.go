package utils

import (
	"fmt"
	"time"
)

func ExpiresAtToStringTime(expiresAt int64) string {
	times := time.Unix(expiresAt, 0)
	return times.Format(time.RFC3339)
}

func CreationTimeFormatKST(creationTime string) (string, error) {
	utcTime, err := time.Parse(time.RFC3339, creationTime)
	if err != nil {
		return "", fmt.Errorf("failed to parse time - CreationTimeFormatKST: %w", err)
	}

	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return "", fmt.Errorf("failed to load location - CreationTimeFormatKST: %w", err)
	}

	kstTime := utcTime.In(location)
	return kstTime.Format(time.RFC3339), nil
}

func CountDays(expiresAt int64) (int, error) {
	expireTimes := time.Unix(expiresAt, 0)

	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return 0, fmt.Errorf("failed to load location: %w", err)
	}

	now := time.Now().In(location)
	duration := expireTimes.Sub(now)
	days := int(duration.Hours()/24) + 1

	return days, nil
}
