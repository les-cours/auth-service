package server

import (
	"database/sql"
	"github.com/les-cours/auth-service/api/auth"
	"github.com/les-cours/auth-service/api/users"
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
