package delivery

type Response struct {
	Quote    string `json:"quote,omitempty"`
	Author   string `json:"author,omitempty"`
	Category string `json:"category,omitempty"`
	Error    string `json:"error,omitempty"`
}
