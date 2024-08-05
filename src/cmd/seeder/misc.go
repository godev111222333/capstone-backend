package seeder

import (
	"fmt"
	"time"
)

const (
	DateTimeLayout = "2006-01-02 15:04:05"
)

type DateTime struct {
	time.Time
}

func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Add(7 * time.Hour).Format(DateTimeLayout), nil
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse(DateTimeLayout, csv)
	date.Time = date.Time.Add(-7 * time.Hour)
	return err
}

func toFilePath(file string) string {
	return fmt.Sprintf("etc/seed/%s", file)
}
