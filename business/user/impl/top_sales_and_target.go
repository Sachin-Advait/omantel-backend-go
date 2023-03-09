package impl

import (
	"math"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetTopTenEmpSalesAndTarget(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesTargetPerStaff, *errors.Error) {

	// cRH, _ := ctx.Get(constant.HeaderInContext)
	// ctxHeader, _ := (cRH).(context.ContextHeader)

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
			Details:    "GetTopTenEmpSalesAndTarget: incorrect manager position: userId-" + user.UserId + " position-" + user.Position,
		}
	}
	user.Position = position

	resp := make([]*models.ProductWiseSalesTargetPerStaff, 0)

	/*
		For Outlet Managers
	*/

	if user.Position == constant.OutletManagerAbb {

		agentList, err := impl.DbService.FetchAllAgentsUserOfOM(userID, quarter, year)
		if err != nil {
			return nil, err
		}

		if len(agentList) == 0 {
			logger.FileLogger.CtxError(ctx).Str("position", user.Position).Msg("GetTopTenEmpSalesAndTarget : no Agent found")
			return resp, nil
		}

		for _, agent := range agentList {
			cacheReport, err := impl.CacheService.GetUserReport(agent.UserId, isQuarterReport, month, quarter, year)
			if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
				resp = append(resp, impl.userSaleTarget(ctx, cacheReport))
			} else {
				_, report, err := impl.updateAgentUserReport(ctx, agent, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}
				resp = append(resp, impl.userSaleTarget(ctx, report))
			}
		}

		sort.Slice(resp, func(i, j int) bool {
			return resp[i].TotalAchievedPercentage > resp[j].TotalAchievedPercentage
		})

		if len(resp) > 10 {
			resp = resp[:10]
		}

		return resp, nil
	}

	/*
		For Other Managers
	*/

	allDistinctMngr, err := impl.DbService.FetchAllMngrReportingToManagerRole(userID, user.Position, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(allDistinctMngr) == 0 {
		logger.FileLogger.CtxError(ctx).Str("position", user.Position).Msg("GetTopTenEmpSalesAndTarget : no Mngr found")
		return resp, nil
	}

	for _, eachMngr := range allDistinctMngr {
		cacheReport, err := impl.CacheService.GetUserReport(eachMngr.UserId, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			resp = append(resp, impl.userSaleTarget(ctx, cacheReport))
		} else {
			report, err := impl.updateManagerReport(ctx, eachMngr, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}
			resp = append(resp, impl.userSaleTarget(ctx, report))
		}
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].TotalAchievedPercentage > resp[j].TotalAchievedPercentage
	})

	if len(resp) > 10 {
		resp = resp[:10]
	}

	return resp, nil
}

func (impl *UserImpl) userSaleTarget(ctx *gin.Context, userReport *models.UserReportDB) *models.ProductWiseSalesTargetPerStaff {

	agentProductWiseRprt := &models.ProductWiseSalesTargetPerStaff{
		Name:                   userReport.UserName,
		HBBSales:               userReport.HBBSalesCount,
		HBBTarget:              userReport.HBBTargetCount,
		HBBAchievedPercent:     userReport.HBBAchievedPercentage,
		HayyakSales:            userReport.HayyakSalesCount,
		HayyakTarget:           userReport.HayyakTargetCount,
		HayyakAchievedPercent:  userReport.HayyakAchievedPercentage,
		BaqatiSales:            userReport.BaqatiSalesCount,
		BaqatiTarget:           userReport.BaqatiTargetCount,
		BaqatiAchievedPercent:  userReport.BaqatiAchievedPercentage,
		UpgradeSales:           userReport.UpgradeSalesCount,
		UpgradeTarget:          userReport.UpgradeTargetCount,
		UpgradeAchievedPercent: userReport.UpradeAchievedPercentage,
	}

	agentProductWiseRprt.TotalSales = agentProductWiseRprt.HBBSales + agentProductWiseRprt.HayyakSales + agentProductWiseRprt.BaqatiSales + agentProductWiseRprt.UpgradeSales
	agentProductWiseRprt.TotalTarget = agentProductWiseRprt.HBBTarget + agentProductWiseRprt.HayyakTarget + agentProductWiseRprt.BaqatiTarget + agentProductWiseRprt.UpgradeTarget

	if agentProductWiseRprt.TotalTarget > 0 {
		agentProductWiseRprt.TotalAchievedPercentage = math.Round(float64(agentProductWiseRprt.TotalSales*100/int(agentProductWiseRprt.TotalTarget))*100) / 100
	}

	return agentProductWiseRprt
}
