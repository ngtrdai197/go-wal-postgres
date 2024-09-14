package error

type DError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewError(statusCode int, message string) *DError {
	return &DError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e *DError) Error() string {
	return e.Message
}
