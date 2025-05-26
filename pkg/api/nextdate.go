package api

import (
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	now := r.FormValue("now")
	var timeNow time.Time
	var err error
	if now == "" {
		timeNow = time.Now()
	} else {
		timeNow, err = time.Parse(DateFormat, now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(timeNow, date, repeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%s", nextDate)

}

func afterNow(date, now time.Time) bool {
	return date.Truncate(24 * time.Hour).After(now.Truncate(24 * time.Hour))
}

func dayCheck(idx []int, day int) bool {
	sort.Ints(idx)
	return slices.Contains(idx, day)
}

func monthCheck(idx []int, month int) bool {
	sort.Ints(idx)
	return slices.Contains(idx, month)
}

func lastDayInMonth(date time.Time) int {
	nextMonth := date.AddDate(0, 1, 0)                   // добавляем один месяц
	lastDay := nextMonth.AddDate(0, 0, -nextMonth.Day()) // получаем первый день следующего месяца и возвращаемся обратно на последний день текущего месяца
	return lastDay.Day()
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {

		return "", fmt.Errorf("rule should not be empty")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date format: %v", err)
	}

	ruleParts := strings.Split(repeat, " ")

	rule := ruleParts[0]

	switch rule {
	case "d":
		if len(ruleParts) != 2 {
			return "", fmt.Errorf("invalid rule format: %v", err)
		}
		interval, err := strconv.Atoi(ruleParts[1])
		if err != nil {
			return "", fmt.Errorf("failed to convert interval to number: %w", err)
		}

		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("invalid interval")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		if len(ruleParts) != 1 {
			return "", fmt.Errorf("invalid rule format: %v", err)
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w":

		if len(ruleParts) < 2 {
			return "", fmt.Errorf("invalid rule format: %v", err)
		}

		daysStr := strings.Split(ruleParts[1], ",")

		var daysIdx []int
		for _, d := range daysStr {
			i, _ := strconv.Atoi(d)
			if i < 1 || i > 7 {
				return "", fmt.Errorf("invalid day: %v", err)
			}
			daysIdx = append(daysIdx, i)
		}

		if !afterNow(date, now) {
			date = now
		}

		for {
			date = date.AddDate(0, 0, 1)
			weekdayNum := int(date.Weekday())
			if weekdayNum == 0 {
				weekdayNum = 7
			}
			if dayCheck(daysIdx, weekdayNum) {
				break
			}
		}

	case "m":

		if len(ruleParts) < 2 {
			return "", fmt.Errorf("invalid rule format: %v", err)
		}

		daysStr := strings.Split(ruleParts[1], ",")
		var daysIdx []int
		for _, d := range daysStr {
			i, err := strconv.Atoi(d)
			if err != nil {
				return "", fmt.Errorf("failed to convert day '%s': %v", d, err)
			}
			if i < -2 || i > 31 {
				return "", fmt.Errorf("дinvalid day: %v", err)
			}
			daysIdx = append(daysIdx, i)
		}

		if !afterNow(date, now) {
			date = now
		}
		for i, d := range daysIdx {
			if d == -1 {
				daysIdx[i] = lastDayInMonth(date)
			}
			if d == -2 {
				daysIdx[i] = lastDayInMonth(date) - 1

			}

		}

		if len(ruleParts) == 2 {

			for {
				date = date.AddDate(0, 0, 1)
				if dayCheck(daysIdx, date.Day()) {
					break
				}
			}

		}

		if len(ruleParts) == 3 {

			monthStr := strings.Split(ruleParts[2], ",")
			var monthIdx []int
			for _, d := range monthStr {
				i, err := strconv.Atoi(d)
				if err != nil {
					return "", fmt.Errorf("failed to convert month '%s': %v", d, err)
				}
				if i < 1 || i > 12 {
					return "", fmt.Errorf("invalid month: %v", err)
				}
				monthIdx = append(monthIdx, i)
			}

			for {
				date = date.AddDate(0, 0, 1)
				if dayCheck(daysIdx, date.Day()) && monthCheck(monthIdx, int(date.Month())) {
					break
				}
			}

		}

	default:
		return "", fmt.Errorf("invalid rule format: %v", err)
	}
	return date.Format(DateFormat), err

}
