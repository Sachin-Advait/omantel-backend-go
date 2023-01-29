package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUsersReportingToGeneralManager(ctx *gin.Context, generalManagerID string) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)
	if ctxHeader.IsAdmin {
		vps, err := impl.PsqlDb.FetchVP(currQuarter, currYear)
		if err != nil {
			return nil, err
		}
		generalManagerID = vps[0]
	}

	response, err := impl.PsqlDb.FetchAllEmpReportingToGivenManagerRole(constant.GeneralManager, generalManagerID, currQuarter, currYear)
	if err != nil {
		return nil, err
	}

	for index, eachSM := range response {

		if eachSM.ID == "" || eachSM.Name == "" || eachSM.Position != constant.SeniorManager {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InvalidDetailsCode,
				Message:    "INVALID SM LIST FOUND",
				Details:    "sm details incorrect - " + eachSM.ID,
			}
		}

		eachSM.Position = constant.EmployeePositionMap[eachSM.Position]
		response[index] = eachSM
	}

	return response, nil
}
