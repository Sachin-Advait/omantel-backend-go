package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetUserSalesCount(ctx *gin.Context, userID, product string, chanel string,
	isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

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

	report, hbbReport, err := impl.CacheService.GetUserSaleCountPerPlanReport(userID, product, chanel, isQuarterReport,
		month, quarter, year)

	if product == constant.HBBCategory && err == nil && hbbReport != nil && objecthelper.IsCurrentDay(hbbReport.ModifiedAt) {
		return hbbReport, nil
	} else if err == nil && report != nil && objecthelper.IsCurrentDay(report.ModifiedAt) {
		return report, nil
	}

	/*
		For Agents
	*/

	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {
		report, hbbReport, err := impl.updateAgentUserSalesCountReport(ctx, user, product, chanel, isQuarterReport,
			quarter, month, year)
		if err != nil {
			return nil, err
		}

		if product == constant.HBBCategory {
			return hbbReport, nil
		} else {
			return report, nil
		}
	}

	/*
		For All Managers
	*/

	mngrReport, hbbMngrReport, err := impl.updateManagerUserSalesCountReport(ctx, user, product, chanel, isQuarterReport,
		quarter, month, year)
	if err != nil {
		return nil, err
	}

	if product == constant.HBBCategory {
		return hbbMngrReport, nil
	} else {
		return mngrReport, nil
	}
}
