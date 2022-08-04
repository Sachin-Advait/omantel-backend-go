package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchRegions(ctx *gin.Context, req *models.DeleteRegionRequest) ([]*models.GetRegionResponse, *errors.Error) {

	resp, err := impl.PsqlDb.FetchRegionRecord(req.RegionCodes)
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
