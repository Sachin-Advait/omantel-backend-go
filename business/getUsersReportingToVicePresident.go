package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToVicePresident(ctx *gin.Context, vicePresidentID int) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	return impl.PsqlDb.FetchAllEmployeesReportingToGivenManagerRole("vice_president", vicePresidentID)
}
