package models

import "time"

type BaseModel struct {
	ID        uint `gorm:"primary_key" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
