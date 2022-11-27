package main

import (
	"encoding/xml"
	"fmt"
)

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

func main() {

	str := `<?xml version="1.0" encoding="utf-8"?><soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"><soapenv:Body><NS1:SendNotificationRs xmlns:NS1="http://www.omantel.om/esb/messages/sendnotification"><NS1:ResponseHeader><FunctionId>17002004</FunctionId><TransactionId>36896</TransactionId><RequestId>OMSIP_1669540763</RequestId><ReturnCode>0</ReturnCode><ReturnDesc>Success/ Successfully validated the request</ReturnDesc></NS1:ResponseHeader></NS1:SendNotificationRs></soapenv:Body></soapenv:Envelope>`

	var i Response

	if err := xml.Unmarshal([]byte(str), &i); err != nil {
		fmt.Println("err: " + err.Error())
		return
	}

	fmt.Printf("success: %+v\n", i)
}
