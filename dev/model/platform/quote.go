package platform

type Quote struct {
	QuoteResponse
	QuoteError
}

type QuoteResponse struct {
	Quote    string `json:"quote"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

type QuoteError struct {
	Error string `json:"error"`
}
