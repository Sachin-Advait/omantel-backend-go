package business

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) GetSaleCountPerPlan(ctx *gin.Context, userID, product string, isQuarterReport bool, quarter, month, year int) ([]*models.SaleCountPerPlan, *errors.Error) {

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
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

	user, err := impl.PsqlDb.FetchUserDetails(userID, quarter, year)
	if err != nil {
		return nil, err
	}

	_, ok = constant.EmployeePositionMap[user.Position]
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

	//TODO: remove at last
	needUpdate = true

	response := make([]*models.SaleCountPerPlan, 0)

	if user.Position == constant.Promoter || user.Position == constant.CSR {

		userReport, err := impl.PsqlDb.FetchAgentUserReport(user.UserId, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil || (err == nil && userReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				ind, resp, err := impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}

				logger.FileLogger.CtxInfo(ctx).Str("agentReport", user.UserId).Str(constant.HBBCategory, ind.HBB.Error).Str(constant.HayyakCategory, ind.Hayyak.Error).
					Str(constant.BaqatiCategory, ind.Baqati.Error).Str(constant.UpgradeCategory, ind.Upgrade.Error).Msg("GetSaleCountPerPlan : trace")
				//TODO: test this out for struct assignment
				userReport = resp
			}
		}

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

	///////

	mngrReport, err := impl.PsqlDb.FetchManagerReport(user.UserId, isQuarterReport, quarter, month, year)
	if !needUpdate {
		if err != nil {
			return nil, err
		}
	} else {
		if err != nil || (err == nil && mngrReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

			resp, err := impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
			if err != nil {
				return nil, err
			}

			logger.FileLogger.CtxInfo(ctx).Interface("mngrReport", resp).Msg("GetSaleCountPerPlan : trace")
			//TODO: test this out for struct assignment
			mngrReport = resp
		}
	}

	var saleCountString string

	switch product {
	case constant.HBBCategory:
		saleCountString = mngrReport.HBBSaleCountPerPlanArray
	case constant.HayyakCategory:
		saleCountString = mngrReport.HayyakSaleCountPerPlanArray
	case constant.BaqatiCategory:
		saleCountString = mngrReport.BaqatiSaleCountPerPlanArray
	case constant.UpgradeCategory:
		saleCountString = mngrReport.UpgradeSaleCountPerPlanArray
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
