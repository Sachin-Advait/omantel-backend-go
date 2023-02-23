package models

type CreateRegionRequest struct {
	RegionCode string `json:"region_code"`
	RegionName string `json:"region_name"`
	Quarter    int    `json:"quarter"`
	Year       int    `json:"year"`
	IsUpdate   bool   `json:"isUpdate"`
}

type CreateOutletRequest struct {
	OutletCode    string  `json:"outlet_code"`
	OutletName    string  `json:"outlet_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	SubRegionName string  `json:"sub_region_name"`
	Lat           float64 `json:"lat"`
	Lang          float64 `json:"lang"`
	Quarter       int     `json:"quarter"`
	Year          int     `json:"year"`
	IsUpdate      bool    `json:"isUpdate"`
}

type CreateProductRequest struct {
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	Weightage   float64 `json:"weightage"`
	Quarter     int     `json:"quarter"`
	Year        int     `json:"year"`
	IsUpdate    bool    `json:"isUpdate"`
}

type CreatePlanRequest struct {
	PlanCode    string  `json:"plan_code"`
	PlanName    string  `json:"plan_name"`
	ProductCode string  `json:"product_code"`
	PlanPrice   float64 `json:"plan_price"`
	SIPPlanRate float64 `json:"sip_plan_rate"`
	Quarter     int     `json:"quarter"`
	Year        int     `json:"year"`
	IsUpdate    bool    `json:"isUpdate"`
}

type CreateTargetRequest struct {
	TargetCode string  `json:"target_code"`
	RegionName string  `json:"region_name"`
	OutletName string  `json:"outlet_name"`
	Month      int     `json:"month"`
	Baqati     float64 `json:"baqati"`
	Hayyak     float64 `json:"hayyak"`
	HBB        float64 `json:"hbb"`
	Upgrade    float64 `json:"upgrade"`
	Quarter    int     `json:"quarter"`
	Year       int     `json:"year"`
	IsUpdate   bool    `json:"isUpdate"`
}

type FetchDeleteRegionRequest struct {
	RegionCodes []string `json:"region_codes"`
	Quarter     int      `json:"quarter"`
	Year        int      `json:"year"`
}

type FetchDeleteOutletRequest struct {
	OutletCodes []string `json:"outlet_codes"`
	Quarter     int      `json:"quarter"`
	Year        int      `json:"year"`
}

type FetchDeleteProductRequest struct {
	ProductCodes []string `json:"product_codes"`
	Quarter      int      `json:"quarter"`
	Year         int      `json:"year"`
}

type FetchDeletePlanRequest struct {
	PlanCodes []string `json:"plan_codes"`
	Quarter   int      `json:"quarter"`
	Year      int      `json:"year"`
}

type FetchDeleteTargetRequest struct {
	TargetCodes []string `json:"target_codes"`
	RegionCode  string   `json:"region_name"`
	OutletCode  []string `json:"outlet_name"`
	Month       int      `json:"month"`
	Quarter     int      `json:"quarter"`
	Year        int      `json:"year"`
}

type RegionResponse struct {
	RegionCode string `json:"region_code"`
	RegionName string `json:"region_name"`
}

type OutletResponse struct {
	OutletCode string  `json:"outlet_code"`
	OutletName string  `json:"outlet_name"`
	RegionCode string  `json:"region_code"`
	RegionName string  `json:"region_name"`
	SubRegion  string  `json:"sub_region_name"`
	Lat        float64 `json:"lat"`
	Lang       float64 `json:"lang"`
}

type ProductResponse struct {
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	Weightage   float64 `json:"weightage"`
}

type PlanResponse struct {
	PlanCode    string  `json:"plan_code"`
	PlanName    string  `json:"plan_name"`
	ProductCode string  `json:"product_code"`
	PlanPrice   float64 `json:"plan_price"`
	SIPPlanRate float64 `json:"sip_plan_rate"`
}

type TargetResponse struct {
	TargetCode string  `json:"target_code"`
	RegionCode string  `json:"region_code"`
	Month      int     `json:"month"`
	Year       int     `json:"year"`
	Baqati     float64 `json:"baqati"`
	Hayyak     float64 `json:"hayyak"`
	HBB        float64 `json:"hbb"`
	Upgrade    float64 `json:"upgrade"`
}

type RegionFailedResponse struct {
	RegionCode string `json:"region_code"`
	Message    string `json:"message"`
}

type OutletFailedResponse struct {
	OutletCode string `json:"outlet_code"`
	Message    string `json:"message"`
}

type ProductFailedResponse struct {
	ProductCode string `json:"product_code"`
	Message     string `json:"message"`
}

type PlanFailedResponse struct {
	PlanCode string `json:"plan_code"`
	Message  string `json:"message"`
}

type TargetFailedResponse struct {
	TargetCode string `json:"target_code"`
	Message    string `json:"message"`
}
