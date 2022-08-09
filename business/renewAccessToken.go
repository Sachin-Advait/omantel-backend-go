package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	jwthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/jwtHelper"
)

func (impl *IncentiveImpl) RenewAccessToken(ctx *gin.Context) (string, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	authUserID, status, err := jwthelper.ValidateRefreshToken(ctxHeader.Authorization, jwthelper.JwtSecretKey)
	if err != nil {
		if status == constant.TokenExpired {
			return "", &errors.Error{
				StatusCode: http.StatusUnauthorized,
				Code:       errors.TokenExpiredCode,
				Message:    errors.TokenExpired,
				Details:    "refreshToken expired - " + err.Error(),
			}
		}

		return "", &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidTokenCode,
			Message:    errors.InvalidToken,
			Details:    "refreshToken invalid - " + err.Error(),
		}
	}

	if authUserID != ctxHeader.XUserId {
		return "", &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.ContradictoryReqParamCode,
			Message:    errors.ContradictoryReqParam,
			Details:    "contradicting userInfo and authInfo",
		}
	}

	accessToken, _, err := jwthelper.RefreshAccessToken(ctxHeader.Authorization, impl.Config.SessionExpiryTime, jwthelper.JwtSecretKey)
	if err != nil {
		return "", &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "refresh accessToken failed - " + err.Error(),
		}
	}

	return accessToken, nil
}
