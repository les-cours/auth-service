package server

import (
	"context"
	"database/sql"
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
				ErrorCode: 4,
			},
		})
		return
	}

	// banned, err := s.checkBannedIP(ReadUserIP(req))
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(Error{
	// 		Message{
	// 			Message: "Server failed to process your request, please try again",
	// 		},
	// 	})
	// 	return
	// }

	// if banned {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	json.NewEncoder(w).Encode(Error{
	// 		Message{
	// 			Message: "You have reach out the limited allowed attempts",
	// 		},
	// 	})
	// 	return
	// }

	var user types.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "Server failed to process your request please try again",
				ErrorCode: 0,
			},
		})
		return
	}

	if user.Osname == "" || user.Password == "" || user.Username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "All fields must be filled, please try again",
				ErrorCode: 2,
			},
		})
		return
	}

	if len(user.Password) > 16 {
		// We don't accept password that take more than 16characters for password
		// because any one can send us 1 million or more characters which will overload
		// the service caused by the encryption the password
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "Username or password wrong, please try again",
				ErrorCode: 1,
			},
		})
		return
	}

	validUser, err := s.userClient.GetUser(ctx.Background(), &users.GetUserRequest{
		Username: strings.ToLower(user.Username),
		Password: user.Password,
	})
	if err != nil {
		log.Printf("Can't getUser: %v", err)
		s.failedLoginAttemptIP(utils.ReadUserIP(req))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "Username or password wrong, please try again",
				ErrorCode: 1,
			},
		})
		return
	}

	s.removeLoginAttempt(utils.ReadUserIP(req))

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
		(account_id, agent_id, ip, timestamp, user_agent, os_name, country, city)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
		`,
		validUser.AccountID,
		validUser.Id,
		utils.ReadUserIP(req),
		loginTimestamp,
		req.UserAgent(),
		user.Osname,
		"Algeria",
		"Algeries") //FIXME: get the from ip the country and city
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

func (s *Server) failedLoginAttemptIP(ip string) {
	attempt := 0
	err := s.DB.QueryRow(
		`
		SELECT
		attempts
		FROM
		login_attempt
		WHERE
		ip = $1
		`, ip).Scan(&attempt)
	if err != nil && err == sql.ErrNoRows {
		s.DB.Exec(
			`
			INSERT INTO
			login_attempt
			(ip, timestamp_seconds, attempts)
			VALUES
			($1, $2, $3)
			`, ip, time.Now().Unix(), 1)
	} else if err == nil {
		s.DB.Exec(
			`
			UPDATE
			login_attempt
			SET
			timestamp_seconds = $1, attempts = $2
			WHERE
			ip = $3
			`, time.Now().Unix(), attempt+1, ip)
	}

}

func (s *Server) removeLoginAttempt(ip string) {
	s.DB.ExecContext(context.Background(),
		`
		DELETE FROM
		login_attempt
		WHERE
		ip = $1
		`, ip)
}
