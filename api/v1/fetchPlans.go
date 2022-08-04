package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func fetchPlans(ctx *gin.Context) {

	logger.FileLogger.CtxInfo(ctx).Interface(constant.BodyLogParam, fmt.Sprintf("%+v", ctx.Request.URL.Query())).Msg("fetchPlans : req body")

	ctx.JSON(adapters.FetchPlans(ctx))
}
