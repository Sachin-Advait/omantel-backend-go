package business

import (
	"bytes"
	"crypto/rand"
	"encoding/xml"
	"fmt"
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

var HttpClient *http.Client

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
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soapenv string   `xml:"soapenv,attr"`
	Body    struct {
		Text               string `xml:",chardata"`
		SendNotificationRs struct {
			Text           string `xml:",chardata"`
			NS1            string `xml:"NS1,attr"`
			ResponseHeader struct {
				Text          string `xml:",chardata"`
				FunctionId    string `xml:"FunctionId"`
				TransactionId string `xml:"TransactionId"`
				RequestId     string `xml:"RequestId"`
				ReturnCode    string `xml:"ReturnCode"`
				ReturnDesc    string `xml:"ReturnDesc"`
			} `xml:"ResponseHeader"`
		} `xml:"SendNotificationRs"`
	} `xml:"Body"`
}

func (impl *IncentiveImpl) GenerateOTP(ctx *gin.Context, clientID string) (map[string]string, *errors.Error) {

	var (
		isUpdate      = true
		contactNumber string
	)

	mobile, otpModel, err := impl.PsqlDb.FetchOTPValidity(clientID)
	if err != nil {
		return nil, err
	}

	contactNumber = mobile

	if otpModel.StaffID == "" {
		name, mobile, err := impl.PsqlDb.FetchUserName(clientID)
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

		if contactNumber == "" {
			contactNumber = mobile
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
	otpModel.OTP = generateOTP(6)
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

	if err := impl.callSOAPClientSteps(clientID, contactNumber, otpModel.OTP); err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "SOAP CALL FAILED",
			Details:    "sms call failed: " + err.Error(),
		}
	}

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

func (impl *IncentiveImpl) callSOAPClientSteps(userID, mobile string, otp int) error {

	req := populateRequest(userID, mobile, otp)

	httpReq, err := impl.generateSOAPRequest(req)
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

	id := strconv.Itoa(int(time.Now().Unix()))

	req := Request{
		UserId:                      userId,
		FunctionId:                  "17002004",
		TransactionId:               id,
		RequestId:                   "OMSIP_" + id,
		RequestorChannelId:          models.RequestorChannelId,
		RequestorUserId:             "OMSIP",
		Timestamp:                   time.Now().Format("2006-01-02T15:04:05"),
		NotificationType:            "SMS",
		NotificationDateTime:        time.Now().Format("2006-01-02T15:04:05"),
		NotificationLanguage:        "1",
		ShortCode:                   models.ShortCode,
		IsDeliveryStatusRequired:    "1",
		NotificaitionValidityPeriod: time.Now().Format("2006-01-02T15:04:05"),
		Body:                        fmt.Sprintf("Dear %v, Your OTP is %v", userId, otp),
		BodyCoding:                  "8",
		BodyContentType:             "1",
		Receiver:                    mobile,
	}

	return &req
}

func (impl *IncentiveImpl) generateSOAPRequest(req *Request) (*http.Request, error) {
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

	r, err := http.NewRequest(http.MethodPost, impl.Config.OTPUrl, bytes.NewBuffer(doc.Bytes()))
	if err != nil {
		return nil, errors.New("request-" + err.Error())
	}

	return r, nil
}

func soapCall(req *http.Request) (*Response, error) {

	resp, err := HttpClient.Do(req)
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

	if r.Body.SendNotificationRs.ResponseHeader.ReturnCode != "0" {
		return nil, errors.New(r.Body.SendNotificationRs.ResponseHeader.ReturnDesc)
	}

	return r, nil
}
