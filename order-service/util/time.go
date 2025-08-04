package util

import (
	"fmt"
	"time"
)

type CustomTime time.Time

const layout = "2006-01-02 15:04:05"

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(ct).Format(layout))
	return []byte(formatted), nil
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(`"`+layout+`"`, string(data))
	if err != nil {
		return err
	}
	*ct = CustomTime(parsed)
	return nil
}

func (ct CustomTime) String() string {
	return time.Time(ct).Format(layout)
}
