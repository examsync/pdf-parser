package errors

// ErrorResponse represents the standardized JSON error envelope returned by all HTTP handlers.
type ErrorResponse struct {
	Success bool         `json:"success"`
	Error   ErrorDetails `json:"error"`
}

// ErrorDetails represents the granular details of an API error payload.
type ErrorDetails struct {
	Code          string      `json:"code"`
	Message       string      `json:"message"`
	Details       interface{} `json:"details,omitempty"`
	Timestamp     string      `json:"timestamp"`
	Path          string      `json:"path"`
	InternalError string      `json:"internal_error,omitempty"`
}
