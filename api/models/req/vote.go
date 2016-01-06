package req

// VoteCreateRequest is sent by the client when voting
type VoteCreateRequest struct {
	Value int `json:"value" binding:"required"`
}
