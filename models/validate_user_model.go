package models

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func ValidateUserMapRequest(req *UserMapRequest, ctxHeader context.ContextHeader) (isQuarterReport bool, quarter, month, year int, err error) {

	if req.OutletName == "" && req.RegionName == "" {
		err = errors.New("blank outlet and region")
		return
	}

	req.OutletName = strings.ToUpper(strings.TrimSpace(req.OutletName))
	req.RegionName = strings.ToUpper(strings.TrimSpace(req.RegionName))

	var date time.Time

	isQuarterReport = true

	if req.Monthly == "" && req.Quarterly == "" {
		err = errors.New("monthy and quartlerly type is blank")
		return
	}

	if req.Monthly != "" {
		date, err = time.Parse("01-2006", strings.TrimSpace(req.Monthly))
		if err != nil {
			err = errors.New("monthly parsing: " + err.Error())
			return
		}

		isQuarterReport = false
		month = int(date.Month())
		year = date.Year()
	}

	if isQuarterReport && req.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(req.Quarterly), "-")
		if len(req) < 2 {
			err = errors.New("quarterly type is incorrect, avoid spaces")
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			err = errors.New("quarter conversion: " + err.Error())
			return

		}

		if quarter < 1 || quarter > 4 {
			err = errors.New("incorrect quarter value - should be between 1 and 4")
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			err = errors.New("quarter year parsing: " + err.Error())
			return
		}

		year = date.Year()
	}

	return
}

func ValidateUserIncentiveRequest(req *UserIncentiveRequest, ctxHeader context.ContextHeader) (isQuarterReport bool, quarter, month, year int, err error) {

	if req.ClientId == "" {
		err = errors.New("blank userId")
		return
	}

	// if !ctxHeader.IsAdmin && ctxHeader.XUserId != req.ClientId {
	// 	err = fmt.Errorf("cannot query for other user, authUser:%v, queryUser:%+v", ctxHeader.XUserId, req.ClientId)
	// 	return
	// }

	var date time.Time

	isQuarterReport = true

	if req.Monthly == "" && req.Quarterly == "" {
		err = errors.New("monthy and quartlerly type is blank")
		return
	}

	if req.Monthly != "" {
		date, err = time.Parse("01-2006", strings.TrimSpace(req.Monthly))
		if err != nil {
			err = errors.New("monthly parsing: " + err.Error())
			return
		}

		isQuarterReport = false
		month = int(date.Month())
		year = date.Year()
	}

	if isQuarterReport && req.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(req.Quarterly), "-")
		if len(req) < 2 {
			err = errors.New("quarterly type is incorrect, avoid spaces")
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			err = errors.New("quarter conversion: " + err.Error())
			return

		}

		if quarter < 1 || quarter > 4 {
			err = errors.New("incorrect quarter value - should be between 1 and 4")
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			err = errors.New("quarter year parsing: " + err.Error())
			return
		}

		year = date.Year()
	}

	return
}

func ValidateUserSalesRequest(req *UserSalesRequest, ctxHeader context.ContextHeader) (isQuarterReport bool, quarter, month, year int, err error) {

	if req.ClientId == "" {
		err = errors.New("blank userId")
		return
	}

	// if !ctxHeader.IsAdmin && ctxHeader.XUserId != req.ClientId {
	// 	err = fmt.Errorf("cannot query for other user, authUser:%v, queryUser:%+v", ctxHeader.XUserId, req.ClientId)
	// 	return
	// }

	product := strings.ToUpper(req.Product)

	if !(product == constant.HBBCategory || product == constant.BaqatiCategory || product == constant.HayyakCategory || product == constant.UpgradeCategory) {
		err = errors.New("incorrect product type")
		return
	}

	var date time.Time

	isQuarterReport = true

	if req.Monthly == "" && req.Quarterly == "" {
		err = errors.New("monthy and quartlerly type is blank")
		return
	}

	if req.Monthly != "" {
		date, err = time.Parse("01-2006", strings.TrimSpace(req.Monthly))
		if err != nil {
			err = errors.New("monthly parsing: " + err.Error())
			return
		}

		isQuarterReport = false
		month = int(date.Month())
		year = date.Year()
	}

	if isQuarterReport && req.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(req.Quarterly), "-")
		if len(req) < 2 {
			err = errors.New("quarterly type is incorrect, avoid spaces")
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			err = errors.New("quarter conversion: " + err.Error())
			return

		}

		if quarter < 1 || quarter > 4 {
			err = errors.New("incorrect quarter value - should be between 1 and 4")
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			err = errors.New("quarter year parsing: " + err.Error())
			return
		}

		year = date.Year()
	}

	return
}

