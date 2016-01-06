package table

// Vote database model for user votes on questions
type Vote struct {
	BaseModel
	PostID   uint   `json:"-" sql:"index;unique_index:idx_postid_userid"`
	PostUUID string `json:"post_uuid" sql:"index"` //TODO this is duplicated from Post table and seems like bad design
	UserID   uint   `json:"-" sql:"index;unique_index:idx_postid_userid"`
	Value    int    `json:"value"` // -1 or 1
}
