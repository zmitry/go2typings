package types

import "fmt"

// WeekDay can be converted to string.
type WeekDay int

const (
	SUNDAY WeekDay = 1
	MONDAY WeekDay = 2
)

type WeekDay2 string

const (
	A  WeekDay2 = "3"
	A2 WeekDay2 = "4"
)

// WeekDay3 has no String method.
type WeekDay3 int

const (
	TUESDAY   WeekDay3 = 5
	WEDNESDAY WeekDay3 = 6
)

func (e WeekDay) String() string {
	switch e {
	case SUNDAY:
		return "sun"
	case MONDAY:
		return "mon"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

type T struct {
	W  WeekDay  `json:"weekday"`
	W2 WeekDay2 `json:"weekday2"`
	W3 WeekDay3 `json:"weekday3"`
}
