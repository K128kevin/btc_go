package main

import (
	"net/http"
	"time"
)

func SaveSession(w http.ResponseWriter, r *http.Request) {
	expire := time.Now().AddDate(0, 0, 1)
    cookie := http.Cookie{"test", "tcookie", "/", "", expire, expire.Format(time.UnixDate), 86400, true, true, "test=tcookie", []string{"test=tcookie"}}
    http.SetCookie(w, &cookie)
}

func ValidateSession(r *http.Request) bool {
	_, err := r.Cookie("test")
	if err != nil {
		return false
	}
	return true
}