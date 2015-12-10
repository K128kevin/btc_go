package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"time"
)

var sessions = make(map[string]Session) // key is token, value is session data
var LoginAttempts = make(map[string]int)

var throttleLimit int = 10
var throttleDuration time.Duration = time.Minute * 5

// handle requests to root ("/")
func Index(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, "Root, nothing here yet")
}

// handle requests to api root ("/api")
func ApiRootHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, "Default api page, will serve a file here to explain APIs")
}

// handle requests to users endpoint ("/users")
func AllUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	queryString := makeUserQueryString("GET", "")
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
	queryString := makeUserQueryString("GET", userId)
	displayString := getUsersFromDB(queryString)
	if displayString == "null" {
		displayString = "No data found :("
	}
	fmt.Fprintf(w, displayString)
}

func UserTokenGenAndEmailLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var userAction UserAction
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Fprintf(w, "Error reading request body")
	}
	if err := json.Unmarshal(body, &userAction); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// generate token with userAction and add it to db
	token := randToken(32)

	// add this token to db
	var status JSONResponse
	// verify action and user
	displayString := getUsersFromDB("select * from users where Email = \"" + userAction.Email + "\";")
	if displayString == "null" {
		status.Error = true
		status.Message = "Could not find email " + userAction.Email
	} else if userAction.Action != "verifyEmail" && userAction.Action != "resetPassword" {
		status.Error = true
		status.Message = "Invalid action provided. Valid actions are: 'verifyEmail', 'resetPassword'"
	} else if AddUserActionToDB(token, userAction) {
		status.Error = false
		status.Message = "Successfully added action " + userAction.Action + " for account " + userAction.Email
	} else {
		status.Error = true
		status.Message = "Failed to add action " + userAction.Action + " for account " + userAction.Email
	}
	retVal, _ := json.Marshal(status)
	fmt.Fprintf(w, string(retVal))
}

func DoAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
//	vars := mux.Vars(r)
//	token := vars["token"]
	var status JSONResponse
	// try to get data from token, return error if it does not exist

	// if token does exist, perform action and remove it from the database

	// return status of action

	retVal, _ := json.Marshal(status)
	fmt.Fprintf(w, string(retVal))
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
	queryString := makeUserQueryString("INSERT", "")
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
	queryString := makeUserQueryString("DELETE", userId)
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
	queryString := makeUserQueryString("UPDATE", userId)
	if (!editUserInDB(queryString, edit)) {
		fmt.Fprintf(w, "Failed to edit user with id: %s", userId)
	} else {
		AllUsers(w, r)
	}
}

// handle CORS and OPTIONS preflight
func CORSOptions(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, AuthToken")
}

func ThrottleLoginAttempts() {
	fmt.Printf("\nStarting throttle updates")
	for _ = range time.Tick(throttleDuration) {
		for k := range LoginAttempts {
			delete(LoginAttempts, k)
		}
	}
}

// handles login attempts to /users/login
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var result JSONResponse
	fmt.Printf("\nADDR: %s", r.RemoteAddr)
	if _, ok := LoginAttempts[r.RemoteAddr]; ok {
		LoginAttempts[r.RemoteAddr]++
		if LoginAttempts[r.RemoteAddr] > throttleLimit {
			result.Error = true
			result.Message = "Too many login attempts - please wait a couple minutes before trying again"
			retVal, _ := json.Marshal(result)
			fmt.Fprintf(w, string(retVal))
			return
		}
	} else {
		LoginAttempts[r.RemoteAddr] = 1
	}
    w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	var login Login
	if err := json.Unmarshal(body, &login); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	email := login.Email
	pass := login.Password
	result = tryToLogIn(w, r, email, pass)
	retVal, err := json.Marshal(result)
//	token := SaveSession(w, r, email, sessions)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(retVal))
	fmt.Printf("\n%s", string(retVal))
}

func CheckSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	token := r.Header.Get(authTokenKey)
	var status JSONResponse
	if token ==  "" {
		status.Error = true
		status.Message = "Session not found"
		fmt.Println("no cookie found")
	} else {
		status.Error = true
		status.Message = "Session cookie found but was not valid"
		if _, ok := sessions[token]; ok {
			status.Error = false;
			status.Message = sessions[token].Email
		}
	}
	retVal, _ := json.Marshal(status)
	fmt.Fprintf(w, string(retVal))
}