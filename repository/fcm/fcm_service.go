package fcm

import "firebase.google.com/go/messaging"

type FCMService interface {
	SendViaTopicMessage(topic string, notif *messaging.Notification, data map[string]string) (string, error)
	SubscribeToTopic(topic string, token []string) (int, error)
	UnSubscribeToTopic(topic string, token []string) (int, error)
}
