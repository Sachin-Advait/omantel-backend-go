package impl

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) FetchVOCCriteria(ctx *gin.Context, req *models.FetchDeleteVOCCriteriaRequest) (*models.VOCCriteriaResponse, *errors.Error) {
	return impl.DbService.FetchVOCCriteriaRecord(req.Year)
}
