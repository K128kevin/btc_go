package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	_ "mysql"
	"io"
	"io/ioutil"
	"strings"
)

// handle requests to root ("/")
func Index(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, "This is a service for accessing user data.")
	fmt.Fprintf(w, "\nYou can access all users by sending a get request to /users")
	fmt.Fprintf(w, "\nYou can delete/edit/read a specific user by sending the appropriate request (delete/put/get) to /users/{userId}")
	fmt.Fprintf(w, "\nYou can create a user by posting JSON data to /users")
}

// handle requests to users endpoint ("/users")
func AllUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	SaveSession(w, r)
	queryString := makeQueryString("GET", "")
	displayString := getUsersFromDB(queryString)
	if displayString == "null" {
		displayString = "No data found :("
	}
	fmt.Fprintf(w, displayString)
}

// handle requests to users/id endpoint ("/users/{userId}")
// returns a single user json object
func SpecificUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	userId := vars["userId"]
	queryString := makeQueryString("GET", userId)
	displayString := getUsersFromDB(queryString)
	if displayString == "null" {
		displayString = "No data found :("
	}
	fmt.Fprintf(w, displayString)
}

// adds given user to database and displays updated json data of all users
func UserCreate(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	var user User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Fprintf(w, "Error reading request body")
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	queryString := makeQueryString("INSERT", "")
	if (!addUserToDB(queryString, user)) {
		fmt.Fprintf(w, "Failed to add user to database")
	} else {
		AllUsers(w, r)
	}
}

// deletes specified user from database and displays updated json data of all users
func UserDelete(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	userId := vars["userId"]
	queryString := makeQueryString("DELETE", userId)
	if (!deleteUserFromDB(queryString)) {
		fmt.Fprintf(w, "Failed to delete user with id: %s", userId)
	} else {
		AllUsers(w, r)
	}
}

// edits specified user by changing specified field to new specified value
// then displays updated json data of all users
func UserEdit(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	userId := vars["userId"]
	var edit User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Fprintf(w, "Error reading request body")
	}
	if err := json.Unmarshal(body, &edit); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		fmt.Printf("error")
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	queryString := makeQueryString("UPDATE", userId)
	if (!editUserInDB(queryString, edit)) {
		fmt.Fprintf(w, "Failed to edit user with id: %s", userId)
	} else {
		AllUsers(w, r)
	}
}

// handle CORS and OPTIONS preflight
func UserOptions(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// handles login attempts
func UserLogin(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Fprintf(w, "Error reading request body")
	}
	parts := strings.Split(string(body), "*SPLITHERE*")
	email := parts[0]
	pass := parts[1]
	if tryToLogIn(email, pass) {
		// save that they are logged in
		// if !SaveSession(w, r) {
		// 	fmt.Printf("\nFailed to save login session")
		// }
		SaveSession(w, r)
		fmt.Printf("\nSuccessful Login")
		fmt.Fprintf(w, "true")
	} else { // if we hit the else statement, that means the login failed
		fmt.Printf("\nLogin Failed")
		fmt.Fprintf(w, "false")
	}
}
