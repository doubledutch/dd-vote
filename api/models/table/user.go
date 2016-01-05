package table

type User struct {
	BaseModel
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`

	// hidden fields
	Email    string `json:"-"`
	Password string `json:"-"`
	ClientID uint   `json:"-" sql:"unique_index" `
}
