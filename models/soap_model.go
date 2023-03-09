package models

import "encoding/xml"

var (
	ShortCode          = "OMSIP"
	RequestorChannelId = "OMSIP"
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

var SMSRequestTemplate = []byte(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:sen="http://www.omantel.om/esb/messages/sendnotification">
<soapenv:Header/>
<soapenv:Body>
   <sen:SendNotificationRq>
	  <sen:RequestHeader>
		 <FunctionId>{{.FunctionId}}</FunctionId>
		 <TransactionId>{{.TransactionId}}</TransactionId>
		 <RequestId>{{.RequestId}}</RequestId>
		 <RequestorChannelId>{{.RequestorChannelId}}</RequestorChannelId>
		 <RequestorUserId>{{.RequestorUserId}}</RequestorUserId>
		 <Timestamp>{{.Timestamp}}</Timestamp>
	  </sen:RequestHeader>
	  <sen:RequestBody>
		 <sen:Notification>
			<NotificationType>SMS</NotificationType>
			<NotificationDateTime>{{.NotificationDateTime}}</NotificationDateTime>
			<NotificationLanguage>1</NotificationLanguage>
			<ShortCode>{{.ShortCode}}</ShortCode>
			<IsDeliveryStatusRequired>1</IsDeliveryStatusRequired>
			<NotificaitionValidityPeriod>{{.NotificaitionValidityPeriod}}</NotificaitionValidityPeriod>
			<Body>{{.Body}}</Body>
			<BodyCoding>8</BodyCoding>
			<BodyContentType>1</BodyContentType>
			<ReceiverList>
			   <Receiver>{{.Receiver}}</Receiver>
			</ReceiverList>
		 </sen:Notification>
	  </sen:RequestBody>
   </sen:SendNotificationRq>
</soapenv:Body>
</soapenv:Envelope>`)

// var SMSResponseTemplate = []byte(`<?xml version="1.0" encoding="UTF-8"?>
// <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
//    <soapenv:Body>
//       <NS1:SendNotificationRs xmlns:NS1="http://www.omantel.om/esb/messages/sendnotification">
//          <NS1:ResponseHeader>
//             <FunctionId>17002004</FunctionId>
//             <TransactionId>36896</TransactionId>
//             <RequestId>OMSIP_1669540763</RequestId>
//             <ReturnCode>0</ReturnCode>
//             <ReturnDesc>Success/ Successfully validated the request</ReturnDesc>
//          </NS1:ResponseHeader>
//       </NS1:SendNotificationRs>
//    </soapenv:Body>
// </soapenv:Envelope>`)

var GetTemplate = `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:sen="http://www.omantel.om/esb/messages/sendnotification">
   <soapenv:Header/>
   <soapenv:Body>
      <sen:SendNotificationRq>
         <sen:RequestHeader>
            <FunctionId>17002004</FunctionId>
            <TransactionId>0</TransactionId>
            <RequestId>{{.HeaderRequestID}}</RequestId>
            <RequestorChannelId>{{.HeaderRequestID}}</RequestorChannelId>
            <RequestorUserId>{{.HeaderStaffID}}</RequestorUserId>
            <Timestamp>{{.HeaderTimeStamp}}</Timestamp>
         </sen:RequestHeader>
         <sen:RequestBody>
            <sen:Notification>
               <NotificationType>SMS</NotificationType>
               <NotificationDateTime>{{.NotificationDateTime}}</NotificationDateTime>
               <NotificationLanguage>1</NotificationLanguage>
               <ShortCode>Omantel</ShortCode>
               <IsDeliveryStatusRequired>1</IsDeliveryStatusRequired>
               <NotificaitionValidityPeriod>{{.NotificaitionValidityPeriod}}</NotificaitionValidityPeriod>
               <Body>{{.Message}}</Body>
               <!--0 for English text 
				8 for other language text
				4 for binary messages like WAP
				16 for flash sms
				and It will be commmon for all SMSs-->
               <BodyCoding>0</BodyCoding>
               <!--0 for English text
						1 for other language text
						2 for binary messages like WAP
						It will be commmon for all SMSs-->
               <BodyContentType>0</BodyContentType>
               <ReceiverList>
                  <!--One receiver only 8 digit Omantel MSISDN without country code-->
                  <Receiver>{{.ContactNumber}}</Receiver>
               </ReceiverList>
            </sen:Notification>
         </sen:RequestBody>
      </sen:SendNotificationRq>
   </soapenv:Body>
</soapenv:Envelope>
`
