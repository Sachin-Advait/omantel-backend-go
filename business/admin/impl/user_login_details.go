package impl

import (
	"net/http"
	"sort"

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

	quarter := getQuarterFromMonth(int(req.From.Month()))
	year := req.From.Year()

	var users []*models.LoginUsers

	if req.Login {
		loginUsers, userErr := impl.DbService.FetchLoginDetails(req.Page, req.Size, req.Region, req.From, req.To)
		if userErr != nil {
			return nil, &errors.Error{StatusCode: userErr.StatusCode, Code: userErr.Code, Err: userErr.Err, Message: userErr.Message, Details: userErr.Details}
		}
		users = append(users, loginUsers...)

	} else {
		masterUsers, userErr := impl.DbService.FetchDistinctUserIdFromOutletAndRegion(req.Region, "", quarter, year)
		if userErr != nil {
			return nil, &errors.Error{StatusCode: userErr.StatusCode, Code: userErr.Code, Err: userErr.Err, Message: userErr.Message, Details: userErr.Details}
		}
		loggedInUsers, userErr := impl.DbService.FetchAllLoginDistinctUsers(req.Region, req.From, req.To)
		if userErr != nil {
			return nil, &errors.Error{StatusCode: userErr.StatusCode, Code: userErr.Code, Err: userErr.Err, Message: userErr.Message, Details: userErr.Details}
		}
	LOOP:
		for _, usermaster := range masterUsers {
			for _, user := range loggedInUsers {
				if usermaster == user {
					continue LOOP
				}
			}
			users = append(users, &models.LoginUsers{StaffId: usermaster})
		}

		sort.Slice(users, func(i, j int) bool {
			return users[i].StaffId < users[j].StaffId
		})

		offset := (req.Page - 1) * req.Size
		if offset < 0 {
			offset = 0
		}

		if offset == 0 {
			if len(users) > req.Size {
				users = users[:req.Size]
			}
		} else {
			if offset >= len(users) {
				users = []*models.LoginUsers{}
			} else {
				users = users[offset:]
			}
		}
	}

	for _, user := range users {
		userDb, userDbErr := impl.DbService.FetchUserDetails(user.StaffId, quarter, year)
		if userDbErr != nil {
			continue
			//return nil, userDbErr
		}
		user.StaffName = userDb.Name
		user.StaffPosition = userDb.Position
		user.StaffMobileNo = userDb.MobileNumber
		user.StaffRegion = userDb.Region
		user.StaffOutlet = userDb.Outlet
		user.SuperVisorId = userDb.ManagerId
		user.SuperVisorPosition = userDb.ManagerPosition
	}

	response := &models.UserLoginDetailsResponse{
		Users: users,
	}

	if len(users) == req.Size {
		response.NextPage = req.Page + 1
		response.Size = req.Size
	} else {
		response.NextPage = -1
		response.Size = len(users)
	}

	return response, nil
}
