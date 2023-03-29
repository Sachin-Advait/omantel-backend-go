package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"firebase.google.com/go/messaging"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (fcm *FCMServiceImpl) SendViaTopicMessage(topic string, notif *messaging.Notification, data map[string]string) (string, error) {

	apnsData := make(map[string]interface{})
	for k, v := range data {
		apnsData[k] = v
	}

	message := &messaging.Message{
		Topic:        topic,
		Notification: notif,
		Data:         data,
		Android: &messaging.AndroidConfig{
			Data: data,
			Notification: &messaging.AndroidNotification{
				Title:      notif.Title,
				Body:       notif.Body,
				ImageURL:   notif.ImageURL,
				Color:      "#FFA500",
				Sticky:     true,
				Priority:   messaging.PriorityHigh,
				Visibility: messaging.VisibilityPublic,
			},
		},
		// APNS: &messaging.APNSConfig{
		// 	//https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CommunicatingwithAPNs.html
		// 	//https://docs.expo.dev/push-notifications/sending-notifications-custom/
		// 	Headers: map[string]string{
		// 		":method":         http.MethodPost,
		// 		"scheme":          "https",
		// 		":path":           "/3/device/<device-token>",
		// 		"host":            "",
		// 		"authorization":   "",
		// 		"apns-id":         "",
		// 		"apns-expiration": "0",
		// 		"apns-priority":   "5",
		// 		"apns-topic":      "YOUR-BUNDLE-IDENTIFIER",
		// 	},
		// 	Payload: &messaging.APNSPayload{
		// 		Aps: &messaging.Aps{
		// 			Alert: &messaging.ApsAlert{
		// 				Title:       notif.Title,
		// 				SubTitle:    "",
		// 				Body:        notif.Body,
		// 				LaunchImage: notif.ImageURL,
		// 			},
		// 			Category:   "",
		// 			CustomData: apnsData,
		// 		},
		// 	},
		// },
	}

	raw, _ := json.Marshal(message)
	fmt.Println("MESSAGE: " + string(raw))

	response, err := fcm.FCMClient.Send(context.TODO(), message)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (fcm *FCMServiceImpl) SubscribeToTopic(topic string, token []string) (int, error) {

	response, err := fcm.FCMClient.SubscribeToTopic(context.TODO(), token, topic)
	if err != nil {
		return 0, err
	}

	raw, _ := objecthelper.MarshalJSON(response.Errors)
	fmt.Println("subscribe errors: " + string(raw))

	if response.FailureCount > 0 {
		return 0, fmt.Errorf("partial failed count - %v", response.FailureCount)
	}

	return response.SuccessCount, nil
}

func (fcm *FCMServiceImpl) UnSubscribeToTopic(topic string, token []string) (int, error) {

	response, err := fcm.FCMClient.UnsubscribeFromTopic(context.TODO(), token, topic)
	if err != nil {
		return 0, err
	}

	raw, _ := objecthelper.MarshalJSON(response.Errors)
	fmt.Println("unsubscribe errors: " + string(raw))

	if response.FailureCount > 0 {
		return 0, fmt.Errorf("partial failed count - %v", response.FailureCount)
	}

	return response.SuccessCount, nil
}
