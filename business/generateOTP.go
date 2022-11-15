package business

import (
	"bytes"
	"crypto/rand"
	"encoding/xml"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

type Request struct {
	//Values are set in below fields as per the request
	UserId                      string
	FunctionId                  string
	TransactionId               string
	RequestId                   string
	RequestorChannelId          string
	RequestorUserId             string
	Timestamp                   string
	NotificationType            string
	NotificationDateTime        string
	NotificationLanguage        string
	ShortCode                   string
	IsDeliveryStatusRequired    string
	NotificaitionValidityPeriod string
	Body                        string
	BodyCoding                  string
	BodyContentType             string
	Receiver                    string
	OTP                         string
}

type Response struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody *SOAPBodyResponse
}

type SOAPBodyResponse struct {
	XMLName      xml.Name `xml:"Body"`
	Resp         *ResponseBody
	FaultDetails *Fault
}

type Fault struct {
	XMLName     xml.Name `xml:"Fault"`
	Faultcode   string   `xml:"faultcode"`
	Faultstring string   `xml:"faultstring"`
}

type ResponseBody struct {
	XMLName   xml.Name `xml:"Response"`
	RequestID string   `xml:"RequestID"`
	Response  *Body
	Status    string `xml:"Status"`
}

type Body struct {
	XMLName       xml.Name `xml:"Body"`
	Status        string   `xml:"Status"`
	Salary        string   `xml:"Salary"`
	Designation   string   `xml:"Designation"`
	Manager       string   `xml:"Manager"`
	Company       string   `xml:"Company"`
	EmployedSince string   `xml:"EmployedSince"`
}

func (impl *IncentiveImpl) GenerateOTP(ctx *gin.Context, clientID string) (map[string]string, *errors.Error) {

	var isUpdate = true

	otpModel, err := impl.PsqlDb.FetchOTPValidity(clientID)
	if err != nil {
		return nil, err
	}

	if otpModel.StaffID == "" {
		name, err := impl.PsqlDb.FetchUserName(clientID)
		if err != nil {
			return nil, err
		}

		if name == "" {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.UserNotFoundCode,
				Message:    errors.UserIDInvalid,
				Details:    "userId is invalid",
			}
		}

		isUpdate = false
	}
	otpModel.StaffID = clientID

	// contactNumber, err := impl.PsqlDb.FetchEmployeeContactNumber(clientID)
	// if err != nil {
	// 	return nil, err
	// }

	// if otpModel.StaffID != 0 && !otpModel.IsActiveUser {
	// 	return nil, &errors.Error{
	// 		StatusCode: http.StatusForbidden,
	// 		Code:       errors.UserBlockedCode,
	// 		Message:    errors.UserBlockedError,
	// 		Details:    "user is blocked",
	// 	}
	// }

	//TODO: reset loginAttemptCount
	// if otpModel.LoginAttemptCount > constant.MaxAllowedLoginAttempt {
	// 	return nil, &errors.Error{
	// 		StatusCode: http.StatusForbidden,
	// 		Code:       errors.RetryLoginLimitReachedCode,
	// 		Message:    errors.RetryLoginLimitReached,
	// 		Details:    "retry login limit reached",
	// 	}
	// }

	//otpLastUpdateTime updated only while generating OTP more than block time
	if time.Since(otpModel.OtpLastUpdateTime).Minutes() <= constant.OtpGenerateBlockTimeInMin {

		//too many times sent within 15 mins
		if otpModel.OTPResendCount > constant.OtpRegenerateAllowedCount {
			return nil, &errors.Error{
				StatusCode: http.StatusForbidden,
				Code:       errors.UserBlockedCode,
				Message:    errors.OTPLimitError,
				Details:    "max otp resend reached",
			}
		}

		//too many failure attempts within 15 mins
		if otpModel.OTPVerifyFailedCount >= constant.MaxAllowedLoginAttempt {
			return nil, &errors.Error{
				StatusCode: http.StatusForbidden,
				Code:       errors.UserBlockedCode,
				Message:    errors.RetryLoginLimitReached,
				Details:    "incorrect otp attempt reached",
			}
		}

		otpModel.OTPResendCount++
	} else {
		//more than 15 mins hence restart counter
		otpModel.OTPResendCount = 0
		otpModel.OTPVerifyFailedCount = 0
		otpModel.OtpLastUpdateTime = time.Now()
	}

	//TODO: generate random num as OTP
	otpModel.OTP = 123456 //generateOTP(6)
	otpModel.OTPUsed = false
	otpModel.OtpCreateTime = time.Now()

	otpModel.UpdateReason = "otp generated"
	otpModel.UpdatedTime = time.Now()

	if otpModel.OTPResendCount > 0 {
		otpModel.UpdateReason = "otp re-generated"
	}

	//store OTP in db
	if err := impl.PsqlDb.InsertStaffOTP(isUpdate, otpModel); err != nil {
		return nil, err
	}

	// reqSMS := &restapi.SMSParams{
	// 	HeaderRequestID:             requestID,
	// 	HeaderStaffID:               strconv.Itoa(clientID),
	// 	HeaderTimeStamp:             otpModel.CreatedTime.Format("2006-01-02 15:04:05"),
	// 	NotificationDateTime:        otpModel.CreatedTime.Format("2006-01-02 15:04:05"),
	// 	NotificaitionValidityPeriod: otpModel.CreatedTime.Add(3 * time.Minute).Format("2006-01-02 15:04:05"),
	// 	Message:                     "Your OTP is: " + strconv.Itoa(otpModel.OTP),
	// 	ContactNumber:               strconv.Itoa(contactNumber),
	// }

	// if err = restapi.SendOTP(reqSMS); err != nil {
	// 	return err
	// }

	// if err := callSOAPClientSteps(clientID, "", otpModel.OTP); err != nil {
	// 	return nil, &errors.Error{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Code:       errors.InternalServerError,
	// 		Message:    "SOAP CALL FAILED",
	// 		Details:    "sms call failed: " + err.Error(),
	// 	}
	// }

	return map[string]string{"error": "nil"}, nil
}

