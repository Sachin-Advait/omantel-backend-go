package business

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToOutletManager(ctx *gin.Context, outletManagerID string) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	omResponse, err := impl.PsqlDb.FetchAllEmpReportingToGivenManagerRole(constant.OutletManager, outletManagerID)
	if err != nil {
		return nil, err
	}

	tlResponse, err := impl.PsqlDb.FetchAllEmpReportingToGivenManagerRole(constant.TeamLeader, outletManagerID)
	if err != nil {
		return nil, err
	}

	omResponse = append(omResponse, tlResponse...)

	return omResponse, nil
}
