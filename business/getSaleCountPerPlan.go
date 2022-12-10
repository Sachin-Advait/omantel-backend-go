package business

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetSaleCountPerPlan(ctx *gin.Context, userID, product string, isQuarterReport bool, quarter, month, year int) ([]*models.SaleCountPerPlan, *errors.Error) {

	_, ok := constant.ProductGroupMap[product]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidReqParameters,
			Message:    "INVALID PRODUCT",
			Details:    "Invalid product",
		}
	}

	user, err := impl.PsqlDb.FetchUserDetails(userID)
	if err != nil {
		return nil, err
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

	var userReport *models.UserReportDB

	response := make([]*models.SaleCountPerPlan, 0)

	if user.Position == constant.Promoter || user.Position == constant.CSR {

		userReport, err = impl.PsqlDb.FetchAgentUserReport(user.UserId, isQuarterReport, quarter, month, year)
		if needUpdate && (err != nil || (err == nil && userReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())))) {
			_, userReport, err = impl.updateAgentUserReport(ctx, user, isQuarterReport, quarter, month, year)
		}

		if err != nil {
			return nil, err
		}
	} else {

		userReport, err = impl.PsqlDb.FetchManagerReport(userID, isQuarterReport, quarter, month, year)
		if needUpdate && (err != nil || (err == nil && userReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())))) {
			userReport, err = impl.updateManagerReport(ctx, user, isQuarterReport, quarter, month, year)
		}

		if err != nil {
			return nil, err
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
