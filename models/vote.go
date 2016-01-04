package models

type Vote struct {
	BaseModel
	PostID uint `json:"post_id" sql:"index;unique_index:idx_postid_userid"`
	UserID uint `json:"user_id" sql:"index;unique_index:idx_postid_userid"`
	Value  int  `json:"value"`
}

type VoteCreateRequest struct {
	PostUUID string `json:"post_id" binding:"required"`
	Value    int    `json:"value" binding:"required"`
}
