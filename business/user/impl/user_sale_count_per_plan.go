package impl

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) GetSaleCountPerPlan(ctx *gin.Context, userID, product string, isQuarterReport bool, quarter, month, year int) ([]*models.SaleCountPerPlan, *errors.Error) {

	// cRH, _ := ctx.Get(constant.HeaderInContext)
	// ctxHeader, _ := (cRH).(context.ContextHeader)

	_, ok := constant.ProductGroupMap[product]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameters,
			Message:    "INVALID PRODUCT - " + product,
			Details:    "invalid product - " + product,
		}
	}

	report, err := impl.CacheService.GetUserReport(userID, isQuarterReport, month, quarter, year)
	if err == nil && report != nil && objecthelper.IsCurrentDay(report.ModifiedAt) {
		if resp, err := impl.userPerPlanSale(ctx, product, report); err == nil {
			return resp, nil
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

		_, report, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		// val, _ := json.Marshal(resp)
		// logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Str("report", string(val)).Msg("GetSaleCountPerPlan : updated")

		return impl.userPerPlanSale(ctx, product, report)
	}

	/*
		For All Managers
	*/

	mngrReport, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}

	// val, _ := json.Marshal(mngrReport)
	// logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Str("report", string(val)).Msg("GetSingleUserIncentive : updated")

	return impl.userPerPlanSale(ctx, product, mngrReport)
}

func (impl *UserImpl) userPerPlanSale(ctx *gin.Context, product string, userReport *models.UserReportDB) ([]*models.SaleCountPerPlan, *errors.Error) {

	response := make([]*models.SaleCountPerPlan, 0)

	var saleCountString string

	switch product {
	case constant.HBBCategory:
		saleCountString = userReport.HBBSaleCountPerPlanArray
	case constant.HayyakCategory:
		saleCountString = userReport.HayyakSaleCountPerPlanArray
	case constant.BaqatiCategory:
		saleCountString = userReport.BaqatiSaleCountPerPlanArray
	case constant.UpgradeCategory:
		saleCountString = userReport.UpgradeSaleCountPerPlanArray
	}

	if saleCountString == "" {
		return response, nil
	}

	if err := json.Unmarshal([]byte(saleCountString), &response); err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "GetSaleCountPerPlan: resp marshal error - " + err.Error(),
		}
	}

	return response, nil
}
