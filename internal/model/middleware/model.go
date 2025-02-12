package middleware

// ErrorResponse defines the structure for API error responses.
//
//	@Description Standard API error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message" example:"Detailed error message"`
}
