package main

import (
	"encoding/json"
	"mApi/pkg/models"
	"net/http"
	"strconv"
)

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	response := app.Addr
	if _, err := w.Write([]byte("Hello, users! Im speak from addr - " + response)); err != nil {
		app.ServeError(w, err)
		return
	}
}

func (app *Application) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	//get json data
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		app.ServeError(w, err)
		return
	}
	nAge, err := u.Age.Int64()
	if err != nil {
		app.ServeError(w, err)
		return
	}
	//insert into the db
	id, err := app.socialDB.InsertUser(u.Name, nAge)
	if err != nil {
		app.ServeError(w, err)
		return
	}
	//get id of user from db
	u.User_ID = id
	//output
	output := "User was created with id: " + strconv.Itoa(u.User_ID)
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(output)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Application) MakeFriends(w http.ResponseWriter, r *http.Request) {
	var u models.Friend

	//read json data
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		app.ServeError(w, err)
		return
	}

	//convert jsonNumber into int64
	nSourceId, err := u.Source_ID.Int64()
	if err != nil {
		app.ServeError(w, err)
		return
	}
	nTargetId, err := u.Target_ID.Int64()
	if err != nil {
		app.ServeError(w, err)
		return
	}
	//insert into the db
	err = app.socialDB.MakeFriends(nSourceId, nTargetId)
	if err != nil {
		app.ServeError(w, err)
		return
	}

	//output
	output := "User with ID " + strconv.Itoa(int(nSourceId)) + " and user with ID " + strconv.Itoa(int(nTargetId)) + " now are friends!!!"
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(output)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var u models.UserID

	//read json data
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		app.ServeError(w, err)
		return
	}
	//convert into int64
	deleteID, err := u.Target_ID.Int64()
	if err != nil {
		app.ServeError(w, err)
		return
	}
	//delete from db
	err = app.socialDB.DeleteUser(deleteID)
	//output
	output := "User with ID " + strconv.Itoa(int(deleteID)) + " was deleted."
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(output)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (app *Application) GetFriends(w http.ResponseWriter, r *http.Request) {
	//get id
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.ServeError(w, err)
		return
	}

	//get all friends from socialDB.friends table
	friends, err := app.socialDB.GetFriends(int64(id))
	if err != nil {
		app.ServeError(w, err)
		return
	}

	//output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

func (app *Application) UpdateAge(w http.ResponseWriter, r *http.Request) {
	var u models.NewAge

	//get user id
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.ServeError(w, err)
		return
	}
	//decode json data
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		app.ServeError(w, err)
		return
	}

	newAge, err := u.NewAge.Int64()
	if err != nil {
		app.ServeError(w, err)
		return
	}

	err = app.socialDB.UpdateAge(id, newAge)
	if err != nil {
		app.ServeError(w, err)
		return
	}

	//output
	output := "Age of user with ID " + strconv.Itoa(id) + " was successfully updated."
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(output)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
