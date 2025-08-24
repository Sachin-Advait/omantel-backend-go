package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetUserMNPSalesCount(ctx *gin.Context, userID string, chanel string, isQuarterReport bool,
	quarter, month, year int) (interface{}, *errors.Error) {

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

	user, err := impl.DbService.FetchUserDetails(userID, quarter, year)
	if err != nil {
		return nil, err
	}
	position, ok := constant.EmployeePositionMap[user.Position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT USER POSITION - " + user.Position,
			Details:    "incorrect user position: userId-" + user.UserId + " position-" + user.Position,
		}
	}
	user.Position = position

	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb || user.Position == constant.OutletManagerAbb {
		chanel = ""
	}

	report, err := impl.CacheService.GetUserMNPSaleCountPerPlanReport(userID, chanel, isQuarterReport, month, quarter, year)
	if err == nil && report != nil {
		modifiedDate := report.(*models.MNPSalesCount).ModifiedAt
		if objecthelper.IsCurrentDay(modifiedDate) {
			return report, nil
		}
	}

	/*
		For Agents
	*/

	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {
		report, err := impl.updateAgentUserMNPSalesCountReport(ctx, user, chanel, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}
		return report, nil
	}

	/*
		For All Managers
	*/

	mngrReport, err := impl.updateManagerUserMNPSalesCountReport(ctx, user, chanel, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}
	return mngrReport, nil
}
