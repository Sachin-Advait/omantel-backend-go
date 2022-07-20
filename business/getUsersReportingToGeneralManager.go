package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
)

func (impl *IncentiveImpl) GetUsersReportingToGeneralManager(ctx *gin.Context, generalManagerID int) ([]*psqldb.ReportingStaffIDAndName, *errors.Error) {

	return impl.PsqlDb.FetchAllEmployeesReportingToGivenManagerRole("general_manager", generalManagerID)
}
