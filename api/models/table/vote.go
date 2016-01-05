package table

type Vote struct {
	BaseModel
	PostID   uint   `json:"-" sql:"index;unique_index:idx_postid_userid"`
	PostUUID string `json:"post_uuid" sql:"index"` //TODO this is duplicated from Post table and seems like bad design
	UserID   uint   `json:"-" sql:"index;unique_index:idx_postid_userid"`
	Value    int    `json:"value"`
}
