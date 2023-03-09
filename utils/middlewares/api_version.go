package middlewares

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/context"
)

// ApiVersionHandler:
//	* sets the version of the api
func ApiVersionHandler(version string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		cRH, _ := ctx.Get(constant.HeaderInContext)
		ctxHeader, _ := (cRH).(context.ContextHeader)

		ctxHeader.XAPIVersion = version

		ctx.Set(constant.HeaderInContext, ctxHeader)
		ctx.Next()
	}
}
