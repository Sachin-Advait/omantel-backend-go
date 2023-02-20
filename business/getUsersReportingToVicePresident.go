package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToVicePresident(ctx *gin.Context, vicePresidentID string, isQuarterReport bool, quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	if !isQuarterReport {
		quarter = impl.getQuarterFromMonth(month)
	}

	response, err := impl.PsqlDb.FetchAllEmpReportingToGivenManagerRole(constant.VicePresidentAbb, vicePresidentID, quarter, year)
	if err != nil {
		return nil, err
	}

	for index, eachGM := range response {

		if eachGM.ID == "" || eachGM.Name == "" || eachGM.Position != constant.GeneralManagerAbb {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InvalidDetailsCode,
				Message:    "INVALID GM LIST FOUND",
				Details:    "gm details incorrect - " + eachGM.ID,
			}
		}

		//eachGM.Position = constant.EmployeePositionMap[eachGM.Position]
		response[index] = eachGM
	}

	return response, nil
}
