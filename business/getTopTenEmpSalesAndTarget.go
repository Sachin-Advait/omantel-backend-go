package business

import (
	"math"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetTopTenEmpSalesAndTarget(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesTargetPerStaff, *errors.Error) {

	position, err := impl.PsqlDb.FetchUserPosition(userId)
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

	empList, err := impl.PsqlDb.FetchAllEmployeesReportingToGivenManagerRole(manger, userId)
	if err != nil {
		return nil, err
	}

	resp := make([]*models.ProductWiseSalesTargetPerStaff, 0)

	for _, v := range empList {

		mngrReport, err := impl.PsqlDb.FetchManagerReport(v.ID, v.Position, isQuarterReport, quarter, month, year)
		if err != nil {
			mngrReport, err = impl.updateManagerReport(ctx, userId, isQuarterReport, quarter, month, year)
		}

		if err != nil || mngrReport == nil {
			return nil, err
		}

		eachUser := &models.ProductWiseSalesTargetPerStaff{
			Name:          v.Name,
			HBBSales:      mngrReport.HBBSalesCount,
			HBBTarget:     float64(mngrReport.HBBTargetCount),
			HayyakSales:   mngrReport.HayyakSalesCount,
			HayyakTarget:  float64(mngrReport.HayyakTargetCount),
			BaqatiSales:   mngrReport.BaqatiSalesCount,
			BaqatiTarget:  float64(mngrReport.BaqatiTargetCount),
			UpgradeSales:  mngrReport.UpgradeSalesCount,
			UpgradeTarget: float64(mngrReport.UpgradeTargetCount),
			TotalSales:    mngrReport.HBBSalesCount + mngrReport.HayyakSalesCount + mngrReport.BaqatiSalesCount + mngrReport.UpgradeSalesCount,
			TotalTarget:   float64(mngrReport.HBBTargetCount + mngrReport.HayyakTargetCount + mngrReport.BaqatiTargetCount + mngrReport.UpgradeTargetCount),
		}

		if mngrReport.HBBTargetCount != 0 {
			eachUser.HBBAchievedPercent = math.Round(float64(mngrReport.HBBSalesCount)*100/float64(mngrReport.HBBTargetCount)) / 100
		}

		if mngrReport.HayyakTargetCount != 0 {
			eachUser.HayyakAchievedPercent = math.Round(float64(mngrReport.HayyakSalesCount)*100/float64(mngrReport.HayyakTargetCount)) / 100
		}

		if mngrReport.BaqatiTargetCount != 0 {
			eachUser.BaqatiAchievedPercent = math.Round(float64(mngrReport.BaqatiSalesCount)*100/float64(mngrReport.BaqatiTargetCount)) / 100
		}

		if mngrReport.UpgradeTargetCount != 0 {
			eachUser.UpgradeAchievedPercent = math.Round(float64(mngrReport.UpgradeSalesCount)*100/float64(mngrReport.UpgradeTargetCount)) / 100
		}

		if eachUser.TotalTarget != 0 {
			eachUser.TotalAchievedPercentage = math.Round(float64(eachUser.TotalSales)*100/float64(eachUser.TotalTarget)) / 100
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
