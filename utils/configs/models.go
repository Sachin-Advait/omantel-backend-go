package configs

import "time"

type IncentiveConfig struct {
	SessionExpiryTime      time.Duration
	RefreshTokenExpiryTime time.Duration
}
