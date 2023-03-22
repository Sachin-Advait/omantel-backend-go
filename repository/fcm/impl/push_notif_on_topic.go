package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"firebase.google.com/go/messaging"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (fcm *FCMServiceImpl) SendViaTopicMessage(topic string, notif *messaging.Notification, data map[string]string) (string, error) {

	message := &messaging.Message{
		Topic:        topic,
		Notification: notif,
		Data:         data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
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
