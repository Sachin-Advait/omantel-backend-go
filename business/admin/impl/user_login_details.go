package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *AdminImpl) GetUserLoginDetails(ctx *gin.Context, req *models.UserLoginDetailsRequest) (*models.UserLoginDetailsResponse, *errors.Error) {

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

	loginUsers, loginUsersErr := impl.DbService.FetchLoginDetails(req.Page, req.Size, req.Region, req.From, req.To)
	if loginUsersErr != nil {

	}

	quarter := getQuarterFromMonth(int(req.From.Month()))
	year := req.From.Year()

	for _, user := range loginUsers {
		userDb, userDbErr := impl.DbService.FetchUserDetails(user.StaffId, quarter, year)
		if userDbErr != nil {
			return nil, userDbErr
		}
		user.StaffName = userDb.UserId
		user.StaffPosition = userDb.Position
		user.StaffMobileNo = userDb.MobileNumber
		user.StaffRegion = userDb.Region
		user.StaffOutlet = userDb.Outlet
		user.SuperVisorId = userDb.ManagerId
		user.SuperVisorPosition = userDb.ManagerPosition
	}

	response := &models.UserLoginDetailsResponse{
		Page:  req.Page + 1,
		Size:  req.Size,
		Users: loginUsers,
	}

	return response, nil
}
