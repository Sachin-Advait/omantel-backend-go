package business

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetTopTenEmpProductSalesCount(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesCountPerStaff, *errors.Error) {

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

	resp := make([]*models.ProductWiseSalesCountPerStaff, 0)

	for _, v := range empList {

		mngrReport, err := impl.PsqlDb.FetchManagerReport(v.ID, v.Position, isQuarterReport, quarter, month, year)
		if err != nil {
			mngrReport, err = impl.updateManagerReport(ctx, userId, isQuarterReport, quarter, month, year)
		}

		if err != nil || mngrReport == nil {
			return nil, err
		}

		eachUser := &models.ProductWiseSalesCountPerStaff{
			Name:    v.Name,
			HBB:     mngrReport.HBBSalesCount,
			Hayyak:  mngrReport.HayyakSalesCount,
			Baqati:  mngrReport.BaqatiSalesCount,
			Upgrade: mngrReport.UpgradeSalesCount,
			Total:   mngrReport.HBBSalesCount + mngrReport.HayyakSalesCount + mngrReport.BaqatiSalesCount + mngrReport.UpgradeSalesCount,
		}

		resp = append(resp, eachUser)
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Total > resp[j].Total
	})

	if len(resp) > 10 {
		resp = resp[:10]
	}

	return resp, nil
}
