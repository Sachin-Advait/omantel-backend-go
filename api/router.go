package api

import (
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gitlab.com/rendezvousPoint/incentive_calculator/api/v1/admin"
	"gitlab.com/rendezvousPoint/incentive_calculator/api/v1/user"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
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
	incentiveV1Router := router.Group(constant.SalesIncentiveRoute, middleware.ApiVersionHandler("1"), middleware.OptionalAuthHandler())
	incentiveRoutes(incentiveV1Router)
	incentiveSwagV1Router := router.Group(constant.SalesIncentiveDocRoute)
	incentiveSwaggerRoutes(incentiveSwagV1Router)

	//admin version 1
	incentiveAdminV1Router := router.Group(constant.SalesIncentiveAdminRoute, middleware.ApiVersionHandler("1"), middleware.MandatoryAdminAuthHandler())
	incentiveAdminRoutes(incentiveAdminV1Router)
	incentiveSwagAdminV1Router := router.Group(constant.SalesIncentiveDocAdminRoute)
	incentiveSwaggerAdminRoutes(incentiveSwagAdminV1Router)

	logger.FileLogger.Info().Str("env", flags.Env()).Int("port", flags.Port()).Msg("Router started...")
	return router
}

// Incentive Routes
func incentiveRoutes(group *gin.RouterGroup) {

	group.GET("/", func(c *gin.Context) {

		pong := models.Pong{
			DT: time.Now(),
		}

		res := models.APIResponse{
			Status:  true,
			Message: errors.SuccessInfoMessage,
			Data:    pong,
		}

		c.JSON(http.StatusOK, res)
	})

	group.GET(constant.GenerateTokenPath, user.GenerateToken)
	group.POST(constant.GenerateOTPPath, user.GenerateOTP)
	group.POST(constant.VerifyOTPPath, user.VerifyOTP)
	group.GET(constant.RenewTokenPath, user.RenewAccessToken)
	group.POST(constant.CaptchaVerificationPath, user.CaptchaVerification)
	group.POST(constant.GetUserPath, user.GetUserFromRegionOrOutlet)
	group.GET(constant.UserLogoutPath, user.UserLogout)

	group.POST(constant.ClientPath, user.GetUserPayout)
	group.POST(constant.UserPayoutPath, user.GetUserPayout)

	group.POST(constant.TargetChartPath, user.GetSingleUserTargetChart)
	group.POST(constant.UserTargetPath, user.GetSingleUserTargetChart)

	group.POST(constant.SalePerPlanPath, user.GetSaleCountPerPlan)
	group.POST(constant.UserSalesPerPlanPath, user.GetSaleCountPerPlan)

	group.POST(constant.ReportToOutletManagerPath, user.GetUsersReportingToOutletManager)
	group.POST(constant.ReportToManagerPath, user.GetUsersReportingToManager)
	group.POST(constant.ReportToSeniorManagerPath, user.GetUsersReportingToSeniorManager)
	group.POST(constant.ReportToGeneralManagerPath, user.GetUsersReportingToGeneralManager)
	group.POST(constant.ReportToVicePresidentPath, user.GetUsersReportingToVicePresident)

	group.POST(constant.OutletSalesCountPath, user.GetOutletTotalSalesCount)
	group.POST(constant.OutletSalesTargetPath, user.GetOutletSalesAndTarget)
	group.POST(constant.RegionSalesCountPath, user.GetRegionTotalSalesCount)
	group.POST(constant.RegionSalesTargetPath, user.GetRegionSalesAndTarget)
	group.POST(constant.TopTenSalesCountPath, user.GetTopTenEmpProductSalesCount)
	group.POST(constant.TopTenSalesTargetPath, user.GetTopTenEmpSalesAndTarget)

	group.POST(constant.UserReportPath, user.GetUserReport)
	group.POST(constant.UserSalesCountPath, user.GetUserSalesCount)
	group.POST(constant.UserPlanPerCategoryPath, user.GetUserPlanCategory)

	group.GET(constant.FetchNotificationPath, user.FetchNotification)

	group.GET(constant.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func incentiveSwaggerRoutes(group *gin.RouterGroup) {

	group.GET("/", func(c *gin.Context) {

		pong := models.Pong{
			DT: time.Now(),
		}

		res := models.APIResponse{
			Status:  true,
			Message: errors.SuccessInfoMessage,
			Data:    pong,
		}

		c.JSON(http.StatusOK, res)
	})

	group.GET(constant.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Incentive Admin Routes
func incentiveAdminRoutes(group *gin.RouterGroup) {

	group.GET("/", func(c *gin.Context) {

		pong := models.Pong{
			DT: time.Now(),
		}

		res := models.APIResponse{
			Status:  true,
			Message: errors.SuccessInfoMessage,
			Data:    pong,
		}

		c.JSON(http.StatusOK, res)
	})

	group.POST(constant.UpsertRegionPath, admin.UpsertNewRegion)
	group.POST(constant.UpsertOutletPath, admin.UpsertNewOutlet)
	group.POST(constant.UpsertProductPath, admin.UpsertNewProduct)
	group.POST(constant.UpsertPlanPath, admin.UpsertNewPlan)
	group.POST(constant.UpsertTargetPath, admin.UpsertNewTarget)
	group.POST(constant.UpsertVOCCriteriaPath, admin.UpsertVOCCriteria)
	group.POST(constant.UpsertVOCMaxPayPath, admin.UpsertVOCMaxPay)
	group.POST(constant.UpsertVOCUserPath, admin.UpsertVOCUser)

	group.POST(constant.FetchRegionPath, admin.FetchRegions)
	group.POST(constant.FetchOutletPath, admin.FetchOutlets)
	group.POST(constant.FetchProductPath, admin.FetchProducts)
	group.POST(constant.FetchPlanPath, admin.FetchPlans)
	group.POST(constant.FetchTargetPath, admin.FetchTargets)
	group.POST(constant.FetchVOCCriteriaPath, admin.FetchVOCCriteria)
	group.POST(constant.FetchVOCMaxPayPath, admin.FetchVOCMaxPay)
	group.POST(constant.FetchVOCUserPath, admin.FetchVOCUser)

	group.POST(constant.DeleteRegionPath, admin.DeleteRegions)
	group.POST(constant.DeleteOutletPath, admin.DeleteOutlets)
	group.POST(constant.DeleteProductPath, admin.DeleteProducts)
	group.POST(constant.DeletePlanPath, admin.DeletePlans)
	group.POST(constant.DeleteTargetPath, admin.DeleteTargets)
	group.POST(constant.DeleteVOCCriteriaPath, admin.DeleteVOCCriteria)
	group.POST(constant.DeleteVOCMaxPayPath, admin.DeleteVOCMaxPay)
	group.POST(constant.DeleteVOCUserPath, admin.DeleteVOCUser)

	group.POST(constant.UpsertNotificationPath, admin.UpsertNotification)
	group.GET(constant.FetchNotificationPath, admin.FetchNotification)
	group.GET(constant.DeleteNotificationPath, admin.DeleteNotification)

	group.GET(constant.LoginChartPath, admin.GetUserLoginCount)
	group.GET(constant.LoginDetailsPath, admin.GetUserLoginDetails)

	group.GET(constant.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func incentiveSwaggerAdminRoutes(group *gin.RouterGroup) {

	group.GET("/", func(c *gin.Context) {

		pong := models.Pong{
			DT: time.Now(),
		}

		res := models.APIResponse{
			Status:  true,
			Message: errors.SuccessInfoMessage,
			Data:    pong,
		}

		c.JSON(http.StatusOK, res)
	})

	group.GET(constant.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))
}
