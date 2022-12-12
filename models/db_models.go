package models

import "time"

type UserReportDB struct {
	UserId                       string
	Position                     string
	ManagerId                    string
	ManagerPosition              string
	Region                       string
	Outlet                       string
	Latitude                     float64
	Longitude                    float64
	HayyakSalesCount             int
	HayyakTargetCount            float64
	HayyakPayout                 float64
	HayyakExtraPayout            float64
	HayyakAchievedPercentage     float64
	HBBSalesCount                int
	HBBTargetCount               float64
	HBBPayout                    float64
	HBBExtraPayout               float64
	HBBAchievedPercentage        float64
	BaqatiSalesCount             int
	BaqatiTargetCount            float64
	BaqatiPayout                 float64
	BaqatiExtraPayout            float64
	BaqatiAchievedPercentage     float64
	UpgradeSalesCount            int
	UpgradeTargetCount           float64
	UpgradePayout                float64
	UpgradeExtraPayout           float64
	UpradeAchievedPercentage     float64
	TotalPayout                  float64
	HBBSaleCountPerPlanArray     string
	HayyakSaleCountPerPlanArray  string
	BaqatiSaleCountPerPlanArray  string
	UpgradeSaleCountPerPlanArray string
	CreatedAt                    time.Time
	ModifiedAt                   time.Time
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
	ModifiedAt         time.Time
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
	ModifiedAt         time.Time
}

type UserBasicDetailsDB struct {
	UserId          string
	Name            string
	Position        string
	MobileNumber    string
	Region          string
	Outlet          string
	ManagerId       string
	ManagerPosition string
	IsAdmin         bool
}
