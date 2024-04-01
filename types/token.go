package types

import "github.com/dgrijalva/jwt-go"

type AuthToken struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"accessToken"`
	ExpiresIn int64  `json:"expiresIn"`
}

type AuthTokenClaim struct {
	*jwt.StandardClaims
	UserToken
}

type RefreshToken struct {
	AccountID string `json:"accountID"`
	AgentID   string `json:"agentID"`
}

type RefreshTokenClaim struct {
	*jwt.StandardClaims
	RefreshToken
}

type UserToken struct {
	AccountID      string        `json:"accountID"`
	ID             string        `json:"id"`
	Username       string        `json:"username"`
	FirstName      string        `json:"firstname"`
	LastName       string        `json:"lastname"`
	Email          string        `json:"email"`
	Avatar         string        `json:"avatar"`
	AccountStatus  AccountStatus `json:"account"`
	Rolename       string        `json:"rolename"`
	RolePersist    bool          `json:"rolePersist"`
	RolePredefined bool          `json:"rolePredefined"`
	CoBrowsing     bool          `json:"coBrowsing"`
	ScreenShare    bool          `json:"screenShare"`
	AudioDownload  bool          `json:"audioDownload"`
	VideoDownload  bool          `json:"videoDownload"`
	Create         Permissions   `json:"create"`
	Update         Permissions   `json:"update"`
	Read           Permissions   `json:"read"`
	Delete         Permissions   `json:"delete"`
	IsAgent        bool          `json:"isAgent"`
}
