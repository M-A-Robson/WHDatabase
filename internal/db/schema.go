package db

import (
	"fmt"
	"time"
)

type Model struct {
	ID           int    `json:"id"`
	Game         string `json:"game"`
	Faction      string `json:"faction"`
	UnitName     string `json:"unit_name"`
	UnitSize     int    `json:"unit_size"`
	Points       int    `json:"points"`
	PurchaseDate string `json:"purchase_date"`
	BuildDate    string `json:"build_date"`
	PaintedDate  string `json:"painted_date"`
	Image        []byte `json:"image"`
}

type Date struct {
	Day   int
	Month int
	Year  int
}

func DateToString(date Date) string {
	return fmt.Sprintf("%04d-%02d-%02d", date.Year, date.Month, date.Day)
}

func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return Date{Day: t.Day(), Month: int(t.Month()), Year: t.Year()}, nil
}

func Today() Date {
	now := time.Now()
	return Date{
		Day:   now.Day(),
		Month: int(now.Month()),
		Year:  now.Year(),
	}
}
