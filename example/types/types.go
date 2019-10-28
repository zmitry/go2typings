package types

// TODO write testcase for this
type M struct {
	Username string `json:"Username2"`
}
type T struct {
	M
	Name []map[string]struct {
		test string
	} `json:"name"`

	// Array<Record<string, string>>
	Lastname []map[string]string `json:"lastname"`
	Surname  []map[string][]*M   `json:"surname"`
}

type UserTag struct {
	Tag string `json:"tag"`
}

type User struct {
	FirstName  string    `json:"firstname"`
	SecondName string    `json:"secondName"`
	Tags       []UserTag `json:"tags"`
}
