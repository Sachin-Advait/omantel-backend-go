package business

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetTopTenEmpProductSalesCount(ctx *gin.Context, clientID string, quarter int, isQuarterReport bool, t time.Time) ([]*models.ProductWiseSalesCountPerStaff, *errors.Error) {

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

	resp := make([]*models.ProductWiseSalesCountPerStaff, 0)

	for _, v := range empList {

		trx, err := impl.PsqlDb.FetchIndividualUserTargetReport(v.ID, int(t.Month()), quarter, isQuarterReport)
		if err != nil {
			return nil, err
		}

		name, err := impl.PsqlDb.FetchUserName(trx.EmpId)
		if err != nil {
			return nil, err
		}

		eachUser := &models.ProductWiseSalesCountPerStaff{
			Name:    name,
			HBB:     trx.Hbb_sales_count,
			Hayyak:  trx.Hayyak_sales_count,
			Baqati:  trx.Baqati_sales_count,
			Upgrade: trx.Upgrade_sales_count,
			Total:   trx.Hbb_sales_count + trx.Hayyak_sales_count + trx.Baqati_sales_count + trx.Upgrade_sales_count,
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
