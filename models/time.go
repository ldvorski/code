package models

import "fmt"

type Time struct {
	Hour   int
	Minute int
}

func NewTime(timeStr string) (Time, error) {
	var t Time
	_, err := fmt.Sscanf(timeStr, "%d:%d", &t.Hour, &t.Minute)
	if err != nil {
		return t, err
	}
	return t, nil
}
