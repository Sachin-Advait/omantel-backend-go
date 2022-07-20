package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToSeniorManager(ctx *gin.Context, seniorManagerID int) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	return impl.PsqlDb.FetchAllEmployeesReportingToGivenManagerRole("senior_manager", seniorManagerID)
}
