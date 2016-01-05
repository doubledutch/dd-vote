package req

type VoteCreateRequest struct {
	PostUUID string `json:"post_id" binding:"required"`
	Value    int    `json:"value" binding:"required"`
}
