package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func invalidDate(value string) error {
	return fmt.Errorf("Invalid date %q. Accepted formats: DD-MM-YYYY, YYYY-MM-DD, MM-YYYY, YYYY-MM, YYYY, with - or / as the separator", value)
}

func invalidTime(value string) error {
	return fmt.Errorf("Invalid time %q. Accepted formats: 7 PM, 07:30 PM, 19:30, 19:30:45", value)
}

func RecordDate(t time.Time) int {
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}

func RecordTime(t time.Time) int {
	return t.Hour()*10000 + t.Minute()*100 + t.Second()
}

func ParseDate(value string) (int, int, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, 0, invalidDate(value)
	}

	separator := ""
	switch {
	case strings.Contains(trimmed, "-"):
		separator = "-"
	case strings.Contains(trimmed, "/"):
		separator = "/"
	}

	if separator == "" {
		year, err := strconv.Atoi(trimmed)
		if err != nil || year < 1 || year > 9999 {
			return 0, 0, invalidDate(value)
		}
		return year * 10000, year*10000 + 9999, nil
	}

	parts := strings.Split(trimmed, separator)
	if len(parts) < 2 || len(parts) > 3 {
		return 0, 0, invalidDate(value)
	}

	numbers := make([]int, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		n, err := strconv.Atoi(part)
		if err != nil {
			return 0, 0, invalidDate(value)
		}
		numbers[i] = n
		parts[i] = part
	}

	last := len(parts) - 1
	yearFirst := len(parts[0]) == 4
	yearLast := len(parts[last]) == 4

	if yearFirst == yearLast {
		return 0, 0, invalidDate(value)
	}

	var year, month, day int
	if yearFirst {
		year, month = numbers[0], numbers[1]
		if len(numbers) == 3 {
			day = numbers[2]
		}
	} else {
		year = numbers[last]
		if len(numbers) == 3 {
			day, month = numbers[0], numbers[1]
		} else {
			month = numbers[0]
		}
	}

	if year < 1 || year > 9999 || month < 1 || month > 12 {
		return 0, 0, invalidDate(value)
	}

	if len(numbers) == 2 {
		base := year*10000 + month*100
		return base, base + 99, nil
	}

	if day < 1 || day > 31 {
		return 0, 0, invalidDate(value)
	}

	exact := year*10000 + month*100 + day
	return exact, exact, nil
}

func ParseTime(value string) (int, int, error) {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		return 0, 0, invalidTime(value)
	}

	isPM := strings.HasSuffix(trimmed, "pm")
	isAM := strings.HasSuffix(trimmed, "am")
	if isPM || isAM {
		trimmed = strings.TrimSpace(trimmed[:len(trimmed)-2])
	}

	parts := strings.Split(trimmed, ":")
	if len(parts) > 3 {
		return 0, 0, invalidTime(value)
	}

	numbers := make([]int, len(parts))
	for i, part := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return 0, 0, invalidTime(value)
		}
		numbers[i] = n
	}

	hour := numbers[0]
	if isPM || isAM {
		if hour < 1 || hour > 12 {
			return 0, 0, invalidTime(value)
		}
		if isPM && hour != 12 {
			hour += 12
		}
		if isAM && hour == 12 {
			hour = 0
		}
	} else if hour < 0 || hour > 23 {
		return 0, 0, invalidTime(value)
	}

	if len(numbers) == 1 {
		base := hour * 10000
		return base, base + 9999, nil
	}

	minute := numbers[1]
	if minute < 0 || minute > 59 {
		return 0, 0, invalidTime(value)
	}

	if len(numbers) == 2 {
		base := hour*10000 + minute*100
		return base, base + 99, nil
	}

	second := numbers[2]
	if second < 0 || second > 59 {
		return 0, 0, invalidTime(value)
	}

	exact := hour*10000 + minute*100 + second
	return exact, exact, nil
}
