package business

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) CreateNewTarget(ctx *gin.Context, req *models.NewTargetRequest) *errors.Error {

	req.TargetCode = req.RegionCode + req.OutletCode + fmt.Sprintf("%d", req.Month)

	return impl.PsqlDb.UpsertTargetRecord(req)
}
