package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchVOCCriteria(ctx *gin.Context, req *models.FetchDeleteVOCCriteriaRequest) (*models.VOCCriteriaResponse, *errors.Error) {
	return impl.PsqlDb.FetchVOCCriteriaRecord(req.Year)
}
