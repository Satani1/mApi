package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type User struct {
	Name    string      `json:"name"`
	Age     json.Number `json:"age"`
	Friends []User      `json:"friends"`
}

type FriendId struct {
	SourceId json.Number `json:"source_id"`
	TargetId json.Number `json:"target_id"`
}

type DeleteUser struct {
	TargetId json.Number `json:"target_id"`
}

type NewAge struct {
	NewAge json.Number `json:"new_age"`
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

// handling all routes
func (s *Server) routes() {
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Hello, users!")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")
	s.HandleFunc("/create", s.createUser()).Methods("POST")
	s.HandleFunc("/get_users", s.listOfUsers()).Methods("GET")
	s.HandleFunc("/make_friends", s.makeFriends()).Methods("POST")
	s.HandleFunc("/user", s.deleteUser()).Methods("DELETE")
	s.HandleFunc("/friends", s.getUserFriends()).Methods("GET")
	s.HandleFunc("/", s.changeUserAge()).Methods("PUT")
}

// handling a creation user
func (s *Server) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		//get all user data from json
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//id for the user
		var id = len(s.usersList) + 1
		//put user in the list
		s.usersList[id] = u
		//output
		//w.Header().Set("Content-Type", "application/json")
		//if err := json.NewEncoder(w).Encode(u); err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		output := "User was created with id: " + strconv.Itoa(id)
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(output)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// get all users
func (s *Server) listOfUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.usersList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// make friends
func (s *Server) makeFriends() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u FriendId
		//get two userIds from json
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var userIndexSource, userIndexTarget int
		//find users in the list of users
		for key := range s.usersList {
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
		//fmt.Fprintf(w, "source and target: %v, %v\n", sourceFriend, targetFriend)
		//adding friends to users
		if thisUser, ok := s.usersList[userIndexSource]; ok {
			thisUser.Friends = append(thisUser.Friends, targetFriend)
			s.usersList[userIndexSource] = thisUser
			//fmt.Fprintf(w, "this user: %v\n from list: %v\n\n", thisUser, s.usersList[userIndexSource])
		}
		if thisUser, ok := s.usersList[userIndexTarget]; ok {
			thisUser.Friends = append(thisUser.Friends, sourceFriend)
			s.usersList[userIndexTarget] = thisUser
			//fmt.Fprintf(w, "this user: %v\n from list: %v\n\n", thisUser, s.usersList[userIndexSource])
		}
		//output
		output := s.usersList[userIndexSource].Name + " and " + s.usersList[userIndexTarget].Name + " are friends now!"
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(output)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// delete user
func (s *Server) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u DeleteUser
		//get an userId from json
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userId, _ := u.TargetId.Int64()
		var deletedUserName string

		//delete user
		if thisUser, ok := s.usersList[int(userId)]; ok {
			deletedUserName = thisUser.Name
			delete(s.usersList, int(userId))
		}
		//output
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(deletedUserName)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// get all friends of a user
func (s *Server) getUserFriends() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get user id from url
		userId, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := s.usersList[userId]
		var userData string
		//get all friends of user
		for _, user := range user.Friends {
			userData += string(user.Name) + " " + string(user.Age) + "\n"
		}
		//output
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(userData)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// change age of a user
func (s *Server) changeUserAge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get user id from url
		userId, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var u NewAge //getting a new age from a json request
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//change age of a user
		if thisUser, ok := s.usersList[userId]; ok {
			thisUser.Age = u.NewAge
			s.usersList[userId] = thisUser
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("User age updated successfully")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
