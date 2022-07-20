package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) CreateNewOutlet(ctx *gin.Context, req *models.NewOutletRequest) *errors.Error {

	return impl.PsqlDb.UpsertOutletRecord(req)
}
