package errors

const (
	SuccessInfoMessage   = "SUCCESS"
	SuccessInfoErrorCode = ""
)

// Error Code
const (
	//common error code
	InternalServerErrorCode   = "CE10001"
	InvalidReqParameterCode   = "CE10002"
	MethodNotAllowedCode      = "CE10003"
	RequestDecodeErrorCode    = "CE10004"
	ContradictoryReqParamCode = "CE10005"

	//user error code
	UserBlockedCode            = "UE11007"
	UserExistsErrorCode        = "UE11001"
	DuplicateUserCode          = "UE11002"
	LoginLimitReachedCode      = "UE11009"
	RetryLoginLimitReachedCode = "UE11010"
	InvalidDetailsCode         = "UE11015"

	UserNotFoundCode        = "UE11043"
	UnauthorizedUserCode    = "UE11042"
	IDPassInvalidCode       = "UE11006"
	InvalidPassCode         = "UE11008"
	InvalidTokenCode        = "UE11011"
	TokenExpiredCode        = "UE11012"
	InvalidTokenDetailsCode = "UE11013"
	SimilarLast3PassCode    = "UE11014"
	OTPUsedCode             = "UE11022"
	OTPInvalidCode          = "UE11019"
	OTPExpiredCode          = "UE11023"

	DataNotFoundCode = "UE11030"

	RenewTokenExpiredCode = "UE11042"

	//other error
	OtherErrorCode = "OE11046"
)

// Error Values
const (
	InternalServerError   = "INTERNAL SERVER ERROR"
	InvalidReqParameters  = "INVALID REQUEST PARAMETERS"
	MethodNotAllowed      = "METHOD NOT ALLOWED"
	RequestDecodeError    = "REQUEST DECODING ERROR"
	ContradictoryReqParam = "CONTRADICTORY REQUEST PARAM"

	UserBlockedError       = "USER IS BLOCKED"
	UserExistsError        = "USER ALREADY EXISTS"
	DuplicateUserError     = "DUPLICATE USER FOUND"
	LoginLimitReached      = "MAX SUCCESS LOGIN LIMIT REACHED"
	RetryLoginLimitReached = "MAX INCORRECT LOGIN LIMIT REACHED, BLOCKED FOR 15 MINS"
	InvalidDetails         = "DETAILS ARE INCORRECT"
	TokenNotFound          = "TOKEN NOT FOUND"
	InvalidToken           = "INVALID TOKEN"
	TokenExpired           = "TOKEN EXPIRED"
	InvalidTokenDetails    = "TOKEN DETAILS ARE INVALID"
	OTPUsedError           = "OTP IS ALREADY USED"
	OTPInvalidError        = "OTP IS INVALID"
	OTPExpiredError        = "OTP IS EXPIRED"
	OTPLimitError          = "TOO MANY TIMES OTP SENT, BLOCKED FOR 15 MINS"

	UserIDInvalid                = "USER ID IS INVALID"
	UserTypeInvalid              = "USER TYPE IS INVALID"
	UserNotFound                 = "USER NOT FOUND"
	UnauthorizedUser             = "UNAUTHORIZED USER"
	UserNameIncorrect            = "USER NAME SHOULD NOT BE BLANK OR SPECIAL CHARS AND LENGTH LESS THAN 15"
	FirstNameIncorrect           = "FIRST NAME SHOULD NOT BE BLANK OR ALPHANUMERIC AND LENGTH LESS THAN 12"
	LastNameIncorrect            = "LAST NAME SHOULD NOT BE BLANK OR ALPHANUMERIC AND LENGTH LESS THAN 12"
	FullNameIncorrect            = "FULL NAME COMBINED HAS INCORRECT LENGTH"
	EmailFormatIncorrect         = "EMAIL FORMAT IS INCORRECT OR LENGTH IS MORE THAN 40"
	EmailIncorrect               = "EMAIL IS INCORRECT"
	ContactNumberFormatIncorrect = "CONTACT NUMBER SHOULD NOT BE BLANK OR ALPHANUMERIC AND LENGTH LESS THAN 15"
	ContactNumberIncorrect       = "CONTACT NUMBER IS INCORRECT"
	PasswordFormatIncorrect      = "PASSWORD SHOULD BE ALPHA NUMERIC AND LENGTH BETWEEN 6 TO 12"
	PasswordIncorrect            = "PASSWORD IS INCORRECT"

	SimilarLast3PassWord = "NEW PASSWORD SHOULD BE DIFFERENT THAN THE LAST 3 PASSWORDS"

	DataNotFound = "DATA NOT FOUND"
)

var ErrorMap = map[int]string{
	1001: "INTERNAL SERVER ERROR",
	1002: "BAD REQUEST",
	1003: "CONTACT NUMBER NOT FOUND",
	1004: "NO RECORDS FOUND",
	1020: "INCORRECT CONTACT NUMBER FORMAT",
	1021: "OTP HAS EXPIRED",
	1022: "OTP IS INVALID",
}

type GenericError struct {
	ErrorMessage string
	Code         int
}

const (
	BAD_REQUEST_CODE      = 400
	UNAUTHORISED_REQUEST  = 401
	INTERNAL_SERVER_ERROR = 500
)
