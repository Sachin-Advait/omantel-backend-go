package impl

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetRegionSalesAndTarget(ctx *gin.Context, userID string, chanel string,
	isQuarterReport bool, quarter, month, year int) (map[string]*models.ProductWiseSalesTargetPerOutlet, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if err := impl.CacheService.IsSessionTokenValid(ctxHeader.XUserId, ctxHeader.Authorization); err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.UnauthorizedUserCode,
			Message:    errors.UnauthorizedUser,
			Details:    err.Details,
		}
	}

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	user, err := impl.DbService.FetchUserBasicDetails(userID, quarter, year)
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

	if user.Position == constant.OutletManagerAbb {
		chanel = ""
	}

	resp := make(map[string]*models.ProductWiseSalesTargetPerOutlet, 0)

	/*
		For Outlet Managers
	*/

	if user.Position == constant.OutletManagerAbb {
		cacheReport, err := impl.CacheService.GetRegionReport(user.Region, chanel, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			rcrd := impl.singleRegionSalesAndTarget(ctx, cacheReport)
			resp[rcrd.Region] = rcrd
			return resp, nil
		}
		regionResp, err := impl.updateRegionReport(ctx, user.Region, chanel, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}
		rcrd := impl.singleRegionSalesAndTarget(ctx, regionResp)
		resp[rcrd.Region] = rcrd
		return resp, nil
	}

	/*
		For Other Managers
	*/

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(user.Position, userID, chanel, quarter, year)
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
		cacheReport, err := impl.CacheService.GetRegionReport(eachOM.Region, chanel, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			rcrd := impl.singleRegionSalesAndTarget(ctx, cacheReport)
			resp[rcrd.Region] = rcrd
		} else {
			regionResp, err := impl.updateRegionReport(ctx, eachOM.Region, chanel, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}
			rcrd := impl.singleRegionSalesAndTarget(ctx, regionResp)
			resp[rcrd.Region] = rcrd
		}
	}

	return resp, nil
}

func (impl *UserImpl) singleRegionSalesAndTarget(ctx *gin.Context, regionReport *models.RegionReportDB) *models.ProductWiseSalesTargetPerOutlet {

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
