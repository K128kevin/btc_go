package main

type User struct {
	Id				int64
	FirstName		string
	LastName		string
	Email			string
	Salt			string
}

type Users []User
