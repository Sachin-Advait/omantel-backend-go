package fcm

type FCMService interface {
	SendViaTopicMessage(topic string, data map[string]string) (string, error)
}
