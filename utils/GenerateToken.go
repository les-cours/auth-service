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

func GenerateTeacherAccessToken(user *users.User) (*types.AuthToken, error) {
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
			Permissions: types.Permissions{
				WriteComment: user.Permissions.WriteComment,
				Live:         user.Permissions.Live,
				Upload:       user.Permissions.Upload,
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