func generateOTP(max int) int {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	otp, _ := strconv.Atoi(string(b))

	return otp
}

func callSOAPClientSteps(userID, mobile string, otp int) error {

	req := populateRequest(userID, mobile, otp)

	httpReq, err := generateSOAPRequest(req)
	if err != nil {
		return err
	}

	_, err = soapCall(httpReq)
	if err != nil {
		return err
	}

	return nil
}

func populateRequest(userId, mobile string, otp int) *Request {

	req := Request{
		UserId:   userId,
		OTP:      strconv.Itoa(otp),
		Receiver: mobile,
	}

	return &req
}

func generateSOAPRequest(req *Request) (*http.Request, error) {
	// Using the var getTemplate to construct request
	template, err := template.New("smsRequest").Parse(string(models.SMSRequestTemplate))
	if err != nil {
		return nil, errors.New("parsing-" + err.Error())
	}

	doc := &bytes.Buffer{}
	// Replacing the doc from template with actual req values
	err = template.Execute(doc, req)
	if err != nil {
		return nil, errors.New("execute-" + err.Error())
	}

	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	err = encoder.Encode(doc.String())
	if err != nil {
		return nil, errors.New("encode-" + err.Error())
	}

	r, err := http.NewRequest(http.MethodPost, "http://www.omantel.om/esb/messages/sendnotification", bytes.NewBuffer([]byte(doc.String())))
	if err != nil {
		return nil, errors.New("request-" + err.Error())
	}

	return r, nil
}

func soapCall(req *http.Request) (*Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, errors.New("call-" + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("readall-" + err.Error())
	}
	defer resp.Body.Close()

	r := &Response{}
	err = xml.Unmarshal(body, &r)
	if err != nil {
		return nil, errors.New("unmarshall-" + err.Error())
	}

	if r.SoapBody.Resp.Status != "200" {
		return nil, err
	}

	return r, nil
}
