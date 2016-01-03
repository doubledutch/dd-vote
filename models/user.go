package models

type User struct {
	BaseModel
	Name     string `json:"name"`
	ClientID uint   `json:"-" sql:"unique_index"`
}
