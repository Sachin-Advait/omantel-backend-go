package psqldb

import (
	"context"
	"net/http"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

// FetchPlans loads the global variable PlanRates with the rate plan for each plan type
func (dbConn *DBImplementation) FetchPlans(quarter, year int) *errors.Error {

	// logger.FileLogger.Info().Time("start", time.Now()).Msg("FetchPlans start")
	// defer logger.FileLogger.Info().Time("end", time.Now()).Msg("FetchPlans done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select plan_code, plan_name, plan_price, sip_plan_rate from plans_table where deleted IS NULL and quarter=$1 and year=$2", quarter, year)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchPlans Query: failed",
			Details:    "FetchPlans Query: " + err.Error(),
		}
	}

	models.PlanRates.Year = year

	if models.PlanRates.PlanRateMap == nil {
		models.PlanRates.PlanRateMap = make(map[int]map[string]*models.ProductPlans)
	}

	if models.PlanRates.PlanRateMap[quarter] == nil {
		models.PlanRates.PlanRateMap[quarter] = make(map[string]*models.ProductPlans)
	}

	var (
		id            string
		plan_name     string
		plan_price    float64
		sip_plan_rate float64
	)

	for rows.Next() {
		if err := rows.Scan(&id, &plan_name, &plan_price, &sip_plan_rate); err != nil {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchPlans Scan: failed",
				Details:    "FetchPlans Scan: " + err.Error(),
			}
		}

		models.PlanRates.Mutex.Lock()

		models.PlanRates.PlanRateMap[quarter][id] = &models.ProductPlans{
			PlanName:    plan_name,
			PlanPrice:   plan_price,
			SIPRatePlan: sip_plan_rate,
		}

		models.PlanRates.Mutex.Unlock()

	}

	return nil
}

// FetchSaleTarget loads the global variable SalesTrgt with the target for each outlet
func (dbConn *DBImplementation) FetchSaleTarget(quarter, year int) *errors.Error {

	// logger.FileLogger.Info().Time("start", time.Now()).Msg("FetchSaleTarget start")
	// defer logger.FileLogger.Info().Time("end", time.Now()).Msg("FetchSaleTarget done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	models.SalesTrgt.Year = year

	if models.SalesTrgt.SaleTargetMap == nil {
		models.SalesTrgt.SaleTargetMap = make(map[string]map[string][12]*models.ProductGrpTarget)
	}

	rows, err := dbConn.Operation.QueryContext(ctx, "select region, outlet, month, baqati, hbb, hayyak, upgrade from targets_table where deleted IS NULL and quarter=$1 and year=$2", quarter, year)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchSaleTarget Query: failed",
			Details:    "FetchSaleTarget Query: " + err.Error(),
		}
	}

	var (
		region  string
		outlet  string
		month   int
		baqati  int
		hBB     int
		hayyak  int
		upgrade int
	)

	for rows.Next() {
		if err := rows.Scan(&region, &outlet, &month, &baqati, &hBB, &hayyak, &upgrade); err != nil {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchSaleTarget Scan: failed",
				Details:    "FetchSaleTarget Scan: " + err.Error(),
			}
		}

		models.SalesTrgt.Mutex.Lock()

		if models.SalesTrgt.SaleTargetMap[region] == nil {
			models.SalesTrgt.SaleTargetMap[region] = make(map[string][12]*models.ProductGrpTarget)
		}

		array, ok := models.SalesTrgt.SaleTargetMap[region][outlet]
		if !ok {
			var temp [12]*models.ProductGrpTarget
			models.SalesTrgt.SaleTargetMap[region][outlet] = temp
		}

		array[month-1] = &models.ProductGrpTarget{
			Baqati:  baqati,
			HBB:     hBB,
			Hayyak:  hayyak,
			Upgrade: upgrade,
		}

		models.SalesTrgt.SaleTargetMap[region][outlet] = array

		models.SalesTrgt.Mutex.Unlock()

	}

	return nil
}

