package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *IncentiveImpl) GetUserReport(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

	// cRH, _ := ctx.Get(constant.HeaderInContext)
	// ctxHeader, _ := (cRH).(context.ContextHeader)

	report, err := impl.RedisCache.GetUserReport(userID, isQuarterReport, month, quarter, year)
	if err == nil && report != nil && objecthelper.IsCurrentDay(report.ModifiedAt) {
		return report, nil
	}

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
	}

	// if ctxHeader.IsAdmin {
	// 	vps, err := impl.PsqlDb.FetchVP(quarter, year)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	userID = vps[0]
	// }

	user, err := impl.PsqlDb.FetchUserDetails(userID, quarter, year)
	if err != nil {
		return nil, err
	}

	_, ok := constant.EmployeePositionMap[user.Position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT USER POSITION - " + user.Position,
			Details:    "incorrect user position: userId-" + user.UserId + " position-" + user.Position,
		}
	}

	/*
		For Agents
	*/

	if user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb {

		_, report, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		// val, _ := json.Marshal(resp)
		// logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Str("report", string(val)).Msg("GetUserReport : updated")

		return report, nil
	}

	/*
		For All Managers
	*/

	mngrReport, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}

	// val, _ := json.Marshal(mngrReport)
	// logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Str("report", string(val)).Msg("GetUserReport : updated")

	return mngrReport, nil
}
