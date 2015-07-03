package main

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
		
	}
	return router
}

var routes = Routes {
	Route {
		"Index",
		"GET",
		"/",
		Index,
	},
	Route {
		"UserIndex",
		"GET",
		"/users",
		AllUsers,
	},
	Route {
		"UserShow",
		"GET",
		"/users/{userId}",
		SpecificUser,
	},
	Route {
		"UserCreate",
		"POST",
		"/users",
		UserCreate,
	},
	Route {
		"UserDelete",
		"DELETE",
		"/users/{userId}",
		UserDelete,
	},
	Route {
		"UserEdit",
		"PUT",
		"/users/{userId}",
		UserEdit,
	},
	Route {
		"Options",
		"OPTIONS",
		"/users/{userId}",
		UserOptions,
	},
	Route {
		"Login",
		"POST",
		"/users/login",
		UserLogin,
	},
}
