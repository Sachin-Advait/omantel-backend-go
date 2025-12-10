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

// GetOutletSalesAndTarget: userID should be a manager
func (impl *UserImpl) GetOutletSalesAndTarget(ctx *gin.Context, userID string, chanel string,
	isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesTargetPerOutlet, *errors.Error) {

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
			Details:    "GetOutletSalesAndTarget: incorrect manager position: userId-" + user.UserId + " position-" + user.Position,
		}
	}
	user.Position = position

	if user.Position == constant.OutletManagerAbb {
		chanel = ""
	}

	resp := make([]*models.ProductWiseSalesTargetPerOutlet, 0)

	/*
		For Outlet Managers
	*/

	if user.Position == constant.OutletManagerAbb {
		cacheReport, err := impl.CacheService.GetOutletReport(user.Outlet, user.Region, chanel, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			return append(resp, impl.singleOutletSalesAndTarget(ctx, cacheReport)), nil
		}
		outletResp, err := impl.updateOutletReport(ctx, user.Region, user.Outlet, chanel, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}
		return append(resp, impl.singleOutletSalesAndTarget(ctx, outletResp)), nil
	}

	/*
		For Other Managers
	*/

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(user.Position, userID, chanel, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("user", userID).Str("position", user.Position).Msg("GetOutletSalesAndTarget : no OM found")
		return resp, nil
	}

	for _, eachOM := range allDistinctOM {
		if eachOM.Position != constant.OutletManagerAbb {
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Message:    eachOM.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + eachOM.Position,
				Details:    "GetOutletSalesAndTarget: incorrect manager position: userId-" + eachOM.UserId + " position-" + eachOM.Position,
			}
		}

		cacheReport, err := impl.CacheService.GetOutletReport(eachOM.Outlet, eachOM.Region, chanel, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			resp = append(resp, impl.singleOutletSalesAndTarget(ctx, cacheReport))
		} else {
			outletResp, err := impl.updateOutletReport(ctx, eachOM.Region, eachOM.Outlet, chanel, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			// val, _ := json.Marshal(outletResp)
			// logger.FileLogger.CtxInfo(ctx).Interface("OtherMngr outletReport", string(val)).Msg("GetOutletSalesAndTarget : trace")

			resp = append(resp, impl.singleOutletSalesAndTarget(ctx, outletResp))
		}
	}

	return resp, nil
}

func (impl *UserImpl) singleOutletSalesAndTarget(ctx *gin.Context,
	outletReport *models.OutletReportDB) *models.ProductWiseSalesTargetPerOutlet {

	eachOutletResp := &models.ProductWiseSalesTargetPerOutlet{
		OutletName:    outletReport.Outlet,
		Region:        outletReport.Region,
		HBBSales:      outletReport.HBBSalesCount,
		HayyakSales:   outletReport.HayyakSalesCount,
		BaqatiSales:   outletReport.BaqatiSalesCount,
		UpgradeSales:  outletReport.UpgradeSalesCount,
		HBBTarget:     outletReport.HBBTargetCount,
		HayyakTarget:  outletReport.HayyakTargetCount,
		BaqatiTarget:  outletReport.BaqatiTargetCount,
		UpgradeTarget: outletReport.UpgradeTargetCount,
	}

	eachOutletResp.TotalSales = outletReport.HBBSalesCount + outletReport.HayyakSalesCount + outletReport.BaqatiSalesCount + outletReport.UpgradeSalesCount
	eachOutletResp.TotalTarget = outletReport.HBBTargetCount + outletReport.HayyakTargetCount + outletReport.BaqatiTargetCount + outletReport.UpgradeTargetCount

	if eachOutletResp.HBBTarget != 0 {
		eachOutletResp.HBBAchievedPercentage = math.Round(float64(eachOutletResp.HBBSales*100/eachOutletResp.HBBTarget)*100) / 100
	}

	if eachOutletResp.HayyakTarget != 0 {
		eachOutletResp.HayyakAchievedPercentage = math.Round(float64(eachOutletResp.HayyakSales*100/eachOutletResp.HayyakTarget)*100) / 100
	}

	if eachOutletResp.BaqatiTarget != 0 {
		eachOutletResp.BaqatiAchievedPercentage = math.Round(float64(eachOutletResp.BaqatiSales*100/eachOutletResp.BaqatiTarget)*100) / 100
	}

	if eachOutletResp.UpgradeTarget != 0 {
		eachOutletResp.UpgradeAchievedPercentage = math.Round(float64(eachOutletResp.UpgradeSales*100/eachOutletResp.UpgradeTarget)*100) / 100
	}

	if eachOutletResp.TotalTarget != 0 {
		eachOutletResp.TotalAchievedPercentage = math.Round(float64(eachOutletResp.TotalSales*100/eachOutletResp.TotalTarget)*100) / 100
	}

	return eachOutletResp
}
