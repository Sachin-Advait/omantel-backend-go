package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchTargets(ctx *gin.Context, req *models.GetTargetRequest) (map[string][]*models.GetTargetResponse, *errors.Error) {

	resp, err := impl.PsqlDb.FetchTargetRecord(req.TargetCode, req.OutletCode, req.Year, req.Month)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, &errors.Error{
			StatusCode: http.StatusOK,
			Code:       errors.DataNotFoundCode,
			Message:    "NO RECORDS FOUND",
			Details:    "NO RECORDS FOUND",
		}
	}

	return resp, nil
}
