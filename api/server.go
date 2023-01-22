package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
	Id      json.Number `json:"id"`
	Name    string      `json:"name"`
	Age     json.Number `json:"age"`
	Friends []User      `json:"friends"`
}

type FriendId struct {
	SourceId json.Number `json:"source_id"`
	TargetId json.Number `json:"target_id"`
}

type Server struct {
	*mux.Router

	usersList map[int]User
}

func NewServer() *Server {
	s := &Server{
		Router:    mux.NewRouter(),
		usersList: map[int]User{},
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, users!"))
	}).Methods("GET")
	s.HandleFunc("/create", s.createUser()).Methods("POST")
	s.HandleFunc("/get_users", s.listOfUsers()).Methods("GET")
	s.HandleFunc("/make_friends", s.makeFriends()).Methods("POST")
}

func (s *Server) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var id int = len(s.usersList) + 1
		s.usersList[id] = u

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User was created"))
	}
}

func (s *Server) listOfUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.usersList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) makeFriends() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u FriendId
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var userIndexSource, userIndexTarget int
		for key, _ := range s.usersList {
			sId, _ := u.SourceId.Int64()
			tId, _ := u.TargetId.Int64()
			if int64(key) == sId {
				userIndexSource = key
				//fmt.Fprintf(w, "sourceFriend with id %v: %v\n", key, s.usersList[key])
			}
			if int64(key) == tId {
				userIndexTarget = key
				//fmt.Fprintf(w, "targetFriend with id %v: %v\n", key, s.usersList[key])
			}
		}

		sourceFriend := s.usersList[userIndexSource]
		targetFriend := s.usersList[userIndexTarget]
		fmt.Fprintf(w, "source and target: %v, %v\n", sourceFriend, targetFriend)
		if thisUser, ok := s.usersList[userIndexSource]; ok {
			thisUser.Friends = append(thisUser.Friends, targetFriend)
			s.usersList[userIndexSource] = thisUser
			fmt.Fprintf(w, "this user: %v\n from list: %v\n\n", thisUser, s.usersList[userIndexSource])
		}
		if thisUser, ok := s.usersList[userIndexTarget]; ok {
			thisUser.Friends = append(thisUser.Friends, sourceFriend)
			s.usersList[userIndexTarget] = thisUser
			fmt.Fprintf(w, "this user: %v\n from list: %v\n\n", thisUser, s.usersList[userIndexSource])
		}

		output := s.usersList[userIndexSource].Name + " and " + s.usersList[userIndexTarget].Name + " are friends now!"
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(output))
	}
}
