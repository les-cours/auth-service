package service

import (
	"database/sql"
	"fmt"
	"github.com/les-cours/auth-service/env"
	_ "github.com/lib/pq"
	"log"
)

func StartDB() (*sql.DB, error) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		env.Settings.Database.PSQLConfig.Host,
		env.Settings.Database.PSQLConfig.Port,
		env.Settings.Database.PSQLConfig.Username,
		env.Settings.Database.PSQLConfig.Password,
		env.Settings.Database.PSQLConfig.DbName,
		env.Settings.Database.PSQLConfig.SslMode,
	)

	db, err := sql.Open("postgres", dbConfig)

	if err != nil {
		log.Fatalf("Failed to connect to postgres database: %v", err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(7)

	return db, nil
}
