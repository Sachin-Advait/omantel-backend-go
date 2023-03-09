package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters/user"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func GenerateToken(ctx *gin.Context) {

	cRH, _ := ctx.Get(constant.HeaderInContext)
	ctxHeader, _ := (cRH).(context.ContextHeader)

	if ctxHeader.XUserId == "" {
		logger.FileLogger.CtxError(ctx).Err(errors.New(errors.UnauthorizedUser)).Msg("generateToken : unauthorized request")
		ctx.JSON(http.StatusBadRequest, models.APIFailureResponse{
			Error:   errors.UnauthorizedUserCode,
			Message: errors.UnauthorizedUser,
		})
		return
	}

	ctx.JSON(user.GenerateToken(ctx))
}
