package business

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
)

type IncentiveImpl struct {
	Config *configs.IncentiveConfig
	PsqlDb *psqldb.DBImplementation
}

func NewIncentiveService(config *configs.IncentiveConfig, dbConfig *models.DbClientConfig) *IncentiveImpl {

	return &IncentiveImpl{
		Config: config,
		PsqlDb: psqldb.NewDatabaseSetup(dbConfig),
	}
}

func (impl *IncentiveImpl) getValueOfPoint(grp string, achievedPercent float64) float64 {
	models.ValOfPoint.Mutex.RLock()
	defer models.ValOfPoint.Mutex.RUnlock()

	switch {
	case achievedPercent < 75:
		return models.ValOfPoint.VopMap[grp][75].Point
	case achievedPercent < 96:
		return models.ValOfPoint.VopMap[grp][95].Point
	case achievedPercent <= 111:
		return models.ValOfPoint.VopMap[grp][110].Point
	case achievedPercent < 140:
		return models.ValOfPoint.VopMap[grp][139].Point
	case achievedPercent >= 140:
		return models.ValOfPoint.VopMap[grp][140].Point
	}

	return 0
}

func (impl *IncentiveImpl) getExtraPay(grp string, achievedPercent float64) float64 {
	models.ValOfPoint.Mutex.RLock()
	defer models.ValOfPoint.Mutex.RUnlock()

	switch {
	case achievedPercent >= 141 && achievedPercent < 160:
		return models.ValOfPoint.VopMap[grp][159].ExtraPay
	case achievedPercent >= 160 && achievedPercent < 180:
		return models.ValOfPoint.VopMap[grp][179].ExtraPay
	case achievedPercent >= 180 && achievedPercent < 190:
		return models.ValOfPoint.VopMap[grp][189].ExtraPay
	case achievedPercent >= 190 && achievedPercent < 200:
		return models.ValOfPoint.VopMap[grp][199].ExtraPay
	case achievedPercent >= 200:
		return models.ValOfPoint.VopMap[grp][200].ExtraPay
	}

	return 0
}

func (impl *IncentiveImpl) getSIPRatePlan(planName string) float64 {

	planID := 0

	models.PlanRates.Mutex.RLock()
	defer models.PlanRates.Mutex.RUnlock()

	for id, v := range models.PlanRates.PlanRateMap {
		if v.PlanName == planName {
			planID = id
			break
		}
	}

	if planID != 0 {
		return models.PlanRates.PlanRateMap[planID].SIPRatePlan
	}

	return 0
}

func (impl *IncentiveImpl) getSaleTargetPerOutlet(t time.Time, quarter int, isQuarterReport bool, grp, region, outlet string) int {

	fmt.Println("check ", t, quarter, isQuarterReport, grp, region, outlet)

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

	models.SalesTrgt.Mutex.RLock()
	defer models.SalesTrgt.Mutex.RUnlock()

	rMap, ok := models.SalesTrgt.SaleTargetMap[region]
	if !ok {
		fmt.Println("region")
		return -1
	}

	if _, ok := rMap[outlet]; !ok {
		fmt.Println("outlet")
		return -1
	}

	if isQuarterReport {
		//quarter target
		result := 0
		switch grp {
		case constant.HBBCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].HBB
			}
			return result
		case constant.HayyakCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].Hayyak
			}
			return result
		case constant.BaqatiCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].Baqati
			}
			return result
		case constant.UpgradeCategory:
			for i := start - 1; i < end; i++ {
				result += models.SalesTrgt.SaleTargetMap[region][outlet][i].Upgrade
			}
			return result
		}
	} else {
		//month target
		switch grp {
		case constant.HBBCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][int(t.Month()-1)].HBB
		case constant.HayyakCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][int(t.Month()-1)].Hayyak
		case constant.BaqatiCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][int(t.Month()-1)].Baqati
		case constant.UpgradeCategory:
			return models.SalesTrgt.SaleTargetMap[region][outlet][int(t.Month()-1)].Upgrade
		}
	}

	fmt.Println("here")
	return -1
}

func (impl *IncentiveImpl) getMaxPayout(grp string) float64 {
	models.MaxPay.Mutex.RLock()
	defer models.MaxPay.Mutex.RUnlock()

	return models.MaxPay.Payoutmap[grp]
}

func (impl *IncentiveImpl) getQuarterFromMonth(month int) int {
	switch {
	case month >= 1 && month <= 3:
		return 1
	case month >= 4 && month <= 6:
		return 2
	case month >= 7 && month <= 9:
		return 3
	case month >= 10 && month <= 12:
		return 4
	}

	return -1
}

func (impl *IncentiveImpl) generateOTP(max int) (string, error) {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if err != nil {
		return "", err
	}

	if n != max {
		return "", errors.New("incorrect otp length")
	}

	for i := 0; i < len(b); i++ {
		b[i] = constant.Table[int(b[i])%len(constant.Table)]
	}

	return string(b), nil
}

func (impl *IncentiveImpl) isRegion(region string) bool {

	if region == "BATINAH NORTH" || region == "BATINAH SOUTH" || region == "DAKHLIA & DAHIRAH" || region == "DHOFAR" ||
		region == "MUSCAT NORTH" || region == "MUSCAT SOUTH" || region == "SHARQIYA" || region == "GRAND TOTAL" {
		return true
	}

	return false
}

func (impl *IncentiveImpl) isProduct(product string) bool {

	if product == constant.HBBCategory || product == constant.BaqatiCategory || product == constant.HayyakCategory || product == constant.UpgradeCategory {
		return true
	}

	if product == "UPGRADE (BAQATI + HBB)" {
		return true
	}

	return false
}
