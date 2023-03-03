package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchVOCMaxPay(ctx *gin.Context, req *models.FetchDeleteVOCCriteriaRequest) ([]*models.VOCMaxPayResponse, *errors.Error) {
	return impl.PsqlDb.FetchVOCMaxPayRecord(req.Year)
}
