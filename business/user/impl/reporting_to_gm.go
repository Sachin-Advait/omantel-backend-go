package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *UserImpl) GetUsersReportingToGeneralManager(ctx *gin.Context, generalManagerID string, chanel string,
	isQuarterReport bool, quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if err := impl.CacheService.IsSessionTokenValid(ctxHeader.XUserId, ctxHeader.Authorization); err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.UnauthorizedUserCode,
			Message:    errors.UnauthorizedUser,
			Details:    err.Details,
		}
	}

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	// if ctxHeader.IsAdmin {
	// 	vps, err := impl.DbService.FetchVP(quarter, year)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	generalManagerID = vps[0]
	// }

	response, err := impl.DbService.FetchAllEmpReportingToGivenManagerRole(constant.GeneralManagerAbb, generalManagerID,
		chanel, quarter, year)
	if err != nil {
		return nil, err
	}

	for index, eachSM := range response {

		if eachSM.ID == "" || eachSM.Name == "" || eachSM.Position != constant.SeniorManagerAbb {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InvalidDetailsCode,
				Message:    "INVALID SM LIST FOUND",
				Details:    "sm details incorrect - " + eachSM.ID,
			}
		}

		//eachSM.Position = constant.EmployeePositionMap[eachSM.Position]
		response[index] = eachSM
	}

	return response, nil
}
