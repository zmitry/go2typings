package types

// TODO write testcase for this
// type M struct {
// 	Username string `json:"Username2"`
// }

type WeekDay int

const (
	SUNDAY WeekDay = iota + 1
	MONDAY
	MONDAY2
	MONDAY3
)

type T struct {
	W WeekDay `json:"weekday"`
}

// type UserTag struct {
// 	Tag string `json:"tag"`
// }

// type User struct {
// 	FirstName  string    `json:"firstname"`
// 	SecondName string    `json:"secondName"`
// 	Tags       []UserTag `json:"tags"`
// }
