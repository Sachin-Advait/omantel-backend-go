package impl

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/fcm"
	fcmClient "gitlab.com/rendezvousPoint/incentive_calculator/utils/fcm"
)

type FCMServiceImpl struct {
	FCMClient fcmClient.FCM
}

func GetFCMServiceInstance(fcmClient fcmClient.FCM) fcm.FCMService {
	return &FCMServiceImpl{
		FCMClient: fcmClient,
	}
}
