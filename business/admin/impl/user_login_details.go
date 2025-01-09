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

	loginUsers, userErr := impl.DbService.FetchLoginDetails(req.Page, req.Size, req.Region, req.From, req.To)
	if userErr != nil {
		return nil, &errors.Error{StatusCode: userErr.StatusCode, Code: userErr.Code, Err: userErr.Err, Message: userErr.Message, Details: userErr.Details}
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
		Users: loginUsers,
	}

	if len(loginUsers) == req.Size {
		response.NextPage = req.Page + 1
		response.Size = req.Size
	} else {
		response.NextPage = -1
		response.Size = len(loginUsers)
	}

	return response, nil
}
