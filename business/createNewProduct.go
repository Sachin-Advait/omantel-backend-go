package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) CreateNewProduct(ctx *gin.Context, req *models.NewProductRequest) *errors.Error {

	return impl.PsqlDb.UpsertProductRecord(req)
}
