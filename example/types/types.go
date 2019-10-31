package types

import "fmt"

// TODO write testcase for this
type M struct {
	Username string `json:"Username2"`
}

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
}

type UserTag struct {
	Tag string `json:"tag"`
}

type User struct {
	FirstName  string    `json:"firstname"`
	SecondName string    `json:"secondName"`
	Tags       []UserTag `json:"tags"`
}
