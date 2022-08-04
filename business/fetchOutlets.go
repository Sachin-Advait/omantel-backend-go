package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchOutlets(ctx *gin.Context, req *models.DeleteOutletRequest) ([]*models.GetOutletResponse, *errors.Error) {

	resp, err := impl.PsqlDb.FetchOutletRecord(req.OutletCodes)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, &errors.Error{
			StatusCode: http.StatusOK,
			Code:       errors.DataNotFoundCode,
			Message:    "NO RECORDS FOUND",
			Details:    "NO RECORDS FOUND",
		}
	}

	return resp, nil
}
