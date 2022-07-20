package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
)

func (impl *IncentiveImpl) GetUsersReportingToManager(ctx *gin.Context, managerID int) ([]*psqldb.ReportingStaffIDAndName, *errors.Error) {

	return impl.PsqlDb.FetchAllEmployeesReportingToGivenManagerRole("manager", managerID)
}
