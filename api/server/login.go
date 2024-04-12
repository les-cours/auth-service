package server

import (
	"encoding/json"
	"github.com/les-cours/auth-service/api/users"
	"github.com/les-cours/auth-service/types"
	"github.com/les-cours/auth-service/utils"
	ctx "golang.org/x/net/context"
	"log"
	"net/http"
	"strings"
	"time"
)

func (s *Server) LoginHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Processing Login ... ")
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "Method not allowed",
				ErrorCode: 403,
			},
		})
		return
	}

	var user types.User
	err := json.NewDecoder(req.Body).Decode(&user)

	switch {
	case err != nil:
		{
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.Error{
				types.Message{
					Message:   "Server failed to process your request please try again",
					ErrorCode: 0,
				},
			})
			return
		}
	case utils.ValidateLoginInput(user.Osname, user.Password, user.Username):
		{
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.Error{
				types.Message{
					Message:   "All fields must be filled, please try again",
					ErrorCode: 2,
				},
			})
			return
		}
	case len(user.Password) > 16:
		{
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.Error{
				types.Message{
					Message:   "Username or password wrong, please try again",
					ErrorCode: 1,
				},
			})
			return
		}
	}

	validUser, err := s.userClient.GetUser(ctx.Background(), &users.GetUserRequest{
		Username: strings.ToLower(user.Username),
		Password: user.Password,
	})
	if err != nil {
		log.Printf("Can't getUser: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "Username or password wrong, please try again",
				ErrorCode: 1,
			},
		})
		return
	}

	accessToken, err := utils.GenerateAccessToken(validUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "Server couldn't process your request please wait",
				ErrorCode: 0,
			},
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(validUser.AccountID, validUser.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "Server couldn't process your request please wait",
				ErrorCode: 0,
			},
		})
	}

	loginTimestamp := time.Now().UnixMilli()
	_, err = s.DB.Exec(
		`
		INSERT INTO
		login_history
		(account_id, ip, timestamp, user_agent, os_name, country, city)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
		`,
		validUser.AccountID,
		utils.ReadUserIP(req),
		loginTimestamp,
		req.UserAgent(),
		user.Osname,
		"Algeria",
		"Setif") //FIXME: get the from ip the country and city
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "Server couldn't process your request please wait",
				ErrorCode: 0,
			},
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken.Token,
		Expires:  time.Unix(refreshToken.ExpiresIn, 0),
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
