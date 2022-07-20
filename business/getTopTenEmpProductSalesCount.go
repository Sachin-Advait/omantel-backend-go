package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
	"net/http"
	"sort"
	"time"
)

func (impl *IncentiveImpl) GetTopTenEmpProductSalesCount(ctx *gin.Context, clientID, quarter int, isQuarterReport bool, t time.Time) ([]*psqldb.ProductWiseSalesCountPerStaff, *errors.Error) {

	position, err := impl.PsqlDb.FetchUserPosition(clientID)
	if err != nil {
		return nil, err
	}

	manger, ok := commons.ManagerMap[position]
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

	resp := make([]*psqldb.ProductWiseSalesCountPerStaff, 0)

	for _, v := range empList {

		trx, err := impl.PsqlDb.FetchIndividualUserTargetReport(v.ID, int(t.Month()), quarter, isQuarterReport)
		if err != nil {
			return nil, err
		}

		name, err := impl.PsqlDb.FetchUserName(trx.EmpId)
		if err != nil {
			return nil, err
		}

		eachUser := &psqldb.ProductWiseSalesCountPerStaff{
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
