package main

import (
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
    "crypto/sha1"
    "encoding/base64"
)

func addUserToDB(qString string, user User) bool {
	// hash password
	hasher := sha1.New()
	hasher.Write([]byte(user.Salt))
    sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
    user.Salt = sha
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

func deleteUserFromDB(qString string) bool {
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/TestDB")
	defer db.Close()
	_, err := db.Exec(qString)
	if err != nil {
		return false
	}
	return true
}

func editUserInDB(qString string, user User) bool {
	// hash password
	hasher := sha1.New()
	hasher.Write([]byte(user.Salt))
    sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
    user.Salt = sha
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
			Salt		string
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

func tryToLogIn(email string, password string) LoginResponse {
	// hash password
	hasher := sha1.New()
	hasher.Write([]byte(password))
    sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
    password = sha
	// compare to db
	queryString := makeUserQueryString("GET", "")
	queryString = queryString + " WHERE Email = '" + email + "'"
	retString := getUsersFromDB(queryString)
	var data Users
	var resp LoginResponse
	if err := json.Unmarshal([]byte(retString), &data); err != nil {
		log.Fatal(err)
		resp.Error = true
		resp.Message = "Invalid JSON provided"
    } else if len(data) < 1 {
		fmt.Printf("\nEmail not found in DB")
		resp.Error = true
		resp.Message = "Email not found"
	} else if password != data[0].Salt {
		resp.Error = true
		resp.Message = "Password was not correct"
	} else {
		resp.Error = false
		resp.Message = "Successfully logged in"
	}
	return resp
}

// qType is "GET", "INSERT", "DELETE", or "UPDATE"
// id will be empty string or id of user to get/delete/put
func makeUserQueryString(qType string, id string) string {
	switch qType {
		case "GET":
			if id == "" {
				return "SELECT * FROM users"
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