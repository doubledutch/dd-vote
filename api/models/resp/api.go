package resp

type ApiResponse struct {
	IsError bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}
