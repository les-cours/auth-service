package types

type Object map[string]interface{}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Osname   string `json:"osname"`
}

type Permissions struct {
	Bots     bool `json:"bots"`
	Triggers bool `json:"triggers"`
	Tickets  bool `json:"tickets"`
	Profiles bool `json:"profiles"`
	Kbas     bool `json:"kbas"`
	Settings bool `json:"settings"`
}

type Plan struct {
	PlanID      string `json:"planID"`
	Name        string `json:"name"`
	PeriodEndAt int64  `json:"periodEndAt"`
	Active      bool   `json:"active"`
	Require     string `json:"require"`
}

type AccountStatus struct {
	AccountID string `json:"accountID"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Plan      Plan   `json:"plan"`
}
type Message struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"code"`
}

type Error struct {
	Error Message `json:"error"`
}
