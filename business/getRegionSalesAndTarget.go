package business

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *IncentiveImpl) GetRegionSalesAndTarget(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesTargetPerOutlet, *errors.Error) {

	// cRH, _ := ctx.Get(constant.HeaderInContext)
	// ctxHeader, _ := (cRH).(context.ContextHeader)

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
	}

	// if ctxHeader.IsAdmin {
	// 	vps, err := impl.PsqlDb.FetchVP(quarter, year)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	userID = vps[0]
	// }

	user, err := impl.PsqlDb.FetchUserBasicDetails(userID, quarter, year)
	if err != nil {
		return nil, err
	}

	position, ok := constant.ManagerMap[user.Position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT MANAGER POSITION - " + user.Position,
			Details:    "GetRegionSalesAndTarget: incorrect manager position: userId-" + user.UserId + " position-" + user.Position,
		}
	}
	user.Position = position

	resp := make([]*models.ProductWiseSalesTargetPerOutlet, 0)

	/*
		For Outlet Managers
	*/

	if user.Position == constant.OutletManagerAbb {

		cacheReport, err := impl.RedisCache.GetRegionReport(user.Region, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			return append(resp, impl.singleRegionSalesAndTarget(ctx, cacheReport)), nil
		}

		regionResp, err := impl.updateRegionReport(ctx, user.Region, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		// val, _ := json.Marshal(outletResp)
		// logger.FileLogger.CtxInfo(ctx).Interface("OM regionReport", string(val)).Msg("GetRegionSalesAndTarget : trace")

		return append(resp, impl.singleRegionSalesAndTarget(ctx, regionResp)), nil
	}

	/*
		For Other Managers
	*/

	allDistinctOM, err := impl.PsqlDb.FetchAllOMReportingToGivenManager(user.Position, userID, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("user", userID).Str("position", user.Position).Msg("GetRegionSalesAndTarget : no OM found")
		return resp, nil
	}

	for _, eachOM := range allDistinctOM {

		if eachOM.Position != constant.OutletManagerAbb {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Message:    eachOM.UserId + ": INCORRECT MANAGER POSITION - " + eachOM.Position,
				Details:    "GetRegionSalesAndTarget: incorrect manager position: userId-" + eachOM.UserId + " position-" + eachOM.Position,
			}
		}

		cacheReport, err := impl.RedisCache.GetRegionReport(eachOM.Region, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			resp = append(resp, impl.singleRegionSalesAndTarget(ctx, cacheReport))
		} else {
			outletResp, err := impl.updateRegionReport(ctx, eachOM.Region, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			// val, _ := json.Marshal(outletResp)
			// logger.FileLogger.CtxInfo(ctx).Interface("OtherMngr regionReport", string(val)).Msg("GetRegionSalesAndTarget : trace")

			resp = append(resp, impl.singleRegionSalesAndTarget(ctx, outletResp))
		}
	}

	return resp, nil
}

func (impl *IncentiveImpl) singleRegionSalesAndTarget(ctx *gin.Context, regionReport *models.RegionReportDB) *models.ProductWiseSalesTargetPerOutlet {

	eachRegionResp := &models.ProductWiseSalesTargetPerOutlet{
		Region:        regionReport.Region,
		HBBSales:      regionReport.HBBSalesCount,
		HayyakSales:   regionReport.HayyakSalesCount,
		BaqatiSales:   regionReport.BaqatiSalesCount,
		UpgradeSales:  regionReport.UpgradeSalesCount,
		HBBTarget:     regionReport.HBBTargetCount,
		HayyakTarget:  regionReport.HayyakTargetCount,
		BaqatiTarget:  regionReport.BaqatiTargetCount,
		UpgradeTarget: regionReport.UpgradeTargetCount,
	}

	eachRegionResp.TotalSales = regionReport.HBBSalesCount + regionReport.HayyakSalesCount + regionReport.BaqatiSalesCount + regionReport.UpgradeSalesCount
	eachRegionResp.TotalTarget = regionReport.HBBTargetCount + regionReport.HayyakTargetCount + regionReport.BaqatiTargetCount + regionReport.UpgradeTargetCount

	if eachRegionResp.HBBTarget != 0 {
		eachRegionResp.HBBAchievedPercentage = math.Round(float64(eachRegionResp.HBBSales*100/eachRegionResp.HBBTarget)*100) / 100
	}

	if eachRegionResp.HayyakTarget != 0 {
		eachRegionResp.HayyakAchievedPercentage = math.Round(float64(eachRegionResp.HayyakSales*100/eachRegionResp.HayyakTarget)*100) / 100
	}

	if eachRegionResp.BaqatiTarget != 0 {
		eachRegionResp.BaqatiAchievedPercentage = math.Round(float64(eachRegionResp.BaqatiSales*100/eachRegionResp.BaqatiTarget)*100) / 100
	}

	if eachRegionResp.UpgradeTarget != 0 {
		eachRegionResp.UpgradeAchievedPercentage = math.Round(float64(eachRegionResp.UpgradeSales*100/eachRegionResp.UpgradeTarget)*100) / 100
	}

	if eachRegionResp.TotalTarget != 0 {
		eachRegionResp.TotalAchievedPercentage = math.Round(float64(eachRegionResp.TotalSales*100/eachRegionResp.TotalTarget)*100) / 100
	}

	return eachRegionResp
}
