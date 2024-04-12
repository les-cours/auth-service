package server

import (
	context "context"
	"github.com/dgrijalva/jwt-go"
	"github.com/les-cours/auth-service/api/auth"
	"github.com/les-cours/auth-service/env"
	"github.com/les-cours/auth-service/types"
	"log"
	"time"
)

func (s *Server) Signup(ctx context.Context, in *auth.SignUpRequest) (*auth.SignUpResponse, error) {

	log.Println("id : ", in.AccountID)
	token, err := s.GenerateAccessToken(ctx, &auth.GenerateAccessTokenRequest{
		AccountID: in.AccountID,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, &auth.GenerateRefreshTokenRequest{
		AccountID: in.AccountID,
	})
	if err != nil {
		return nil, err
	}

	signupTokenExpiresAt := time.Now().Add(time.Second * time.Duration(env.Settings.SignupLinkLife)).Unix()
	signupTokenHash := jwt.New(jwt.SigningMethodHS256)
	signupTokenHash.Claims = &types.SignupTokenClaim{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: signupTokenExpiresAt,
		},
		SignupToken: types.SignupToken{
			AccountID: in.AccountID,
			Email:     in.Email,
		}}

	signupToken, err := signupTokenHash.SignedString([]byte(env.Settings.JWTSignupTokenSecret))

	return &auth.SignUpResponse{
		AccessToken: &auth.AuthToken{
			Token:     token.Token.Token,
			ExpiresAt: token.Token.ExpiresAt,
			TokenType: token.Token.TokenType,
		},
		RefreshToken: &auth.RefreshToken{
			Token:     refreshToken.Token.Token,
			ExpiresAt: refreshToken.Token.ExpiresAt,
		},
		SignupToken: &auth.SignupToken{
			Token:     signupToken,
			ExpiresAt: signupTokenExpiresAt,
		},
	}, nil

}
