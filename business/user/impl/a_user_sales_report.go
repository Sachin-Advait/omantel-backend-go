package impl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) updateAgentUserSalesCountReport(ctx *gin.Context, user *models.UserBasicDetailsDB, product string,
	chanel string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("month", month).Int("year", year).Msg("ASalesCountReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "updateAgentUserSalesCountReport: invalid agent details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("ASalesCountReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": INVALID AGENT DETAILS",
			Details:    "updateAgentUserSalesCountReport: invalid agent details - userId:" + user.UserId,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("ASalesCountReport : info")

	if !(user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb) {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("position", user.Position).Msg("ASalesCountReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details:    "ASalesCountReport: incorrect agent position - userId:" + user.UserId + " position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).
			Msg("ASalesCountReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + user.Region + " outlet:" + user.Outlet,
			Details:    "ASalesCountReport: agent blank region or outlet - userId: " + user.UserId + " region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	allUserSales, err := impl.DbService.FetchUserSalesCountPerPlan(user.UserId, isQuarterReport, quarter, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserSalesCountReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("ASalesCountReport : fetch sales failed")
		return nil, err
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("saleGroupCount", len(allUserSales)).
		Msg("ASalesCountReport : total agents and salesGroup")

	if product == constant.HBBCategory {
		response := models.HbbUserSalesCountPerPlanDB{
			Product:    constant.HBBCategory,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		for _, v := range allUserSales {
			if v.ProductGroup != constant.HBBCategory {
				continue
			}
			switch v.BusinessUnit {
			case "GOVERNMENT":
				if response.BusinessUnit.Government == nil {
					response.BusinessUnit.Government = &models.SaleCalculation{}
				}
				response.BusinessUnit.Government.Count = response.BusinessUnit.Government.Count + v.Quantity
			case "CONSUMER":
				if response.BusinessUnit.Consumer == nil {
					response.BusinessUnit.Consumer = &models.SaleCalculation{}
				}
				response.BusinessUnit.Consumer.Count = response.BusinessUnit.Consumer.Count + v.Quantity

				if v.LineType == "FIXED WFBB SERVICE" {
					//category
					if response.Category.CBUWFBB == nil {
						response.Category.CBUWFBB = &models.SaleCalculation{}
					}
					response.Category.CBUWFBB.Count = response.Category.CBUWFBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.CBUWFBB == nil {
						response.SubCategory.CBUWFBB = make(map[string]interface{})
					}
					if response.SubCategory.CBUWFBB["total"] == nil {
						response.SubCategory.CBUWFBB["total"] = 0
					}
					total := response.SubCategory.CBUWFBB["total"].(int)
					response.SubCategory.CBUWFBB["total"] = total + v.Quantity
					if response.SubCategory.CBUWFBB[v.SubCategory] == nil {
						response.SubCategory.CBUWFBB[v.SubCategory] = &models.SaleCalculation{}
					}
					subcategory := response.SubCategory.CBUWFBB[v.SubCategory].(*models.SaleCalculation)
					subcategory.Count = subcategory.Count + v.Quantity
					response.SubCategory.CBUWFBB[v.SubCategory] = subcategory

					//serviceCategory
					if response.ServiceCategory.CBUWFBB == nil {
						response.ServiceCategory.CBUWFBB = make(map[string]interface{})
					}
					if response.ServiceCategory.CBUWFBB["total"] == nil {
						response.ServiceCategory.CBUWFBB["total"] = 0
					}
					total = response.ServiceCategory.CBUWFBB["total"].(int)
					response.ServiceCategory.CBUWFBB["total"] = total + v.Quantity

					//plans
					if response.Plans.CBUWFBB == nil {
						response.Plans.CBUWFBB = make(map[string]int)
					}
					total = response.Plans.CBUWFBB["total"]
					response.Plans.CBUWFBB["total"] = total + v.Quantity
					total = response.Plans.CBUWFBB[v.CurrentPlan]
					response.Plans.CBUWFBB[v.CurrentPlan] = total + v.Quantity
				}
				if v.LineType == "ILL SERVICE" || v.LineType == "MPLS SERVICE" || v.LineType == "P2P LEASED LINE SERVICE" || v.LineType == "PSTN VOICE SERVICE" {
					//category
					if response.Category.CBUHBB == nil {
						response.Category.CBUHBB = &models.SaleCalculation{}
					}
					response.Category.CBUHBB.Count = response.Category.CBUHBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.CBUHBB == nil {
						response.SubCategory.CBUHBB = make(map[string]interface{})
					}
					if response.SubCategory.CBUHBB["total"] == nil {
						response.SubCategory.CBUHBB["total"] = 0
					}

					total := response.SubCategory.CBUHBB["total"].(int)
					response.SubCategory.CBUHBB["total"] = total + v.Quantity
					if response.SubCategory.CBUHBB[v.SubCategory] == nil {
						response.SubCategory.CBUHBB[v.SubCategory] = &models.SaleCalculation{}
					}
					subcategory := response.SubCategory.CBUHBB[v.SubCategory].(*models.SaleCalculation)
					subcategory.Count = subcategory.Count + v.Quantity
					response.SubCategory.CBUHBB[v.SubCategory] = subcategory

					//serviceCategory
					if response.ServiceCategory.CBUHBB == nil {
						response.ServiceCategory.CBUHBB = make(map[string]interface{})
					}
					if response.ServiceCategory.CBUHBB["total"] == nil {
						response.ServiceCategory.CBUHBB["total"] = 0
					}
					total = response.ServiceCategory.CBUHBB["total"].(int)
					response.ServiceCategory.CBUHBB["total"] = total + v.Quantity
					if response.ServiceCategory.CBUHBB[v.FinalServiceCategory] == nil {
						response.ServiceCategory.CBUHBB[v.FinalServiceCategory] = &models.SaleCalculation{}
					}
					srvcCat := response.ServiceCategory.CBUHBB[v.FinalServiceCategory].(*models.SaleCalculation)
					srvcCat.Count = srvcCat.Count + v.Quantity
					response.ServiceCategory.CBUHBB[v.FinalServiceCategory] = srvcCat

					//plans
					if response.Plans.CBUHBB == nil {
						response.Plans.CBUHBB = make(map[string]int)
					}
					count := response.Plans.CBUHBB["total"]
					response.Plans.CBUHBB["total"] = count + v.Quantity
					count = response.Plans.CBUHBB[v.CurrentPlan]
					response.Plans.CBUHBB[v.CurrentPlan] = count + v.Quantity
				}
			case "CORPORATE":
				if response.BusinessUnit.Corporate == nil {
					response.BusinessUnit.Corporate = &models.SaleCalculation{}
				}
				response.BusinessUnit.Corporate.Count = response.BusinessUnit.Corporate.Count + v.Quantity

				if v.LineType == "FIXED WFBB SERVICE" {
					//category
					if response.Category.EBUWFBB == nil {
						response.Category.EBUWFBB = &models.SaleCalculation{}
					}
					response.Category.EBUWFBB.Count = response.Category.EBUWFBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.EBUWFBB == nil {
						response.SubCategory.EBUWFBB = make(map[string]interface{})
					}
					if response.SubCategory.EBUWFBB["total"] == nil {
						response.SubCategory.EBUWFBB["total"] = 0
					}
					total := response.SubCategory.EBUWFBB["total"].(int)
					response.SubCategory.EBUWFBB["total"] = total + v.Quantity
					if response.SubCategory.EBUWFBB[v.SubCategory] == nil {
						response.SubCategory.EBUWFBB[v.SubCategory] = &models.SaleCalculation{}
					}
					subcategory := response.SubCategory.EBUWFBB[v.SubCategory].(*models.SaleCalculation)
					subcategory.Count = subcategory.Count + v.Quantity
					response.SubCategory.EBUWFBB[v.SubCategory] = subcategory

					//serviceCategory
					if response.ServiceCategory.EBUWFBB == nil {
						response.ServiceCategory.EBUWFBB = make(map[string]interface{})
					}
					if response.ServiceCategory.EBUWFBB["total"] == nil {
						response.ServiceCategory.EBUWFBB["total"] = 0
					}
					total = response.ServiceCategory.EBUWFBB["total"].(int)
					response.ServiceCategory.EBUWFBB["total"] = total + v.Quantity

					//plans
					if response.Plans.EBUWFBB == nil {
						response.Plans.EBUWFBB = make(map[string]int)
					}
					total = response.Plans.EBUWFBB["total"]
					response.Plans.EBUWFBB["total"] = total + v.Quantity
					total = response.Plans.EBUWFBB[v.CurrentPlan]
					response.Plans.EBUWFBB[v.CurrentPlan] = total + v.Quantity
				}
				if v.LineType == "ILL SERVICE" || v.LineType == "MPLS SERVICE" || v.LineType == "P2P LEASED LINE SERVICE" ||
					v.LineType == "PSTN VOICE SERVICE" {
					//category
					if response.Category.EBUHBB == nil {
						response.Category.EBUHBB = &models.SaleCalculation{}
					}
					response.Category.EBUHBB.Count = response.Category.EBUHBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.EBUHBB == nil {
						response.SubCategory.EBUHBB = make(map[string]interface{})
					}
					if response.SubCategory.EBUHBB["total"] == nil {
						response.SubCategory.EBUHBB["total"] = 0
					}
					total := response.SubCategory.EBUHBB["total"].(int)
					response.SubCategory.EBUHBB["total"] = total + v.Quantity
					if response.SubCategory.EBUHBB[v.SubCategory] == nil {
						response.SubCategory.EBUHBB[v.SubCategory] = &models.SaleCalculation{}
					}
					subcategory := response.SubCategory.EBUHBB[v.SubCategory].(*models.SaleCalculation)
					subcategory.Count = subcategory.Count + v.Quantity
					response.SubCategory.EBUHBB[v.SubCategory] = subcategory

					//serviceCategory
					if response.ServiceCategory.EBUHBB == nil {
						response.ServiceCategory.EBUHBB = make(map[string]interface{})
					}
					if response.ServiceCategory.EBUHBB["total"] == nil {
						response.ServiceCategory.EBUHBB["total"] = 0
					}
					total = response.ServiceCategory.EBUHBB["total"].(int)
					response.ServiceCategory.EBUHBB["total"] = total + v.Quantity
					if response.ServiceCategory.EBUHBB[v.FinalServiceCategory] == nil {
						response.ServiceCategory.EBUHBB[v.FinalServiceCategory] = &models.SaleCalculation{}
					}
					srvcCat := response.ServiceCategory.EBUHBB[v.FinalServiceCategory].(*models.SaleCalculation)
					srvcCat.Count = srvcCat.Count + v.Quantity
					response.ServiceCategory.EBUHBB[v.FinalServiceCategory] = srvcCat

					//plans
					if response.Plans.EBUHBB == nil {
						response.Plans.EBUHBB = make(map[string]int)
					}
					count := response.Plans.EBUHBB["total"]
					response.Plans.EBUHBB["total"] = count + v.Quantity
					count = response.Plans.EBUHBB[v.CurrentPlan]
					response.Plans.EBUHBB[v.CurrentPlan] = count + v.Quantity
				}
			default:
			}
		}

		//businesssUnit
		if response.BusinessUnit.Government == nil {
			response.BusinessUnit.Government = &models.SaleCalculation{}
		}
		if response.BusinessUnit.Corporate == nil {
			response.BusinessUnit.Corporate = &models.SaleCalculation{}
		}
		if response.BusinessUnit.Consumer == nil {
			response.BusinessUnit.Consumer = &models.SaleCalculation{}
		}
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//category
		if response.Category.CBUWFBB == nil {
			response.Category.CBUWFBB = &models.SaleCalculation{}
		}
		if response.Category.CBUHBB == nil {
			response.Category.CBUHBB = &models.SaleCalculation{}
		}
		if response.Category.EBUWFBB == nil {
			response.Category.EBUWFBB = &models.SaleCalculation{}
		}
		if response.Category.EBUHBB == nil {
			response.Category.EBUHBB = &models.SaleCalculation{}
		}
		response.Category.CBUWFBB.Percentage = 100
		response.Category.CBUHBB.Percentage = 100
		response.Category.EBUWFBB.Percentage = 100
		response.Category.EBUHBB.Percentage = 100

		//subCategory
		if response.SubCategory.CBUWFBB == nil {
			response.SubCategory.CBUWFBB = make(map[string]interface{})
		}
		if response.SubCategory.CBUHBB == nil {
			response.SubCategory.CBUHBB = make(map[string]interface{})
		}
		if response.SubCategory.EBUWFBB == nil {
			response.SubCategory.EBUWFBB = make(map[string]interface{})
		}
		if response.SubCategory.EBUHBB == nil {
			response.SubCategory.EBUHBB = make(map[string]interface{})
		}
		var subCategoryCBUWFBBTotal int
		for subcategory, v := range response.SubCategory.CBUWFBB {
			if subcategory == "total" {
				subCategoryCBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUWFBBTotal)
			}
			response.SubCategory.CBUWFBB[subcategory] = obj
		}
		var subCategoryCBUHBBTotal int
		for subcategory, v := range response.SubCategory.CBUHBB {
			if subcategory == "total" {
				subCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUHBBTotal)
			}
			response.SubCategory.CBUHBB[subcategory] = obj
		}
		var subCategoryEBUWFBBTotal int
		for subcategory, v := range response.SubCategory.EBUWFBB {
			if subcategory == "total" {
				subCategoryEBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUWFBBTotal)
			}
			response.SubCategory.EBUWFBB[subcategory] = obj
		}
		var subCategoryEBUHBBTotal int
		for subcategory, v := range response.SubCategory.EBUHBB {
			if subcategory == "total" {
				subCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUHBBTotal)
			}
			response.SubCategory.EBUHBB[subcategory] = obj
		}

		//serviceCategory
		if response.ServiceCategory.CBUWFBB == nil {
			response.ServiceCategory.CBUWFBB = make(map[string]interface{})
		}
		if response.ServiceCategory.CBUHBB == nil {
			response.ServiceCategory.CBUHBB = make(map[string]interface{})
		}
		if response.ServiceCategory.EBUWFBB == nil {
			response.ServiceCategory.EBUWFBB = make(map[string]interface{})
		}
		if response.ServiceCategory.EBUHBB == nil {
			response.ServiceCategory.EBUHBB = make(map[string]interface{})
		}
		var srvcCategoryCBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.CBUHBB {
			if srvccategory == "total" {
				srvcCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryCBUHBBTotal)
			}
			response.ServiceCategory.CBUHBB[srvccategory] = obj
		}
		var srvcCategoryEBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.EBUHBB {
			if srvccategory == "total" {
				srvcCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryEBUHBBTotal)
			}
			response.ServiceCategory.EBUHBB[srvccategory] = obj
		}

		go func(report *models.HbbUserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(user.UserId, report.Product, chanel, isQuarterReport, month, quarter, year, report,
				impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("ASalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("product", product).Msg("ASalesCountReport : cache update success")
			}
		}(&response)

		return &response, nil
	} else {
		response := models.UserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		for _, v := range allUserSales {
			if v.ProductGroup != product {
				continue
			}
			switch v.BusinessUnit {
			case "GOVERNMENT":
				if response.BusinessUnit.Government == nil {
					response.BusinessUnit.Government = &models.SaleCalculation{}
				}
				response.BusinessUnit.Government.Count = response.BusinessUnit.Government.Count + v.Quantity
			case "CONSUMER":
				if response.BusinessUnit.Consumer == nil {
					response.BusinessUnit.Consumer = &models.SaleCalculation{}
				}
				response.BusinessUnit.Consumer.Count = response.BusinessUnit.Consumer.Count + v.Quantity

				//subcategory
				if response.SubCategory.CBU == nil {
					response.SubCategory.CBU = make(map[string]interface{})
				}
				if response.SubCategory.CBU["total"] == nil {
					response.SubCategory.CBU["total"] = 0
				}
				total := response.SubCategory.CBU["total"].(int)
				total = total + v.Quantity
				response.SubCategory.CBU["total"] = total
				if response.SubCategory.CBU[v.SubCategory] == nil {
					response.SubCategory.CBU[v.SubCategory] = &models.SaleCalculation{}
				}
				subcategory := response.SubCategory.CBU[v.SubCategory].(*models.SaleCalculation)
				subcategory.Count = subcategory.Count + v.Quantity
				response.SubCategory.CBU[v.SubCategory] = subcategory

				//plans
				if response.Plans.CBU == nil {
					response.Plans.CBU = make(map[string]int)
				}
				total = response.Plans.CBU["total"]
				response.Plans.CBU["total"] = total + v.Quantity
				total = response.Plans.CBU[v.CurrentPlan]
				response.Plans.CBU[v.CurrentPlan] = total + v.Quantity
			case "CORPORATE":
				if response.BusinessUnit.Corporate == nil {
					response.BusinessUnit.Corporate = &models.SaleCalculation{}
				}
				response.BusinessUnit.Corporate.Count = response.BusinessUnit.Corporate.Count + v.Quantity

				//subcategory
				if response.SubCategory.EBU == nil {
					response.SubCategory.EBU = make(map[string]interface{})
				}
				if response.SubCategory.EBU["total"] == nil {
					response.SubCategory.EBU["total"] = 0
				}
				total := response.SubCategory.EBU["total"].(int)
				total = total + v.Quantity
				response.SubCategory.EBU["total"] = total
				if response.SubCategory.EBU[v.SubCategory] == nil {
					response.SubCategory.EBU[v.SubCategory] = &models.SaleCalculation{}
				}
				subcategory := response.SubCategory.EBU[v.SubCategory].(*models.SaleCalculation)
				subcategory.Count = subcategory.Count + v.Quantity
				response.SubCategory.EBU[v.SubCategory] = subcategory

				//plans
				if response.Plans.EBU == nil {
					response.Plans.EBU = make(map[string]int)
				}
				total = response.Plans.EBU["total"]
				response.Plans.EBU["total"] = total + v.Quantity
				total = response.Plans.EBU[v.CurrentPlan]
				response.Plans.EBU[v.CurrentPlan] = total + v.Quantity
			default:
			}
		}

		//businesssUnit
		if response.BusinessUnit.Government == nil {
			response.BusinessUnit.Government = &models.SaleCalculation{}
		}
		if response.BusinessUnit.Corporate == nil {
			response.BusinessUnit.Corporate = &models.SaleCalculation{}
		}
		if response.BusinessUnit.Consumer == nil {
			response.BusinessUnit.Consumer = &models.SaleCalculation{}
		}
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//subCategory
		if response.SubCategory.CBU == nil {
			response.SubCategory.CBU = make(map[string]interface{})
		}
		if response.SubCategory.EBU == nil {
			response.SubCategory.EBU = make(map[string]interface{})
		}

		var subCategoryCBUTotal int
		for subcategory, v := range response.SubCategory.CBU {
			if subcategory == "total" {
				subCategoryCBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUTotal)
			}
			response.SubCategory.CBU[subcategory] = obj
		}
		var subCategoryEBUTotal int
		for subcategory, v := range response.SubCategory.EBU {
			if subcategory == "total" {
				subCategoryEBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUTotal)
			}
			response.SubCategory.EBU[subcategory] = obj
		}

		go func(report *models.UserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(user.UserId, report.Product, chanel, isQuarterReport,
				month, quarter, year, report, impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("product", product).
					Err(errors.New(err.Details.(string))).Msg("ASalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("product", product).
					Msg("ASalesCountReport : cache update success")
			}
		}(&response)

		return &response, nil
	}
}

