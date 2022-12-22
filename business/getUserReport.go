package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) GetUserReport(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) (*models.UserReportDB, *errors.Error) {

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

	var needUpdate bool

	currTime := time.Now()
	currYear, currMonth, currDay := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if isQuarterReport && year == currYear && quarter <= currQuarter {
		needUpdate = true
	} else if !isQuarterReport && year == currYear && month <= int(currMonth) {
		needUpdate = true
	}

	if user.Position == constant.Promoter || user.Position == constant.CSR {

		userReport, err := impl.PsqlDb.FetchAgentUserReport(user.UserId, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}

			logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetUserReport : precalculated")

			userReport.Position = constant.EmployeePositionMap[userReport.Position]
			userReport.ManagerPosition = constant.EmployeePositionMap[userReport.ManagerPosition]

			return userReport, nil
		}

		if err != nil || (err == nil && userReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

			_, resp, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetUserReport : updated")

			resp.Position = constant.EmployeePositionMap[resp.Position]
			resp.ManagerPosition = constant.EmployeePositionMap[resp.ManagerPosition]

			return resp, nil
		}

		logger.FileLogger.CtxInfo(ctx).Str("agent", user.UserId).Msg("GetUserReport : trace")

		userReport.Position = constant.EmployeePositionMap[userReport.Position]
		userReport.ManagerPosition = constant.EmployeePositionMap[userReport.ManagerPosition]

		return userReport, nil
	}

	/////

	mngrReport, err := impl.PsqlDb.FetchManagerReport(user.UserId, isQuarterReport, quarter, month, year)
	if !needUpdate {
		if err != nil {
			return nil, err
		}

		logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetUserReport : precalculated")

		mngrReport.Position = constant.EmployeePositionMap[mngrReport.Position]
		mngrReport.ManagerPosition = constant.EmployeePositionMap[mngrReport.ManagerPosition]

		return mngrReport, nil
	}

	if err != nil || (err == nil && mngrReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

		resp, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetUserReport : updated")

		resp.Position = constant.EmployeePositionMap[resp.Position]
		resp.ManagerPosition = constant.EmployeePositionMap[resp.ManagerPosition]

		return resp, nil
	}

	logger.FileLogger.CtxInfo(ctx).Str("manger", user.UserId).Msg("GetUserReport : trace")

	mngrReport.Position = constant.EmployeePositionMap[mngrReport.Position]
	mngrReport.ManagerPosition = constant.EmployeePositionMap[mngrReport.ManagerPosition]

	return mngrReport, nil
}
