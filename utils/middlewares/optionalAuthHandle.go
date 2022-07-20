package middlewares

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	jwthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/jwtHelper"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
	"net/http"
)

type APIResponse struct {
	Status    bool        `json:"status" example:"false" enums:"true,false"`
	Message   string      `json:"message" example:"CE10001"`
	ErrorCode string      `json:"errorcode" example:"INTERNAL SERVER ERROR"`
	Data      interface{} `json:"data"`
}

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
			logger.FileLogger.CtxError(ctx).Err(errors.New("MIDDLEWARE AUTH FAILED")).Stack().Msg("AuthMiddleware : incorrect access token length")
			ctx.JSON(http.StatusBadRequest, APIResponse{
				ErrorCode: errors.InvalidTokenCode,
				Message:   errors.InvalidToken,
				Data:      "MIDDLEWARE AUTH FAILED",
			})
			ctx.Abort()
			return
		}

		if len(ctxHeader.Authorization) > 7 {

			authUser, userType, source, status, err := jwthelper.ValidateAccessToken(ctxHeader.Authorization[7:], jwthelper.JwtSecretKey)

			if err != nil {
				if status == constant.TokenExpired {
					logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("AuthMiddleware : expired access token")
					ctx.JSON(http.StatusUnauthorized, APIResponse{
						ErrorCode: errors.TokenExpiredCode,
						Message:   errors.TokenExpired,
					})
					ctx.Abort()
					return
				}

				logger.FileLogger.CtxError(ctx).Err(err).Stack().Msg("AuthMiddleware : invalid access token")
				ctx.JSON(http.StatusBadRequest, APIResponse{
					ErrorCode: errors.InvalidTokenCode,
					Message:   errors.InvalidToken,
				})
				ctx.Abort()
				return
			}

			if !objecthelper.IsStrictlyAlphNumeric(authUser) || (ctxHeader.XUserId != "" && authUser != ctxHeader.XUserId) {
				logger.FileLogger.CtxError(ctx).Err(errors.New("USERID AND AUTHTOKEN ARE CONTRADICTORY")).Stack().Msg("AuthMiddleware : userID contradicting with access token")
				ctx.JSON(http.StatusBadRequest, APIResponse{
					ErrorCode: errors.ContradictoryReqParamCode,
					Message:   errors.ContradictoryReqParam,
					Data:      "USERID AND AUTHTOKEN ARE CONTRADICTORY",
				})
				ctx.Abort()
				return
			}

			if ctxHeader.XUserType != "" && userType != ctxHeader.XUserType {
				logger.FileLogger.CtxError(ctx).Err(errors.New("USERTYPE AND AUTHTOKEN ARE CONTRADICTORY")).Stack().Msg("AuthMiddleware : userType contradicting with access token")
				ctx.JSON(http.StatusBadRequest, APIResponse{
					ErrorCode: errors.ContradictoryReqParamCode,
					Message:   errors.ContradictoryReqParam,
					Data:      "USERTYPE AND AUTHTOKEN ARE CONTRADICTORY",
				})
				ctx.Abort()
				return
			}

			if ctxHeader.XSourceID != "" && source != ctxHeader.XSourceID {
				logger.FileLogger.CtxError(ctx).Err(errors.New("SOURCE AND AUTHTOKEN ARE CONTRADICTORY")).Stack().Msg("AuthMiddleware : source contradicting with access token")
				ctx.JSON(http.StatusBadRequest, APIResponse{
					ErrorCode: errors.ContradictoryReqParamCode,
					Message:   errors.ContradictoryReqParam,
					Data:      "SOURCE AND AUTHTOKEN ARE CONTRADICTORY",
				})
				ctx.Abort()
				return
			}

			ctxHeader.XUserId = authUser
			ctxHeader.XUserType = userType
			ctxHeader.XSourceID = source
			ctxHeader.Authorization = ctxHeader.Authorization[7:]
		}

		ctx.Set(constant.HeaderInContext, ctxHeader)
		ctx.Next()
	}
}
