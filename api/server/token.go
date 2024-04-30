package server

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/les-cours/auth-service/env"
	"github.com/les-cours/auth-service/types"
	"log"
	"net/http"
	"strings"
)

func (s *Server) TokenHealthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.Message{
			Message:   "Http method not allowed",
			ErrorCode: 403,
		})
		return
	}

	authorizationHeader := req.Header.Get("Authorization")
	if authorizationHeader != "" {
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) == 2 {
			if bearerToken[0] != env.Settings.TokenType {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(
					types.Message{
						Message: "Wrong type of token",
					},
				)
				return
			}

			token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("there was an error")
				}

				return []byte(env.Settings.JWTAccessTokenSecret), nil
			})

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(
					types.Message{
						Message: "You are not authorized",
					},
				)
				return
			}

			if token.Valid {
				json.NewEncoder(w).Encode(types.Message{Message: "Token Valid"})
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(types.Message{
				Message: "You are not authorized",
			})
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.Message{
			Message: "You are not authorized",
		})
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(types.Message{
		Message: "You are not authorized",
	})
}
