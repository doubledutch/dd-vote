package resp

// APIResponse is the wrapper around all api responses. If there is an error,
// then a Message is sent. If there is no error, then Value is sent.
type APIResponse struct {
	IsError bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}
