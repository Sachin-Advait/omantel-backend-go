package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToManager(ctx *gin.Context, managerID string) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	response, err := impl.PsqlDb.FetchAllEmpReportingToGivenManagerRole(constant.Manager, managerID, currQuarter, currYear)
	if err != nil {
		return nil, err
	}

	for index, eachOM := range response {

		if eachOM.ID == "" || eachOM.Name == "" || !(eachOM.Position == constant.OutletManager || eachOM.Position == constant.TeamLeader) {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InvalidDetailsCode,
				Message:    "INVALID OM LIST FOUND",
				Details:    "om details incorrect - " + eachOM.ID,
			}
		}

		eachOM.Position = constant.EmployeePositionMap[eachOM.Position]
		response[index] = eachOM
	}

	return response, nil
}
