package utils

import (
	"github.com/les-cours/auth-service/env"
	"github.com/les-cours/auth-service/types"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/les-cours/auth-service/api/users"
)

func GenerateAccessToken(validUser *users.User) (*types.AuthToken, error) {
	accessTokenExpiresAt := time.Now().Add(time.Minute * time.Duration(env.Settings.AccessTokenLife)).Unix()
	accessTokenHash := jwt.New(jwt.SigningMethodHS256)
	accessTokenHash.Claims = &types.AuthTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: accessTokenExpiresAt,
		},
		types.UserToken{
			AccountID: validUser.AccountID,
			ID:        validUser.Id,
			Username:  validUser.Username,
			FirstName: validUser.FirstName,
			LastName:  validUser.LastName,
			Avatar:    validUser.Avatar,
			AccountStatus: types.AccountStatus{
				AccountID: validUser.AccountID,
				Name:      validUser.Account.Name,
				Status:    validUser.Account.Status,
				Plan: types.Plan{
					PlanID:      validUser.Account.Plan.PlanID,
					Name:        validUser.Account.Plan.Name,
					PeriodEndAt: validUser.Account.Plan.PeriodEndAt,
					Active:      validUser.Account.Plan.Active,
					Require:     validUser.Account.Plan.Require,
				},
			},
			Rolename:       validUser.Role.Name,
			RolePersist:    validUser.Role.Persist,
			RolePredefined: validUser.Role.Predefined,
			Email:          validUser.Email,
			CoBrowsing:     validUser.Role.CoBrowsing,
			ScreenShare:    validUser.Role.ScreenShare,
			AudioDownload:  validUser.Role.AudioDownload,
			VideoDownload:  validUser.Role.VideoDownload,
			Create:         PermissionsRPCToCode(validUser.Role.Create),
			Update:         PermissionsRPCToCode(validUser.Role.Update),
			Read:           PermissionsRPCToCode(validUser.Role.Read),
			Delete:         PermissionsRPCToCode(validUser.Role.Delete),
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

func GenerateVisitorToken(accountID string, visitorID string) (*types.AuthToken, error) {
	visitorTokenExpiresAt := time.Now().Add(time.Second * time.Duration(env.Settings.AccessTokenLife)).Unix()
	visitorTokenHash := jwt.New(jwt.SigningMethodHS256)
	visitorTokenHash.Claims = &types.AuthTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: visitorTokenExpiresAt,
		},
		types.UserToken{
			AccountID: accountID,
			ID:        visitorID,
			IsAgent:   false,
		},
	}
	visitorToken, err := visitorTokenHash.SignedString([]byte(env.Settings.JWTAccessTokenSecret))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	token := &types.AuthToken{
		Token:     visitorToken,
		ExpiresIn: visitorTokenExpiresAt,
		TokenType: env.Settings.TokenType,
	}

	return token, nil

}
