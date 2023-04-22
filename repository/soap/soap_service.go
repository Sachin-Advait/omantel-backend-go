package soap

type SOAPService interface {
	CallSOAPClientSteps(userID, userName, mobile, url string, otp int) error
}
