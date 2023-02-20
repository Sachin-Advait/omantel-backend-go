package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *IncentiveImpl) GetRegionTotalSalesCount(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesCountPerOutlet, *errors.Error) {

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
			Details:    "GetRegionTotalSalesCount: incorrect manager position: userId-" + user.UserId + " position-" + user.Position,
		}
	}
	user.Position = position

	resp := make([]*models.ProductWiseSalesCountPerOutlet, 0)

	/*
		For Outlet Managers
	*/

	if user.Position == constant.OutletManagerAbb {

		cacheReport, err := impl.RedisCache.GetRegionReport(user.Region, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			return append(resp, impl.singleRegionSalesCount(ctx, cacheReport)), nil
		}

		regionResp, err := impl.updateRegionReport(ctx, user.Region, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		// val, _ := json.Marshal(outletResp)
		// logger.FileLogger.CtxInfo(ctx).Interface("OM regionReport", string(val)).Msg("GetRegionTotalSalesCount : trace")

		return append(resp, impl.singleRegionSalesCount(ctx, regionResp)), nil
	}

	/*
		For Other Managers
	*/

	allDistinctOM, err := impl.PsqlDb.FetchAllOMReportingToGivenManager(user.Position, userID, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("user", userID).Str("position", user.Position).Msg("GetRegionTotalSalesCount : no OM found")
		return resp, nil
	}

	for _, eachOM := range allDistinctOM {

		if eachOM.Position != constant.OutletManagerAbb {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Message:    eachOM.UserId + ": INCORRECT MANAGER POSITION - " + eachOM.Position,
				Details:    "GetRegionTotalSalesCount: incorrect manager position: userId-" + eachOM.UserId + " position-" + eachOM.Position,
			}
		}

		cacheReport, err := impl.RedisCache.GetRegionReport(eachOM.Region, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			resp = append(resp, impl.singleRegionSalesCount(ctx, cacheReport))
		} else {
			outletResp, err := impl.updateRegionReport(ctx, eachOM.Region, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			// val, _ := json.Marshal(outletResp)
			// logger.FileLogger.CtxInfo(ctx).Interface("OtherMngr regionReport", string(val)).Msg("GetRegionTotalSalesCount : trace")

			resp = append(resp, impl.singleRegionSalesCount(ctx, outletResp))
		}
	}

	return resp, nil
}

func (impl *IncentiveImpl) singleRegionSalesCount(ctx *gin.Context, regionReport *models.RegionReportDB) *models.ProductWiseSalesCountPerOutlet {

	eachRegionResp := &models.ProductWiseSalesCountPerOutlet{
		Region:  regionReport.Region,
		HBB:     regionReport.HBBSalesCount,
		Hayyak:  regionReport.HayyakSalesCount,
		Baqati:  regionReport.BaqatiSalesCount,
		Upgrade: regionReport.UpgradeSalesCount,
	}

	eachRegionResp.Total = regionReport.HBBSalesCount + regionReport.HayyakSalesCount + regionReport.BaqatiSalesCount + regionReport.UpgradeSalesCount

	return eachRegionResp
}
