package api

import (
	"github.com/gin-gonic/gin"
	goActuator "github.com/sinhashubham95/go-actuator"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
)

var (
	actuatorHandler = goActuator.GetActuatorHandler(&goActuator.Config{
		Env:     flags.Env(),
		Name:    constant.AppName,
		Port:    flags.Port(),
		Version: constant.AppVersion,
	})
)

func actuator(ctx *gin.Context) {
	actuatorHandler(ctx.Writer, ctx.Request)
}
