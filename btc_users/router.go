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

var apiRoot = "/api"

var routes = Routes {
	Route {
		"Index",
		"GET",
		"/",
		Index,
	},
	Route {
		"APIRoot",
		"GET",
		apiRoot,
		ApiRootHandler,
	},
	Route {
		"UserIndex",
		"GET",
		apiRoot + "/users",
		AllUsers,
	},
	Route {
		"UserShow",
		"GET",
		apiRoot + "/users/{userId}",
		SpecificUser,
	},
	Route {
		"UserCreate",
		"POST",
		apiRoot + "/users",
		UserCreate,
	},
	Route {
		"DoAction",
		"POST",
		apiRoot + "/doaction",
		UserTokenGenAndEmailLink,
	},
	Route {
		"DoAction",
		"GET",
		apiRoot + "/doaction/{token}",
		DoAction,
	},
	Route {
		"UserDelete",
		"DELETE",
		apiRoot + "/users/{userId}",
		UserDelete,
	},
	Route {
		"UserEdit",
		"PUT",
		apiRoot + "/users/{userId}",
		UserEdit,
	},
	Route {
		"Options",
		"OPTIONS",
		apiRoot + "/users/{userId}",
		CORSOptions,
	},
	Route {
		"Options",
		"OPTIONS",
		apiRoot + "/data/predictions",
		CORSOptions,
	},
	Route {
		"Options",
		"OPTIONS",
		apiRoot + "/sessions",
		CORSOptions,
	},
	Route {
		"Login",
		"POST",
		apiRoot + "/users/login",
		UserLogin,
	},
	Route {
		"Logout",
		"POST",
		apiRoot + "/users/logout",
		UserLogout,
	},
	Route {
		"ChangePassword",
		"POST",
		apiRoot + "/users/changepassword",
		ChangePassword,
	},
	Route {
		"CheckSession",
		"GET",
		apiRoot + "/sessions",
		CheckSession,
	},
	Route {
		"GetPredictionData",
		"GET",
		apiRoot + "/data/predictions",
		PredictionGet,
	},
	Route {
		"AddPredictionData",
		"POST",
		apiRoot + "/data/predictions",
		PredictionAdd,
	},
	Route {
		"GetPriceData",
		"GET",
		apiRoot + "/data/prices",
		PriceGet,
	},
	Route {
		"AddPriceData",
		"POST",
		apiRoot + "/data/prices",
		PriceAdd,
	},
}
