package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToGeneralManager(ctx *gin.Context, generalManagerID string) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	return impl.PsqlDb.FetchAllEmpReportingToGivenManagerRole(constant.GeneralManager, generalManagerID)
}
