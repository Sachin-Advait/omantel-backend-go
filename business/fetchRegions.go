package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchRegions(ctx *gin.Context, req *models.RegionRequest) ([]*models.RegionResponse, *errors.Error) {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if req.Quarter == 0 {
		req.Quarter = currQuarter
	}

	if req.Year == 0 {
		req.Year = currYear
	}

	resp, err := impl.PsqlDb.FetchRegionRecord(req.RegionCodes, req.Quarter, req.Year)
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
