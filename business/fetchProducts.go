package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchProducts(ctx *gin.Context, req *models.FetchDeleteProductRequest) ([]*models.ProductResponse, *errors.Error) {
	return impl.PsqlDb.FetchProductRecord(req.ProductCodes, req.Quarter, req.Year)
}
