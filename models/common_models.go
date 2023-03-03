package models

import "time"

var (
	CaptchaSitesEnableKey  string
	CaptchaSecretKey       string
	VerifyCaptchaGoogleApi string
)

type Pong struct {
	DT time.Time `json:"time"`
}

type APIResponse struct {
	Status    bool        `json:"status" example:"false" enums:"true,false"`
	Message   string      `json:"message" example:"CE10001"`
	ErrorCode string      `json:"errorcode" example:"INTERNAL SERVER ERROR"`
	Data      interface{} `json:"data"`
}

type APISuccessResponse struct {
	Status bool        `json:"status" example:"false" enums:"true,false"`
	Data   interface{} `json:"data"`
}

type APIFailureResponse struct {
	Status  bool   `json:"status" example:"false" enums:"true,false"`
	Error   string `json:"error" example:"INTERNAL SERVER ERROR"`
	Message string `json:"message" example:"failed reason"`
}

type SuccessMessage struct {
	Status string `json:"status"`
}
