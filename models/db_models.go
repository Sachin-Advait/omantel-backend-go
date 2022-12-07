package models

import "time"

type UserReportDB struct {
	UserId                       string
	ManagerId                    string
	Isquarterly                  bool
	Month                        int
	Quarter                      int
	Year                         int
	Hayyak_sales_count           int
	Hayyak_target                float64
	Hayyak_payout                float64
	Hayyak_extra_payout          float64
	Hayyak_achieved_percentage   float64
	Hbb_sales_count              int
	Hbb_target                   float64
	Hbb_payout                   float64
	Hbb_extra_payout             float64
	Hbb_achieved_percentage      float64
	Baqati_sales_count           int
	Baqati_target                float64
	Baqati_payout                float64
	Baqati_extra_payout          float64
	Baqati_achieved_percentage   float64
	Upgrade_sales_count          int
	Upgrade_target               float64
	Upgrade_payout               float64
	Upgrade_extra_payout         float64
	Uprade_achieved_percentage   float64
	Total_payout                 float64
	HBBSaleCountPerPlanArray     string
	HayyakSaleCountPerPlanArray  string
	BaqatiSaleCountPerPlanArray  string
	UpgradeSaleCountPerPlanArray string
	CreateAt                     time.Time
}

type ManagerReportDB struct {
	UserId             string
	Position           string
	ManagerId          string
	Region             string
	Outlet             string
	Latitude           float64
	Longitude          float64
	HBBSalesCount      int
	HBBTargetCount     int
	HayyakSalesCount   int
	HayyakTargetCount  int
	BaqatiSalesCount   int
	BaqatiTargetCount  int
	UpgradeSalesCount  int
	UpgradeTargetCount int
	CreatedAt          time.Time
}

type OutletReportDB struct {
	Outlet             string
	Region             string
	Latitude           float64
	Longitude          float64
	HBBSalesCount      int
	HBBTargetCount     int
	HayyakSalesCount   int
	HayyakTargetCount  int
	BaqatiSalesCount   int
	BaqatiTargetCount  int
	UpgradeSalesCount  int
	UpgradeTargetCount int
	CreatedAt          time.Time
}

type RegionReportDB struct {
	Region             string
	HBBSalesCount      int
	HBBTargetCount     int
	HayyakSalesCount   int
	HayyakTargetCount  int
	BaqatiSalesCount   int
	BaqatiTargetCount  int
	UpgradeSalesCount  int
	UpgradeTargetCount int
	CreatedAt          time.Time
}

type UserBasicDetailsDB struct {
	UserId         string
	Name           string
	Position       string
	MobileNumber   string
	Region         string
	Outlet         string
	OutletManageId string
}

type ManagerBasicDetailsDB struct {
	UserId          string
	Name            string
	Position        string
	MobileNumber    string
	Region          string
	Outlet          string
	ManagerId       string
	ManagerPosition string
}
