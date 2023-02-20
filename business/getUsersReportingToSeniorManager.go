package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToSeniorManager(ctx *gin.Context, seniorManagerID string, isQuarterReport bool, quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error) {

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
		seniorManagerID = vps[0]
	}

	response, err := impl.PsqlDb.FetchAllEmpReportingToGivenManagerRole(constant.SeniorManagerAbb, seniorManagerID, quarter, year)
	if err != nil {
		return nil, err
	}

	for index, eachMngr := range response {

		if eachMngr.ID == "" || eachMngr.Name == "" || eachMngr.Position != constant.ManagerAbb {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InvalidDetailsCode,
				Message:    "INVALID MANAGER LIST FOUND",
				Details:    "mngr details incorrect - " + eachMngr.ID,
			}
		}

		//eachMngr.Position = constant.EmployeePositionMap[eachMngr.Position]
		response[index] = eachMngr
	}

	return response, nil
}
