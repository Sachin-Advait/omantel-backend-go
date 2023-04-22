package models

import "time"

type UserSessionToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	LoginTime    time.Time `json:"login_time"`
}
