package req

import "github.com/doubledutch/dd-vote/api/models/table"

// UserRequest is sent by the client when logging in
type UserRequest struct {
	Firstname string `json:"firstName" binding:"required"`
	Lastname  string `json:"lastName" binding:"required"`
	ClientID  uint   `json:"userId" binding:"required"`
}

// AdminLoginRequest is sent by the client when logging in as an admin
type AdminLoginRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	GroupUUID string `json:"groupId" binding:"required"`
}

// IsValid returns whether the request has valid data
func (req *UserRequest) IsValid() bool {
	return len(req.Firstname) > 0 && len(req.Lastname) > 0
}

// ToUser converts a UserRequest object to User
func (req *UserRequest) ToUser() table.User {
	var user table.User
	user.ClientID = req.ClientID
	user.Firstname = req.Firstname
	user.Lastname = req.Lastname
	return user
}

// ToUser converts an AdminLoginRequest object to User
func (req *AdminLoginRequest) ToUser() table.User {
	var user table.User
	user.Email = req.Email
	user.Password = req.Password
	return user
}
