package models

import "encoding/json"

type User struct {
	User_ID int         `json:"user_ID"`
	Name    string      `json:"name"`
	Age     json.Number `json:"age"`
}
type UserForDB struct {
	User_ID int    `json:"user_ID"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
}

type Friend struct {
	Source_ID json.Number `json:"source_id"`
	Target_ID json.Number `json:"target_id"`
}

type UserID struct {
	Target_ID json.Number `json:"target_id"`
}

type NewAge struct {
	NewAge json.Number `json:"newAge"`
}
