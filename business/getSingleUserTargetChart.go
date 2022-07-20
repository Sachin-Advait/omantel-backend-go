package business

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
	"math"
	"net/http"
	"time"
)

func (impl *IncentiveImpl) GetSingleUserTargetChart(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) (*psqldb.ChartDataResponse, *errors.Error) {

	trx, err := impl.PsqlDb.FetchUserSales(clientID, int(t.Month()), t.Year(), quarter, isQuarterReport)
	if err != nil {
		return nil, err
	}

	resp := &psqldb.ChartDataResponse{}

	if trx.Region == "" || trx.Outlet == "" {
		ind := &psqldb.TargetAchievementData{Error: "no sales record found for this period"}
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
		ind := &psqldb.TargetAchievementData{Error: "no employee found for this outlet - " + trx.Outlet}
		resp.HBB = ind
		resp.Baqati = ind
		resp.Hayyak = ind
		resp.Upgrade = ind
		return resp, nil
	}

	totalEmp--

	logger.FileLogger.CtxInfo(ctx).Int("clientID", clientID).Str("outlet", trx.Outlet).Str("region", trx.Region).
		Int("totalEmp", totalEmp).Msg("GetSingleUserTargetChart : sales info")

	for grp, rcds := range trx.SaleRecords {

		data := &psqldb.TargetAchievementData{}

		for _, v := range rcds {
			data.SalesCount += v.Quantity
		}

		target := float64(impl.getSaleTargetPerOutlet(t, quarter, isQuarterReport, grp, trx.Region, trx.Outlet)) / float64(totalEmp)
		if target <= 0 {
			data.Error = fmt.Errorf("per employee target found is %f", target).Error()
		}

		data.Target = math.Round(target*100) / 100
		data.AchievedPercentage = math.Round((float64(data.SalesCount)*100/target)*100) / 100

		switch grp {
		case commons.HBBCategory:
			resp.HBB = data
		case commons.HayyakCategory:
			resp.Hayyak = data
		case commons.BaqatiCategory:
			resp.Baqati = data
		case commons.UpgradeCategory:
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
