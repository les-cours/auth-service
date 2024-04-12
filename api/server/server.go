package server

import (
	"context"
	"database/sql"
	"log"

	"github.com/les-cours/auth-service/api/auth"
	"github.com/les-cours/auth-service/api/users"
	"github.com/les-cours/auth-service/utils"
)

type Server struct {
	DB         *sql.DB
	userClient users.UserServiceClient
	auth.UnimplementedAuthServiceServer
}

var instance *Server

func GetInstance(userClient users.UserServiceClient, db *sql.DB) *Server {
	if instance != nil {
		return instance
	}

	instance = &Server{
		DB:         db,
		userClient: userClient,
	}
	return instance
}

func (s *Server) GenerateAccessToken(ctx context.Context, in *auth.GenerateAccessTokenRequest) (*auth.GenerateAccessTokenResponse, error) {

	var validUser, err = s.userClient.GetUserByID(ctx, &users.GetUserByIDRequest{
		AccountID: in.AccountID,
	})
	if err != nil {
		log.Println("err server.go GenerateAccessToken",err)
		return nil, err
	}

	token, err := utils.GenerateAccessToken(validUser)
	if err != nil {
		log.Println("err server.go utils.GenerateAccessToken ",err)
		return nil, err
	}

	var response = &auth.GenerateAccessTokenResponse{
		Token: &auth.AuthToken{
			Token:     token.Token,
			TokenType: token.TokenType,
			ExpiresAt: token.ExpiresIn,
		},
	}
	return response, nil
}

func (s *Server) GenerateRefreshToken(ctx context.Context, in *auth.GenerateRefreshTokenRequest) (*auth.GenerateRefreshTokenResponse, error) {

	var validUser, err = s.userClient.GetUserByID(ctx, &users.GetUserByIDRequest{
		AccountID: in.AccountID,
	})
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateAccessToken(validUser)
	if err != nil {
		return nil, err
	}

	var response = &auth.GenerateRefreshTokenResponse{
		Token: &auth.AuthToken{
			Token:     token.Token,
			ExpiresAt: token.ExpiresIn,
		},
	}

	return response, nil
}
