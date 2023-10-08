package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
)

func authMiddleware(handler http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handle..", r.Method, r.RequestURI)

		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value != generateDailyToken(username, password) {
			w.Header().Set("Content-Type", "text/html")

			if r.Method == "POST" && r.RequestURI == "/login" {
				user := r.PostFormValue("user")
				pass := r.PostFormValue("pass")

				// Check if they match the preset username and password
				if user == username && pass == password {
					// Set the cookie
					log.Println("Set cookie")
					http.SetCookie(w, &http.Cookie{
						Name:  "token",
						Value: generateDailyToken(user, pass),
					})
				} else {
					log.Println("Invalid credentials")
					time.Sleep(time.Second)
				}
				// Redirect back to the referer
				http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
			}

			// If method is not POST or cookie does not exist/is invalid, show the form
			w.Write([]byte(`
				<form method="POST" action="/login">
					<input name="user" placeholder="user"> 
					<input type="password" name="pass" type="password" placeholder="password">
					<button type="submit">Sign in</button>
				</form>
				`))
			return
		}

		// If the cookie is valid, continue handling
		handler.ServeHTTP(w, r)
	})
}

func generateDailyToken(user, pass string) string {
	currentDate := time.Now().Format("2006-01-02")
	data := fmt.Sprintf("%s:%s:%s", currentDate, user, pass)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
