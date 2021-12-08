package restapi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
)

type SMSParams struct {
	HeaderRequestID             string
	HeaderStaffID               string
	HeaderTimeStamp             string
	NotificationDateTime        string
	NotificaitionValidityPeriod string
	Message                     string
	ContactNumber               string
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

func SendOTP(sms *SMSParams) *commons.GenericError {

	req, err := generateSOAPRequest(sms)
	if err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: err.Error()}
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: "external api cal failed " + err.Error()}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: "external api readBody failed " + err.Error()}
	}

	defer resp.Body.Close()

	r := &Response{}
	err = xml.Unmarshal(body, &r)
	if err != nil {
		return &commons.GenericError{Code: 1001, ErrorMessage: "external api " + err.Error()}
	}

	if r.Body.SendNotificationRs.ResponseHeader.ReturnCode != "200" {
		return &commons.GenericError{Code: 1001, ErrorMessage: "external api response invalid"}
	}

	return nil
}

func generateSOAPRequest(smsReq *SMSParams) (*http.Request, error) {
	// Using the var getTemplate to construct request
	template, err := template.New("sms_" + smsReq.HeaderStaffID).Parse(commons.GetTemplate)
	if err != nil {
		return nil, fmt.Errorf("error marshling template object. %s", err.Error())
	}

	doc := &bytes.Buffer{}
	// Replacing the doc from template with actual req values
	err = template.Execute(doc, smsReq)
	if err != nil {
		return nil, fmt.Errorf("template execute error %s", err.Error())
	}

	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	err = encoder.Encode(doc.String())
	if err != nil {
		return nil, fmt.Errorf("encoder Encode error. %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, "https://pki.omantel.om:17862/services/SendNotification", bytes.NewBuffer([]byte(doc.String())))
	if err != nil {
		return nil, fmt.Errorf("error making a request. %s", err.Error())
	}

	req.SetBasicAuth("SA_ZoneAlert", "Omantel@123")

	// ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	// defer cancel()

	req.Header.Set("user-agent", "sampleTest")
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")

	// req = req.WithContext(ctx)

	return req, nil
}

//https://pki.omantel.om:17862/services/SendNotification

//credentials SA_ZoneAlert:Omantel@123

// In soap request basically we are passing xml data (it have phoneNumber, message). soap require encoded credentials instead of plain text. it also require agent(const agent = new https.Agent({
// 	rejectUnauthorized: false
// });)

// // this is our soap request after creating xml, agent and encoded credentials.
// const { response } = await soapRequest({
// 		url: url, headers: {
// 			'user-agent': 'sampleTest',
// 			'Content-Type': 'text/xml;charset=UTF-8',
// 			Authorization: `Basic ${encodedCredentials}`,
// 			SOAPAction: 'SendNotification',
// 		}, xml: xml, extraOpts: {
// 			httpsAgent: agent
// 		}
// 	});
