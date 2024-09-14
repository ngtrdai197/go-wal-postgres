package error

// Define global errors
var (
	ErrBadRequest = NewError(400, "Bad Request")
	ErrNotFound   = NewError(404, "Not Found")
	ErrInternal   = NewError(500, "Internal Server Error")

	ErrUUIDInvalid = NewError(400, "Invalid UUID")
)
