package utils

import (
	"fmt"
	"strconv"
	"time"
)

type Count struct {
	LeftDays int
	Err      error
}

func (d *Count) Validate() string {
	if d.Err != nil {
		return d.Err.Error()
	}

	switch d.LeftDays {
	case 0:
		return "Expired"
	case -1:
		return "Never Expires"
	}

	return strconv.Itoa(d.LeftDays)
}

func ExpiresAtToStringTime(expiresAt int64) string {
	if expiresAt == -1 {
		return "Never Expires"
	}
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

func CountDays(expiresAt int64) *Count {
	expireTimes := time.Unix(expiresAt, 0)

	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return &Count{
			0, fmt.Errorf("failed to load location: %w", err),
		}
	}

	now := time.Now().In(location)
	duration := expireTimes.Sub(now)
	days := int(duration.Hours() / 24)

	if days < 0 {
		days = -1
	}

	return &Count{LeftDays: days, Err: nil}
}
