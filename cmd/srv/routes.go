package main

import (
	"github.com/gorilla/mux"
)

func (app *Application) Routes() *mux.Router {
	rMux := mux.NewRouter()

	rMux.HandleFunc("/", app.Home)

	rMux.HandleFunc("/create-user", app.CreateUser).Methods("POST")
	rMux.HandleFunc("/make-friends", app.MakeFriends).Methods("POST")
	rMux.HandleFunc("/delete-user", app.DeleteUser).Methods("DELETE")
	rMux.HandleFunc("/friends", app.GetFriends).Methods("GET")
	rMux.HandleFunc("/user_id", app.UpdateAge).Methods("PUT")

	return rMux
}
