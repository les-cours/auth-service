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
	ID          string      `json:"id"`
	UserType    string      `json:"userType"`
	AccountID   string      `json:"accountID"`
	Username    string      `json:"username"`
	FirstName   string      `json:"firstname"`
	LastName    string      `json:"lastname"`
	Email       string      `json:"email"`
	Avatar      string      `json:"avatar"`
	Plan        Plan        `json:"plan"`
	Permissions Permissions `json:"permissions"`
}

type SignupTokenClaim struct {
	*jwt.StandardClaims
	SignupToken
}
type SignupToken struct {
	AccountID string `json:"accountID"`
	Email     string `json:"email"`
}
