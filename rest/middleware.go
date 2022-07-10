package rest

import (
	"crypto/subtle"
	"net/http"
	"smart-awning/config"
)

func basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter username and password."`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		usernamesMatch := subtle.ConstantTimeCompare([]byte(username), []byte(config.HTTP.Username)) == 1
		passwordsMatch := subtle.ConstantTimeCompare([]byte(password), []byte(config.HTTP.Password)) == 1

		if usernamesMatch && passwordsMatch {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
