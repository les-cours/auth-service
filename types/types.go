package types

type Object map[string]interface{}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Osname   string `json:"osname"`
}

type Permissions struct {
	WriteComment bool `json:"writeComment"`
	Live         bool `json:"live"`
	Upload       bool `json:"upload"`
}

type Message struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"code"`
}

type Error struct {
	Error Message `json:"error"`
}
