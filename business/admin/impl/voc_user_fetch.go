package impl

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) FetchVOCUser(ctx *gin.Context, req *models.FetchDeleteVOCUserRequest) ([]*models.VOCUserResponse, *errors.Error) {
	return impl.DbService.FetchVOCUserRecord(req.UserIds, req.Year)
}
