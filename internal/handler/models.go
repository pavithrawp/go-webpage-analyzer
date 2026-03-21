package handler

// JSON request body for the analyze endpoint
type analyzeRequest struct {
	URL string `json:"url"`
}
