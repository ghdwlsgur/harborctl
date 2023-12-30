package utils

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type Expiration int

const (
	Expired      Expiration = 0
	NeverExpired Expiration = -1

	location            = "Asia/Seoul"
	ExpiredMessage      = "Expired"
	NeverExpiredMessage = "Never Expires"
)

type Count struct {
	LeftDays int
	Err      error
}

func (d *Count) Validate() string {
	if d.Err != nil {
		return d.Err.Error()
	}

	switch Expiration(d.LeftDays) {
	case Expired:
		return ExpiredMessage
	case NeverExpired:
		return NeverExpiredMessage
	}

	return strconv.Itoa(d.LeftDays)
}

func ExpiresAtToStringTime(expiresAt int64) string {
	if Expiration(expiresAt) == NeverExpired {
		return NeverExpiredMessage
	}
	return time.Unix(expiresAt, 0).Format(time.RFC3339)
}

func CreationTimeFormatKST(creationTime string) (string, error) {
	utcTime, err := time.Parse(time.RFC3339, creationTime)
	if err != nil {
		return "", fmt.Errorf("time.Parse - CreationTimeFormatKST: %w", err)
	}

	location, err := time.LoadLocation(location)
	if err != nil {
		return "", fmt.Errorf("time.LoadLocation - CreationTimeFormatKST: %w", err)
	}

	kstTime := utcTime.In(location)
	return kstTime.Format(time.RFC3339), nil
}

func CountDays(expiresAt int64) *Count {
	if expiresAt < 0 {
		return &Count{LeftDays: int(NeverExpired), Err: nil}
	}

	expireTimes := time.Unix(expiresAt, 0)
	location, err := time.LoadLocation(location)
	if err != nil {
		return &Count{
			0, fmt.Errorf("time.LoadLocation - CountDays: %w", err),
		}
	}

	now := time.Now().In(location)
	duration := expireTimes.Sub(now)

	if days := int(math.Ceil(duration.Hours() / 24)); days < 0 {
		return &Count{LeftDays: int(Expired), Err: nil}
	} else {
		return &Count{LeftDays: days, Err: nil}
	}
}
