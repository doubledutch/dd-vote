package models

type Permission struct {
	BaseModel
	UserID   uint
	Role     uint // see roles in users package
	Metadata string

	// associations
	User User
}
