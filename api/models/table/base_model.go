package table

import "time"

// BaseModel includes columns that should be in every table, including a
// primary key, timestamps for creation and updating, and soft deletes
type BaseModel struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}
