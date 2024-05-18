package service

import (
	"net/http"
)

func cors(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCounter.Inc()
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", r.Method)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		h.ServeHTTP(w, r)

	})
}
