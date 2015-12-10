package main

type User struct {
	Id				int64
	FirstName		string
	LastName		string
	Email			string
	Salt			string
}

type Users []User

type Login struct {
	Email			string
	Password		string
}

type JSONResponse struct {
	Error 			bool
	Message 		string
}

type UserAction struct {
	Email			string
	Action			string
}