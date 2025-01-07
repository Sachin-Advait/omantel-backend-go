package impl

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetSaleCountPerPlan(ctx *gin.Context, userID, product string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

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

	_, ok := constant.ProductGroupMap[product]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameters,
			Message:    "INVALID PRODUCT - " + product,
			Details:    "invalid product - " + product,
		}
	}

	report, err := impl.CacheService.GetUserSaleCountPerPlanReport(userID, product, isQuarterReport, month, quarter, year)
	if err == nil && report != nil {
		var modifiedDate time.Time
		if product == constant.HBBCategory {
			modifiedDate = report.(*models.HbbUserSalesCountPerPlanDB).ModifiedAt
		} else {
			modifiedDate = report.(*models.UserSalesCountPerPlanDB).ModifiedAt
		}

		if objecthelper.IsCurrentDay(modifiedDate) {
			return report, nil
		}
	}

	// Compute real time

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

	/*
		For Agents
	*/

	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {

		report, err := impl.updateAgentUserSaleCountPerPlanReport(ctx, user, product, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}
		return report, nil
	}

	/*
		For All Managers
	*/

	mngrReport, err := impl.updateManagerUserSaleCountPerPlanReport(ctx, user, product, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}
	return mngrReport, nil
}
