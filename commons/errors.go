package commons

type GenericError struct {
	ErrorMessage string
	Code         int
}

const (
	BAD_REQUEST_CODE      = 400
	UNAUTHORISED_REQUEST  = 401
	INTERNAL_SERVER_ERROR = 500
)

var BadRequestError = GenericError{Code: BAD_REQUEST_CODE, ErrorMessage: "Invalid arguments to Api"}
var AuthenticationError = GenericError{Code: UNAUTHORISED_REQUEST, ErrorMessage: "Authorisation Failed"}
var JsonDeserializationFailed = GenericError{Code: BAD_REQUEST_CODE, ErrorMessage: "JSON deserialization of arguments failed"}
var CacheMissError = GenericError{ErrorMessage: "Cache miss"}
var KeyRefreshFailed = GenericError{ErrorMessage: "Refresh key operation failed"}
