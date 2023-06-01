package models

import "fmt"

type Date struct {
	Day   int
	Month int
	Year  int
}

func NewDate(dateStr string) (Date, error) {
	var d Date
	_, err := fmt.Sscanf(dateStr, "%d-%d-%d", &d.Day, &d.Month, &d.Year)
	if err != nil {
		return d, err
	}
	return d, nil
}
