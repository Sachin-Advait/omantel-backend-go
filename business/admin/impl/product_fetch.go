package impl

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) FetchProducts(ctx *gin.Context, req *models.FetchDeleteProductRequest) ([]*models.ProductResponse, *errors.Error) {
	return impl.DbService.FetchProductRecord(req.ProductCodes, req.Quarter, req.Year)
}
