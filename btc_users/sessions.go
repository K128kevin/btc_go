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

var refreshRate time.Duration = time.Second * 10
var sessionDuration time.Duration = time.Hour
var authTokenKey string = "AuthToken"
var tokenSize int = 16

func SaveSession(w http.ResponseWriter, r *http.Request, email string, sessions map[string]Session) string {

	token := randToken(tokenSize)
	prevToken, err := r.Cookie(authTokenKey)

	// found cookie
	if err == nil {
		// found session
		if sessions[prevToken.Value].Email == email {
			fmt.Printf("\n%s logged in already", email)
			renewSession(prevToken.Value, sessions, w)
			return "logged in already"
		}
	}

	// no session found or session found but no matches
//	cookie := http.Cookie{Name: authTokenKey, Value: token, Expires: time.Now().Add(sessionDuration), HttpOnly: true}
	var tempSession Session
	tempSession.Email = email
	tempSession.Expiration = time.Now().Add(sessionDuration)
	sessions[token] = tempSession
//	http.SetCookie(w, &cookie)
	return token;
}

func randToken(size int) string {
	token := make([]byte, size)
	rand.Read(token)
	return fmt.Sprintf("%x", token)
}

func renewSession(token string, sessions map[string]Session, w http.ResponseWriter) {
	var tempSession Session
	tempSession.Expiration = time.Now().Add(sessionDuration)
	tempSession.Email = sessions[token].Email

	sessions[token] = tempSession
	cookie := http.Cookie{Name: authTokenKey, Value: token, Expires: time.Now().Add(time.Hour), HttpOnly: true}
	http.SetCookie(w, &cookie)
}

func StartSessionUpdates() {
	fmt.Printf("\nStarting session updates")
	for _ = range time.Tick(refreshRate) {
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