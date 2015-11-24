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

type LoginResponse struct {
	Error 			bool
	Message 		string
}