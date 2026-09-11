package dto

// APIResponse standard success response format
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// APIErrorResponse standard error response format
type APIErrorResponse struct {
	Success bool        `json:"success" example:"false"`
	Message string      `json:"message" example:"Validation error"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string, errors interface{}) APIErrorResponse {
	return APIErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}
