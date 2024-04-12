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
}

type Plan struct {
	PlanID      string `json:"planID"`
	Name        string `json:"name"`
	PeriodEndAt int64  `json:"periodEndAt"`
	Active      bool   `json:"active"`
	Require     string `json:"require"`
}

type Message struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"code"`
}

type Error struct {
	Error Message `json:"error"`
}
