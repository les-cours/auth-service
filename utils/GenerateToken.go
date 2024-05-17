package utils

import (
	"github.com/les-cours/auth-service/env"
	"github.com/les-cours/auth-service/types"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/les-cours/auth-service/api/users"
)

var AccessTokenExpiresAt = time.Now().Add(time.Minute * 60 * 24 * 360).Unix() // time.Now().Add(time.Minute * time.Duration(env.Settings.AccessTokenLife)).Unix()
func GenerateAccessToken(user *users.User) (*types.AuthToken, error) {
	accessTokenHash := jwt.New(jwt.SigningMethodHS256)
	accessTokenHash.Claims = &types.AuthTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: AccessTokenExpiresAt,
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
		ExpiresIn: AccessTokenExpiresAt,
		TokenType: env.Settings.TokenType,
	}

	return token, nil
}

func GenerateTeacherAccessToken(user *users.User) (*types.AuthToken, error) {

	accessTokenHash := jwt.New(jwt.SigningMethodHS256)
	accessTokenHash.Claims = &types.AuthTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: AccessTokenExpiresAt,
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
			Create: types.Permissions{
				USER:     user.CREATE.Users,
				LEARNING: user.CREATE.Learning,
				ORGS:     user.CREATE.Orgs,
			},
			Update: types.Permissions{
				USER:     user.UPDATE.Users,
				LEARNING: user.UPDATE.Learning,
				ORGS:     user.UPDATE.Orgs,
			},
			Read: types.Permissions{
				USER:     user.READ.Users,
				LEARNING: user.READ.Learning,
				ORGS:     user.READ.Orgs,
			},
			Delete: types.Permissions{
				USER:     user.DELETE.Users,
				LEARNING: user.DELETE.Learning,
				ORGS:     user.DELETE.Orgs,
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
		ExpiresIn: AccessTokenExpiresAt,
		TokenType: env.Settings.TokenType,
	}

	return token, nil
}

//func GenerateAdminAccessToken(user *users.User) (*types.AuthToken, error) {
//	AccessTokenExpiresAt := time.Now().Add(time.Minute * time.Duration(env.Settings.AccessTokenLife)).Unix()
//	accessTokenHash := jwt.New(jwt.SigningMethodHS256)
//	accessTokenHash.Claims = &types.AuthTokenClaim{
//		&jwt.StandardClaims{
//			ExpiresAt: AccessTokenExpiresAt,
//		},
//		types.UserToken{
//			ID:        user.Id,
//			UserType:  user.UserType,
//			AccountID: user.AccountID,
//			Username:  user.Username,
//			FirstName: user.FirstName,
//			LastName:  user.LastName,
//			Email:     user.Email,
//			Avatar:    user.Avatar,
//			Create: types.Permissions{
//				USER:     true,
//				LEARNING: true,
//				ORGS:     true,
//			},
//			Update: types.Permissions{
//				USER:     true,
//				LEARNING: true,
//				ORGS:     true,
//			},
//			Read: types.Permissions{
//				USER:     true,
//				LEARNING: true,
//				ORGS:     true,
//			},
//			Delete: types.Permissions{
//				USER:     true,
//				LEARNING: true,
//				ORGS:     true,
//			},
//		},
//	}
//	accessToken, err := accessTokenHash.SignedString([]byte(env.Settings.JWTAccessTokenSecret))
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//
//	token := &types.AuthToken{
//		Token:     accessToken,
//		ExpiresIn: AccessTokenExpiresAt,
//		TokenType: env.Settings.TokenType,
//	}
//
//	return token, nil
//}
