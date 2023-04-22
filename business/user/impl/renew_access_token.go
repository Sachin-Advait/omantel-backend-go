package impl

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	jwthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/jwtHelper"
)

func (impl *UserImpl) RenewAccessToken(ctx *gin.Context) (string, *errors.Error) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	var (
		isAdmin    bool
		authUserID string
		status     int
		err        error
	)

	authUserID, status, err = jwthelper.ValidateRefreshToken(ctxHeader.SuperAuthorization, jwthelper.JwtSecretKey)
	if err != nil {
		invalidErr, ok := err.(*jwt.ValidationError)
		if ok && invalidErr.Errors == jwt.ValidationErrorSignatureInvalid {
			isAdmin = true
		}
	}

	if isAdmin {
		authUserID, status, err = jwthelper.ValidateAdminRefreshToken(ctxHeader.SuperAuthorization, jwthelper.AdminJwtSecretKey)
	}

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

	var accessToken string

	if isAdmin {
		accessToken, _, err = jwthelper.RefreshAccessToken(ctxHeader.SuperAuthorization, impl.Config.SessionExpiryTime, jwthelper.AdminJwtSecretKey, true)
	} else {
		accessToken, _, err = jwthelper.RefreshAccessToken(ctxHeader.SuperAuthorization, impl.Config.SessionExpiryTime, jwthelper.JwtSecretKey, false)
	}

	if err != nil {
		return "", &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "refresh accessToken failed - " + err.Error(),
		}
	}

	session := &models.UserSessionToken{
		AccessToken:  accessToken,
		RefreshToken: ctxHeader.SuperAuthorization,
		LoginTime:    time.Now(),
	}

	//multiple login validation check
	if cacheErr := impl.CacheService.InsertUseSessionToken(authUserID, session, impl.Config.SessionExpiryTime); cacheErr != nil {
		return "", cacheErr
	}

	return accessToken, nil
}
