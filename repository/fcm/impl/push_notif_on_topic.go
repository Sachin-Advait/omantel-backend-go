package impl

import (
	"context"

	"firebase.google.com/go/messaging"
)

func (fcm *FCMServiceImpl) SendViaTopicMessage(topic string, data map[string]string) (string, error) {

	message := &messaging.Message{
		Data:  data,
		Topic: topic,
	}

	response, err := fcm.FCMClient.Send(context.TODO(), message)
	if err != nil {
		return "", err
	}

	return response, nil
}
