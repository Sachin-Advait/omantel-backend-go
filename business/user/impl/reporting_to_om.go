package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *UserImpl) GetUsersReportingToOutletManager(ctx *gin.Context, outletManagerID string, isQuarterReport bool, quarter, month, year int) ([]*models.ReportingStaffIDAndName, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if err := impl.CacheService.IsSessionTokenValid(ctxHeader.XUserId, ctxHeader.Authorization[7:]); err != nil {
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
		outletManagerID = vps[0]
	}

	response, err := impl.DbService.FetchAllEmpReportingToGivenManagerRole(constant.OutletManagerAbb, outletManagerID, quarter, year)
	if err != nil {
		return nil, err
	}

	for index, agents := range response {

		if agents.ID == "" || agents.Name == "" || !(agents.Position == constant.PromoterAbb || agents.Position == constant.CSRAbb) {
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InvalidDetailsCode,
				Message:    "INVALID AGENT LIST FOUND",
				Details:    "agent details incorrect - " + agents.ID,
			}
		}

		//agents.Position = constant.EmployeePositionMap[agents.Position]
		response[index] = agents
	}

	return response, nil
}
