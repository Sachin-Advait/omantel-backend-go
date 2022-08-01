package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) DeleteProducts(ctx *gin.Context, req *models.DeleteProductRequest) *errors.Error {

	return impl.PsqlDb.DeleteProductRecord(req.ProductCodes)
}
