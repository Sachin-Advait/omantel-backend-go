package impl

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) FetchVOCMaxPay(ctx *gin.Context, req *models.FetchDeleteVOCCriteriaRequest) ([]*models.VOCMaxPayResponse, *errors.Error) {
	return impl.DbService.FetchVOCMaxPayRecord(req.Year)
}
