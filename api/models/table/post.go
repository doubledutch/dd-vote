package table

// Post (Question) database model for questions asked in a group
type Post struct {
	BaseModel
	Name      string `json:"name" sql:"unique_index:idx_name_groupid" binding:"required"`
	UUID      string `json:"uuid" sql:"unique_index"`
	Upvotes   uint   `json:"upvotes"`
	Downvotes uint   `json:"downvotes"`

	// hidden fields
	GroupID   uint `json:"-" sql:"index;unique_index:idx_name_groupid"`
	CreatedBy uint `json:"-"`

	// associations
	Comments []Comment `json:"comments"`
	User     User      `json:"-" gorm:"ForeignKey:CreatedBy"`
	Group    Group     `json:"-"`
}

// IsValidForCreate returns whether the Post object has valid data
func (post *Post) IsValidForCreate() bool {
	return len(post.Name) > 0 && len(post.Name) <= 140
}
