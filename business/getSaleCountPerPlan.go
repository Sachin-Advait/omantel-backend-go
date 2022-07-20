package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
	"sort"
	"time"
)

func (impl *IncentiveImpl) GetSaleCountPerPlan(ctx *gin.Context, clientID, quarter int, product string, isQuarterReport bool, t time.Time) ([]*psqldb.SaleCountPerPlan, *errors.Error) {

	trx, err := impl.PsqlDb.FetchUserSales(clientID, int(t.Month()), t.Year(), quarter, isQuarterReport)
	if err != nil {
		return nil, err
	}

	m := make(map[string]int)

	for grp, rcds := range trx.SaleRecords {

		if product != commons.UpgradeCategory && grp != product {
			continue
		}

		for _, each := range rcds {

			v, ok := m[each.CurrentProduct]
			if !ok {
				m[each.CurrentProduct] = each.Quantity
			}

			if grp == commons.UpgradeCategory && each.PreviousProduct != "" {
				m[each.CurrentProduct] = v + each.Quantity
			}

			if grp != commons.UpgradeCategory && each.PreviousProduct == "" {
				m[each.CurrentProduct] = v + each.Quantity
			}
		}
	}

	resp := make([]*psqldb.SaleCountPerPlan, 0)

	for k, v := range m {
		resp = append(resp, &psqldb.SaleCountPerPlan{
			CurrentPlan: k,
			Quantity:    v,
		})
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Quantity > resp[j].Quantity
	})

	return resp, nil
}
