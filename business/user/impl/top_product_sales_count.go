package impl

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetTopTenEmpProductSalesCount(ctx *gin.Context, userID string, chanel string,
	isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesCountPerStaff, *errors.Error) {

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
			Details: "GetTopTenEmpProductSalesCount: incorrect manager position: userId-" +
				user.UserId + " position-" + user.Position,
		}
	}
	user.Position = position

	if user.Position == constant.OutletManagerAbb {
		chanel = ""
	}

	resp := make([]*models.ProductWiseSalesCountPerStaff, 0)

	/*
		For Outlet Managers
	*/

	if user.Position == constant.OutletManagerAbb {

		agentList, err := impl.DbService.FetchAllAgentsUserOfOM(userID, quarter, year)
		if err != nil {
			return nil, err
		}

		if len(agentList) == 0 {
			logger.FileLogger.CtxError(ctx).Str("position", user.Position).Msg("GetTopTenEmpProductSalesCount : no Agent found")
			return resp, nil
		}

		for _, agent := range agentList {
			cacheReport, err := impl.CacheService.GetUserReport(agent.UserId, chanel, isQuarterReport, month, quarter, year)
			if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
				resp = append(resp, impl.userSaleCount(ctx, cacheReport))
			} else {
				_, report, err := impl.updateAgentUserReport(ctx, agent, chanel, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}
				resp = append(resp, impl.userSaleCount(ctx, report))
			}
		}

		sort.Slice(resp, func(i, j int) bool {
			return resp[i].Total > resp[j].Total
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
		logger.FileLogger.CtxError(ctx).Str("position", user.Position).Msg("GetTopTenEmpProductSalesCount : no Mngr found")
		return resp, nil
	}

	for _, eachMngr := range allDistinctMngr {
		cacheReport, err := impl.CacheService.GetUserReport(eachMngr.UserId, chanel, isQuarterReport, month, quarter, year)
		if err == nil && cacheReport != nil && objecthelper.IsCurrentDay(cacheReport.ModifiedAt) {
			resp = append(resp, impl.userSaleCount(ctx, cacheReport))
		} else {
			report, err := impl.updateManagerReport(ctx, eachMngr, chanel, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}
			resp = append(resp, impl.userSaleCount(ctx, report))
		}
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Total > resp[j].Total
	})

	if len(resp) > 10 {
		resp = resp[:10]
	}

	return resp, nil
}

func (impl *UserImpl) userSaleCount(ctx *gin.Context, userReport *models.UserReportDB) *models.ProductWiseSalesCountPerStaff {

	agentProductWiseRprt := &models.ProductWiseSalesCountPerStaff{
		Name:    userReport.UserName,
		HBB:     userReport.HBBSalesCount,
		Hayyak:  userReport.HayyakSalesCount,
		Baqati:  userReport.BaqatiSalesCount,
		Upgrade: userReport.UpgradeSalesCount,
	}

	agentProductWiseRprt.Total = agentProductWiseRprt.HBB + agentProductWiseRprt.Hayyak + agentProductWiseRprt.Baqati + agentProductWiseRprt.Upgrade

	return agentProductWiseRprt
}
