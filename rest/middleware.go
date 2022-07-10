package rest

import (
	"crypto/subtle"
	"golang.org/x/crypto/bcrypt"
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

		passwordMatch := bcrypt.CompareHashAndPassword([]byte(config.Get("password")), []byte(password)) == nil
		usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(config.Get("username"))) == 1

		if usernameMatch && passwordMatch {
			next.ServeHTTP(w, r)
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
