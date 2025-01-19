package impl

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetUserPlanCategory(ctx *gin.Context, userID, product string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

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

	// if product != constant.HBBCategory {
	// 	return nil, &errors.Error{
	// 		StatusCode: http.StatusBadRequest,
	// 		Code:       errors.InvalidReqParameters,
	// 		Message:    "INVALID PRODUCT - " + product,
	// 		Details:    "invalid product - " + product,
	// 	}
	// }

	report, err := impl.CacheService.GetUserPlanCategoryReport(userID, product, isQuarterReport, month, quarter, year)
	if err == nil && report != nil {
		planReport, ok := report.(map[string]interface{})
		if !ok {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    errors.InternalServerError,
				Details:    "getUserPlanCategoryReport - invalid cache data",
			}
		}
		modifiedDate, _ := planReport["ModifiedAt"].(time.Time)
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

		report, err := impl.updateAgentUserPlanPerCategoryReport(ctx, user, product, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}
		return report, nil
	}

	/*
		For All Managers
	*/

	mngrReport, err := impl.updateManagerUserPlanPerCategoryReport(ctx, user, product, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}
	return mngrReport, nil
}
