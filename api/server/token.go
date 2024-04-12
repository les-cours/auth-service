package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/les-cours/auth-service/api/users"
	"github.com/les-cours/auth-service/env"
	"github.com/les-cours/auth-service/types"
	"github.com/les-cours/auth-service/utils"
	ctx "golang.org/x/net/context"
)

func (s *Server) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.Message{
			Message: "Method not allowed",
		})
		return
	}

	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message: "Bad request",
			},
		})
		return
	}

	refreshToken := cookie.Value
	user := types.RefreshTokenClaim{}
	_, err = jwt.ParseWithClaims(refreshToken, &user, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.Settings.JWTRefreshTokenSecret), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    "",
			MaxAge:   -1,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message: "You are not authorized",
			},
		})
		return
	}

	validUser, err := s.userClient.GetUserByID(ctx.Background(), &users.GetUserByIDRequest{
		AccountID: user.AccountID,
	})
	if err != nil && err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message: "Something goes wrong with the cookie",
			},
		})
		return
	}

	if err != nil {
		log.Printf("Can't getUser: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message: "Server failed to process your request, please try again",
			},
		})
		return
	}

	accessToken, err := utils.GenerateAccessToken(validUser)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message: "Server failed to process your request please try again",
			},
		})
		return
	}

	refreshTokenExpiresAt := time.Now().Add(time.Second * time.Duration(env.Settings.RefreshTokenLife)).Unix()
	refreshTokenHash := jwt.New(jwt.SigningMethodHS512)
	refreshTokenHash.Claims = &types.RefreshTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiresAt,
		},
		types.RefreshToken{
			AccountID: validUser.AccountID,
			AgentID:   validUser.Id,
		},
	}
	newRefreshToken, err := refreshTokenHash.SignedString([]byte(env.Settings.JWTRefreshTokenSecret))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Message{
			Message: "Server failed to process your request please try again",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(time.Second * time.Duration(env.Settings.RefreshTokenLife)),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	json.NewEncoder(w).Encode(types.Object{
		"user": types.AuthToken{
			Token:     accessToken.Token,
			TokenType: accessToken.TokenType,
			ExpiresIn: accessToken.ExpiresIn,
		},
	})
}

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
