package table

// Group database model for groups/pages in the /g/:group naming scheme
type Group struct {
	BaseModel
	Name string `json:"name" sql:"unique_index" form:"name" binding:"required"`

	// associations
	Posts []Post
}

// IsValidForCreate returns whether the Group object has valid data
func (group *Group) IsValidForCreate() bool {
	return len(group.Name) > 0 && len(group.Name) <= 200
}
