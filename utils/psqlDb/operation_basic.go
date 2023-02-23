package psqldb

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

// FetchPlans loads the global variable PlanRates with the rate plan for each plan type
func (dbConn *DBImplementation) FetchPlans(year int) *errors.Error {

	// logger.FileLogger.Info().Time("start", time.Now()).Msg("FetchPlans start")
	// defer logger.FileLogger.Info().Time("end", time.Now()).Msg("FetchPlans done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select plan_code, plan_name, plan_price, sip_plan_rate, quarter from plans_table where deleted IS NULL and year=$1", year)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchPlans Query: failed",
			Details:    "FetchPlans Query: " + err.Error(),
		}
	}

	models.PlanRates.Mutex.Lock()
	defer models.PlanRates.Mutex.Unlock()

	if models.PlanRates.PlanRateMap == nil {
		models.PlanRates.PlanRateMap = make(map[int]map[int]map[string]*models.ProductPlans)
	}

	if models.PlanRates.PlanRateMap[year] == nil {
		models.PlanRates.PlanRateMap[year] = make(map[int]map[string]*models.ProductPlans)
	}

	var (
		id            string
		plan_name     string
		plan_price    float64
		sip_plan_rate float64
		quarter       int
	)

	for rows.Next() {
		if err := rows.Scan(&id, &plan_name, &plan_price, &sip_plan_rate, &quarter); err != nil {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchPlans Scan: failed",
				Details:    "FetchPlans Scan: " + err.Error(),
			}
		}

		if models.PlanRates.PlanRateMap[year][quarter] == nil {
			models.PlanRates.PlanRateMap[year][quarter] = make(map[string]*models.ProductPlans)
		}

		if models.PlanRates.PlanRateMap[year][quarter][id] == nil {
			models.PlanRates.PlanRateMap[year][quarter][id] = &models.ProductPlans{
				PlanName:    plan_name,
				PlanPrice:   plan_price,
				SIPRatePlan: sip_plan_rate,
			}
		}
	}

	return nil
}

// FetchSaleTarget loads the global variable SalesTrgt with the target for each outlet
func (dbConn *DBImplementation) FetchSaleTarget(year int) *errors.Error {

	// logger.FileLogger.Info().Time("start", time.Now()).Msg("FetchSaleTarget start")
	// defer logger.FileLogger.Info().Time("end", time.Now()).Msg("FetchSaleTarget done")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select region, outlet, month, baqati, hbb, hayyak, upgrade from targets_table where deleted IS NULL and year=$1", year)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchSaleTarget Query: failed",
			Details:    "FetchSaleTarget Query: " + err.Error(),
		}
	}

	models.SalesTrgt.Mutex.Lock()
	defer models.SalesTrgt.Mutex.Unlock()

	if models.SalesTrgt.SaleTargetMap == nil {
		models.SalesTrgt.SaleTargetMap = make(map[int]map[string]map[string][12]*models.ProductGrpTarget)
	}

	if models.SalesTrgt.SaleTargetMap[year] == nil {
		models.SalesTrgt.SaleTargetMap[year] = make(map[string]map[string][12]*models.ProductGrpTarget)
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

		if models.SalesTrgt.SaleTargetMap[year][region] == nil {
			models.SalesTrgt.SaleTargetMap[year][region] = make(map[string][12]*models.ProductGrpTarget)
		}

		array, ok := models.SalesTrgt.SaleTargetMap[year][region][outlet]
		if !ok {
			var temp [12]*models.ProductGrpTarget
			models.SalesTrgt.SaleTargetMap[year][region][outlet] = temp
		}

		array[month-1] = &models.ProductGrpTarget{
			Baqati:  baqati,
			HBB:     hBB,
			Hayyak:  hayyak,
			Upgrade: upgrade,
		}

		models.SalesTrgt.SaleTargetMap[year][region][outlet] = array
	}

	return nil
}

