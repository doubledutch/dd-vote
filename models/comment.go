package models

type Comment struct {
	BaseModel
	PostID  uint   `json:"-" sql:"index"`
	UserID  uint   `json:"-" sql:"index"`
	Comment string `json:"comment"`
}
