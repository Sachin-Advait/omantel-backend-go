package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToOutletManager(ctx *gin.Context, outletManagerID string, isQuarterReport bool, quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
	}

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)
	if ctxHeader.IsAdmin {
		vps, err := impl.PsqlDb.FetchVP(quarter, year)
		if err != nil {
			return nil, err
		}
		outletManagerID = vps[0]
	}

	response, err := impl.PsqlDb.FetchAllEmpReportingToGivenManagerRole(constant.OutletManager, outletManagerID, quarter, year)
	if err != nil {
		return nil, err
	}

	for index, agents := range response {

		if agents.ID == "" || agents.Name == "" || !(agents.Position == constant.Promoter || agents.Position == constant.CSR) {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InvalidDetailsCode,
				Message:    "INVALID AGENT LIST FOUND",
				Details:    "agent details incorrect - " + agents.ID,
			}
		}

		agents.Position = constant.EmployeePositionMap[agents.Position]
		response[index] = agents
	}

	return response, nil
}
