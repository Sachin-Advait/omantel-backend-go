package models

type SingleUserIncentiveRequest struct {
	ClientId  string `json:"client_id"`
	Monthly   string `json:"monthly"`
	Quarterly string `json:"quarterly"`
}

type SalesCountRequest struct {
	ClientId  string `json:"client_id"`
	Monthly   string `json:"monthly"`
	Quarterly string `json:"quarterly"`
	Product   string `json:"product"`
}

type GenerateOTP struct {
	ClientId string `json:"client_id"`
}

type VerifyOTP struct {
	ClientId string `json:"client_id"`
	Otp      int    `json:"otp"`
}

type VerifyOTPResponse struct {
	Authenticated bool   `json:"authenticated"`
	StaffName     string `json:"staff_name"`
	StaffPosition string `json:"staff_position"`
	Error         string `json:"error"`
}

type UserReportRequest struct {
	ClientId string `json:"client_id"`
	Month    int    `json:"month"`
	Quarter  int    `json:"quarter"`
}

type UserReporting struct {
	ClientId string `json:"client_id"`
}
