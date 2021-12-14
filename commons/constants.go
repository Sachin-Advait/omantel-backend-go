package commons

// Config Viper
const (
	DEV         = "DEV"
	QA          = "QA"
	PROD        = "PROD"
	CONFIG_PATH = "local_bkp"
	CONFIG_DIR  = "./config"
)

// common
const (
	EMPTY       = ""
	EnvVariable = "ENVIRONMENT"
)

const (
	LOG_FILE_NAME = "SalesIncentive.Log"
)

const (
	ClientIDHeader = "client_id"
	QuarterHeader  = "quarter"
	DateHeader     = "date"
)

const (
	HBBCategory     = "HBB"
	BaqatiCategory  = "BAQATI"
	HayyakCategory  = "HAYYAK"
	UpgradeCategory = "UPGRADE"
)

const (
	VicePresident  = "VICE PRESIDENT"
	GeneralManager = "GENERAl MANAGER"
	SeniorManager  = "SENIOR MANAGER"
	Manager        = "MANAGER"
	OutletManager  = "OUTLET MANAGER"
	Promoter       = "PROMOTER"
)

const (
	BatinahNorth   = "Batinah North"
	BatinahSouth   = "Batinah South"
	DakhliaDahirah = "Dakhlia & Dahirah"
	Dhofar         = "Dhofar"
	MuscatNorth    = "Muscat North"
	MuscatSouth    = "Muscat South"
	Sharqiya       = "Sharqiya"
)

var (
	Table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

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
