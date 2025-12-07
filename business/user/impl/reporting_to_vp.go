package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *UserImpl) GetUsersReportingToVicePresident(ctx *gin.Context, vicePresidentID string, chanel string,
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

	if ctxHeader.IsAdmin {
		vps, err := impl.DbService.FetchVP(quarter, year)
		if err != nil {
			return nil, err
		}
		vicePresidentID = vps[0]
	}

	response, err := impl.DbService.FetchAllEmpReportingToGivenManagerRole(constant.VicePresidentAbb, vicePresidentID,
		chanel, quarter, year)
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
