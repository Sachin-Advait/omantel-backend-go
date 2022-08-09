package v1

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/adapters"
)

func renewAccessToken(ctx *gin.Context) {

	ctx.JSON(adapters.RenewAccessToken(ctx))
}