// FetchValueOfPoint loads the global variable ValOfPoint with the value of point for each product type
func (dbConn *DBImplementation) FetchValueOfPoint(year int) *errors.Error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	rows, err := dbConn.Operation.QueryContext(ctx, "select product, achieved_percent, vop_value, extra_pay, quarter from vop_table where deleted IS NULL and year=$1", year)
	if err != nil {
		return &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerError,
			Message:    "FetchValueOfPoint Query: failed",
			Details:    "FetchValueOfPoint Query: " + err.Error(),
		}
	}

	models.ValOfPoint.Mutex.Lock()
	defer models.ValOfPoint.Mutex.Unlock()

	if models.ValOfPoint.VopMap == nil {
		models.ValOfPoint.VopMap = make(map[int]map[int]map[string]map[int]*models.VOP)
	}

	if models.ValOfPoint.VopMap[year] == nil {
		models.ValOfPoint.VopMap[year] = make(map[int]map[string]map[int]*models.VOP)
	}

	for rows.Next() {

		var (
			product         string
			achievedPercent int
			vopValue        float64
			extraPay        float64
			quarter         int
		)

		if err := rows.Scan(&product, &achievedPercent, &vopValue, &extraPay, &quarter); err != nil {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchValueOfPoint Scan: failed",
				Details:    "FetchValueOfPoint Scan: " + err.Error(),
			}
		}

		if quarter < 1 || quarter > 4 {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchValueOfPoint: invalid quarter record",
				Details:    "FetchValueOfPoint: invalid quarter record - " + fmt.Sprintf("year:%v, quarter:%v", year, quarter),
			}
		}

		if !(product == constant.HBBCategory || product == constant.HayyakCategory || product == constant.BaqatiCategory || product == constant.UpgradeCategory) {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchValueOfPoint: invalid product record",
				Details:    "FetchValueOfPoint: invalid product record - " + fmt.Sprintf("year:%v, quarter:%v, product:%v", year, quarter, product),
			}
		}

		if !(achievedPercent == 75 || achievedPercent == 95 || achievedPercent == 110 || achievedPercent == 139 || achievedPercent == 140 || achievedPercent == 159 ||
			achievedPercent == 179 || achievedPercent == 189 || achievedPercent == 199 || achievedPercent == 200) {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchValueOfPoint: invalid achievedPercent record",
				Details: "FetchValueOfPoint: invalid achievedPercent record - " + fmt.Sprintf("year:%v, quarter:%v, product:%v, achievedPercent:%v",
					year, quarter, product, achievedPercent),
			}
		}

		if models.ValOfPoint.VopMap[year][quarter] == nil {
			models.ValOfPoint.VopMap[year][quarter] = make(map[string]map[int]*models.VOP)
		}

		if models.ValOfPoint.VopMap[year][quarter][product] == nil {
			models.ValOfPoint.VopMap[year][quarter][product] = make(map[int]*models.VOP)
		}

		if models.ValOfPoint.VopMap[year][quarter][product] == nil {
			models.ValOfPoint.VopMap[year][quarter][product] = make(map[int]*models.VOP)
		}

		if models.ValOfPoint.VopMap[year][quarter][product][achievedPercent] == nil {
			models.ValOfPoint.VopMap[year][quarter][product][achievedPercent] = &models.VOP{Point: vopValue, ExtraPay: extraPay}
		}
	}

	//validation
	for _, quarterMap := range models.ValOfPoint.VopMap[year] {
		if len(quarterMap) != 4 {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchValueOfPoint: all products not found",
				Details:    "FetchValueOfPoint: all products not found - " + fmt.Sprintf("%v", len(quarterMap)),
			}
		}

		for _, productMap := range quarterMap {
			if len(productMap) != 10 {
				return &errors.Error{
					StatusCode: http.StatusInternalServerError,
					Code:       errors.InternalServerError,
					Message:    "FetchValueOfPoint: all achievedPercent not found",
					Details:    "FetchValueOfPoint: all achievedPercent not found - " + fmt.Sprintf("%v", len(productMap)),
				}
			}

		}
	}

	return nil
}

