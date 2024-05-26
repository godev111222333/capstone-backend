package misc

import (
	"time"
)

type DateTime struct {
	time.Time
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(time.DateTime)), nil
}

func (t *DateTime) UnmarshalJSON(b []byte) error {
	p, err := time.Parse(time.DateTime, string(b))
	if err != nil {
		return err
	}
	t.Time = p
	return nil
}
