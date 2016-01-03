package models

type Error struct {
	IsError bool   `json:"error"`
	Message string `json:"message"`
}
