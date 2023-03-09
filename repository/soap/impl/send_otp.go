package impl

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
)

func (impl *SOAPServiceImpl) CallSOAPClientSteps(userID, mobile, url string, otp int) error {

	req := impl.populateRequest(userID, mobile, otp)

	httpReq, err := impl.generateSOAPRequest(url, req)
	if err != nil {
		return err
	}

	_, err = impl.soapCall(httpReq)
	if err != nil {
		return err
	}

	return nil
}

func (impl *SOAPServiceImpl) generateSOAPRequest(url string, req *models.Request) (*http.Request, error) {
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

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(doc.Bytes()))
	if err != nil {
		return nil, errors.New("request-" + err.Error())
	}

	r.Header.Set("user-agent", "sampleTest")
	r.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	r.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("SA-OmantelSIP:Des!B@buEnglishMem@7533")))
	r.Header.Set("SOAPAction", "SendNotification")

	return r, nil
}

func (impl *SOAPServiceImpl) populateRequest(userId, mobile string, otp int) *models.Request {

	id := strconv.Itoa(int(time.Now().Unix()))

	req := models.Request{
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

func (impl *SOAPServiceImpl) soapCall(req *http.Request) (*models.Response, error) {

	resp, err := impl.HttpClient.Do(req)
	if err != nil {
		return nil, errors.New("call-" + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("readall-" + err.Error())
	}
	defer resp.Body.Close()

	r := &models.Response{}
	err = xml.Unmarshal(body, &r)
	if err != nil {
		return nil, errors.New("unmarshall-" + err.Error())
	}

	if r.Body.SendNotificationRs.ResponseHeader.ReturnCode != "0" {
		return nil, errors.New(r.Body.SendNotificationRs.ResponseHeader.ReturnDesc)
	}

	return r, nil
}
