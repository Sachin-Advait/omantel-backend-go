package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	jwthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/jwtHelper"
)

func (impl *IncentiveImpl) GenerateToken(ctx *gin.Context) (string, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	user, err := impl.PsqlDb.FetchUserBasicDetails(ctxHeader.XUserId, 0, 0, true)
	if err != nil {
		return "", err
	}

	if !user.IsActive {
		return "", &errors.Error{
			StatusCode: http.StatusUnauthorized,
			Code:       errors.UnauthorizedUserCode,
			Message:    "YOU ARE NOT AUTHORIZED TO VIEW REPORTS",
			Details:    "position:" + user.Position + ", contact:" + user.MobileNumber,
		}
	}

	accessTokenExpiry := impl.Config.SessionExpiryTime

	accessToken, _, jwtErr := jwthelper.CreateAccessToken(ctxHeader.XUserId, accessTokenExpiry, jwthelper.JwtSecretKey, false)
	if jwtErr != nil {
		return "", &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "create at failed - " + err.Error(),
		}
	}

	return accessToken, nil
}
