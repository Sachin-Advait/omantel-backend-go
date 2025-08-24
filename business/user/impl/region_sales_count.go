package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetRegionTotalSalesCount(ctx *gin.Context, userID string, chanel string,
	isQuarterReport bool, quarter, month, year int) (map[string]*models.ProductWiseSalesCountPerOutlet, *errors.Error) {

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
			Details:    "GetRegionTotalSalesCount: incorrect manager position: userId-" + user.UserId + " position-" + user.Position,
		}
	}
	user.Position = position

	if user.Position == constant.OutletManagerAbb {
		chanel = ""
	}

	resp := make(map[string]*models.ProductWiseSalesCountPerOutlet, 0)

	/*
		For Outlet Managers
	*/

	if user.Position == constant.OutletManagerAbb {
		cacheReport, err := impl.CacheService.GetRegionReport(user.Region, chanel, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			rcrd := impl.singleRegionSalesCount(ctx, cacheReport)
			resp[rcrd.Region] = rcrd
			return resp, nil
		}
		regionResp, err := impl.updateRegionReport(ctx, user.Region, chanel, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}
		rcrd := impl.singleRegionSalesCount(ctx, regionResp)
		resp[rcrd.Region] = rcrd
		return resp, nil
	}

	/*
		For Other Managers
	*/

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(user.Position, userID, quarter, year)
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

		cacheReport, err := impl.CacheService.GetRegionReport(eachOM.Region, chanel, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			rcrd := impl.singleRegionSalesCount(ctx, cacheReport)
			resp[rcrd.Region] = rcrd
		} else {
			regionResp, err := impl.updateRegionReport(ctx, eachOM.Region, chanel, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			// val, _ := json.Marshal(outletResp)
			// logger.FileLogger.CtxInfo(ctx).Interface("OtherMngr regionReport", string(val)).Msg("GetRegionTotalSalesCount : trace")

			rcrd := impl.singleRegionSalesCount(ctx, regionResp)
			resp[rcrd.Region] = rcrd
		}
	}

	return resp, nil
}

func (impl *UserImpl) singleRegionSalesCount(ctx *gin.Context, regionReport *models.RegionReportDB) *models.ProductWiseSalesCountPerOutlet {

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
