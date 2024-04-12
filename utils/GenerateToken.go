package utils

import (
	"github.com/les-cours/auth-service/env"
	"github.com/les-cours/auth-service/types"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/les-cours/auth-service/api/users"
)

func GenerateAccessToken(user *users.User) (*types.AuthToken, error) {
	accessTokenExpiresAt := time.Now().Add(time.Minute * time.Duration(env.Settings.AccessTokenLife)).Unix()
	accessTokenHash := jwt.New(jwt.SigningMethodHS256)
	accessTokenHash.Claims = &types.AuthTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: accessTokenExpiresAt,
		},
		types.UserToken{
			ID:        user.Id,
			UserType:  user.UserType,
			AccountID: user.AccountID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Plan: types.Plan{
				PlanID:      user.Plan.PlanID,
				Name:        user.Plan.Name,
				PeriodEndAt: user.Plan.PeriodEndAt,
				Active:      user.Plan.Active,
				Require:     user.Plan.Require,
			},
			Permissions: types.Permissions{
				WriteComment: user.Permissions.WriteComment,
				Live:         user.Permissions.Live,
			},
		},
	}
	accessToken, err := accessTokenHash.SignedString([]byte(env.Settings.JWTAccessTokenSecret))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	token := &types.AuthToken{
		Token:     accessToken,
		ExpiresIn: accessTokenExpiresAt,
		TokenType: env.Settings.TokenType,
	}

	return token, nil
}

func GenerateRefreshToken(accountID string, agentID string) (*types.AuthToken, error) {

	refreshTokenExpiresAt := time.Now().Add(time.Second * time.Duration(env.Settings.RefreshTokenLife)).Unix()
	refreshTokenHash := jwt.New(jwt.SigningMethodHS512)
	refreshTokenHash.Claims = &types.RefreshTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiresAt,
		},
		types.RefreshToken{
			AccountID: accountID,
			AgentID:   agentID,
		},
	}
	refreshToken, err := refreshTokenHash.SignedString([]byte(env.Settings.JWTRefreshTokenSecret))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	token := &types.AuthToken{
		Token:     refreshToken,
		ExpiresIn: refreshTokenExpiresAt,
	}

	return token, nil

}
