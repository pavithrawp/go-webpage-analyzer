package handler

import (
	"crypto/subtle"
	"net/http"
	"os"
)

const (
	pprofUsernameEnv      = "PPROF_USERNAME"
	pprofPasswordEnv      = "PPROF_PASSWORD"
	pprofRealm            = `Basic realm="restricted"`
	errUnauthorized       = "unauthorized"
	headerWWWAuthenticate = "WWW-Authenticate"
)

// PprofAuth is a middleware that protects pprof endpoints with basic authentication
func PprofAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set(headerWWWAuthenticate, pprofRealm)
			http.Error(w, errUnauthorized, http.StatusUnauthorized)
			return
		}

		expectedUsername := os.Getenv(pprofUsernameEnv)
		expectedPassword := os.Getenv(pprofPasswordEnv)

		// use subtle.ConstantTimeCompare to prevent timing attacks
		usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(expectedUsername)) == 1
		passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) == 1

		if !usernameMatch || !passwordMatch {
			w.Header().Set(headerWWWAuthenticate, pprofRealm)
			http.Error(w, errUnauthorized, http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
