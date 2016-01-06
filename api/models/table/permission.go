package table

// Permission database model for
type Permission struct {
	BaseModel
	UserID   uint   `sql:"index"`
	Role     uint   // (SuperAdmin, EventAdmin, GroupAdmin) - see roles in users package
	Metadata string // could be a specific group name for GroupAdmin, or pattern for EventAdmin

	// associations
	User User
}
