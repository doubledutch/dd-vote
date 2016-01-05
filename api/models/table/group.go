package table

type Group struct {
	BaseModel
	Name string `json:"name" sql:"unique_index" form:"name" binding:"required"`

	// associations
	Posts []Post
}
