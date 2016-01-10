package table

// Comment database model for comments on questions
type Comment struct {
	BaseModel
	Comment string `json:"comment" binding:"required"`

	// hidden fields
	PostID uint `json:"-" sql:"index"`
	UserID uint `json:"-" sql:"index"`

	// associations
	User User `json:"user"`
}

// IsValidForCreate returns whether the Comment object has valid data
func (comment *Comment) IsValidForCreate() bool {
	return len(comment.Comment) > 0 && len(comment.Comment) <= 500
}
