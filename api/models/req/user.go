package req

import "github.com/jordanjoz/dd-vote/api/models/table"

type UserRequest struct {
	Firstname string `json:"firstName" binding:"required"`
	Lastname  string `json:"lastName" binding:"required"`
	ClientID  uint   `json:"userId" binding:"required"`
}

type AdminLoginRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	GroupUUID string `json:"groupId" binding:"required"`
}

func (req *UserRequest) ToUser() table.User {
	var user table.User
	user.ClientID = req.ClientID
	user.Firstname = req.Firstname
	user.Lastname = req.Lastname
	return user
}

func (req *AdminLoginRequest) ToUser() table.User {
	var user table.User
	user.Email = req.Email
	user.Password = req.Password
	return user
}
