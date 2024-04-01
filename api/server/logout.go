package server

import (
	"encoding/json"
	"net/http"

	"github.com/les-cours/auth-service/types"
)

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(types.Message{
			Message: "HTTP Method not allowed",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
	json.NewEncoder(w).Encode(types.Message{
		Message: "logout successfuly",
	})
}
