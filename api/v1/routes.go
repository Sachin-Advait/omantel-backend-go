package v1

import (
	"net/http"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Incentive Routes
func IncentiveRoutes(group *gin.RouterGroup) {

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

	group.GET(constant.GenerateTokenPath, generateToken)

	group.POST(constant.GenerateOTPPath, generateOTP)
	group.POST(constant.VerifyOTPPath, verifyOTP)
	group.POST(constant.CaptchaVerificationPath, captchaVerification)
	//group.OPTIONS(constant.VerifyOTPPath, verifyOTP)

	group.POST(constant.ClientPath, getSingleUserIncentive)
	group.POST(constant.TargetChartPath, getSingleUserTargetChart)
	group.POST(constant.SalePerPlanPath, getSaleCountPerPlan)
	group.POST(constant.UserReportPath, getUserReport)
	group.POST(constant.ReportToOutletManagerPath, getUsersReportingToOutletManager)
	group.POST(constant.ReportToManagerPath, getUsersReportingToManager)
	group.POST(constant.ReportToSeniorManagerPath, getUsersReportingToSeniorManager)
	group.POST(constant.ReportToGeneralManagerPath, getUsersReportingToGeneralManager)
	group.POST(constant.ReportToVicePresidentPath, getUsersReportingToVicePresident)
	group.POST(constant.UserPayoutPath, getUserPayout)
	group.POST(constant.UserTargetPath, getUserTarget)
	group.POST(constant.UsrSalesPerPlanPath, getSaleCountPerPlan)
	group.POST(constant.OutletSalesCountPath, getOutletTotalSalesCount)
	group.POST(constant.OutletSalesTargetPath, getOutletSalesAndTarget)
	group.POST(constant.RegionSalesCountPath, getRegionTotalSalesCount)
	group.POST(constant.RegionSalesTargetPath, getRegionSalesAndTarget)
	group.POST(constant.TopTenSalesCountPath, getTopTenEmpProductSalesCount)
	group.POST(constant.TopTenSalesTargetPath, getTopTenEmpSalesAndTarget)
	group.GET(constant.RenewTokenPath, renewAccessToken)

	group.POST(constant.UpsertRegionPath, upsertNewRegion)
	group.POST(constant.UpsertOutletPath, upsertNewOutlet)
	group.POST(constant.UpsertProductPath, upsertNewProduct)
	group.POST(constant.UpsertPlanPath, upsertNewPlan)
	group.POST(constant.UpsertTargetPath, upsertNewTarget)

	group.POST(constant.FetchRegionPath, fetchRegions)
	group.POST(constant.FetchOutletPath, fetchOutlets)
	group.POST(constant.FetchProductPath, fetchProducts)
	group.POST(constant.FetchPlanPath, fetchPlans)
	group.POST(constant.FetchTargetPath, fetchTargets)

	group.POST(constant.DeleteRegionPath, deleteRegions)
	group.POST(constant.DeleteOutletPath, deleteOutlets)
	group.POST(constant.DeleteProductPath, deleteProducts)
	group.POST(constant.DeletePlanPath, deletePlans)
	group.POST(constant.DeleteTargetPath, deleteTargets)

	group.GET(constant.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func IncentiveSwaggerRoutes(group *gin.RouterGroup) {

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
