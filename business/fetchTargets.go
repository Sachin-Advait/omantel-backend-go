package business

import (
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchTargets(ctx *gin.Context, req *models.FetchDeleteTargetRequest) (map[string][]*models.TargetResponse, *errors.Error) {

	var (
		err  *errors.Error
		resp map[string][]*models.TargetResponse
	)

	switch {
	case len(req.TargetCodes) > 0: //targetCodes of one year
		resp, err = impl.PsqlDb.FetchTargetRecord(req.TargetCodes, nil, 0, req.Quarter, req.Year)
	case len(req.OutletCode) > 0: //outletCodes of one quarter and of one year
		resp, err = impl.PsqlDb.FetchTargetRecord(nil, req.OutletCode, 0, req.Quarter, req.Year)
	case req.Month != 0: //all records of one month and of one year
		resp, err = impl.PsqlDb.FetchTargetRecord(nil, nil, req.Month, req.Quarter, req.Year)
	case req.Quarter != 0: //all records of one quarter and of one year
		resp, err = impl.PsqlDb.FetchTargetRecord(nil, nil, 0, req.Quarter, req.Year)
	default: //all records of one year
		resp, err = impl.PsqlDb.FetchTargetRecord(nil, nil, 0, 0, req.Year)
	}

	if err != nil {
		return nil, err
	}

	for k, v := range resp {

		sort.Slice(v, func(i, j int) bool {
			return v[i].Month > v[j].Month
		})

		resp[k] = v
	}

	return resp, nil
}
