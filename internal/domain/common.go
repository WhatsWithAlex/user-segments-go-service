package domain

type CommonResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func ErrorResponse(message string) CommonResponse {
	return CommonResponse{
		Message: message,
		Success: false,
	}
}

func SuccessResponse(message string) CommonResponse {
	return CommonResponse{
		Message: message,
		Success: true,
	}
}
