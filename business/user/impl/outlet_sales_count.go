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

func (impl *UserImpl) GetOutletTotalSalesCount(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesCountPerOutlet, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if err := impl.CacheService.IsSessionTokenValid(ctxHeader.XUserId, ctxHeader.Authorization[7:]); err != nil {
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

	// if ctxHeader.IsAdmin {
	// 	vps, err := impl.DbService.FetchVP(quarter, year)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	userID = vps[0]
	// }

	user, err := impl.DbService.FetchUserBasicDetails(userID, quarter, year)
	if err != nil {
		return nil, err
	}

	position, ok := constant.ManagerMap[user.Position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT MANAGER POSITION - " + user.Position,
			Details:    "GetOutletTotalSalesCount: incorrect manager position: userId-" + user.UserId + " position-" + user.Position,
		}
	}
	user.Position = position

	resp := make([]*models.ProductWiseSalesCountPerOutlet, 0)

	/*
		For Outlet Managers
	*/

	if user.Position == constant.OutletManagerAbb {

		cacheReport, err := impl.CacheService.GetOutletReport(user.Outlet, user.Region, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			return append(resp, impl.singleOutletSalesCount(ctx, cacheReport)), nil
		}

		outletResp, err := impl.updateOutletReport(ctx, user.Region, user.Outlet, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		// val, _ := json.Marshal(outletResp)
		// logger.FileLogger.CtxInfo(ctx).Interface("OM outletReport", string(val)).Msg("GetOutletTotalSalesCount : trace")

		return append(resp, impl.singleOutletSalesCount(ctx, outletResp)), nil
	}

	/*
		For Other Managers
	*/

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(user.Position, userID, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("user", userID).Str("position", user.Position).Msg("GetOutletTotalSalesCount : no OM found")
		return resp, nil
	}

	for _, eachOM := range allDistinctOM {

		if eachOM.Position != constant.OutletManagerAbb {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Message:    eachOM.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + eachOM.Position,
				Details:    "GetOutletTotalSalesCount: incorrect manager position: userId-" + eachOM.UserId + " position-" + eachOM.Position,
			}
		}

		cacheReport, err := impl.CacheService.GetOutletReport(eachOM.Outlet, eachOM.Region, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			resp = append(resp, impl.singleOutletSalesCount(ctx, cacheReport))
		} else {
			outletResp, err := impl.updateOutletReport(ctx, eachOM.Region, eachOM.Outlet, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			// val, _ := json.Marshal(outletResp)
			// logger.FileLogger.CtxInfo(ctx).Interface("OtherMngr outletReport", string(val)).Msg("GetOutletTotalSalesCount : trace")

			resp = append(resp, impl.singleOutletSalesCount(ctx, outletResp))
		}
	}

	return resp, nil
}

func (impl *UserImpl) singleOutletSalesCount(ctx *gin.Context, outletReport *models.OutletReportDB) *models.ProductWiseSalesCountPerOutlet {

	eachOutletResp := &models.ProductWiseSalesCountPerOutlet{
		OutletName: outletReport.Outlet,
		Region:     outletReport.Region,
		Latitude:   outletReport.Latitude,
		Longitude:  outletReport.Longitude,
		HBB:        outletReport.HBBSalesCount,
		Hayyak:     outletReport.HayyakSalesCount,
		Baqati:     outletReport.BaqatiSalesCount,
		Upgrade:    outletReport.UpgradeSalesCount,
	}

	eachOutletResp.Total = outletReport.HBBSalesCount + outletReport.HayyakSalesCount + outletReport.BaqatiSalesCount + outletReport.UpgradeSalesCount

	return eachOutletResp
}
