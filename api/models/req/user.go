package req

type UserRequest struct {
	Firstname string `json:"firstName" binding:"required"`
	Lastname  string `json:"lastName" binding:"required"`
	ClientID  uint   `json:"userId" binding:"required"`
}

type AdminLoginRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	GroupUUID string `json:"groupId" binding:"required"`
}
