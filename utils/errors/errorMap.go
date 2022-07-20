package errors

var ErrorCodeMap = map[string]string{
	InternalServerErrorCode:   InternalServerError,
	InvalidReqParameterCode:   InvalidReqParameters,
	RequestDecodeErrorCode:    RequestDecodeError,
	ContradictoryReqParamCode: ContradictoryReqParam,
}
