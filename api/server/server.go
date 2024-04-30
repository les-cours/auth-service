package server

import (
	"database/sql"
	"github.com/les-cours/auth-service/api/auth"
	"github.com/les-cours/auth-service/api/users"
	"go.uber.org/zap"
)

type Server struct {
	DB         *sql.DB
	userClient users.UserServiceClient
	Logger     *zap.Logger
	auth.UnimplementedAuthServiceServer
}

var instance *Server

func GetInstance(userClient users.UserServiceClient, db *sql.DB, logger *zap.Logger) *Server {
	if instance != nil {
		return instance
	}

	instance = &Server{
		DB:         db,
		userClient: userClient,
		Logger:     logger,
	}
	return instance
}
