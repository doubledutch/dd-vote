package req

// VoteCreateRequest is sent by the client when voting
type VoteCreateRequest struct {
	Value int `json:"value" binding:"required"`
}

// IsValid returns whether the request has valid data
func (req *VoteCreateRequest) IsValid() bool {
	return req.Value == 1 || req.Value == -1
}