// FetchValueOfPoint loads the global variable ValOfPoint with the value of point for each product type
func (dbConn *DBImplementation) FetchValueOfPoint(quarter, year int) *errors.Error {

	// rows, err := dbConn.Operation.Query("select group, top_percentage, vop, extra_pay from Product_Plans_Table where deleted IS NULL and quarter=$1 and year=$2", quarter, year)
	// if err != nil {
	// 	return err
	// }

	// var (
	// 	group          string
	// 	top_percentage int
	// 	vop            float64
	// 	extra_pay      float64
	// )

	// for rows.Next() {
	// 	if err := rows.Scan(&group, &top_percentage, &vop, &extra_pay); err != nil {
	// 		return err
	// 	}

	// 	ValOfPoint.Mutex.Lock()
	// 	if ValOfPoint.VopMap == nil {
	// 		ValOfPoint.VopMap = make(map[string]map[int]*VOP)
	// 	}

	// 	if ValOfPoint.VopMap[group] == nil {
	// 		ValOfPoint.VopMap[group] = make(map[int]*VOP)
	// 	}

	// 	ValOfPoint.VopMap[group][top_percentage] = &VOP{
	// 		Point:    vop,
	// 		ExtraPay: extra_pay,
	// 	}

	// 	ValOfPoint.Mutex.Unlock()

	// }

	models.ValOfPoint.Year = year

	if models.ValOfPoint.VopMap == nil {
		models.ValOfPoint.VopMap = make(map[int]map[string]map[int]*models.VOP)
	}

	if models.ValOfPoint.VopMap[quarter] == nil {
		models.ValOfPoint.VopMap[quarter] = make(map[string]map[int]*models.VOP)
	}

	models.ValOfPoint.VopMap[quarter]["HBB"] = make(map[int]*models.VOP)
	models.ValOfPoint.VopMap[quarter]["BAQATI"] = make(map[int]*models.VOP)
	models.ValOfPoint.VopMap[quarter]["HAYYAK"] = make(map[int]*models.VOP)
	models.ValOfPoint.VopMap[quarter]["UPGRADE"] = make(map[int]*models.VOP)

	models.ValOfPoint.Mutex.Lock()

	models.ValOfPoint.VopMap[quarter]["HBB"][75] = &models.VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["HBB"][95] = &models.VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["HBB"][110] = &models.VOP{
		Point:    0.25,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["HBB"][139] = &models.VOP{
		Point:    0.45,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["HBB"][140] = &models.VOP{
		Point:    0.75,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["HBB"][159] = &models.VOP{
		Point:    0,
		ExtraPay: 17.5,
	}
	models.ValOfPoint.VopMap[quarter]["HBB"][179] = &models.VOP{
		Point:    0,
		ExtraPay: 35.0,
	}
	models.ValOfPoint.VopMap[quarter]["HBB"][189] = &models.VOP{
		Point:    0,
		ExtraPay: 52.5,
	}
	models.ValOfPoint.VopMap[quarter]["HBB"][199] = &models.VOP{
		Point:    0,
		ExtraPay: 70.0,
	}
	models.ValOfPoint.VopMap[quarter]["HBB"][200] = &models.VOP{
		Point:    0,
		ExtraPay: 87.5,
	}

	models.ValOfPoint.VopMap[quarter]["BAQATI"][75] = &models.VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["BAQATI"][95] = &models.VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["BAQATI"][110] = &models.VOP{
		Point:    0.25,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["BAQATI"][139] = &models.VOP{
		Point:    0.45,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["BAQATI"][140] = &models.VOP{
		Point:    0.75,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["BAQATI"][159] = &models.VOP{
		Point:    0,
		ExtraPay: 17.5,
	}
	models.ValOfPoint.VopMap[quarter]["BAQATI"][179] = &models.VOP{
		Point:    0,
		ExtraPay: 35.0,
	}
	models.ValOfPoint.VopMap[quarter]["BAQATI"][189] = &models.VOP{
		Point:    0,
		ExtraPay: 53,
	}
	models.ValOfPoint.VopMap[quarter]["BAQATI"][199] = &models.VOP{
		Point:    0,
		ExtraPay: 70.0,
	}
	models.ValOfPoint.VopMap[quarter]["BAQATI"][200] = &models.VOP{
		Point:    0,
		ExtraPay: 87.5,
	}

	models.ValOfPoint.VopMap[quarter]["HAYYAK"][75] = &models.VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["HAYYAK"][95] = &models.VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["HAYYAK"][110] = &models.VOP{
		Point:    0.25,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["HAYYAK"][139] = &models.VOP{
		Point:    0.45,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["HAYYAK"][140] = &models.VOP{
		Point:    0.75,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["HAYYAK"][159] = &models.VOP{
		Point:    0,
		ExtraPay: 11,
	}
	models.ValOfPoint.VopMap[quarter]["HAYYAK"][179] = &models.VOP{
		Point:    0,
		ExtraPay: 22,
	}
	models.ValOfPoint.VopMap[quarter]["HAYYAK"][189] = &models.VOP{
		Point:    0,
		ExtraPay: 33,
	}
	models.ValOfPoint.VopMap[quarter]["HAYYAK"][199] = &models.VOP{
		Point:    0,
		ExtraPay: 44,
	}
	models.ValOfPoint.VopMap[quarter]["HAYYAK"][200] = &models.VOP{
		Point:    0,
		ExtraPay: 55,
	}

	models.ValOfPoint.VopMap[quarter]["UPGRADE"][75] = &models.VOP{
		Point:    0.0,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["UPGRADE"][95] = &models.VOP{
		Point:    0.1,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["UPGRADE"][110] = &models.VOP{
		Point:    0.15,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["UPGRADE"][139] = &models.VOP{
		Point:    0.3,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["UPGRADE"][140] = &models.VOP{
		Point:    0.5,
		ExtraPay: 0,
	}
	models.ValOfPoint.VopMap[quarter]["UPGRADE"][159] = &models.VOP{
		Point:    0,
		ExtraPay: 4,
	}
	models.ValOfPoint.VopMap[quarter]["UPGRADE"][179] = &models.VOP{
		Point:    0,
		ExtraPay: 8,
	}
	models.ValOfPoint.VopMap[quarter]["UPGRADE"][189] = &models.VOP{
		Point:    0,
		ExtraPay: 12,
	}
	models.ValOfPoint.VopMap[quarter]["UPGRADE"][199] = &models.VOP{
		Point:    0,
		ExtraPay: 16,
	}
	models.ValOfPoint.VopMap[quarter]["UPGRADE"][200] = &models.VOP{
		Point:    0,
		ExtraPay: 20,
	}

	models.ValOfPoint.Mutex.Unlock()

	return nil
}

// FetchMaxPayout loads the gloabal variable MaxPay with the maximum payout limit
func (dbConn *DBImplementation) FetchMaxPayout(quarter, year int) *errors.Error {

	// rows, err := dbConn.Operation.Query("select group, percentage, max_pay from Max_Payout_Table where deleted IS NULL and quarter=$1 and year=$2", quarter, year)
	// if err != nil {
	// 	return err
	// }

	// var (
	// 	group      string
	// 	percentage int
	// 	max_pay    int
	// )

	// for rows.Next() {
	// 	if err := rows.Scan(&group, &percentage, &max_pay); err != nil {
	// 		return err
	// 	}

	// 	MaxPay.Mutex.Lock()
	// 	if MaxPay.Payoutmap == nil {
	// 		MaxPay.Payoutmap = make(map[string]float64)
	// 	}

	// 	MaxPay.Payoutmap[group] = float64(percentage) * float64(max_pay) / 100

	// 	MaxPay.Mutex.Unlock()

	// }

	models.MaxPay.Year = year

	if models.MaxPay.Payoutmap == nil {
		models.MaxPay.Payoutmap = make(map[int]map[string]float64)
	}

	if models.MaxPay.Payoutmap[quarter] == nil {
		models.MaxPay.Payoutmap[quarter] = make(map[string]float64)
	}

	models.MaxPay.Mutex.Lock()

	models.MaxPay.Payoutmap[quarter]["HBB"] = 350
	models.MaxPay.Payoutmap[quarter]["BAQATI"] = 350
	models.MaxPay.Payoutmap[quarter]["HAYYAK"] = 220
	models.MaxPay.Payoutmap[quarter]["UPGRADE"] = 80
	models.MaxPay.Mutex.Unlock()

	return nil
}
