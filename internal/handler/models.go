package handler

// JSON request body for the analyze endpoint
type analyzeRequest struct {
	URL string `json:"url"`
}

// errorResponse represents a standard JSON error response
type errorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
