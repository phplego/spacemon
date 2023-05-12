package main

import (
	"net/http"
)

func basicAuthMiddleware(handler http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendUnauthorized := func() {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			sendUnauthorized()
			return
		}
		handler.ServeHTTP(w, r)
	})
}
