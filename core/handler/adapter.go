package handler

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/db"
	"time"
)

func (impl *SalesIncentiveImpl) getValueOfPoint(grp string, achievedPercent float64) float64 {
	db.ValOfPoint.Mutex.RLock()
	defer db.ValOfPoint.Mutex.RUnlock()

	switch {
	case achievedPercent < 75:
		return db.ValOfPoint.VopMap[grp][75].Point
	case achievedPercent < 96:
		return db.ValOfPoint.VopMap[grp][95].Point
	case achievedPercent <= 111:
		return db.ValOfPoint.VopMap[grp][110].Point
	case achievedPercent < 140:
		return db.ValOfPoint.VopMap[grp][139].Point
	case achievedPercent >= 140:
		return db.ValOfPoint.VopMap[grp][140].Point
	}

	return 0
}

func (impl *SalesIncentiveImpl) getExtraPay(grp string, achievedPercent float64) float64 {
	db.ValOfPoint.Mutex.RLock()
	defer db.ValOfPoint.Mutex.RUnlock()

	switch {
	case achievedPercent >= 141 && achievedPercent < 160:
		return db.ValOfPoint.VopMap[grp][159].ExtraPay
	case achievedPercent >= 160 && achievedPercent < 180:
		return db.ValOfPoint.VopMap[grp][179].ExtraPay
	case achievedPercent >= 180 && achievedPercent < 190:
		return db.ValOfPoint.VopMap[grp][189].ExtraPay
	case achievedPercent >= 190 && achievedPercent < 200:
		return db.ValOfPoint.VopMap[grp][199].ExtraPay
	case achievedPercent >= 200:
		return db.ValOfPoint.VopMap[grp][200].ExtraPay
	}

	return 0
}

func (impl *SalesIncentiveImpl) getSIPRatePlan(planName string) float64 {

	planID := 0

	db.PlanRates.Mutex.RLock()
	defer db.PlanRates.Mutex.RUnlock()

	for id, v := range db.PlanRates.PlanRateMap {
		if v.PlanName == planName {
			planID = id
			break
		}
	}

	if planID != 0 {
		return db.PlanRates.PlanRateMap[planID].SIPRatePlan
	}

	return 0
}

func (impl *SalesIncentiveImpl) getSaleTargetPerOutlet(t time.Time, quarter int, isQuarterReport bool, grp, region, outlet string) int {

	var (
		start int
		end   int
	)

	if isQuarterReport {
		switch quarter {
		case 1:
			start = 1
			end = 3
		case 2:
			start = 4
			end = 6
		case 3:
			start = 7
			end = 9
		case 4:
			start = 10
			end = 12
		default:
			return -1
		}
	}

	db.SalesTrgt.Mutex.RLock()
	defer db.SalesTrgt.Mutex.RUnlock()

	rMap, ok := db.SalesTrgt.SaleTargetMap[region]
	if !ok {
		return -1
	}

	if _, ok := rMap[outlet]; !ok {
		return -1
	}

	if isQuarterReport {
		//quarter target
		result := 0
		switch grp {
		case commons.HBBCategory:
			for i:=start-1; i < end; i++ {
				result += db.SalesTrgt.SaleTargetMap[region][outlet][i].HBB
			}
			return result
		case commons.HayyakCategory:
			for i:=start-1; i < end; i++  {
				result += db.SalesTrgt.SaleTargetMap[region][outlet][i].Hayyak
			}
			return result
		case commons.BaqatiCategory:
			for i:=start-1; i < end; i++  {
				result += db.SalesTrgt.SaleTargetMap[region][outlet][i].Baqati
			}
			return result
		case commons.UpgradeCategory:
			for i:=start-1; i < end; i++  {
				result += db.SalesTrgt.SaleTargetMap[region][outlet][i].Upgrade
			}
			return result
		}
	} else {
		//month target
		switch grp {
		case commons.HBBCategory:
			return db.SalesTrgt.SaleTargetMap[region][outlet][int(t.Month()-1)].HBB
		case commons.HayyakCategory:
			return db.SalesTrgt.SaleTargetMap[region][outlet][int(t.Month()-1)].Hayyak
		case commons.BaqatiCategory:
			return db.SalesTrgt.SaleTargetMap[region][outlet][int(t.Month()-1)].Baqati
		case commons.UpgradeCategory:
			return db.SalesTrgt.SaleTargetMap[region][outlet][int(t.Month()-1)].Upgrade
		}
	}

	return -1
}

func (impl *SalesIncentiveImpl) getMaxPayout(grp string) float64 {
	db.MaxPay.Mutex.RLock()
	defer db.MaxPay.Mutex.RUnlock()

	return db.MaxPay.Payoutmap[grp]
}
