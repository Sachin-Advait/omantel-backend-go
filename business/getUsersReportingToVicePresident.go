package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	psqldb "gitlab.com/rendezvousPoint/incentive_calculator/utils/psqlDb"
)

func (impl *IncentiveImpl) GetUsersReportingToVicePresident(ctx *gin.Context, vicePresidentID int) ([]*psqldb.ReportingStaffIDAndName, *errors.Error) {

	return impl.PsqlDb.FetchAllEmployeesReportingToGivenManagerRole("vice_president", vicePresidentID)
}
