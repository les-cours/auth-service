package types

type Object map[string]interface{}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Osname   string `json:"osname"`
}

type Permissions struct {
	USER     bool `json:"user"`
	LEARNING bool `json:"learning"`
	ORGS     bool `json:"orgs"`
	PAYMENT  bool `json:"payment"`
}

type Message struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"code"`
}

type Error struct {
	Error Message `json:"error"`
}
