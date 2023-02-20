package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToManager(ctx *gin.Context, managerID string, isQuarterReport bool, quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error) {

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
		managerID = vps[0]
	}

	response, err := impl.PsqlDb.FetchAllEmpReportingToGivenManagerRole(constant.ManagerAbb, managerID, quarter, year)
	if err != nil {
		return nil, err
	}

	for index, eachOM := range response {

		if eachOM.ID == "" || eachOM.Name == "" || eachOM.Position != constant.OutletManagerAbb { //!(eachOM.Position == constant.OutletManager || eachOM.Position == constant.TeamLeader)
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InvalidDetailsCode,
				Message:    "INVALID OM LIST FOUND",
				Details:    "om details incorrect - " + eachOM.ID,
			}
		}

		//eachOM.Position = constant.EmployeePositionMap[eachOM.Position]
		response[index] = eachOM
	}

	return response, nil
}
