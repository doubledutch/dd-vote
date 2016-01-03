package models

type User struct {
	BaseModel
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	ClientID  uint   `json:"-" sql:"unique_index" `
}

type UserRequest struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	ClientID  uint   `json:"client_id" binding:"required"`
}
