package main

import (
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
    "crypto/sha1"
    "encoding/base64"
	"net/http"
)

var UserActionTableName string = "one_time_actions"

func AddUserToDB(qString string, user User) bool {
	// hash password
    user.Salt = HashPassword(user.Salt)
    // connect to db
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	defer db.Close()
	_, err := db.Exec(qString + "'" +
						user.FirstName + "','" +
						user.LastName + "','" +
						user.Email + "','" +
						user.Salt + "');")
	if err != nil {
		return false
	}
	return true
}

func DeleteUserFromDB(qString string) bool {
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	defer db.Close()
	_, err := db.Exec(qString)
	if err != nil {
		return false
	}
	return true
}

func EditUserInDB(qString string, user User) bool {
	// hash password
    user.Salt = HashPassword(user.Salt)
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	defer db.Close()
	_, err := db.Exec(qString + "'" +
						user.FirstName + "','" +
						user.LastName + "','" +
						user.Email + "','" +
						user.Salt + "');")
	if err != nil {
		return false
	}
	return true
}

func GetActionFromDB(token string) (string, string, error) {
	fmt.Printf("\nGetActionFromDB called\n")
	var err error = nil

	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	rows, _ := db.Query("select email, action from one_time_actions where token = \"" + token + "\";")
	defer db.Close()
	defer rows.Close()

	var (
		email		string
		action		string
	)
	for rows.Next() {
		err = rows.Scan(&email, &action)
	}

	return email, action, err
}

func DeleteActionFromDB(token string) error {
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	defer db.Close()

	_, err := db.Exec("delete from one_time_actions where token = \"" + token + "\";")
	return err
}

func getUsersFromDB(qString string) string {
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	rows, _ := db.Query(qString)
	defer db.Close()
	defer rows.Close()
	var tempUsers Users
	for rows.Next() {
		var (
			id			int64
			FirstName	string
			LastName	string
			Email		string
			Salt 		string
		)
		err := rows.Scan(&id, &FirstName, &LastName, &Email, &Salt)
		if err != nil {
			log.Fatal(err)
		}
		tempUsers = append(tempUsers, User{id, FirstName, LastName, Email, Salt})
	}
	retVal, err := json.Marshal(tempUsers)
	if err != nil {
		return "An error occurred converting returned data to json"
	}
	return string(retVal)
}

func tryToChangePassword(newPass NewPass) JSONResponse {
	oldPass := HashPassword(newPass.OldPass)
	// compare to db
	queryString := MakeUserQueryString("GET", "")
	queryString = queryString + " WHERE Email = '" + newPass.Email + "'"
	retString := getUsersFromDB(queryString)
	var data Users
	var status JSONResponse
	if err := json.Unmarshal([]byte(retString), &data); err != nil {
		log.Fatal(err)
		status.Error = true
		status.Message = "Invalid JSON provided"
	} else if len(data) < 1 {
		fmt.Printf("\nEmail not found in DB")
		status.Error = true
		status.Message = "Email not found"
	} else if oldPass != data[0].Salt {
		status.Error = true
		status.Message = "Password was not correct"
	} else {
		status.Error = false
		status.Message = "success"
		err := SetPassword(newPass.Email, newPass.NewPass)
		if err != nil {
			status.Error = true
			status.Message = "Failed to set new password"
		}
	}

	return status
}

func TryToLogIn(w http.ResponseWriter, r *http.Request, email string, password string) JSONResponse {
	// hash password
    password = HashPassword(password)
	// compare to db
	queryString := MakeUserQueryString("GET", "")
	queryString = queryString + " WHERE Email = '" + email + "'"
	retString := getUsersFromDB(queryString)
	var data Users
	var status JSONResponse
	if err := json.Unmarshal([]byte(retString), &data); err != nil {
		log.Fatal(err)
		status.Error = true
		status.Message = "Invalid JSON provided"
    } else if len(data) < 1 {
		fmt.Printf("\nEmail not found in DB")
		status.Error = true
		status.Message = "Email not found"
	} else if password != data[0].Salt {
		status.Error = true
		status.Message = "Password was not correct"
	} else {
		status.Error = false
		status.Message = SaveSession(w, r, email, sessions)
	}
	return status
}

func AddUserActionToDB(token string, userAction UserAction) bool {
	fmt.Printf("\nAddUserActionToDB called\n")
	qString := "INSERT INTO " + UserActionTableName + " (Token,Email,Action) values(\"" + token + "\",\"" + userAction.Email + "\",\"" + userAction.Action + "\");"
	fmt.Printf("\nQuery String: %s\n", qString)
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	defer db.Close()
	_, err := db.Exec(qString)
	if err != nil {
		fmt.Printf("\nQuery failed\n")
		return false
	}
	return true
}

// qType is "GET", "INSERT", "DELETE", or "UPDATE"
// id will be empty string or id of user to get/delete/put
func MakeUserQueryString(qType string, id string) string {
	switch qType {
		case "GET":
			if id == "" {
				return "SELECT id, FirstName, LastName, Email, Salt FROM users"
			} else {
				return "SELECT * FROM users WHERE id = " + id
			}
		case "DELETE":
			return "DELETE FROM users WHERE id = " + id
		case "INSERT":
			return "INSERT INTO users (FirstName, LastName, Email, Salt) VALUES("
		case "UPDATE":
			return "REPLACE INTO users (id, FirstName, LastName, Email, Salt) VALUES('" + id + "',"
		default:
			panic("unrecognized query type")
	}
}

func IsVerified(email string) bool {
	fmt.Printf("\nisVerified called\n")
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	rows, _ := db.Query("select * from users where Verified = true AND Email = \"" + email + "\"")
	defer db.Close()
	defer rows.Close()
	next := rows.Next()
	fmt.Printf("Next: %t", next)
	return next
}

func MakeVerified(email string) error {
	fmt.Printf("\nmakeVerified called\n")
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	defer db.Close()

	fmt.Printf("\nEmail: %s\n", email)
	_, err := db.Exec("update users set Verified = true where Email = \"" + email + "\";")
	return err
}

func SetPassword(email string, password string) error {
	fmt.Printf("\nsetPassword called\n")
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	defer db.Close()

	// hash password
	newPW := HashPassword(password)

	query := "update users set Salt = \"" + newPW + "\" where Email = \"" + email + "\";"
	fmt.Printf("\nQuery: %s", query)
	_, err := db.Exec(query)
	return err
}

func HashPassword(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}