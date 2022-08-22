package v2

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

	group.POST(constant.UpsertRegionPath, createNewRegionViaFile)
	group.POST(constant.UpsertOutletPath, createNewOutletViaFile)
	group.POST(constant.UpsertProductPath, createNewProductViaFile)
	group.POST(constant.UpsertPlanPath, createNewPlanViaFile)
	group.POST(constant.UpsertTargetPath, createNewTargetViaFile)

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
