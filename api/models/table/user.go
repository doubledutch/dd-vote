package table

// User database model includes users from logged in with client credentials
// (no password) and admins (with emails and passwords)
type User struct {
	BaseModel
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`

	// hidden fields
	Email    string `json:"-"`
	Password string `json:"-"`
	ClientID uint   `json:"-" sql:"unique_index" `
}
