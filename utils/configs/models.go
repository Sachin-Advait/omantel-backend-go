package configs

import "time"

type IncentiveConfig struct {
	OTPUrl                 string
	SessionExpiryTime      time.Duration
	RefreshTokenExpiryTime time.Duration
}
