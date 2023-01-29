package user

import (
	"encoding/json"
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

type NewAge struct {
	NewAge json.Number `json:"new_age"`
}

type DeleteUser struct {
	TargetId json.Number `json:"target_id"`
}
