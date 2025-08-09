package impl

import (
	"fmt"
	"math"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/database"
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/fcm"
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/redis"
	"gitlab.com/rendezvousPoint/incentive_calculator/repository/soap"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

var (
	getQuarterFromMonth = objecthelper.GetQuarterFromMonth
)

type UserImpl struct {
	Config       *configs.UserConfig
	DbService    database.DbService
	CacheService redis.CacheService
	FCMService   fcm.FCMService
	SOAPService  soap.SOAPService
}

func GetUserService(config *configs.UserConfig, dbService database.DbService, cacheService redis.CacheService, fcmService fcm.FCMService,
	soapService soap.SOAPService) *UserImpl {

	return &UserImpl{
		Config:       config,
		DbService:    dbService,
		CacheService: cacheService,
		FCMService:   fcmService,
		SOAPService:  soapService,
	}
}

func (impl *UserImpl) SessionValidation(userId, currentAccessToken string) bool {

	session, err := impl.CacheService.GetUserSessionToken(userId)
	if err != nil && err.Code != errors.UserNotFoundCode {
		return false
	}

	if err != nil && err.Code == errors.UserNotFoundCode {
		return true
	}

	if session.AccessToken == currentAccessToken {
		return true
	}

	return false
}

func (impl *UserImpl) getValueOfPoint(quarter, year int, grp string, achievedPercent float64) (float64, error) {
	models.ValOfPoint.Mutex.RLock()
	defer models.ValOfPoint.Mutex.RUnlock()

	if models.ValOfPoint.VopMap[year] == nil {
		return 0, errors.New("VOP YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.ValOfPoint.VopMap[year][quarter] == nil {
		return 0, errors.New("VOP QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.ValOfPoint.VopMap[year][quarter][grp] == nil {
		return 0, errors.New("VOP GROUP INVALID - " + grp)
	}

	switch {
	case achievedPercent < 75:
		return models.ValOfPoint.VopMap[year][quarter][grp][75].Point, nil
	case achievedPercent < 96:
		return models.ValOfPoint.VopMap[year][quarter][grp][95].Point, nil
	case achievedPercent <= 111:
		return models.ValOfPoint.VopMap[year][quarter][grp][110].Point, nil
	case achievedPercent < 140:
		return models.ValOfPoint.VopMap[year][quarter][grp][139].Point, nil
	case achievedPercent >= 140:
		return models.ValOfPoint.VopMap[year][quarter][grp][140].Point, nil
	}

	return 0, errors.New("VOP ACHIEVEMENT PERCENT INVALID")
}

func (impl *UserImpl) getExtraPay(quarter, year int, grp string, achievedPercent float64) (float64, error) {
	models.ValOfPoint.Mutex.RLock()
	defer models.ValOfPoint.Mutex.RUnlock()

	if models.ValOfPoint.VopMap[year] == nil {
		return 0, errors.New("VOEP YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.ValOfPoint.VopMap[year][quarter] == nil {
		return 0, errors.New("VOEP QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.ValOfPoint.VopMap[year][quarter][grp] == nil {
		return 0, errors.New("VOEP GROUP INVALID - " + grp)
	}

	switch {
	case achievedPercent >= 141 && achievedPercent < 160:
		return models.ValOfPoint.VopMap[year][quarter][grp][159].ExtraPay, nil
	case achievedPercent >= 160 && achievedPercent < 180:
		return models.ValOfPoint.VopMap[year][quarter][grp][179].ExtraPay, nil
	case achievedPercent >= 180 && achievedPercent < 190:
		return models.ValOfPoint.VopMap[year][quarter][grp][189].ExtraPay, nil
	case achievedPercent >= 190 && achievedPercent < 200:
		return models.ValOfPoint.VopMap[year][quarter][grp][199].ExtraPay, nil
	case achievedPercent >= 200:
		return models.ValOfPoint.VopMap[year][quarter][grp][200].ExtraPay, nil
	}

	return 0, nil
}

func (impl *UserImpl) getSIPRatePlan(quarter, year int, planName string) (float64, error) {

	planID := ""

	models.PlanRates.Mutex.RLock()
	defer models.PlanRates.Mutex.RUnlock()

	if models.PlanRates.PlanRateMap[year] == nil {
		return 0, errors.New("PLAN-NAME YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.PlanRates.PlanRateMap[year][quarter] == nil {
		return 0, errors.New("PLAN-NAME QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	for id, v := range models.PlanRates.PlanRateMap[year][quarter] {
		if v.PlanName == planName {
			planID = id
			break
		}
	}

	if planID != "" {
		return models.PlanRates.PlanRateMap[year][quarter][planID].SIPRatePlan, nil
	}

	return 0, errors.New("PLAN-NAME INVALID - " + planName)
}

func (impl *UserImpl) getSaleTargetPerOutlet(year, month, quarter int, isQuarterReport bool, grp, region, outlet string) (int, error) {

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
			return 0, errors.New("TARGET QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
		}
	}

	models.SalesTrgt.Mutex.RLock()
	defer models.SalesTrgt.Mutex.RUnlock()

	if models.SalesTrgt.SaleTargetMap[year] == nil {
		return 0, errors.New("TARGET YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	rMap, ok := models.SalesTrgt.SaleTargetMap[year][region]
	if !ok {
		return 0, errors.New("TARGET REGION NOT FOUND - " + region)
	}

	if _, ok := rMap[outlet]; !ok {
		return 0, errors.New("TARGET OUTLET NOT FOUND - " + outlet)
	}

	if isQuarterReport {
		//quarter target
		result := 0
		switch grp {
		case constant.HBBCategory:
			for i := start - 1; i < end; i++ {
				if models.SalesTrgt.SaleTargetMap[year][region][outlet][i] == nil {
					return 0, fmt.Errorf("TARGET HBB BLANK FOR MONTH - %v", i)
				}
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].HBB
			}
			return result, nil

		case constant.HayyakCategory:
			for i := start - 1; i < end; i++ {
				if models.SalesTrgt.SaleTargetMap[year][region][outlet][i] == nil {
					return 0, fmt.Errorf("TARGET HAYYAK BLANK FOR MONTH - %v", i)
				}
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].Hayyak
			}
			return result, nil

		case constant.BaqatiCategory:
			for i := start - 1; i < end; i++ {
				if models.SalesTrgt.SaleTargetMap[year][region][outlet][i] == nil {
					return 0, fmt.Errorf("TARGET BAQATI BLANK FOR MONTH - %v", i)
				}
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].Baqati
			}
			return result, nil

		case constant.UpgradeCategory:
			for i := start - 1; i < end; i++ {
				if models.SalesTrgt.SaleTargetMap[year][region][outlet][i] == nil {
					return 0, fmt.Errorf("TARGET UPGRADE BLANK FOR MONTH - %v", i)
				}
				result += models.SalesTrgt.SaleTargetMap[year][region][outlet][i].Upgrade
			}
			return result, nil

		}
	} else {
		//month target
		switch grp {
		case constant.HBBCategory:
			return models.SalesTrgt.SaleTargetMap[year][region][outlet][month-1].HBB, nil
		case constant.HayyakCategory:
			return models.SalesTrgt.SaleTargetMap[year][region][outlet][month-1].Hayyak, nil
		case constant.BaqatiCategory:
			return models.SalesTrgt.SaleTargetMap[year][region][outlet][month-1].Baqati, nil
		case constant.UpgradeCategory:
			return models.SalesTrgt.SaleTargetMap[year][region][outlet][month-1].Upgrade, nil
		}
	}

	return 0, errors.New("TARGET GROUP INVALID - " + grp)
}

func (impl *UserImpl) getMaxPayout(quarter, year int, position, grp string) (float64, error) {

	pos := constant.EmployeePositionMap[position]

	models.MaxPay.Mutex.RLock()
	defer models.MaxPay.Mutex.RUnlock()

	if models.MaxPay.Payoutmap[year] == nil {
		return 0, errors.New("MAX PAYOUT YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.MaxPay.Payoutmap[year][quarter] == nil {
		return 0, errors.New("MAX PAYOUT QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.MaxPay.Payoutmap[year][quarter][pos] == nil {
		return 0, errors.New("MAX PAYOUT POSITION INVALID - " + fmt.Sprintf("%v", position))
	}

	val, ok := models.MaxPay.Payoutmap[year][quarter][pos][grp]
	if !ok {
		return 0, errors.New("MAX PAYOUT GROUP INVALID - " + grp)
	}

	return val, nil
}

func (impl *UserImpl) getMaxExtraPayout(quarter, year int, position, grp string) (float64, error) {

	pos := constant.EmployeePositionMap[position]

	models.MaxExtraPay.Mutex.RLock()
	defer models.MaxExtraPay.Mutex.RUnlock()

	if models.MaxExtraPay.ExtraPayoutmap[year] == nil {
		return 0, errors.New("MAX EXTRA PAYOUT YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.MaxExtraPay.ExtraPayoutmap[year][quarter] == nil {
		return 0, errors.New("MAX EXTRA PAYOUT QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.MaxExtraPay.ExtraPayoutmap[year][quarter][pos] == nil {
		return 0, errors.New("MAX EXTRA PAYOUT POSITION INVALID - " + fmt.Sprintf("%v", position))
	}

	val, ok := models.MaxExtraPay.ExtraPayoutmap[year][quarter][pos][grp]
	if !ok {
		return 0, errors.New("MAX EXTRA PAYOUT GROUP INVALID - " + grp)
	}

	return val, nil
}

func (impl *UserImpl) getVOCCriteria(year int, achievedPercent float64) (float64, error) {

	models.VOCDetail.Mutex.RLock()
	defer models.VOCDetail.Mutex.RUnlock()

	if models.VOCDetail.Criteria[year] == nil {
		return 0, errors.New("VOC CRITERIA YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	switch {
	case achievedPercent < 85:
		return models.VOCDetail.Criteria[year][constant.VOCUnsatisfactory], nil
	case achievedPercent >= 85 && achievedPercent < 89:
		return models.VOCDetail.Criteria[year][constant.VOCAlmostAchieved], nil
	case achievedPercent >= 89 && achievedPercent < 93:
		return models.VOCDetail.Criteria[year][constant.VOCAchieved], nil
	case achievedPercent >= 93 && achievedPercent <= 97:
		return models.VOCDetail.Criteria[year][constant.VOCWellAchieved], nil
	default:
		return models.VOCDetail.Criteria[year][constant.VOCOutstanding], nil
	}

}

func (impl *UserImpl) getVOCMaxPay(year int, position string) (float64, error) {

	models.VOCDetail.Mutex.RLock()
	defer models.VOCDetail.Mutex.RUnlock()

	if models.VOCDetail.MaxPay[year] == nil {
		return 0, errors.New("VOC MAX-PAY YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	return models.VOCDetail.MaxPay[year][position], nil
}

func (impl *UserImpl) getThresholdPayout(quarter, year int, position, grp string, achievedPercent float64) (float64, error) {

	pos := constant.EmployeePositionMap[position]
	achieved := int(math.Round(achievedPercent))

	models.ThresholdCheck.Mutex.RLock()
	defer models.ThresholdCheck.Mutex.RUnlock()

	if achieved < 75 {
		return 0, nil
	}

	if achieved > 140 {
		achieved = 140
	}

	if models.ThresholdCheck.Limit[year] == nil {
		return 0, errors.New("THRESHOLD PAYOUT YEAR INVALID - " + fmt.Sprintf("%v", year))
	}

	if models.ThresholdCheck.Limit[year][quarter] == nil {
		return 0, errors.New("THRESHOLD PAYOUT QUARTER INVALID - " + fmt.Sprintf("%v", quarter))
	}

	if models.ThresholdCheck.Limit[year][quarter][pos] == nil {
		return 0, errors.New("THRESHOLD PAYOUT POSITION INVALID - " + fmt.Sprintf("%v", position))
	}

	if models.ThresholdCheck.Limit[year][quarter][pos][grp] == nil {
		return 0, errors.New("THRESHOLD PAYOUT GROUP INVALID - " + fmt.Sprintf("%v", grp))
	}

	val, ok := models.ThresholdCheck.Limit[year][quarter][pos][grp][achieved]
	if !ok {
		return 0, errors.New("THRESHOLD PAYOUT ACHIEVED PERCENT INVALID - " + grp)
	}

	return val, nil
}
