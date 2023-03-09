package impl

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
)

// FetchPlans loads the global variable PlanRates with the rate plan for each plan type
func (dbConn *DbServiceImpl) FetchPlans(year int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	rows, err := dbConn.DBClient.QueryContext(ctx, "select plan_code, plan_name, plan_price, sip_plan_rate, quarter from plans_table where deleted IS NULL and year=$1", year)
	if err != nil {
		return fmt.Errorf("FetchPlans: query - " + err.Error())
	}

	models.PlanRates.Mutex.Lock()
	defer models.PlanRates.Mutex.Unlock()

	if models.PlanRates.PlanRateMap == nil {
		models.PlanRates.PlanRateMap = make(map[int]map[int]map[string]*models.ProductPlans)
	}

	if models.PlanRates.PlanRateMap[year] == nil {
		models.PlanRates.PlanRateMap[year] = make(map[int]map[string]*models.ProductPlans)
	}

	for rows.Next() {

		var (
			planCode    string
			planName    string
			planPrice   float64
			sipPlanRate float64
			quarter     int
		)

		if err := rows.Scan(&planCode, &planName, &planPrice, &sipPlanRate, &quarter); err != nil {
			return fmt.Errorf("FetchPlans: scan - " + err.Error())
		}

		if models.PlanRates.PlanRateMap[year][quarter] == nil {
			models.PlanRates.PlanRateMap[year][quarter] = make(map[string]*models.ProductPlans)
		}

		if models.PlanRates.PlanRateMap[year][quarter][planCode] == nil {
			models.PlanRates.PlanRateMap[year][quarter][planCode] = &models.ProductPlans{
				PlanName:    planName,
				PlanPrice:   planPrice,
				SIPRatePlan: sipPlanRate,
			}
		}
	}

	return nil
}