func ValidateUserSaleCountRequest(req *UserSalesRequest, ctxHeader context.ContextHeader) (isQuarterReport bool, quarter, month, year int, err error) {

	if req.ClientId == "" {
		err = errors.New("blank userId")
		return
	}

	// if !ctxHeader.IsAdmin && ctxHeader.XUserId != req.ClientId {
	// 	err = fmt.Errorf("cannot query for other user, authUser:%v, queryUser:%+v", ctxHeader.XUserId, req.ClientId)
	// 	return
	// }

	// product := strings.ToUpper(req.Product)

	// if !(product == constant.HBBCategory || product == constant.BaqatiCategory || product == constant.HayyakCategory || product == constant.UpgradeCategory) {
	// 	err = errors.New("incorrect product type")
	// 	return
	// }

	var date time.Time

	isQuarterReport = true

	if req.Monthly == "" && req.Quarterly == "" {
		err = errors.New("monthy and quartlerly type is blank")
		return
	}

	if req.Monthly != "" {
		date, err = time.Parse("01-2006", strings.TrimSpace(req.Monthly))
		if err != nil {
			err = errors.New("monthly parsing: " + err.Error())
			return
		}

		isQuarterReport = false
		month = int(date.Month())
		year = date.Year()
	}

	if isQuarterReport && req.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(req.Quarterly), "-")
		if len(req) < 2 {
			err = errors.New("quarterly type is incorrect, avoid spaces")
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			err = errors.New("quarter conversion: " + err.Error())
			return

		}

		if quarter < 1 || quarter > 4 {
			err = errors.New("incorrect quarter value - should be between 1 and 4")
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			err = errors.New("quarter year parsing: " + err.Error())
			return
		}

		year = date.Year()
	}

	return
}

func ValidateSingleUserTargetChartResponse(res *ChartDataResponse) {

	if res.HBB == nil {
		res.HBB = &TargetAchievementData{}
	}
	if res.Hayyak == nil {
		res.Hayyak = &TargetAchievementData{}
	}
	if res.Baqati == nil {
		res.Baqati = &TargetAchievementData{}
	}
	if res.Upgrade == nil {
		res.Upgrade = &TargetAchievementData{}
	}

	res.HBB.Target = math.Round(res.HBB.Target*1000) / 1000
	res.HBB.AchievedPercentage = math.Round(res.HBB.AchievedPercentage*1000) / 1000
	res.Hayyak.Target = math.Round(res.Hayyak.Target*1000) / 1000
	res.Hayyak.AchievedPercentage = math.Round(res.Hayyak.AchievedPercentage*1000) / 1000
	res.Baqati.Target = math.Round(res.Baqati.Target*1000) / 1000
	res.Baqati.AchievedPercentage = math.Round(res.Baqati.AchievedPercentage*1000) / 1000
	res.Upgrade.Target = math.Round(res.Upgrade.Target*1000) / 1000
	res.Upgrade.AchievedPercentage = math.Round(res.Upgrade.AchievedPercentage*1000) / 1000
}

func ValidateTopTenSalesResponse(res []*ProductWiseSalesTargetPerStaff) {

	for index, val := range res {
		res[index].HBBTarget = math.Round(val.HBBTarget*1000) / 1000
		res[index].HBBAchievedPercent = math.Round(val.HBBAchievedPercent*1000) / 1000
		res[index].HayyakTarget = math.Round(val.HayyakTarget*1000) / 1000
		res[index].HayyakAchievedPercent = math.Round(val.HayyakAchievedPercent*1000) / 1000
		res[index].BaqatiTarget = math.Round(val.BaqatiTarget*1000) / 1000
		res[index].BaqatiAchievedPercent = math.Round(val.BaqatiAchievedPercent*1000) / 1000
		res[index].UpgradeTarget = math.Round(val.UpgradeTarget*1000) / 1000
		res[index].UpgradeAchievedPercent = math.Round(val.UpgradeAchievedPercent*1000) / 1000
		res[index].TotalTarget = math.Round(val.TotalTarget*1000) / 1000
		res[index].TotalAchievedPercentage = math.Round(val.TotalAchievedPercentage*1000) / 1000
	}
}

func ValidateUserPayoutResponse(ctx *gin.Context, res *TotalIncentive) {

	var achievedCount int

	if res.HBB == nil {
		res.HBB = &IndividualIncentive{}
	}
	if res.Hayyak == nil {
		res.Hayyak = &IndividualIncentive{}
	}
	if res.Baqati == nil {
		res.Baqati = &IndividualIncentive{}
	}
	if res.Upgrade == nil {
		res.Upgrade = &IndividualIncentive{}
	}
	if res.VOC == nil {
		res.VOC = &IndividualIncentive{}
	}

	if res.HBB.Payout != 0 {
		achievedCount++
	} else {
		//logger.FileLogger.CtxInfo(ctx).Str("HBB", "not achieved").Msg("GetUserPayout : response validation")
	}

	if res.Hayyak.Payout != 0 {
		achievedCount++
	} else {
		//logger.FileLogger.CtxInfo(ctx).Str("HAYYAK", "not achieved").Msg("GetUserPayout : response validation")
	}

	if res.Baqati.Payout != 0 {
		achievedCount++
	} else {
		//logger.FileLogger.CtxInfo(ctx).Str("BAQATI", "not achieved").Msg("GetUserPayout : response validation")
	}

	if res.Upgrade.Payout != 0 {
		achievedCount++
	} else {
		//logger.FileLogger.CtxInfo(ctx).Str("UPGRADE", "not achieved").Msg("GetUserPayout : response validation")
	}

	if achievedCount < 3 {
		res.TotalPay = 0
		res.HBB.Payout = 0
		res.HBB.ExtraPayout = 0
		res.Hayyak.Payout = 0
		res.Hayyak.ExtraPayout = 0
		res.Baqati.Payout = 0
		res.Baqati.ExtraPayout = 0
		res.Upgrade.Payout = 0
		res.Upgrade.ExtraPayout = 0
		res.VOC.Payout = 0
		res.VOC.ExtraPayout = 0

		return
	}

	res.TotalPay = math.Round(res.TotalPay*1000) / 1000
	res.HBB.Payout = math.Round(res.HBB.Payout*1000) / 1000
	res.HBB.ExtraPayout = math.Round(res.HBB.ExtraPayout*1000) / 1000
	res.Hayyak.Payout = math.Round(res.Hayyak.Payout*1000) / 1000
	res.Hayyak.ExtraPayout = math.Round(res.Hayyak.ExtraPayout*1000) / 1000
	res.Baqati.Payout = math.Round(res.Baqati.Payout*1000) / 1000
	res.Baqati.ExtraPayout = math.Round(res.Baqati.ExtraPayout*1000) / 1000
	res.Upgrade.Payout = math.Round(res.Upgrade.Payout*1000) / 1000
	res.Upgrade.ExtraPayout = math.Round(res.Upgrade.ExtraPayout*1000) / 1000
}