func (impl *UserImpl) updateOutletManagerUserSaleCountReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB,
	product string, chanel string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {

		if omngr == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
				Int("month", month).Int("year", year).Msg("OMSalesCountReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "updateOMUserSalesReport: invalid omngr details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMSalesCountReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "updateOMUserSalesReport: invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if omngr.Position != constant.OutletManagerAbb { //!(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader)
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("position", omngr.Position).Msg("OMSalesCountReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details:    "updateOMUserSalesReport: incorrect omngr position - userId:" + omngr.UserId + " position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Msg("OMSalesCountReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + omngr.Region + " outlet:" + omngr.Outlet,
			Details:    "updateOMUserSalesReport: agent blank region or outlet - userId: " + omngr.UserId + " region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMSalesCountReport : info")

	userAgents, err := impl.DbService.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(errors.New(err.Details.(string))).Msg("OMSalesCountReport : fetch agents failed")
		return nil, err
	}

	if product == constant.HBBCategory {
		response := &models.HbbUserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		if len(userAgents) == 0 {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMSalesCountReport : zero agents")
			return response, nil
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).Msg("OMSalesCountReport : total agents")
		}
		for _, eachAgent := range userAgents {
			agentReport, err := impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, chanel, isQuarterReport, month, quarter, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMSalesCountReport : cache report failed")
				resp, err := impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, chanel, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("OMSalesCountReport : agent log")
				agentReport = resp
			} else if agentReport != nil {
				userSaleCount, _ := agentReport.(*models.HbbUserSalesCountPerPlanDB)
				if !objecthelper.IsCurrentDay(userSaleCount.ModifiedAt) {
					resp, err := impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, chanel, isQuarterReport, quarter, month, year)
					if err != nil {
						logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMSalesCountReport : agent report failed")
						continue
					}
					rawReport, _ := json.Marshal(resp)
					logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("OMSalesCountReport : agent log")
					agentReport = resp
				}
			}

			if agentReport == nil {
				continue
			}

			userSaleCount := agentReport.(*models.HbbUserSalesCountPerPlanDB)

			//businessUnit
			if response.BusinessUnit.Government == nil {
				response.BusinessUnit.Government = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Corporate == nil {
				response.BusinessUnit.Corporate = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Consumer == nil {
				response.BusinessUnit.Consumer = &models.SaleCalculation{}
			}
			response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
			response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
			response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

			//category
			if response.Category.CBUWFBB == nil {
				response.Category.CBUWFBB = &models.SaleCalculation{}
			}
			if response.Category.CBUHBB == nil {
				response.Category.CBUHBB = &models.SaleCalculation{}
			}
			if response.Category.EBUWFBB == nil {
				response.Category.EBUWFBB = &models.SaleCalculation{}
			}
			if response.Category.EBUHBB == nil {
				response.Category.EBUHBB = &models.SaleCalculation{}
			}
			response.Category.CBUWFBB.Count += userSaleCount.Category.CBUWFBB.Count
			response.Category.CBUHBB.Count += userSaleCount.Category.CBUHBB.Count
			response.Category.EBUWFBB.Count += userSaleCount.Category.EBUWFBB.Count
			response.Category.EBUHBB.Count += userSaleCount.Category.EBUHBB.Count

			//subCategory
			if response.SubCategory.CBUHBB == nil {
				response.SubCategory.CBUHBB = make(map[string]interface{})
			}
			if response.SubCategory.CBUWFBB == nil {
				response.SubCategory.CBUWFBB = make(map[string]interface{})
			}
			if response.SubCategory.EBUWFBB == nil {
				response.SubCategory.EBUWFBB = make(map[string]interface{})
			}
			if response.SubCategory.EBUHBB == nil {
				response.SubCategory.EBUHBB = make(map[string]interface{})
			}

			for subcategory, val := range userSaleCount.SubCategory.CBUWFBB {
				if subcategory == "total" {
					if response.SubCategory.CBUWFBB["total"] == nil {
						response.SubCategory.CBUWFBB["total"] = 0
					}
					response.SubCategory.CBUWFBB["total"] = val.(int) + userSaleCount.SubCategory.CBUWFBB["total"].(int)
					continue
				}
				if response.SubCategory.CBUWFBB[subcategory] == nil {
					response.SubCategory.CBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBUWFBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBUWFBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.CBUHBB {
				if subcategory == "total" {
					if response.SubCategory.CBUHBB["total"] == nil {
						response.SubCategory.CBUHBB["total"] = 0
					}
					response.SubCategory.CBUHBB["total"] = val.(int) + userSaleCount.SubCategory.CBUHBB["total"].(int)
					continue
				}
				if response.SubCategory.CBUHBB[subcategory] == nil {
					response.SubCategory.CBUHBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBUHBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBUHBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBUWFBB {
				if subcategory == "total" {
					if response.SubCategory.EBUWFBB["total"] == nil {
						response.SubCategory.EBUWFBB["total"] = 0
					}
					response.SubCategory.EBUWFBB["total"] = val.(int) + userSaleCount.SubCategory.EBUWFBB["total"].(int)
					continue
				}
				if response.SubCategory.EBUWFBB[subcategory] == nil {
					response.SubCategory.EBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBUWFBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBUWFBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBUHBB {
				if subcategory == "total" {
					if response.SubCategory.EBUHBB["total"] == nil {
						response.SubCategory.EBUHBB["total"] = 0
					}
					response.SubCategory.EBUHBB["total"] = val.(int) + userSaleCount.SubCategory.EBUHBB["total"].(int)
					continue
				}
				if response.SubCategory.EBUHBB[subcategory] == nil {
					response.SubCategory.EBUHBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBUHBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBUHBB[subcategory] = reportSubCat
			}

			//serviceCategory
			if response.ServiceCategory.CBUHBB == nil {
				response.ServiceCategory.CBUHBB = make(map[string]interface{})
			}
			if response.ServiceCategory.CBUWFBB == nil {
				response.ServiceCategory.CBUWFBB = make(map[string]interface{})
			}
			if response.ServiceCategory.EBUWFBB == nil {
				response.ServiceCategory.EBUWFBB = make(map[string]interface{})
			}
			if response.ServiceCategory.EBUHBB == nil {
				response.ServiceCategory.EBUHBB = make(map[string]interface{})
			}

			for srvccategory, val := range userSaleCount.ServiceCategory.CBUWFBB {
				if srvccategory == "total" {
					if response.ServiceCategory.CBUWFBB["total"] == nil {
						response.ServiceCategory.CBUWFBB["total"] = 0
					}
					response.ServiceCategory.CBUWFBB["total"] = val.(int) + userSaleCount.ServiceCategory.CBUWFBB["total"].(int)
				}
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.CBUHBB {
				if srvccategory == "total" {
					if response.ServiceCategory.CBUHBB["total"] == nil {
						response.ServiceCategory.CBUHBB["total"] = 0
					}
					response.ServiceCategory.CBUHBB["total"] = val.(int) + userSaleCount.ServiceCategory.CBUHBB["total"].(int)
					continue
				}
				if response.ServiceCategory.CBUHBB[srvccategory] == nil {
					response.ServiceCategory.CBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.ServiceCategory.CBUHBB[srvccategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.ServiceCategory.CBUHBB[srvccategory] = reportSubCat
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.EBUWFBB {
				if srvccategory == "total" {
					if response.ServiceCategory.EBUWFBB["total"] == nil {
						response.ServiceCategory.EBUWFBB["total"] = 0
					}
					response.ServiceCategory.EBUWFBB["total"] = val.(int) + userSaleCount.ServiceCategory.EBUWFBB["total"].(int)
				}
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.EBUHBB {
				if srvccategory == "total" {
					if response.ServiceCategory.EBUHBB["total"] == nil {
						response.ServiceCategory.EBUHBB["total"] = 0
					}
					response.ServiceCategory.EBUHBB["total"] = val.(int) + userSaleCount.ServiceCategory.EBUHBB["total"].(int)
					continue
				}
				if response.ServiceCategory.EBUHBB[srvccategory] == nil {
					response.ServiceCategory.EBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.ServiceCategory.EBUHBB[srvccategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.ServiceCategory.EBUHBB[srvccategory] = reportSubCat
			}

			//plans
			if response.Plans.CBUWFBB == nil {
				response.Plans.CBUWFBB = make(map[string]int)
			}
			if response.Plans.CBUHBB == nil {
				response.Plans.CBUHBB = make(map[string]int)
			}
			if response.Plans.EBUWFBB == nil {
				response.Plans.EBUWFBB = make(map[string]int)
			}
			if response.Plans.EBUHBB == nil {
				response.Plans.EBUHBB = make(map[string]int)
			}

			for plan, count := range userSaleCount.Plans.CBUWFBB {
				response.Plans.CBUWFBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.CBUHBB {
				response.Plans.CBUHBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBUWFBB {
				response.Plans.EBUWFBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBUHBB {
				response.Plans.EBUHBB[plan] += count
			}
		}

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//category
		response.Category.CBUWFBB.Percentage = 100
		response.Category.CBUHBB.Percentage = 100
		response.Category.EBUWFBB.Percentage = 100
		response.Category.EBUHBB.Percentage = 100

		//subCategory
		var subCategoryCBUWFBBTotal int
		for subcategory, v := range response.SubCategory.CBUWFBB {
			if subcategory == "total" {
				subCategoryCBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUWFBBTotal)
			}
			response.SubCategory.CBUWFBB[subcategory] = obj
		}
		var subCategoryCBUHBBTotal int
		for subcategory, v := range response.SubCategory.CBUHBB {
			if subcategory == "total" {
				subCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUHBBTotal)
			}
			response.SubCategory.CBUHBB[subcategory] = obj
		}
		var subCategoryEBUWFBBTotal int
		for subcategory, v := range response.SubCategory.EBUWFBB {
			if subcategory == "total" {
				subCategoryEBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUWFBBTotal)
			}
			response.SubCategory.EBUWFBB[subcategory] = obj
		}
		var subCategoryEBUHBBTotal int
		for subcategory, v := range response.SubCategory.EBUHBB {
			if subcategory == "total" {
				subCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUHBBTotal)
			}
			response.SubCategory.EBUHBB[subcategory] = obj
		}

		//serviceCategory
		var srvcCategoryCBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.CBUHBB {
			if srvccategory == "total" {
				srvcCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryCBUHBBTotal)
			}
			response.ServiceCategory.CBUHBB[srvccategory] = obj
		}
		var srvcCategoryEBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.EBUHBB {
			if srvccategory == "total" {
				srvcCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryEBUHBBTotal)
			}
			response.ServiceCategory.EBUHBB[srvccategory] = obj
		}

		go func(report *models.HbbUserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(omngr.UserId, report.Product, chanel, isQuarterReport,
				month, quarter, year, report, impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("OMSalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("product", product).Msg("OMSalesCountReport : cache update success")
			}
		}(response)

		return response, nil
	} else {
		response := &models.UserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		if len(userAgents) == 0 {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMSalesCountReport : zero agents")
			return response, nil
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).Msg("OMSalesCountReport : total agents")
		}
		for _, eachAgent := range userAgents {
			agentReport, err := impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, chanel, isQuarterReport, month, quarter, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMSalesCountReport : cache report failed")
				resp, err := impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, chanel, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("OMSalesCountReport : agent log")
				agentReport = resp
			} else if agentReport != nil {
				userSaleCount, _ := agentReport.(*models.UserSalesCountPerPlanDB)
				if !objecthelper.IsCurrentDay(userSaleCount.ModifiedAt) {
					resp, err := impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, chanel, isQuarterReport, quarter, month, year)
					if err != nil {
						logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("OMSalesCountReport : agent report failed")
						continue
					}
					rawReport, _ := json.Marshal(resp)
					logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("OMSalesCountReport : agent log")
					agentReport = resp
				}
			}

			if agentReport == nil {
				continue
			}

			userSaleCount := agentReport.(*models.UserSalesCountPerPlanDB)

			//businessUnit
			if response.BusinessUnit.Government == nil {
				response.BusinessUnit.Government = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Corporate == nil {
				response.BusinessUnit.Corporate = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Consumer == nil {
				response.BusinessUnit.Consumer = &models.SaleCalculation{}
			}
			response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
			response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
			response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

			//subcategory
			if response.SubCategory.CBU == nil {
				response.SubCategory.CBU = make(map[string]interface{})
			}
			if response.SubCategory.EBU == nil {
				response.SubCategory.EBU = make(map[string]interface{})
			}
			if response.SubCategory.CBU["total"] == nil {
				response.SubCategory.CBU["total"] = 0
			}
			if response.SubCategory.EBU["total"] == nil {
				response.SubCategory.EBU["total"] = 0
			}

			for subcategory, val := range userSaleCount.SubCategory.CBU {
				if val == nil {
					continue
				}
				if subcategory == "total" {
					response.SubCategory.CBU["total"] = response.SubCategory.CBU["total"].(int) + val.(int)
					continue
				}
				if response.SubCategory.CBU[subcategory] == nil {
					response.SubCategory.CBU[subcategory] = &models.SaleCalculation{}
				}
				userSubCat, ok := val.(*models.SaleCalculation)
				if !ok {
					logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Interface("type", reflect.TypeOf(val)).
						Msg("OMSalesCountReport : cbu invalid type conversion")
					continue
				}
				reportSubCat := response.SubCategory.CBU[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBU[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBU {
				if subcategory == "total" {
					response.SubCategory.EBU["total"] = response.SubCategory.EBU["total"].(int) + val.(int)
					continue
				}
				if response.SubCategory.EBU[subcategory] == nil {
					response.SubCategory.EBU[subcategory] = &models.SaleCalculation{}
				}
				userSubCat, ok := val.(*models.SaleCalculation)
				if !ok {
					logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Interface("type", reflect.TypeOf(val)).
						Msg("OMSalesCountReport : ebu invalid type conversion")
					continue
				}
				reportSubCat := response.SubCategory.EBU[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBU[subcategory] = reportSubCat
			}

			//plan
			if response.Plans.CBU == nil {
				response.Plans.CBU = make(map[string]int)
			}
			if response.Plans.EBU == nil {
				response.Plans.EBU = make(map[string]int)
			}
			for plan, count := range userSaleCount.Plans.CBU {
				response.Plans.CBU[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBU {
				response.Plans.EBU[plan] += count
			}
		}

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//subCategory
		var subCategoryCBUTotal int
		for subcategory, v := range response.SubCategory.CBU {
			if subcategory == "total" {
				subCategoryCBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUTotal)
			}
			response.SubCategory.CBU[subcategory] = obj
		}
		var subCategoryEBUTotal int
		for subcategory, v := range response.SubCategory.EBU {
			if subcategory == "total" {
				subCategoryEBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUTotal)
			}
			response.SubCategory.EBU[subcategory] = obj
		}

		go func(report *models.UserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(omngr.UserId, report.Product, chanel, isQuarterReport, month, quarter, year, report,
				impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("OMSalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("product", product).Msg("OMSalesCountReport : cache update success")
			}
		}(response)

		return response, nil
	}
}

func (impl *UserImpl) updateManagerUserSalesCountReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB, product string,
	chanel string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if mngr == nil || mngr.UserId == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("MngrSalesCountReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "updateMngrUserSalesCountReport: invalid mngr details",
		}
	}

	position, ok := constant.ManagerMap[mngr.Position]
	if !ok {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrSalesCountReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details:    "updateMngrUserSalesCountReport: incorrect mngr position - userId:" + mngr.UserId + " position:" + mngr.Position,
		}
	}
	mngr.Position = position

	rawDetails, _ := json.Marshal(mngr)
	logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("user", string(rawDetails)).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("MngrSalesCountReport : info")

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManagerAbb { // || mngr.Position == constant.TeamLeader
		report, err := impl.updateOutletManagerUserSaleCountReport(ctx, mngr, product, chanel, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	if mngr.Position != constant.VicePresidentAbb && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrSalesCountReport : user's manager info missing")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    mngr.UserId + ": INVALID MANAGER'S MANAGER DETAILS",
			Details:    mngr.UserId + ": invalid manager's manager details",
		}
	}

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(mngr.Position, mngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Err(errors.New(err.Details.(string))).Msg("MngrSalesCountReport : fetch om list failed")
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrSalesCountReport : zero om")
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Int("omCount", len(allDistinctOM)).Msg("MngrSalesCountReport : total om")
	}

	if product == constant.HBBCategory {
		response := &models.HbbUserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		for _, eachAgent := range allDistinctOM {
			agentReport, err := impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, chanel, isQuarterReport, month, quarter, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("MngrSalesCountReport : cache report failed")
				resp, err := impl.updateOutletManagerUserSaleCountReport(ctx, eachAgent, product, chanel, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("MngrSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("MngrSalesCountReport : agent log")
				agentReport = resp
			} else if agentReport != nil {
				userSaleCount, _ := agentReport.(*models.HbbUserSalesCountPerPlanDB)
				if !objecthelper.IsCurrentDay(userSaleCount.ModifiedAt) {
					resp, err := impl.updateOutletManagerUserSaleCountReport(ctx, eachAgent, product, chanel, isQuarterReport, quarter, month, year)
					if err != nil {
						logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("MngrSalesCountReport : agent report failed")
						continue
					}
					rawReport, _ := json.Marshal(resp)
					logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("MngrSalesCountReport : agent log")
					agentReport = resp
				}
			}

			if agentReport == nil {
				continue
			}
			userSaleCount := agentReport.(*models.HbbUserSalesCountPerPlanDB)

			//businessUnit
			if response.BusinessUnit.Government == nil {
				response.BusinessUnit.Government = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Corporate == nil {
				response.BusinessUnit.Corporate = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Consumer == nil {
				response.BusinessUnit.Consumer = &models.SaleCalculation{}
			}
			response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
			response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
			response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

			//category
			if response.Category.CBUWFBB == nil {
				response.Category.CBUWFBB = &models.SaleCalculation{}
			}
			if response.Category.CBUHBB == nil {
				response.Category.CBUHBB = &models.SaleCalculation{}
			}
			if response.Category.EBUWFBB == nil {
				response.Category.EBUWFBB = &models.SaleCalculation{}
			}
			if response.Category.EBUHBB == nil {
				response.Category.EBUHBB = &models.SaleCalculation{}
			}
			response.Category.CBUWFBB.Count += userSaleCount.Category.CBUWFBB.Count
			response.Category.CBUHBB.Count += userSaleCount.Category.CBUHBB.Count
			response.Category.EBUWFBB.Count += userSaleCount.Category.EBUWFBB.Count
			response.Category.EBUHBB.Count += userSaleCount.Category.EBUHBB.Count

			//subCategory
			if response.SubCategory.CBUWFBB == nil {
				response.SubCategory.CBUWFBB = make(map[string]interface{})
			}
			if response.SubCategory.CBUHBB == nil {
				response.SubCategory.CBUHBB = make(map[string]interface{})
			}
			if response.SubCategory.EBUWFBB == nil {
				response.SubCategory.EBUWFBB = make(map[string]interface{})
			}
			if response.SubCategory.EBUHBB == nil {
				response.SubCategory.EBUHBB = make(map[string]interface{})
			}

			for subcategory, val := range userSaleCount.SubCategory.CBUWFBB {
				if subcategory == "total" {
					if response.SubCategory.CBUWFBB["total"] == nil {
						response.SubCategory.CBUWFBB["total"] = 0
					}
					response.SubCategory.CBUWFBB["total"] = val.(int) + userSaleCount.SubCategory.CBUWFBB["total"].(int)
					continue
				}
				if response.SubCategory.CBUWFBB[subcategory] == nil {
					response.SubCategory.CBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBUWFBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBUWFBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.CBUHBB {
				if subcategory == "total" {
					if response.SubCategory.CBUHBB["total"] == nil {
						response.SubCategory.CBUHBB["total"] = 0
					}
					response.SubCategory.CBUHBB["total"] = val.(int) + userSaleCount.SubCategory.CBUHBB["total"].(int)
					continue
				}
				if response.SubCategory.CBUHBB[subcategory] == nil {
					response.SubCategory.CBUHBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBUHBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBUHBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBUWFBB {
				if subcategory == "total" {
					if response.SubCategory.EBUWFBB["total"] == nil {
						response.SubCategory.EBUWFBB["total"] = 0
					}
					response.SubCategory.EBUWFBB["total"] = val.(int) + userSaleCount.SubCategory.EBUWFBB["total"].(int)
					continue
				}
				if response.SubCategory.EBUWFBB[subcategory] == nil {
					response.SubCategory.EBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBUWFBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBUWFBB[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBUHBB {
				if subcategory == "total" {
					if response.SubCategory.EBUHBB["total"] == nil {
						response.SubCategory.EBUHBB["total"] = 0
					}
					response.SubCategory.EBUHBB["total"] = val.(int) + userSaleCount.SubCategory.EBUHBB["total"].(int)
					continue
				}
				if response.SubCategory.EBUHBB[subcategory] == nil {
					response.SubCategory.EBUHBB[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBUHBB[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBUHBB[subcategory] = reportSubCat
			}

			//serviceCategory
			if response.ServiceCategory.CBUWFBB == nil {
				response.ServiceCategory.CBUWFBB = make(map[string]interface{})
			}
			if response.ServiceCategory.CBUHBB == nil {
				response.ServiceCategory.CBUHBB = make(map[string]interface{})
			}
			if response.ServiceCategory.EBUWFBB == nil {
				response.ServiceCategory.EBUWFBB = make(map[string]interface{})
			}
			if response.ServiceCategory.EBUHBB == nil {
				response.ServiceCategory.EBUHBB = make(map[string]interface{})
			}

			for srvccategory, val := range userSaleCount.ServiceCategory.CBUWFBB {
				if srvccategory == "total" {
					if response.ServiceCategory.CBUWFBB["total"] == nil {
						response.ServiceCategory.CBUWFBB["total"] = 0
					}
					response.ServiceCategory.CBUWFBB["total"] = val.(int) + userSaleCount.ServiceCategory.CBUWFBB["total"].(int)
				}
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.CBUHBB {
				if srvccategory == "total" {
					if response.ServiceCategory.CBUHBB["total"] == nil {
						response.ServiceCategory.CBUHBB["total"] = 0
					}
					response.ServiceCategory.CBUHBB["total"] = val.(int) + userSaleCount.ServiceCategory.CBUHBB["total"].(int)
					continue
				}
				if response.ServiceCategory.CBUHBB[srvccategory] == nil {
					response.ServiceCategory.CBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.ServiceCategory.CBUHBB[srvccategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.ServiceCategory.CBUHBB[srvccategory] = reportSubCat
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.EBUWFBB {
				if srvccategory == "total" {
					if response.ServiceCategory.EBUWFBB["total"] == nil {
						response.ServiceCategory.EBUWFBB["total"] = 0
					}
					response.ServiceCategory.EBUWFBB["total"] = val.(int) + userSaleCount.ServiceCategory.EBUWFBB["total"].(int)
				}
			}
			for srvccategory, val := range userSaleCount.ServiceCategory.EBUHBB {
				if srvccategory == "total" {
					if response.ServiceCategory.EBUHBB["total"] == nil {
						response.ServiceCategory.EBUHBB["total"] = 0
					}
					response.ServiceCategory.EBUHBB["total"] = val.(int) + userSaleCount.ServiceCategory.EBUHBB["total"].(int)
					continue
				}
				if response.ServiceCategory.EBUHBB[srvccategory] == nil {
					response.ServiceCategory.EBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.ServiceCategory.EBUHBB[srvccategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.ServiceCategory.EBUHBB[srvccategory] = reportSubCat
			}

			//plans
			if response.Plans.CBUWFBB == nil {
				response.Plans.CBUWFBB = make(map[string]int)
			}
			if response.Plans.CBUHBB == nil {
				response.Plans.CBUHBB = make(map[string]int)
			}
			if response.Plans.EBUWFBB == nil {
				response.Plans.EBUWFBB = make(map[string]int)
			}
			if response.Plans.EBUHBB == nil {
				response.Plans.EBUHBB = make(map[string]int)
			}

			for plan, count := range userSaleCount.Plans.CBUWFBB {
				response.Plans.CBUWFBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.CBUHBB {
				response.Plans.CBUHBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBUWFBB {
				response.Plans.EBUWFBB[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBUHBB {
				response.Plans.EBUHBB[plan] += count
			}
		}

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//category
		response.Category.CBUWFBB.Percentage = 100
		response.Category.CBUHBB.Percentage = 100
		response.Category.EBUWFBB.Percentage = 100
		response.Category.EBUHBB.Percentage = 100

		//subCategory
		var subCategoryCBUWFBBTotal int
		for subcategory, v := range response.SubCategory.CBUWFBB {
			if subcategory == "total" {
				subCategoryCBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUWFBBTotal)
			}
			response.SubCategory.CBUWFBB[subcategory] = obj
		}
		var subCategoryCBUHBBTotal int
		for subcategory, v := range response.SubCategory.CBUHBB {
			if subcategory == "total" {
				subCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUHBBTotal)
			}
			response.SubCategory.CBUHBB[subcategory] = obj
		}
		var subCategoryEBUWFBBTotal int
		for subcategory, v := range response.SubCategory.EBUWFBB {
			if subcategory == "total" {
				subCategoryEBUWFBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUWFBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUWFBBTotal)
			}
			response.SubCategory.EBUWFBB[subcategory] = obj
		}
		var subCategoryEBUHBBTotal int
		for subcategory, v := range response.SubCategory.EBUHBB {
			if subcategory == "total" {
				subCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUHBBTotal)
			}
			response.SubCategory.EBUHBB[subcategory] = obj
		}

		//serviceCategory
		var srvcCategoryCBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.CBUHBB {
			if srvccategory == "total" {
				srvcCategoryCBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryCBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryCBUHBBTotal)
			}
			response.ServiceCategory.CBUHBB[srvccategory] = obj
		}
		var srvcCategoryEBUHBBTotal int
		for srvccategory, v := range response.ServiceCategory.EBUHBB {
			if srvccategory == "total" {
				srvcCategoryEBUHBBTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if srvcCategoryEBUHBBTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(srvcCategoryEBUHBBTotal)
			}
			response.ServiceCategory.EBUHBB[srvccategory] = obj
		}

		go func(report *models.HbbUserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(mngr.UserId, report.Product, chanel, isQuarterReport, month, quarter, year, report,
				impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("MngrSalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("product", product).Msg("MngrSalesCountReport : cache update success")
			}
		}(response)

		return response, nil
	} else {
		response := &models.UserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		for _, eachAgent := range allDistinctOM {
			agentReport, err := impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, chanel, isQuarterReport, month, quarter, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("MngrSalesCountReport : cache report failed")
				resp, err := impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, chanel, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("MngrSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("MngrSalesCountReport : agent log")
				agentReport = resp
			} else if agentReport != nil {
				userSaleCount, _ := agentReport.(*models.UserSalesCountPerPlanDB)
				if !objecthelper.IsCurrentDay(userSaleCount.ModifiedAt) {
					resp, err := impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, chanel, isQuarterReport, quarter, month, year)
					if err != nil {
						logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("MngrSalesCountReport : agent report failed")
						continue
					}
					rawReport, _ := json.Marshal(resp)
					logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("MngrSalesCountReport : agent log")
					agentReport = resp
				}
			}

			if agentReport == nil {
				continue
			}
			userSaleCount := agentReport.(*models.UserSalesCountPerPlanDB)

			//businessUnit
			if response.BusinessUnit.Government == nil {
				response.BusinessUnit.Government = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Corporate == nil {
				response.BusinessUnit.Corporate = &models.SaleCalculation{}
			}
			if response.BusinessUnit.Consumer == nil {
				response.BusinessUnit.Consumer = &models.SaleCalculation{}
			}
			response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
			response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
			response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

			//category
			if response.SubCategory.CBU == nil {
				response.SubCategory.CBU = make(map[string]interface{})
			}
			if response.SubCategory.EBU == nil {
				response.SubCategory.EBU = make(map[string]interface{})
			}
			for subcategory, val := range userSaleCount.SubCategory.CBU {
				if subcategory == "total" {
					if response.SubCategory.CBU["total"] == nil {
						response.SubCategory.CBU["total"] = 0
					}
					response.SubCategory.CBU["total"] = val.(int) + userSaleCount.SubCategory.CBU["total"].(int)
					continue
				}
				if response.SubCategory.CBU[subcategory] == nil {
					response.SubCategory.CBU[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.CBU[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.CBU[subcategory] = reportSubCat
			}
			for subcategory, val := range userSaleCount.SubCategory.EBU {
				if subcategory == "total" {
					if response.SubCategory.EBU["total"] == nil {
						response.SubCategory.EBU["total"] = 0
					}
					response.SubCategory.EBU["total"] = val.(int) + userSaleCount.SubCategory.EBU["total"].(int)
					continue
				}
				if response.SubCategory.EBU[subcategory] == nil {
					response.SubCategory.EBU[subcategory] = &models.SaleCalculation{}
				}
				userSubCat := val.(*models.SaleCalculation)
				reportSubCat := response.SubCategory.EBU[subcategory].(*models.SaleCalculation)
				reportSubCat.Count = reportSubCat.Count + userSubCat.Count
				response.SubCategory.EBU[subcategory] = reportSubCat
			}

			//plan
			if response.Plans.CBU == nil {
				response.Plans.CBU = make(map[string]int)
			}
			if response.Plans.EBU == nil {
				response.Plans.EBU = make(map[string]int)
			}
			for plan, count := range userSaleCount.Plans.CBU {
				response.Plans.CBU[plan] += count
			}
			for plan, count := range userSaleCount.Plans.EBU {
				response.Plans.EBU[plan] += count
			}
		}

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//subCategory
		var subCategoryCBUTotal int
		for subcategory, v := range response.SubCategory.CBU {
			if subcategory == "total" {
				subCategoryCBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryCBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryCBUTotal)
			}
			response.SubCategory.CBU[subcategory] = obj
		}
		var subCategoryEBUTotal int
		for subcategory, v := range response.SubCategory.EBU {
			if subcategory == "total" {
				subCategoryEBUTotal = v.(int)
				continue
			}
			obj := v.(*models.SaleCalculation)
			if subCategoryEBUTotal != 0 {
				obj.Percentage = float64(obj.Count) * 100 / float64(subCategoryEBUTotal)
			}
			response.SubCategory.EBU[subcategory] = obj
		}

		go func(report *models.UserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(mngr.UserId, report.Product, chanel, isQuarterReport, month, quarter, year, report,
				impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
					Msg("MngrSalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("product", product).Msg("MngrSalesCountReport : cache update success")
			}
		}(response)

		return response, nil
	}
}
