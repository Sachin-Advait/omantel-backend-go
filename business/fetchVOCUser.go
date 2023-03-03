package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) FetchVOCUser(ctx *gin.Context, req *models.FetchDeleteVOCUserRequest) ([]*models.VOCUserResponse, *errors.Error) {
	return impl.PsqlDb.FetchVOCUserRecord(req.UserIds, req.Year)
}
