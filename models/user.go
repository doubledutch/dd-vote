package models

type User struct {
	BaseModel
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`

	// hidden fields
	Email    string `json:"-"`
	Password string `json:"-"`
	ClientID uint   `json:"-" sql:"unique_index" `
}

type UserRequest struct {
	Firstname string `json:"firstName" binding:"required"`
	Lastname  string `json:"lastName" binding:"required"`
	ClientID  uint   `json:"userId" binding:"required"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
