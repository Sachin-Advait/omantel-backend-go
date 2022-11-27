const moment = require('moment-timezone');
const soapRequest = require('easy-soap-request');
const https = require('https');
const base64 = require('base-64');
const SentMessage = require('../models/SentMessage');
const { queueLogger } = require('../config/logger');
//const smsQueue = require('../queues/sendSmsQueue');

function getXml(message, phoneNumber) {
    let xml = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:sen="http://www.omantel.om/esb/messages/sendnotification">
    <soapenv:Header/>
    <soapenv:Body>
       <sen:SendNotificationRq>
          <sen:RequestHeader>
             <FunctionId>17002004</FunctionId>
             <TransactionId>0</TransactionId>
             <RequestId>ZONE_${(new Date()).getTime()}</RequestId>
             <RequestorChannelId>ZONE</RequestorChannelId>
             <RequestorUserId>ZONE</RequestorUserId>
             <Timestamp> ${moment().add(2, 'years').toISOString()}</Timestamp>
          </sen:RequestHeader>
          <sen:RequestBody>
             <sen:Notification>
                <NotificationType>SMS</NotificationType>
                <NotificationDateTime>${moment().tz('Asia/Muscat').format("YYYY-MM-DD HH:mm:ss")}</NotificationDateTime>
                <NotificationLanguage>1</NotificationLanguage>
                <ShortCode>Hassin</ShortCode>
                <IsDeliveryStatusRequired>1</IsDeliveryStatusRequired>
                <NotificaitionValidityPeriod>${moment().tz('Asia/Muscat').toISOString()}</NotificaitionValidityPeriod>
                <Body>${message}</Body>
                <BodyCoding>8</BodyCoding>
                <BodyContentType>1</BodyContentType>
                <ReceiverList>
                   <Receiver>${phoneNumber}</Receiver>
                </ReceiverList>
             </sen:Notification>
          </sen:RequestBody>
       </sen:SendNotificationRq>
    </soapenv:Body>
 </soapenv:Envelope>`;
    return xml;
}

module.exports.sendSms = async function (phoneNumber, message, queue = true) {
    const url = 'https://pki.omantel.om:17862/services/SendNotification';

    if (phoneNumber.length === 11) {
        phoneNumber = phoneNumber.substring(3);
    }

    let xml = getXml(message, phoneNumber);

    const agent = new https.Agent({
        rejectUnauthorized: false
    });

    let sentMessageData = new SentMessage({
        phoneNumber: phoneNumber,
        message: message,
        data: xml
    });

    sentMessageData = await sentMessageData.save();

    console.log(" ******** xml data ********* ", xml);

    let credentials = "SA_ZoneAlert:Omantel@123";
    let encodedCredentials = base64.encode(credentials);

    if (process.env.NODE_ENV === 'development') {
	console.log("in developement")
        return;
    }

    try {
	console.log("in production ****")
        const { response } = await soapRequest({
            url: url, headers: {
                'user-agent': 'sampleTest',
                'Content-Type': 'text/xml;charset=UTF-8',
                Authorization: `Basic ${encodedCredentials}`,
                SOAPAction: 'SendNotification',
            }, xml: xml, extraOpts: {
                httpsAgent: agent
            }
        });
        const { body, statusCode } = response;
        console.log(" ******* statusCode ***** ", statusCode);
        sentMessageData.statusCode = statusCode;
        sentMessageData.response = body;
        await sentMessageData.save();

    } catch (err) {
	console.log("error",err)
        if (queue) {
            // this.pushMessageToQueue(phoneNumber, message);
        } else {
            queueLogger.info(" ***** error while seding sms ****** ", { err });
            throw new Error("Can not send sms");
        }
    }

}

// module.exports.pushMessageToQueue = function pushMessageToQueue(phoneNumber, message) {
//     smsQueue.add({
//         phoneNumber: phoneNumber,
//         message: message
//     }, {
//         attempts: 20
//     });
// }
