package business

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) GetSingleUserTargetChart(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*models.ChartDataResponse, *errors.Error) {

	trx, err := impl.PsqlDb.FetchUserSales(userId, month, year, quarter, isQuarterReport)
	if err != nil {
		return nil, err
	}

	resp := &models.ChartDataResponse{}

	if trx.Region == "" || trx.Outlet == "" {
		ind := &models.TargetAchievementData{Error: "no sales record found for this period"}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	totalEmp, err := impl.PsqlDb.FetchTotalNonManagerialEmployees(trx.Region, trx.Outlet)
	if err != nil {
		return nil, err
	} else if totalEmp == 0 {
		ind := &models.TargetAchievementData{Error: "no employee found for this outlet - " + trx.Outlet}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	totalEmp--

	logger.FileLogger.CtxInfo(ctx).Str("clientID", userId).Str("outlet", trx.Outlet).Str("region", trx.Region).
		Int("totalEmp", totalEmp).Msg("GetSingleUserTargetChart : sales info")

	for grp, rcds := range trx.SaleRecords {

		data := &models.TargetAchievementData{}

		for _, v := range rcds {
			data.SalesCount += v.Quantity
		}

		target := float64(impl.getSaleTargetPerOutlet(month, quarter, isQuarterReport, grp, trx.Region, trx.Outlet)) / float64(totalEmp)
		if target <= 0 {
			data.Error = fmt.Errorf("per employee target found is %f", target).Error()
		}

		data.Target = math.Round(target*100) / 100
		data.AchievedPercentage = math.Round((float64(data.SalesCount)*100/target)*100) / 100

		switch grp {
		case constant.HBBCategory:
			resp.HBB = data
		case constant.HayyakCategory:
			resp.Hayyak = data
		case constant.BaqatiCategory:
			resp.Baqati = data
		case constant.UpgradeCategory:
			resp.Upgrade = data
		default:
			return nil, &errors.Error{
				StatusCode: http.StatusBadRequest,
				Code:       errors.InvalidReqParameters,
				Message:    "no category found",
				Details:    "no category found"}
		}
	}

	return resp, nil
}