// FetchSaleTarget loads the global variable SalesTrgt with the target for each outlet
func (dbConn *DbServiceImpl) FetchSaleTarget(year int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	rows, err := dbConn.DBClient.QueryContext(ctx, "select region, outlet, month, baqati, hbb, hayyak, upgrade from targets_table where deleted IS NULL and year=$1", year)
	if err != nil {
		return fmt.Errorf("FetchSaleTarget: query - " + err.Error())
	}

	models.SalesTrgt.Mutex.Lock()
	defer models.SalesTrgt.Mutex.Unlock()

	if models.SalesTrgt.SaleTargetMap == nil {
		models.SalesTrgt.SaleTargetMap = make(map[int]map[string]map[string][12]*models.ProductGrpTarget)
	}

	if models.SalesTrgt.SaleTargetMap[year] == nil {
		models.SalesTrgt.SaleTargetMap[year] = make(map[string]map[string][12]*models.ProductGrpTarget)
	}

	for rows.Next() {

		var (
			region  string
			outlet  string
			month   int
			baqati  int
			hBB     int
			hayyak  int
			upgrade int
		)

		if err := rows.Scan(&region, &outlet, &month, &baqati, &hBB, &hayyak, &upgrade); err != nil {
			return fmt.Errorf("FetchSaleTarget: scan - " + err.Error())
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
func (dbConn *DbServiceImpl) FetchValueOfPoint(year int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	rows, err := dbConn.DBClient.QueryContext(ctx, "select product, achieved_percent, vop_value, extra_pay, quarter from vop_table where deleted IS NULL and year=$1", year)
	if err != nil {
		return fmt.Errorf("FetchValueOfPoint: query - " + err.Error())
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
			return fmt.Errorf("FetchValueOfPoint: scan - " + err.Error())
		}

		if quarter < 1 || quarter > 4 {
			return fmt.Errorf("FetchValueOfPoint: invalid quarter - " + fmt.Sprintf("year:%v, quarter:%v", year, quarter))
		}

		if !(product == constant.HBBCategory || product == constant.HayyakCategory || product == constant.BaqatiCategory || product == constant.UpgradeCategory) {
			return fmt.Errorf("FetchValueOfPoint: invalid product - " + fmt.Sprintf("year:%v, quarter:%v, product:%v", year, quarter, product))
		}

		if !(achievedPercent == 75 || achievedPercent == 95 || achievedPercent == 110 || achievedPercent == 139 || achievedPercent == 140 || achievedPercent == 159 ||
			achievedPercent == 179 || achievedPercent == 189 || achievedPercent == 199 || achievedPercent == 200) {
			return fmt.Errorf("FetchValueOfPoint: invalid achievedPercent - " + fmt.Sprintf("year:%v, quarter:%v, product:%v, achievedPercent:%v",
				year, quarter, product, achievedPercent))
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
	for quarter, quarterMap := range models.ValOfPoint.VopMap[year] {
		if len(quarterMap) != 4 {
			return fmt.Errorf("FetchValueOfPoint: invalid product length - " + fmt.Sprintf("quarter:%v, productLen:%v", quarter, len(quarterMap)))
		}

		for product, productMap := range quarterMap {
			if len(productMap) != 10 {
				return fmt.Errorf("FetchValueOfPoint: invalid achievedPercent length - " + fmt.Sprintf("quarter:%v, product:%v, achievedPercentLen:%v",
					quarter, product, len(productMap)))
			}
		}
	}

	return nil
}

// FetchMaxPayout loads the gloabal variable MaxPay with the maximum payout limit of A,OM,M,SM
func (dbConn *DbServiceImpl) FetchMaxPayout(year int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	rows, err := dbConn.DBClient.QueryContext(ctx, "select position, product, max_pay, quarter from max_payout_table where deleted IS NULL and year=$1", year)
	if err != nil {
		return fmt.Errorf("FetchMaxPayout: query - " + err.Error())
	}

	models.MaxPay.Mutex.Lock()
	defer models.MaxPay.Mutex.Unlock()

	if models.MaxPay.Payoutmap == nil {
		models.MaxPay.Payoutmap = make(map[int]map[int]map[string]map[string]float64)
	}

	if models.MaxPay.Payoutmap[year] == nil {
		models.MaxPay.Payoutmap[year] = make(map[int]map[string]map[string]float64)
	}

	for rows.Next() {

		var (
			position string
			product  string
			pay      float64
			quarter  int
		)

		if err := rows.Scan(&position, &product, &pay, &quarter); err != nil {
			return fmt.Errorf("FetchMaxPayout: scan - " + err.Error())
		}

		if quarter < 1 || quarter > 4 {
			return fmt.Errorf("FetchMaxPayout: invalid quarter - " + fmt.Sprintf("year:%v, quarter:%v", year, quarter))
		}

		if !(position == constant.PromoterCC || position == constant.OutletManagerCC || position == constant.ManagerCC || position == constant.SeniorManagerCC) {
			return fmt.Errorf("FetchMaxPayout: invalid position - " + fmt.Sprintf("year:%v, quarter:%v, position:%v, product:%v", year, quarter, position, product))
		}

		if !(product == constant.HBBCategory || product == constant.HayyakCategory || product == constant.BaqatiCategory || product == constant.UpgradeCategory) {
			return fmt.Errorf("FetchMaxPayout: invalid product - " + fmt.Sprintf("year:%v, quarter:%v, position:%v, product:%v", year, quarter, position, product))
		}

		if pay == 0 {
			return fmt.Errorf("FetchMaxPayout: invalid pay - " + fmt.Sprintf("year:%v, quarter:%v, position:%v, product:%v, pay:%v",
				year, quarter, position, product, pay))
		}

		if models.MaxPay.Payoutmap[year][quarter] == nil {
			models.MaxPay.Payoutmap[year][quarter] = make(map[string]map[string]float64)
		}

		if models.MaxPay.Payoutmap[year][quarter][constant.EmployeePositionMap[position]] == nil {
			models.MaxPay.Payoutmap[year][quarter][constant.EmployeePositionMap[position]] = make(map[string]float64)
		}

		models.MaxPay.Payoutmap[year][quarter][constant.EmployeePositionMap[position]][product] = pay
	}

	//validation
	for quarter, quarterMap := range models.MaxPay.Payoutmap[year] {
		if len(quarterMap) != 4 {
			return fmt.Errorf("FetchMaxPayout: invalid position length - " + fmt.Sprintf("quarter:%v, positionLen:%v", quarter, len(quarterMap)))
		}

		for position, positionMap := range quarterMap {
			if len(positionMap) != 4 {
				return fmt.Errorf("FetchMaxPayout: invalid product length - " + fmt.Sprintf("quarter:%v, position:%v, productLen:%v", quarter, position, len(positionMap)))
			}
		}
	}

	return nil
}

// FetchMaxExtraPayout loads the gloabal variable MaxExtraPay with the maximum extra payout limit of A
func (dbConn *DbServiceImpl) FetchMaxExtraPayout(year int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	rows, err := dbConn.DBClient.QueryContext(ctx, "select position, product, max_pay, quarter from max_extra_payout_table where deleted IS NULL and year=$1", year)
	if err != nil {
		return fmt.Errorf("FetchMaxExtraPayout: query - " + err.Error())
	}

	models.MaxExtraPay.Mutex.Lock()
	defer models.MaxExtraPay.Mutex.Unlock()

	if models.MaxExtraPay.ExtraPayoutmap == nil {
		models.MaxExtraPay.ExtraPayoutmap = make(map[int]map[int]map[string]map[string]float64)
	}

	if models.MaxExtraPay.ExtraPayoutmap[year] == nil {
		models.MaxExtraPay.ExtraPayoutmap[year] = make(map[int]map[string]map[string]float64)
	}

	for rows.Next() {

		var (
			position string
			product  string
			extraPay float64
			quarter  int
		)

		if err := rows.Scan(&position, &product, &extraPay, &quarter); err != nil {
			return fmt.Errorf("FetchMaxExtraPayout: scan - " + err.Error())
		}

		if quarter < 1 || quarter > 4 {
			return fmt.Errorf("FetchMaxExtraPayout: invalid quarter - " + fmt.Sprintf("year:%v, quarter:%v", year, quarter))
		}

		if position != constant.PromoterCC {
			return fmt.Errorf("FetchMaxExtraPayout: invalid position - " + fmt.Sprintf("year:%v, quarter:%v, position:%v, product:%v", year, quarter, position, product))
		}

		if !(product == constant.HBBCategory || product == constant.HayyakCategory || product == constant.BaqatiCategory || product == constant.UpgradeCategory) {
			return fmt.Errorf("FetchMaxExtraPayout: invalid product - " + fmt.Sprintf("year:%v, quarter:%v, position:%v, product:%v", year, quarter, position, product))
		}

		if extraPay == 0 {
			return fmt.Errorf("FetchMaxExtraPayout: invalid extraPay - " + fmt.Sprintf("year:%v, quarter:%v, position:%v, product:%v, extraPay:%v",
				year, quarter, position, product, extraPay))
		}

		if models.MaxExtraPay.ExtraPayoutmap[year][quarter] == nil {
			models.MaxExtraPay.ExtraPayoutmap[year][quarter] = make(map[string]map[string]float64)
		}

		if models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.EmployeePositionMap[position]] == nil {
			models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.EmployeePositionMap[position]] = make(map[string]float64)
		}

		models.MaxExtraPay.ExtraPayoutmap[year][quarter][constant.EmployeePositionMap[position]][product] = extraPay
	}

	//validation
	for quarter, quarterMap := range models.MaxExtraPay.ExtraPayoutmap[year] {
		if len(quarterMap) != 1 {
			return fmt.Errorf("FetchMaxExtraPayout: invalid position length - " + fmt.Sprintf("quarter:%v, positionLen:%v", quarter, len(quarterMap)))
		}

		for position, positionMap := range quarterMap {
			if len(positionMap) != 4 {
				return fmt.Errorf("FetchMaxExtraPayout: invalid product length - " + fmt.Sprintf("quarter:%v, position:%v, productLen:%v", quarter, position, len(positionMap)))
			}
		}
	}

	return nil
}

// FetchMaxThreshold loads the gloabal variable ThresholdCheck with the threshold pay limit of A,OM,M,SM
func (dbConn *DbServiceImpl) FetchMaxThreshold(year int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	rows, err := dbConn.DBClient.QueryContext(ctx, "select position, product, achieved_percent, payout, quarter from threshold_table where deleted IS NULL and year=$1", year)
	if err != nil {
		return fmt.Errorf("FetchMaxThreshold: query - " + err.Error())
	}

	models.ThresholdCheck.Mutex.Lock()
	defer models.ThresholdCheck.Mutex.Unlock()

	if models.ThresholdCheck.Limit == nil {
		models.ThresholdCheck.Limit = make(map[int]map[int]map[string]map[string]map[int]float64)
	}

	if models.ThresholdCheck.Limit[year] == nil {
		models.ThresholdCheck.Limit[year] = make(map[int]map[string]map[string]map[int]float64)
	}

	for rows.Next() {

		var (
			position        string
			product         string
			achievedPercent int
			threshold       float64
			quarter         int
		)

		if err := rows.Scan(&position, &product, &achievedPercent, &threshold, &quarter); err != nil {
			return fmt.Errorf("FetchMaxThreshold: scan - " + err.Error())
		}

		if quarter < 1 || quarter > 4 {
			return fmt.Errorf("FetchMaxThreshold: invalid quarter - " + fmt.Sprintf("year:%v, quarter:%v", year, quarter))
		}

		if !(position == constant.PromoterCC || position == constant.OutletManagerCC || position == constant.ManagerCC || position == constant.SeniorManagerCC) {
			return fmt.Errorf("FetchMaxThreshold: invalid position - " + fmt.Sprintf("year:%v, quarter:%v, position:%v, product:%v", year, quarter, position, product))
		}

		if !(product == constant.HBBCategory || product == constant.HayyakCategory || product == constant.BaqatiCategory || product == constant.UpgradeCategory) {
			return fmt.Errorf("FetchMaxThreshold: invalid product - " + fmt.Sprintf("year:%v, quarter:%v, position:%v, product:%v", year, quarter, position, product))
		}

		if achievedPercent == 0 {
			return fmt.Errorf("FetchMaxThreshold: invalid achievedPercent - " + fmt.Sprintf("year:%v, quarter:%v, position:%v, product:%v, achieved:%v",
				year, quarter, position, product, achievedPercent))
		}

		if threshold == 0 {
			return fmt.Errorf("FetchMaxThreshold: invalid threshold - " + fmt.Sprintf("year:%v, quarter:%v, position:%v, product:%v, achieved:%v, threshold:%v,",
				year, quarter, position, product, achievedPercent, threshold))
		}

		if models.ThresholdCheck.Limit[year][quarter] == nil {
			models.ThresholdCheck.Limit[year][quarter] = make(map[string]map[string]map[int]float64)
		}

		if models.ThresholdCheck.Limit[year][quarter][constant.EmployeePositionMap[position]] == nil {
			models.ThresholdCheck.Limit[year][quarter][constant.EmployeePositionMap[position]] = make(map[string]map[int]float64)
		}

		if models.ThresholdCheck.Limit[year][quarter][constant.EmployeePositionMap[position]][product] == nil {
			models.ThresholdCheck.Limit[year][quarter][constant.EmployeePositionMap[position]][product] = make(map[int]float64)
		}

		models.ThresholdCheck.Limit[year][quarter][constant.EmployeePositionMap[position]][product][achievedPercent] = threshold
	}

	//validation
	for quarter, quarterMap := range models.ThresholdCheck.Limit[year] {
		if len(quarterMap) != 3 {
			return fmt.Errorf("FetchMaxThreshold: invalid position length - " + fmt.Sprintf("quarter:%v, positionLen:%v", quarter, len(quarterMap)))
		}

		for position, positionMap := range quarterMap {
			if len(positionMap) != 4 {
				return fmt.Errorf("FetchMaxThreshold: invalid product length - " + fmt.Sprintf("quarter:%v, position:%v, productLen:%v", quarter, position, len(positionMap)))
			}

			for product, productMap := range positionMap {
				if len(productMap) != 66 {
					return fmt.Errorf("FetchMaxThreshold: invalid achievedPercent length - " + fmt.Sprintf("quarter:%v, position:%v, product:%v, achievedPercentLen:%v",
						quarter, position, product, len(productMap)))
				}
			}
		}
	}

	return nil
}

// FetchVOCCriteria loads the gloabal variable VOCDetail with the voc criteria value
func (dbConn *DbServiceImpl) FetchVOCCriteria(year int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var (
		unsatisfactory float64
		almostAchieved float64
		achieved       float64
		wellAchieved   float64
		outstanding    float64
	)

	models.VOCDetail.Mutex.Lock()
	defer models.VOCDetail.Mutex.Unlock()

	if models.VOCDetail.Criteria == nil {
		models.VOCDetail.Criteria = make(map[int]map[string]float64)
	}

	if models.VOCDetail.Criteria[year] == nil {
		models.VOCDetail.Criteria[year] = make(map[string]float64)
	}

	err := dbConn.DBClient.QueryRowContext(ctx, "select unsatisfactory, almost_achieved, achieved, well_achieved, outstanding from voc_criteria_table where deleted IS NULL and year=$1",
		year).Scan(&unsatisfactory, &almostAchieved, &achieved, &wellAchieved, &outstanding)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("FetchVOCCriteria: query - " + err.Error())
	}

	models.VOCDetail.Criteria[year] = map[string]float64{
		constant.VOCUnsatisfactory: unsatisfactory,
		constant.VOCAlmostAchieved: almostAchieved,
		constant.VOCAchieved:       achieved,
		constant.VOCWellAchieved:   wellAchieved,
		constant.VOCOutstanding:    outstanding,
	}

	return nil
}

// FetchVOCMaxPay loads the gloabal variable VOCDetail with the voc maxPay value
func (dbConn *DbServiceImpl) FetchVOCMaxPay(year int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	models.VOCDetail.Mutex.Lock()
	defer models.VOCDetail.Mutex.Unlock()

	if models.VOCDetail.MaxPay == nil {
		models.VOCDetail.MaxPay = make(map[int]map[string]float64)
	}

	if models.VOCDetail.MaxPay[year] == nil {
		models.VOCDetail.MaxPay[year] = make(map[string]float64)
	}

	rows, err := dbConn.DBClient.QueryContext(ctx, "select position, max_pay from voc_max_pay_table where deleted IS NULL and year=$1", year)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("FetchVOCMaxPay: query - " + err.Error())
	}

	for rows.Next() {

		var (
			position string
			maxPay   float64
		)

		if err := rows.Scan(&position, &maxPay); err != nil {
			return fmt.Errorf("FetchVOCMaxPay: scan - " + err.Error())
		}

		if !(position == constant.PromoterCC || position == constant.OutletManagerCC || position == constant.ManagerCC || position == constant.SeniorManagerCC) {
			return fmt.Errorf("FetchVOCMaxPay: invalid position - " + fmt.Sprintf("year:%v, position:%v", year, position))
		}

		if maxPay == 0 {
			return fmt.Errorf("FetchVOCMaxPay: invalid maxPay - " + fmt.Sprintf("year:%v, position:%v, maxPay:%v", year, position, maxPay))
		}

		models.VOCDetail.MaxPay[year][constant.EmployeePositionMap[position]] = maxPay
	}

	//validation
	// if len(models.VOCDetail.MaxPay[year]) != 4 {
	// 	return fmt.Errorf("FetchVOCMaxPay: invalid position length - " + fmt.Sprintf("year:%v, positionLen:%v", year, len(models.VOCDetail.MaxPay[year])))
	// }

	return nil
}