func ValidateUserReportRequest(req *UserReportRequest, ctxHeader context.ContextHeader) (isQuarterReport bool, err error) {

	if req.ClientId == "" {
		err = errors.New("blank userId")
		return
	}

	// if req.ClientId == "" || ctxHeader.XUserId != req.ClientId {
	// 	err = errors.New("incorrect user token")
	// 	return
	// }

	if !((req.Month >= 1 && req.Month <= 12) || (req.Quarter >= 1 && req.Quarter <= 4)) {
		err = errors.New("monthly or quarterly type is incorrect")
		return
	}

	if req.Quarter >= 1 && req.Quarter <= 4 {
		isQuarterReport = true
	}

	return isQuarterReport, nil
}

func ValidateUserReportResponse(ctx *gin.Context, res *UserReportDB) {

	var achievedCount int

	if res.HBBPayout != 0 {
		achievedCount++
	} else {
		//logger.FileLogger.CtxInfo(ctx).Str("HBB", "not achieved").Msg("GetUserReport : response validation")
	}

	if res.HayyakPayout != 0 {
		achievedCount++
	} else {
		//logger.FileLogger.CtxInfo(ctx).Str("HAYYAK", "not achieved").Msg("GetUserReport : response validation")
	}

	if res.BaqatiPayout != 0 {
		achievedCount++
	} else {
		//logger.FileLogger.CtxInfo(ctx).Str("BAQATI", "not achieved").Msg("GetUserReport : response validation")
	}

	if res.UpgradePayout != 0 {
		achievedCount++
	} else {
		//logger.FileLogger.CtxInfo(ctx).Str("UPGRADE", "not achieved").Msg("GetUserReport : response validation")
	}

	if achievedCount < 3 {
		res.HBBPayout = 0
		res.HBBExtraPayout = 0
		res.HayyakPayout = 0
		res.HayyakExtraPayout = 0
		res.BaqatiPayout = 0
		res.BaqatiExtraPayout = 0
		res.UpgradePayout = 0
		res.UpgradeExtraPayout = 0
		res.TotalPayout = 0

		return
	}

	res.HBBTargetCount = math.Round(res.HBBTargetCount*1000) / 1000
	res.HBBPayout = math.Round(res.HBBPayout*1000) / 1000
	res.HBBExtraPayout = math.Round(res.HBBExtraPayout*1000) / 1000
	res.HBBAchievedPercentage = math.Round(res.HBBAchievedPercentage*1000) / 1000

	res.HayyakTargetCount = math.Round(res.HayyakTargetCount*1000) / 1000
	res.HayyakPayout = math.Round(res.HayyakPayout*1000) / 1000
	res.HayyakExtraPayout = math.Round(res.HayyakExtraPayout*1000) / 1000
	res.HayyakAchievedPercentage = math.Round(res.HayyakAchievedPercentage*1000) / 1000

	res.BaqatiTargetCount = math.Round(res.BaqatiTargetCount*1000) / 1000
	res.BaqatiPayout = math.Round(res.BaqatiPayout*1000) / 1000
	res.BaqatiExtraPayout = math.Round(res.BaqatiExtraPayout*1000) / 1000
	res.BaqatiAchievedPercentage = math.Round(res.BaqatiAchievedPercentage*1000) / 1000

	res.UpgradeTargetCount = math.Round(res.UpgradeTargetCount*1000) / 1000
	res.UpgradePayout = math.Round(res.UpgradePayout*1000) / 1000
	res.UpgradeExtraPayout = math.Round(res.UpgradeExtraPayout*1000) / 1000
	res.UpradeAchievedPercentage = math.Round(res.UpradeAchievedPercentage*1000) / 1000

	res.TotalPayout = math.Round(res.TotalPayout*1000) / 1000
}
