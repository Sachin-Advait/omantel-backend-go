package soap

type SOAPService interface {
	CallSOAPClientSteps(userID, mobile, url string, otp int) error
}
