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

var root = "/api"

var routes = Routes {
	Route {
		"Index",
		"GET",
		root,
		Index,
	},
	Route {
		"UserIndex",
		"GET",
		root + "/users",
		AllUsers,
	},
	Route {
		"UserShow",
		"GET",
		root + "/users/{userId}",
		SpecificUser,
	},
	Route {
		"UserCreate",
		"POST",
		root + "/users",
		UserCreate,
	},
	Route {
		"UserDelete",
		"DELETE",
		root + "/users/{userId}",
		UserDelete,
	},
	Route {
		"UserEdit",
		"PUT",
		root + "/users/{userId}",
		UserEdit,
	},
	Route {
		"Options",
		"OPTIONS",
		root + "/users/{userId}",
		UserOptions,
	},
	Route {
		"Login",
		"POST",
		root + "/users/login",
		UserLogin,
	},
	Route {
		"GetPredictionData",
		"GET",
		root + "/data/predictions",
		PredictionGet,
	},
	Route {
		"AddPredictionData",
		"POST",
		root + "/data/predictions",
		PredictionAdd,
	},
	Route {
		"GetPriceData",
		"GET",
		root + "/data/prices",
		PriceGet,
	},
	Route {
		"AddPriceData",
		"POST",
		root + "/data/prices",
		PriceAdd,
	},
}
