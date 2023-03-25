package admin

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

type AdminProvider interface {
	FetchRegions(ctx *gin.Context, req *models.FetchDeleteRegionRequest) ([]*models.RegionResponse, *errors.Error)
	UpsertRegion(ctx *gin.Context, req []*models.CreateRegionRequest) ([]*models.RegionFailedResponse, *errors.Error)
	DeleteRegions(ctx *gin.Context, req *models.FetchDeleteRegionRequest) *errors.Error

	FetchOutlets(ctx *gin.Context, req *models.FetchDeleteOutletRequest) ([]*models.OutletResponse, *errors.Error)
	UpsertOutlet(ctx *gin.Context, req []*models.CreateOutletRequest) ([]*models.OutletFailedResponse, *errors.Error)
	DeleteOutlets(ctx *gin.Context, req *models.FetchDeleteOutletRequest) *errors.Error

	FetchProducts(ctx *gin.Context, req *models.FetchDeleteProductRequest) ([]*models.ProductResponse, *errors.Error)
	UpsertProduct(ctx *gin.Context, req []*models.CreateProductRequest) ([]*models.ProductFailedResponse, *errors.Error)
	DeleteProducts(ctx *gin.Context, req *models.FetchDeleteProductRequest) *errors.Error

	FetchPlans(ctx *gin.Context, req *models.FetchDeletePlanRequest) ([]*models.PlanResponse, *errors.Error)
	UpsertPlan(ctx *gin.Context, req []*models.CreatePlanRequest) ([]*models.PlanFailedResponse, *errors.Error)
	DeletePlans(ctx *gin.Context, req *models.FetchDeletePlanRequest) *errors.Error

	FetchTargets(ctx *gin.Context, req *models.FetchDeleteTargetRequest) (map[string][]*models.TargetResponse, *errors.Error)
	UpsertTarget(ctx *gin.Context, req []*models.CreateTargetRequest) ([]*models.TargetFailedResponse, *errors.Error)
	DeleteTargets(ctx *gin.Context, req *models.FetchDeleteTargetRequest) *errors.Error

	FetchVOCCriteria(ctx *gin.Context, req *models.FetchDeleteVOCCriteriaRequest) (*models.VOCCriteriaResponse, *errors.Error)
	UpsertVOCCriteria(ctx *gin.Context, req []*models.CreateVOCCriteriaRequest) ([]*models.VOCCriteriaFailedResponse, *errors.Error)
	DeleteVOCCriteria(ctx *gin.Context, req *models.FetchDeleteVOCCriteriaRequest) *errors.Error

	FetchVOCMaxPay(ctx *gin.Context, req *models.FetchDeleteVOCCriteriaRequest) ([]*models.VOCMaxPayResponse, *errors.Error)
	UpsertVOCMaxPay(ctx *gin.Context, req []*models.CreateVOCMaxPayRequest) ([]*models.VOCMaxPayFailedResponse, *errors.Error)
	DeleteVOCMaxPay(ctx *gin.Context, req []*models.FetchDeleteVOCMaxPayRequest) ([]*models.VOCMaxPayFailedResponse, *errors.Error)

	FetchVOCUser(ctx *gin.Context, req *models.FetchDeleteVOCUserRequest) ([]*models.VOCUserResponse, *errors.Error)
	UpsertVOCUser(ctx *gin.Context, req []*models.CreateVOCUserRequest) ([]*models.VOCCUserFailedResponse, *errors.Error)
	DeleteVOCUer(ctx *gin.Context, req *models.FetchDeleteVOCUserRequest) *errors.Error

	UpsertNotification(ctx *gin.Context, req *models.Notification) (*models.APISuccessResponse, *errors.Error)
	FetchNotification(ctx *gin.Context, category, id string) (interface{}, *errors.Error)
	DeleteNotification(ctx *gin.Context, id string) *errors.Error
}