// FetchMaxPayout loads the gloabal variable MaxPay with the maximum payout limit
func (dbConn *DBImplementation) FetchMaxPayout(year int) *errors.Error {

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

	models.MaxPay.Mutex.Lock()
	defer models.MaxPay.Mutex.Unlock()

	if models.MaxPay.Payoutmap == nil {
		models.MaxPay.Payoutmap = make(map[int]map[int]map[string]map[string]float64)
	}

	if models.MaxPay.Payoutmap[year] == nil {
		models.MaxPay.Payoutmap[year] = make(map[int]map[string]map[string]float64)
	}

	for quarter := 1; quarter <= 4; quarter++ {
		if models.MaxPay.Payoutmap[year][quarter] == nil {
			models.MaxPay.Payoutmap[year][quarter] = make(map[string]map[string]float64)
		}

		if models.MaxPay.Payoutmap[year][quarter][constant.PromoterAbb] == nil {
			models.MaxPay.Payoutmap[year][quarter][constant.PromoterAbb] = make(map[string]float64)
		}

		if models.MaxPay.Payoutmap[year][quarter][constant.OutletManagerAbb] == nil {
			models.MaxPay.Payoutmap[year][quarter][constant.OutletManagerAbb] = make(map[string]float64)
		}

		if models.MaxPay.Payoutmap[year][quarter][constant.ManagerAbb] == nil {
			models.MaxPay.Payoutmap[year][quarter][constant.ManagerAbb] = make(map[string]float64)
		}

		if models.MaxPay.Payoutmap[year][quarter][constant.SeniorManagerAbb] == nil {
			models.MaxPay.Payoutmap[year][quarter][constant.SeniorManagerAbb] = make(map[string]float64)
		}

		models.MaxPay.Payoutmap[year][quarter][constant.PromoterAbb]["HBB"] = 350
		models.MaxPay.Payoutmap[year][quarter][constant.PromoterAbb]["BAQATI"] = 350
		models.MaxPay.Payoutmap[year][quarter][constant.PromoterAbb]["HAYYAK"] = 220
		models.MaxPay.Payoutmap[year][quarter][constant.PromoterAbb]["UPGRADE"] = 80

		models.MaxPay.Payoutmap[year][quarter][constant.OutletManagerAbb]["HBB"] = 420
		models.MaxPay.Payoutmap[year][quarter][constant.OutletManagerAbb]["BAQATI"] = 420
		models.MaxPay.Payoutmap[year][quarter][constant.OutletManagerAbb]["HAYYAK"] = 264
		models.MaxPay.Payoutmap[year][quarter][constant.OutletManagerAbb]["UPGRADE"] = 96

		models.MaxPay.Payoutmap[year][quarter][constant.ManagerAbb]["HBB"] = 656.25
		models.MaxPay.Payoutmap[year][quarter][constant.ManagerAbb]["BAQATI"] = 656.25
		models.MaxPay.Payoutmap[year][quarter][constant.ManagerAbb]["HAYYAK"] = 412.5
		models.MaxPay.Payoutmap[year][quarter][constant.ManagerAbb]["UPGRADE"] = 150

		models.MaxPay.Payoutmap[year][quarter][constant.SeniorManagerAbb]["HBB"] = 985
		models.MaxPay.Payoutmap[year][quarter][constant.SeniorManagerAbb]["BAQATI"] = 985
		models.MaxPay.Payoutmap[year][quarter][constant.SeniorManagerAbb]["HAYYAK"] = 619
		models.MaxPay.Payoutmap[year][quarter][constant.SeniorManagerAbb]["UPGRADE"] = 225
	}

	return nil
}

