package core_http_response

type ErrorResponse struct {
	Error   string `json:"error"   example:"system error message: error type"`
	Message string `json:"message" example:"human readable message"`
}
