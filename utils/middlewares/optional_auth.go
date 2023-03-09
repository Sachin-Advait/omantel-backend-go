package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	jwthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/jwtHelper"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

// OptionalAuthHandler:
//	* fetches values from header and sets into contextHeader in upperCase
//	* if token is present then it validates if its correct
//	* if userID is passed in header then it validates if it corresponds to the right token
//	* if token is not present then it is allowed to bypass
//	* context authorization field is reset with only the access token part
func OptionalAuthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cRH, _ := ctx.Get(constant.HeaderInContext)
		ctxHeader, _ := (cRH).(context.ContextHeader)

		tokenLength := len(ctxHeader.Authorization)

		if tokenLength > 0 && tokenLength <= 7 {
			logger.FileLogger.CtxError(ctx).Err(errors.New("MIDDLEWARE AUTH FAILED")).Msg("AuthMiddleware : incorrect token length")
			ctx.JSON(http.StatusBadRequest, models.APIResponse{
				ErrorCode: errors.InvalidTokenCode,
				Message:   errors.InvalidToken,
			})
			ctx.Abort()
			return
		}

		if len(ctxHeader.Authorization) > 7 {

			//first validate admin token
			authUser, status, err := jwthelper.ValidateAdminAccessToken(ctxHeader.Authorization[7:], jwthelper.AdminJwtSecretKey)
			if err == nil {
				ctxHeader.IsAdmin = true
				ctxHeader.XUserId = authUser
				ctxHeader.Authorization = ctxHeader.Authorization[7:]
			} else {
				authUser, status, err = jwthelper.ValidateAccessToken(ctxHeader.Authorization[7:], jwthelper.JwtSecretKey)
				if err != nil {
					if status == constant.TokenExpired {
						logger.FileLogger.CtxError(ctx).Err(err).Msg("AuthMiddleware : expired token")
						ctx.JSON(http.StatusUnauthorized, models.APIResponse{
							ErrorCode: errors.TokenExpiredCode,
							Message:   errors.TokenExpired,
						})
						ctx.Abort()
						return
					}

					logger.FileLogger.CtxError(ctx).Err(err).Msg("AuthMiddleware : invalid token")
					ctx.JSON(http.StatusBadRequest, models.APIResponse{
						ErrorCode: errors.InvalidTokenCode,
						Message:   errors.InvalidToken,
					})
					ctx.Abort()
					return
				}

				if ctxHeader.XUserId != "" && authUser != ctxHeader.XUserId {
					logger.FileLogger.CtxError(ctx).Err(errors.New("USERID AND AUTHTOKEN ARE CONTRADICTORY")).Msg("AuthMiddleware : userID contradicting with token")
					ctx.JSON(http.StatusBadRequest, models.APIResponse{
						ErrorCode: errors.ContradictoryReqParamCode,
						Message:   errors.ContradictoryReqParam,
						Data:      "USERID AND AUTHTOKEN ARE CONTRADICTORY",
					})
					ctx.Abort()
					return
				}

				ctxHeader.IsAdmin = false
				ctxHeader.XUserId = authUser
				ctxHeader.Authorization = ctxHeader.Authorization[7:]
			}
		}

		ctx.Set(logger.UserIDLogParam, ctxHeader.XUserId)
		ctx.Set(constant.HeaderInContext, ctxHeader)
		ctx.Next()
	}
}
