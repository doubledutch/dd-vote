package req

// VoteCreateRequest is sent by the client when voting
type VoteCreateRequest struct {
	PostUUID string `json:"post_id" binding:"required"`
	Value    int    `json:"value" binding:"required"`
}
