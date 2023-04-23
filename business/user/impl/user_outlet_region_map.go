package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

// GetUserFromRegionOrOutlet - user from outlet and region
func (impl *UserImpl) GetUserFromRegionOrOutlet(ctx *gin.Context, outlet, region string, quarter, year int) ([]*models.UserMapResponse, *errors.Error) {

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

	users, err := impl.DbService.FetchUserIdFromOutletAndRegion(region, outlet, quarter, year)
	if err != nil {
		return nil, err
	}

	respArray := make([]*models.UserMapResponse, 0)

	for _, user := range users {

		if outlet != "" && user.Position == constant.OutletManagerAbb {
			respArray = append(respArray, &models.UserMapResponse{
				ClientId:   user.UserId,
				Name:       user.Name,
				Position:   user.Position,
				RegionName: user.Region,
				OutletName: user.Outlet,
				IsAdmin:    user.IsAdmin,
			})
		}

		if outlet == "" && user.Position == constant.ManagerAbb {
			respArray = append(respArray, &models.UserMapResponse{
				ClientId:   user.UserId,
				Name:       user.Name,
				Position:   user.Position,
				RegionName: user.Region,
				OutletName: user.Outlet,
				IsAdmin:    user.IsAdmin,
			})
		}
	}

	return respArray, nil
}
