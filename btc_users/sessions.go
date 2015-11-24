package main

import (
	"net/http"
	"time"
	"crypto/rand"
	"fmt"
)

type Session struct {
	Email			string
	Expiration		time.Time
}

func SaveSession(w http.ResponseWriter, r *http.Request, email string, sessions map[string]Session) {

	token := randToken()
	prevToken, err := r.Cookie("AuthToken")

	// found cookie
	if err == nil {
		// found session
		if sessions[prevToken.Value].Email == email {
			fmt.Printf("\n%s logged in already", email)
			renewSession(prevToken.Value, sessions, w)
			return
		}
	}

	// no session found or session found but no matches
	cookie := http.Cookie{Name: "AuthToken", Value: token, Expires: time.Now().Add(time.Hour), HttpOnly: true}
	var tempSession Session
	tempSession.Email = email
	tempSession.Expiration = time.Now().Add(time.Hour)
	sessions[token] = tempSession
	http.SetCookie(w, &cookie)

}

func randToken() string {
	token := make([]byte, 16)
	rand.Read(token)
	return fmt.Sprintf("%x", token)
}

func renewSession(token string, sessions map[string]Session, w http.ResponseWriter) {
	var tempSession Session
	tempSession.Expiration = time.Now().Add(time.Hour)
	tempSession.Email = sessions[token].Email

	sessions[token] = tempSession
	cookie := http.Cookie{Name: "AuthToken", Value: token, Expires: time.Now().Add(time.Hour), HttpOnly: true}
	http.SetCookie(w, &cookie)
}

func startUpdates() {
	fmt.Printf("\nStarting session updates")
	for _ = range time.Tick(10 * time.Second) {
		updateSessions(sessions)
	}
}

func updateSessions(sessions map[string]Session) {
	for token, state := range sessions {
		if state.Expiration.Before(time.Now()) {
			delete(sessions, token)
			fmt.Printf("\nSession %s deleted for email %s", token, state.Email)
		}
	}
}