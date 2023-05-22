package configs

import "time"

type AdminConfig struct {
	FCMOfferTopic        string
	FCMTrainingTopic     string
	FCMNotificationTopic string
}

type UserConfig struct {
	OTPUrl                 string
	FCMOfferTopic          string
	FCMTrainingTopic       string
	FCMNotificationTopic   string
	ReportCacheTTL         time.Duration
	SessionExpiryTime      time.Duration
	RefreshTokenExpiryTime time.Duration
}

type IncentiveConfig struct {
	OTPUrl                 string
	SessionExpiryTime      time.Duration
	RefreshTokenExpiryTime time.Duration
	FCMOfferTopic          string
	FCMTrainingTopic       string
	FCMNotificationTopic   string
}
