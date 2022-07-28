package api

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	v1 "gitlab.com/rendezvousPoint/incentive_calculator/api/v1"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	middleware "gitlab.com/rendezvousPoint/incentive_calculator/utils/middlewares"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func GetRouter(middlewares ...gin.HandlerFunc) *gin.Engine {

	router := gin.New()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middlewares...)
	router.Use(gin.Recovery())

	//configure swagger
	router.GET(constant.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))

	//configure actuator
	router.GET(constant.ActuatorRoute, actuator)

	//version 1
	incentiveV1Router := router.Group(constant.SalesIncentiveRoute, middleware.HeaderLogging())
	v1.IncentiveRoutes(incentiveV1Router)

	incentiveSwagV1Router := router.Group(constant.SalesIncentiveDocRoute)
	v1.IncentiveSwaggerRoutes(incentiveSwagV1Router)

	logger.FileLogger.Info().Str("env", flags.Env()).Int("port", flags.Port()).Msg("Router started...")
	return router
}
