package api

import (
	"encoding/json"
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

	usersList []User
}

func NewServer() *Server {
	s := &Server{
		Router:    mux.NewRouter(),
		usersList: []User{},
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

		s.usersList = append(s.usersList, u)

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

		//sourceId, err := u.SourceId.Int64()
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//targetId, err := u.SourceId.Int64()
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}

		var userIndexSource, userIndexTarget int
		for index, user := range s.usersList {
			if user.Id == u.SourceId {
				userIndexSource = index
			}
			if user.Id == u.TargetId {
				userIndexTarget = index
			}
		}
		//fmt.Fprint(w, s.usersList[userIndexSource], s.usersList[userIndexTarget])

		sourceFriend := s.usersList[userIndexSource]
		targetFriend := s.usersList[userIndexTarget]
		s.usersList[userIndexTarget].Friends = append(s.usersList[userIndexTarget].Friends, sourceFriend)
		s.usersList[userIndexSource].Friends = append(s.usersList[userIndexSource].Friends, targetFriend)

		//w.Header().Set("Content-Type", "application/json")
		//if err := json.NewEncoder(w).Encode(s.usersList); err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}

		output := s.usersList[userIndexSource].Name + " and " + s.usersList[userIndexTarget].Name + " are friends now!"
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(output))
	}
}
