package business

import (
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetTopTenEmpSalesAndTarget(ctx *gin.Context, clientID string, quarter int, isQuarterReport bool, t time.Time) ([]*models.ProductWiseSalesTargetPerStaff, *errors.Error) {

	position, err := impl.PsqlDb.FetchUserPosition(clientID)
	if err != nil {
		return nil, err
	}

	manger, ok := constant.ManagerMap[position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "incorrect position",
			Details:    "incorrect position",
		}
	}

	empList, err := impl.PsqlDb.FetchAllEmployeesReportingToGivenManagerRole(manger, clientID)
	if err != nil {
		return nil, err
	}

	resp := make([]*models.ProductWiseSalesTargetPerStaff, 0)

	for _, v := range empList {

		trx, err := impl.PsqlDb.FetchIndividualUserTargetReport(v.ID, int(t.Month()), quarter, isQuarterReport)
		if err != nil {
			return nil, err
		}

		name, _, err := impl.PsqlDb.FetchUserName(trx.EmpId)
		if err != nil {
			return nil, err
		}

		eachUser := &models.ProductWiseSalesTargetPerStaff{
			Name:                    name,
			HBBSales:                trx.Hbb_sales_count,
			HBBTarget:               math.Round(trx.Hbb_target*100) / 100,
			HBBAchievedPercent:      math.Round(trx.Hbb_achieved_percentage*100) / 100,
			HayyakSales:             trx.Hayyak_sales_count,
			HayyakTarget:            math.Round(trx.Hayyak_target*100) / 100,
			HayyakAchievedPercent:   math.Round(trx.Hayyak_achieved_percentage*100) / 100,
			BaqatiSales:             trx.Baqati_sales_count,
			BaqatiTarget:            math.Round(trx.Baqati_target*100) / 100,
			BaqatiAchievedPercent:   math.Round(trx.Baqati_achieved_percentage*100) / 100,
			UpgradeSales:            trx.Upgrade_sales_count,
			UpgradeTarget:           math.Round(trx.Upgrade_target*100) / 100,
			UpgradeAchievedPercent:  math.Round(trx.Uprade_achieved_percentage*100) / 100,
			TotalSales:              trx.Hbb_sales_count + trx.Hayyak_sales_count + trx.Baqati_sales_count + trx.Upgrade_sales_count,
			TotalTarget:             math.Round((trx.Hbb_target+trx.Hayyak_target+trx.Baqati_target+trx.Upgrade_target)*100) / 100,
			TotalAchievedPercentage: math.Round(((trx.Hbb_achieved_percentage+trx.Hayyak_achieved_percentage+trx.Baqati_achieved_percentage+trx.Uprade_achieved_percentage)/4)*100) / 100,
		}

		resp = append(resp, eachUser)
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].TotalAchievedPercentage > resp[j].TotalAchievedPercentage
	})

	if len(resp) > 10 {
		resp = resp[:10]
	}

	return resp, nil
}
