package models

type Comment struct {
	BaseModel
	Comment string `json:"comment" binding:"required"`

	// hidden fields
	PostID uint `json:"-" sql:"index"`
	UserID uint `json:"-" sql:"index"`

	// associations
	User User
}
