package models

type Post struct {
	BaseModel
	Name      string `json:"name" sql:"unique_index:idx_name_groupid"`
	GroupID   uint   `json:"-" sql:"index;unique_index:idx_name_groupid"`
	CreatedBy uint   `json:"-"`
	Upvotes   uint   `json:"upvotes"`
	Downvotes uint   `json:"downvotes"`
}
