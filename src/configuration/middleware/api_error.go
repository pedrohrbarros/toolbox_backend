package api_error

import "net/http"

type ApiError struct {
	Message string     `json:"message"`
	Err     string     `json:"error"`
	Code    int        `json:"code"`
	Causes  []Causes `json:"causes"`
}

type Causes struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r*ApiError) Error() string {
	return r.Message
}

func NewApiError(message string, err string, code int, causes []Causes) *ApiError {
	return &ApiError{
		Message: message,
		Err:     err,
		Code:    code,
		Causes:  causes,
	}
}

func NewBadRequestError(message string) *ApiError {
	return &ApiError{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
	}
}

func NewBadRequestValidationError(message string, causes []Causes) *ApiError {
	return &ApiError {
		Message: message,
		Err: "bad_request",
		Code: http.StatusBadRequest,
		Causes: causes,
	}
}

func NewInternalServerError(message string) *ApiError {
	return &ApiError{
		Message: message,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	}	
}

func NewTooManyRequestsError(message string) *ApiError {
	return &ApiError {
		Message: message,
		Err: "too_many_requests",
		Code: http.StatusTooManyRequests,
	}
}