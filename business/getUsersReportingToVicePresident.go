package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToVicePresident(ctx *gin.Context, vicePresidentID string) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	response, err := impl.PsqlDb.FetchAllEmpReportingToGivenManagerRole(constant.VicePresident, vicePresidentID, currQuarter, currYear)
	if err != nil {
		return nil, err
	}

	for index, eachGM := range response {

		if eachGM.ID == "" || eachGM.Name == "" || eachGM.Position != constant.GeneralManager {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InvalidDetailsCode,
				Message:    "INVALID GM LIST FOUND",
				Details:    "gm details incorrect - " + eachGM.ID,
			}
		}

		eachGM.Position = constant.EmployeePositionMap[eachGM.Position]
		response[index] = eachGM
	}

	return response, nil
}
