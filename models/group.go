package models

type Group struct {
	BaseModel
	Name  string `json:"name" sql:"unique_index" form:"name"`
	Posts []Post
}
