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

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    accessToken.Token,
		Expires:  time.Unix(accessToken.ExpiresIn, 0),
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

func (s *Server) LoginTeacherHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Processing Login ... %v", req.Method)

	if req.Method == "OPTIONS" {
		w.WriteHeader(200)

		return
	}

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
		Username:  strings.ToLower(user.Username),
		Password:  user.Password,
		IsTeacher: true,
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

	accessToken, err := utils.GenerateTeacherAccessToken(validUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "Server couldn't process your request please wait",
				ErrorCode: 0,
			},
		})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    accessToken.Token,
		Expires:  time.Unix(accessToken.ExpiresIn, 0),
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

func (s *Server) LoginAdminHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == "OPTIONS" {
		w.WriteHeader(200)

		return
	}

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
		IsAdmin:  true,
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

	accessToken, err := utils.GenerateAdminAccessToken(validUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.Error{
			types.Message{
				Message:   "Server couldn't process your request please wait",
				ErrorCode: 0,
			},
		})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    accessToken.Token,
		Expires:  time.Unix(accessToken.ExpiresIn, 0),
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
