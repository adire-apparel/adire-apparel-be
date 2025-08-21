package lib

type ApiError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func NewApiErrror(message string, status int) *ApiError {
	return &ApiError{
		Message: message,
		Status:  status,
	}
}

func (s *ApiError) Error() string {
	return s.Message
}
