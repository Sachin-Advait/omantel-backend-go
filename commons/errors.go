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

// var BadRequestError = GenericError{Code: BAD_REQUEST_CODE, ErrorMessage: "Invalid arguments to Api"}
// var AuthenticationError = GenericError{Code: UNAUTHORISED_REQUEST, ErrorMessage: "Authorisation Failed"}
// var JsonDeserializationFailed = GenericError{Code: BAD_REQUEST_CODE, ErrorMessage: "JSON deserialization of arguments failed"}
// var CacheMissError = GenericError{ErrorMessage: "Cache miss"}
// var KeyRefreshFailed = GenericError{ErrorMessage: "Refresh key operation failed"}

var ErrorMap = map[int]string{
	1001: "INTERNAL SERVER ERROR",
	1002: "BAD REQUEST",
	1003: "CONTACT NUMBER NOT FOUND",
	1020: "INCORRECT CONTACT NUMBER FORMAT",
	1021: "OTP HAS EXPIRED",
	1022: "OTP IS INVALID",
}