// FetchMaxExtraPayout loads the gloabal variable MaxExtraPay with the maximum extra payout limit
func (dbConn *DBImplementation) FetchMaxExtraPayout(year int) *errors.Error {

	models.MaxExtraPay.Mutex.Lock()
	defer models.MaxExtraPay.Mutex.Unlock()

	if models.MaxExtraPay.ExtraPayoutmap == nil {
		models.MaxExtraPay.ExtraPayoutmap = make(map[int]map[int]map[string]map[string]float64)
	}

	if models.MaxExtraPay.ExtraPayoutmap[year] == nil {
		models.MaxExtraPay.ExtraPayoutmap[year] = make(map[int]map[string]map[string]float64)
	}

	for quarter := 1; quarter <= 4; quarter++ {
		if models.MaxExtraPay.ExtraPayoutmap[year][quarter] == nil {
			models.MaxExtraPay.ExtraPayoutmap[year][quarter] = make(map[string]map[string]float64)
		}

		if models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.PromoterAbb] == nil {
			models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.PromoterAbb] = make(map[string]float64)
		}

		if models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.OutletManagerAbb] == nil {
			models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.OutletManagerAbb] = make(map[string]float64)
		}

		if models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.ManagerAbb] == nil {
			models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.ManagerAbb] = make(map[string]float64)
		}

		if models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.SeniorManagerAbb] == nil {
			models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.SeniorManagerAbb] = make(map[string]float64)
		}

		models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.PromoterAbb]["HBB"] = 88
		models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.PromoterAbb]["BAQATI"] = 88
		models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.PromoterAbb]["HAYYAK"] = 55
		models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.PromoterAbb]["UPGRADE"] = 20
	}

	return nil
}

// FetchMaxThresholdCheck loads the gloabal variable ThresholdCheck with the threshold pay limit
func (dbConn *DBImplementation) FetchMaxThresholdCheck(year int) *errors.Error {

	// logger.FileLogger.Info().Time("start", time.Now()).Msg("FetchPlans start")
	// defer logger.FileLogger.Info().Time("end", time.Now()).Msg("FetchPlans done")

	models.ThresholdCheck.Mutex.Lock()
	defer models.ThresholdCheck.Mutex.Unlock()

	if models.ThresholdCheck.Limit == nil {
		models.ThresholdCheck.Limit = make(map[int]map[int]map[string]map[string]map[int]float64)
	}

	if models.ThresholdCheck.Limit[year] == nil {
		models.ThresholdCheck.Limit[year] = make(map[int]map[string]map[string]map[int]float64)
	}

	positions := []string{constant.PromoterCC, constant.OutletManagerCC, constant.ManagerCC, constant.SeniorManagerAbb}

	for _, position := range positions {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		rows, err := dbConn.Operation.QueryContext(ctx, "select product, achieved_percent, payout, quarter from threshold_table where deleted IS NULL and position=$1 and year=$2", position, year)
		if err != nil {
			return &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerError,
				Message:    "FetchPlans Query: failed",
				Details:    "FetchPlans Query: " + err.Error(),
			}
		}

		for rows.Next() {

			var (
				product         string
				achievedPercent int
				payout          float64
				quarter         int
			)

			if err := rows.Scan(&product, &achievedPercent, &payout, &quarter); err != nil {
				return &errors.Error{
					StatusCode: http.StatusInternalServerError,
					Code:       errors.InternalServerError,
					Message:    "FetchMaxThresholdCheck Scan: failed",
					Details:    "FetchMaxThresholdCheck Scan: " + err.Error(),
				}
			}

			pos := constant.EmployeePositionMap[position]
			if pos == "" {
				return &errors.Error{
					StatusCode: http.StatusInternalServerError,
					Code:       errors.InternalServerError,
					Message:    "blank user position",
					Details:    "FetchMaxThresholdCheck: blank position",
				}
			}

			if models.ThresholdCheck.Limit[year][quarter] == nil {
				models.ThresholdCheck.Limit[year][quarter] = make(map[string]map[string]map[int]float64)
			}

			if models.ThresholdCheck.Limit[year][quarter][pos] == nil {
				models.ThresholdCheck.Limit[year][quarter][pos] = make(map[string]map[int]float64)
			}

			if models.ThresholdCheck.Limit[year][quarter][pos][product] == nil {
				models.ThresholdCheck.Limit[year][quarter][pos][product] = make(map[int]float64)
			}

			models.ThresholdCheck.Limit[year][quarter][pos][product][achievedPercent] = payout
		}
	}

	return nil
}
