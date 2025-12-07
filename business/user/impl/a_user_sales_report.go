package impl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
	objecthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/objectHelper"
)

func (impl *UserImpl) updateAgentUserSalesCountReport(ctx *gin.Context, user *models.UserBasicDetailsDB, product string,
	chanel string, isQuarterReport bool, quarter, month, year int) (*models.UserSalesCountPerPlanDB, *models.HbbUserSalesCountPerPlanDB,
	*errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("month", month).Int("year", year).Msg("ASalesCountReport : nil info")
			return nil, nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "updateAgentUserSalesCountReport: invalid agent details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("ASalesCountReport : invalid info")
		return nil, nil, &errors.Error{
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
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details:    "ASalesCountReport: incorrect agent position - userId:" + user.UserId + " position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).
			Msg("ASalesCountReport : invalid outlet info")
		return nil, nil, &errors.Error{
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
		return nil, nil, err
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("saleGroupCount", len(allUserSales)).
		Msg("ASalesCountReport : total agents and salesGroup")

	if product == constant.HBBCategory {
		response := &models.HbbUserSalesCountPerPlanDB{
			Product:    constant.HBBCategory,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		response.BusinessUnit.Consumer = &models.SaleCalculation{}
		response.BusinessUnit.Corporate = &models.SaleCalculation{}
		response.BusinessUnit.Government = &models.SaleCalculation{}
		response.Category.CBUWFBB = &models.SaleCalculation{Percentage: 100}
		response.Category.CBUHBB = &models.SaleCalculation{Percentage: 100}
		response.Category.EBUWFBB = &models.SaleCalculation{Percentage: 100}
		response.Category.EBUHBB = &models.SaleCalculation{Percentage: 100}
		response.SubCategory.CBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.SubCategory.CBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.SubCategory.EBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.SubCategory.EBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.CBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.CBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.EBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.EBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.Plans.CBUWFBB = map[string]int{"total": 0}
		response.Plans.CBUHBB = map[string]int{"total": 0}
		response.Plans.EBUWFBB = map[string]int{"total": 0}
		response.Plans.EBUHBB = map[string]int{"total": 0}

		for _, v := range allUserSales {
			if v.ProductGroup != constant.HBBCategory {
				continue
			}
			switch v.BusinessUnit {
			case "GOVERNMENT":
				response.BusinessUnit.Government.Count = response.BusinessUnit.Government.Count + v.Quantity
			case "CONSUMER":
				response.BusinessUnit.Consumer.Count = response.BusinessUnit.Consumer.Count + v.Quantity

				if v.LineType == "FIXED WFBB SERVICE" {
					//category
					response.Category.CBUWFBB.Count = response.Category.CBUWFBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.CBUWFBB[v.SubCategory] == nil {
						response.SubCategory.CBUWFBB[v.SubCategory] = &models.SaleCalculation{}
					}
					response.SubCategory.CBUWFBB["total"].Count = response.SubCategory.CBUWFBB["total"].Count + v.Quantity
					response.SubCategory.CBUWFBB[v.SubCategory].Count = response.SubCategory.CBUWFBB[v.SubCategory].Count + v.Quantity

					//serviceCategory
					response.ServiceCategory.CBUWFBB["total"].Count = response.ServiceCategory.CBUWFBB["total"].Count + v.Quantity

					//plans
					response.Plans.CBUWFBB["total"] = response.Plans.CBUWFBB["total"] + v.Quantity
					response.Plans.CBUWFBB[v.CurrentPlan] = response.Plans.CBUWFBB[v.CurrentPlan] + v.Quantity
				}
				if v.LineType == "ILL SERVICE" || v.LineType == "MPLS SERVICE" || v.LineType == "P2P LEASED LINE SERVICE" ||
					v.LineType == "PSTN VOICE SERVICE" {
					//category
					response.Category.CBUHBB.Count = response.Category.CBUHBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.CBUHBB[v.SubCategory] == nil {
						response.SubCategory.CBUHBB[v.SubCategory] = &models.SaleCalculation{}
					}
					response.SubCategory.CBUHBB["total"].Count = response.SubCategory.CBUHBB["total"].Count + v.Quantity
					response.SubCategory.CBUHBB[v.SubCategory].Count = response.SubCategory.CBUHBB[v.SubCategory].Count + v.Quantity

					//serviceCategory
					if response.ServiceCategory.CBUHBB[v.FinalServiceCategory] == nil {
						response.ServiceCategory.CBUHBB[v.FinalServiceCategory] = &models.SaleCalculation{}
					}
					response.ServiceCategory.CBUHBB["total"].Count = response.ServiceCategory.CBUHBB["total"].Count + v.Quantity
					response.ServiceCategory.CBUHBB[v.FinalServiceCategory].Count = response.ServiceCategory.CBUHBB[v.FinalServiceCategory].Count + v.Quantity

					//plans
					response.Plans.CBUHBB["total"] = response.Plans.CBUHBB["total"] + v.Quantity
					response.Plans.CBUHBB[v.CurrentPlan] = response.Plans.CBUHBB[v.CurrentPlan] + v.Quantity
				}
			case "CORPORATE":
				response.BusinessUnit.Corporate.Count = response.BusinessUnit.Corporate.Count + v.Quantity

				if v.LineType == "FIXED WFBB SERVICE" {
					//category
					response.Category.EBUWFBB.Count = response.Category.EBUWFBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.EBUWFBB[v.SubCategory] == nil {
						response.SubCategory.EBUWFBB[v.SubCategory] = &models.SaleCalculation{}
					}
					response.SubCategory.EBUWFBB["total"].Count = response.SubCategory.EBUWFBB["total"].Count + v.Quantity
					response.SubCategory.EBUWFBB[v.SubCategory].Count = response.SubCategory.EBUWFBB[v.SubCategory].Count + v.Quantity

					//serviceCategory
					response.ServiceCategory.EBUWFBB["total"].Count = response.ServiceCategory.EBUWFBB["total"].Count + v.Quantity

					//plans
					response.Plans.EBUWFBB["total"] = response.Plans.EBUWFBB["total"] + v.Quantity
					response.Plans.EBUWFBB[v.CurrentPlan] = response.Plans.EBUWFBB[v.CurrentPlan] + v.Quantity
				}
				if v.LineType == "ILL SERVICE" || v.LineType == "MPLS SERVICE" || v.LineType == "P2P LEASED LINE SERVICE" ||
					v.LineType == "PSTN VOICE SERVICE" {
					//category
					response.Category.EBUHBB.Count = response.Category.EBUHBB.Count + v.Quantity

					//subcategory
					if response.SubCategory.EBUHBB[v.SubCategory] == nil {
						response.SubCategory.EBUHBB[v.SubCategory] = &models.SaleCalculation{}
					}
					response.SubCategory.EBUHBB["total"].Count = response.SubCategory.EBUHBB["total"].Count + v.Quantity
					response.SubCategory.EBUHBB[v.SubCategory].Count = response.SubCategory.EBUHBB[v.SubCategory].Count + v.Quantity

					//serviceCategory
					if response.ServiceCategory.EBUHBB[v.FinalServiceCategory] == nil {
						response.ServiceCategory.EBUHBB[v.FinalServiceCategory] = &models.SaleCalculation{}
					}
					response.ServiceCategory.EBUHBB["total"].Count = response.ServiceCategory.EBUHBB["total"].Count + v.Quantity
					response.ServiceCategory.EBUHBB[v.FinalServiceCategory].Count = response.ServiceCategory.EBUHBB[v.FinalServiceCategory].Count + v.Quantity

					//plans
					response.Plans.EBUHBB["total"] = response.Plans.EBUHBB["total"] + v.Quantity
					response.Plans.EBUHBB[v.CurrentPlan] = response.Plans.EBUHBB[v.CurrentPlan] + v.Quantity
				}
			default:
			}
		}

		// HBB CALCULATION

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//category

		//subCategory
		subCategoryCBUWFBBTotal := response.SubCategory.CBUWFBB["total"].Count
		for subcategory, v := range response.SubCategory.CBUWFBB {
			if subcategory == "total" || subCategoryCBUWFBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryCBUWFBBTotal)
			response.SubCategory.CBUWFBB[subcategory] = v
		}

		subCategoryCBUHBBTotal := response.SubCategory.CBUHBB["total"].Count
		for subcategory, v := range response.SubCategory.CBUHBB {
			if subcategory == "total" || subCategoryCBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryCBUHBBTotal)
			response.SubCategory.CBUHBB[subcategory] = v
		}

		subCategoryEBUWFBBTotal := response.SubCategory.EBUWFBB["total"].Count
		for subcategory, v := range response.SubCategory.EBUWFBB {
			if subcategory == "total" || subCategoryEBUWFBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryEBUWFBBTotal)
			response.SubCategory.EBUWFBB[subcategory] = v
		}

		subCategoryEBUHBBTotal := response.SubCategory.EBUHBB["total"].Count
		for subcategory, v := range response.SubCategory.EBUHBB {
			if subcategory == "total" || subCategoryEBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryEBUHBBTotal)
			response.SubCategory.EBUHBB[subcategory] = v
		}

		//serviceCategory
		srvcCategoryCBUHBBTotal := response.ServiceCategory.CBUHBB["total"].Count
		for srvccategory, v := range response.ServiceCategory.CBUHBB {
			if srvccategory == "total" || srvcCategoryCBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(srvcCategoryCBUHBBTotal)
			response.ServiceCategory.CBUHBB[srvccategory] = v
		}

		srvcCategoryEBUHBBTotal := response.ServiceCategory.EBUHBB["total"].Count
		for srvccategory, v := range response.ServiceCategory.EBUHBB {
			if srvccategory == "total" || srvcCategoryEBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(srvcCategoryEBUHBBTotal)
			response.ServiceCategory.EBUHBB[srvccategory] = v
		}

		go func(report *models.HbbUserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(user.UserId, report.Product, chanel,
				isQuarterReport, month, quarter, year, nil, report, impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("product", product).
					Err(errors.New(err.Details.(string))).Msg("ASalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("product", product).
					Msg("ASalesCountReport : cache update success")
			}
		}(response)

		return nil, response, nil
	}

	response := &models.UserSalesCountPerPlanDB{
		Product:    product,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	response.BusinessUnit.Consumer = &models.SaleCalculation{}
	response.BusinessUnit.Corporate = &models.SaleCalculation{}
	response.BusinessUnit.Government = &models.SaleCalculation{}
	response.SubCategory.CBU = map[string]*models.SaleCalculation{"total": {}}
	response.SubCategory.EBU = map[string]*models.SaleCalculation{"total": {}}
	response.Plans.CBU = map[string]int{"total": 0}
	response.Plans.EBU = map[string]int{"total": 0}

	for _, v := range allUserSales {
		if v.ProductGroup != product {
			continue
		}
		switch v.BusinessUnit {
		case "GOVERNMENT":
			response.BusinessUnit.Government.Count = response.BusinessUnit.Government.Count + v.Quantity
		case "CONSUMER":
			response.BusinessUnit.Consumer.Count = response.BusinessUnit.Consumer.Count + v.Quantity

			//subcategory
			if response.SubCategory.CBU[v.SubCategory] == nil {
				response.SubCategory.CBU[v.SubCategory] = &models.SaleCalculation{}
			}
			response.SubCategory.CBU["total"].Count = response.SubCategory.CBU["total"].Count + v.Quantity
			response.SubCategory.CBU[v.SubCategory].Count = response.SubCategory.CBU[v.SubCategory].Count + v.Quantity

			//plans
			response.Plans.CBU["total"] = response.Plans.CBU["total"] + v.Quantity
			response.Plans.CBU[v.CurrentPlan] = response.Plans.CBU[v.CurrentPlan] + v.Quantity
		case "CORPORATE":
			response.BusinessUnit.Corporate.Count = response.BusinessUnit.Corporate.Count + v.Quantity

			//subcategory
			if response.SubCategory.EBU[v.SubCategory] == nil {
				response.SubCategory.EBU[v.SubCategory] = &models.SaleCalculation{}
			}
			response.SubCategory.EBU["total"].Count = response.SubCategory.EBU["total"].Count + v.Quantity
			response.SubCategory.EBU[v.SubCategory].Count = response.SubCategory.EBU[v.SubCategory].Count + v.Quantity

			//plans
			response.Plans.EBU["total"] = response.Plans.EBU["total"] + v.Quantity
			response.Plans.EBU[v.CurrentPlan] = response.Plans.EBU[v.CurrentPlan] + v.Quantity
		default:
		}
	}

	// OTHER CALCULATION

	//businesssUnit
	sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
		float64(response.BusinessUnit.Consumer.Count))
	if sum != 0 {
		response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
		response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
		response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
	}

	//subCategory
	subCategoryCBUTotal := response.SubCategory.CBU["total"].Count
	for subcategory, v := range response.SubCategory.CBU {
		if subcategory == "total" || subCategoryCBUTotal == 0 {
			continue
		}
		v.Percentage = float64(v.Count) * 100 / float64(subCategoryCBUTotal)
		response.SubCategory.CBU[subcategory] = v
	}

	subCategoryEBUTotal := response.SubCategory.EBU["total"].Count
	for subcategory, v := range response.SubCategory.EBU {
		if subcategory == "total" || subCategoryEBUTotal == 0 {
			continue
		}
		v.Percentage = float64(v.Count) * 100 / float64(subCategoryEBUTotal)
		response.SubCategory.EBU[subcategory] = v
	}

	go func(report *models.UserSalesCountPerPlanDB) {
		err = impl.CacheService.InsertUserSaleCountPerPlanReport(user.UserId, report.Product, chanel, isQuarterReport,
			month, quarter, year, report, nil, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", user.UserId).Str("product", product).
				Err(errors.New(err.Details.(string))).Msg("ASalesCountReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("product", product).
				Msg("ASalesCountReport : cache update success")
		}
	}(response)

	return response, nil, nil
}

func (impl *UserImpl) updateOutletManagerUserSaleCountReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB,
	product string, chanel string, isQuarterReport bool, quarter, month, year int) (*models.UserSalesCountPerPlanDB,
	*models.HbbUserSalesCountPerPlanDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {
		if omngr == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").
				Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).Int("year", year).
				Msg("OMSalesCountReport : nil info")
			return nil, nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "updateOMUserSalesReport: invalid omngr details",
			}
		}
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
			Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).
			Int("year", year).Msg("OMSalesCountReport : invalid info")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "updateOMUserSalesReport: invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if omngr.Position != constant.OutletManagerAbb { //!(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader)
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("position", omngr.Position).
			Msg("OMSalesCountReport : invalid position")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details: "updateOMUserSalesReport: incorrect omngr position - userId:" + omngr.UserId +
				" position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
			Str("outlet", omngr.Outlet).Msg("OMSalesCountReport : invalid outlet info")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + omngr.Region +
				" outlet:" + omngr.Outlet,
			Details: "updateOMUserSalesReport: agent blank region or outlet - userId: " + omngr.UserId +
				" region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).
		Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).Int("year", year).
		Msg("OMSalesCountReport : info")

	userAgents, err := impl.DbService.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(errors.New(err.Details.(string))).
			Msg("OMSalesCountReport : fetch agents failed")
		return nil, nil, err
	}

	if product == constant.HBBCategory {
		response := &models.HbbUserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		response.BusinessUnit.Consumer = &models.SaleCalculation{}
		response.BusinessUnit.Corporate = &models.SaleCalculation{}
		response.BusinessUnit.Government = &models.SaleCalculation{}
		response.Category.CBUWFBB = &models.SaleCalculation{Percentage: 100}
		response.Category.CBUHBB = &models.SaleCalculation{Percentage: 100}
		response.Category.EBUWFBB = &models.SaleCalculation{Percentage: 100}
		response.Category.EBUHBB = &models.SaleCalculation{Percentage: 100}
		response.SubCategory.CBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.SubCategory.CBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.SubCategory.EBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.SubCategory.EBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.CBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.CBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.EBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.EBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.Plans.CBUWFBB = map[string]int{"total": 0}
		response.Plans.CBUHBB = map[string]int{"total": 0}
		response.Plans.EBUWFBB = map[string]int{"total": 0}
		response.Plans.EBUHBB = map[string]int{"total": 0}

		if len(userAgents) == 0 {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMSalesCountReport : zero agents")
			return nil, response, nil
		}
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).
			Msg("OMSalesCountReport : total agents")

		for _, eachAgent := range userAgents {
			var agentReport *models.HbbUserSalesCountPerPlanDB

			_, agentReport, err = impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product,
				chanel, isQuarterReport, month, quarter, year)
			if err != nil {
				logger.FileLogger.CtxWarn(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
					Msg("OMSalesCountReport : cache report failed")
				_, agentReport, err = impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, chanel, isQuarterReport,
					quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
						Err(err).Msg("OMSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(agentReport)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
					Str("record", string(rawReport)).Msg("OMSalesCountReport : agent log")

			} else if agentReport != nil {
				if !objecthelper.IsCurrentDay(agentReport.ModifiedAt) {
					_, agentReport, err = impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, chanel, isQuarterReport,
						quarter, month, year)
					if err != nil {
						logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
							Err(err).Msg("OMSalesCountReport : agent report failed")
						continue
					}
					rawReport, _ := json.Marshal(agentReport)
					logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
						Str("record", string(rawReport)).Msg("OMSalesCountReport : agent log")
				}
			}

			if agentReport == nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
					Err(errors.New("nil record")).Msg("OMSalesCountReport : agent report failed")
				continue
			}

			userSaleCount := agentReport

			//businessUnit
			response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
			response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
			response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

			//category
			response.Category.CBUWFBB.Count += userSaleCount.Category.CBUWFBB.Count
			response.Category.CBUHBB.Count += userSaleCount.Category.CBUHBB.Count
			response.Category.EBUWFBB.Count += userSaleCount.Category.EBUWFBB.Count
			response.Category.EBUHBB.Count += userSaleCount.Category.EBUHBB.Count

			//subCategory
			for subcategory, val := range userSaleCount.SubCategory.CBUWFBB {
				if subcategory == "total" {
					response.SubCategory.CBUWFBB["total"].Count = response.SubCategory.CBUWFBB["total"].Count + val.Count
					continue
				}
				if response.SubCategory.CBUWFBB[subcategory] == nil {
					response.SubCategory.CBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				response.SubCategory.CBUWFBB[subcategory].Count += val.Count
			}

			for subcategory, val := range userSaleCount.SubCategory.CBUHBB {
				if subcategory == "total" {
					response.SubCategory.CBUHBB["total"].Count = response.SubCategory.CBUHBB["total"].Count + val.Count
					continue
				}
				if response.SubCategory.CBUHBB[subcategory] == nil {
					response.SubCategory.CBUHBB[subcategory] = &models.SaleCalculation{}
				}
				response.SubCategory.CBUHBB[subcategory].Count += val.Count
			}

			for subcategory, val := range userSaleCount.SubCategory.EBUWFBB {
				if subcategory == "total" {
					response.SubCategory.EBUWFBB["total"].Count = response.SubCategory.EBUWFBB["total"].Count + val.Count
					continue
				}
				if response.SubCategory.EBUWFBB[subcategory] == nil {
					response.SubCategory.EBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				response.SubCategory.EBUWFBB[subcategory].Count += val.Count
			}

			for subcategory, val := range userSaleCount.SubCategory.EBUHBB {
				if subcategory == "total" {
					response.SubCategory.EBUHBB["total"].Count = response.SubCategory.EBUHBB["total"].Count + val.Count
					continue
				}
				if response.SubCategory.EBUHBB[subcategory] == nil {
					response.SubCategory.EBUHBB[subcategory] = &models.SaleCalculation{}
				}
				response.SubCategory.EBUHBB[subcategory].Count += val.Count
			}

			//serviceCategory
			for srvccategory, val := range userSaleCount.ServiceCategory.CBUWFBB {
				if srvccategory == "total" {
					response.ServiceCategory.CBUWFBB["total"].Count = response.ServiceCategory.CBUWFBB["total"].Count + val.Count
				}
			}

			for srvccategory, val := range userSaleCount.ServiceCategory.CBUHBB {
				if srvccategory == "total" {
					response.ServiceCategory.CBUHBB["total"].Count = response.ServiceCategory.CBUHBB["total"].Count + val.Count
					continue
				}
				if response.ServiceCategory.CBUHBB[srvccategory] == nil {
					response.ServiceCategory.CBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				response.ServiceCategory.CBUHBB[srvccategory].Count += val.Count
			}

			for srvccategory, val := range userSaleCount.ServiceCategory.EBUWFBB {
				if srvccategory == "total" {
					response.ServiceCategory.EBUWFBB["total"].Count = response.ServiceCategory.EBUWFBB["total"].Count + val.Count
				}
			}

			for srvccategory, val := range userSaleCount.ServiceCategory.EBUHBB {
				if srvccategory == "total" {
					response.ServiceCategory.EBUHBB["total"].Count = response.ServiceCategory.EBUHBB["total"].Count + val.Count
					continue
				}
				if response.ServiceCategory.EBUHBB[srvccategory] == nil {
					response.ServiceCategory.EBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				response.ServiceCategory.EBUHBB[srvccategory].Count += val.Count
			}

			//plans
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

		// HBB CALCULATION

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//category

		//subCategory
		subCategoryCBUWFBBTotal := response.SubCategory.CBUWFBB["total"].Count
		for subcategory, v := range response.SubCategory.CBUWFBB {
			if subcategory == "total" || subCategoryCBUWFBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryCBUWFBBTotal)
			response.SubCategory.CBUWFBB[subcategory] = v
		}

		subCategoryCBUHBBTotal := response.SubCategory.CBUHBB["total"].Count
		for subcategory, v := range response.SubCategory.CBUHBB {
			if subcategory == "total" || subCategoryCBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryCBUHBBTotal)
			response.SubCategory.CBUHBB[subcategory] = v
		}

		subCategoryEBUWFBBTotal := response.SubCategory.EBUWFBB["total"].Count
		for subcategory, v := range response.SubCategory.EBUWFBB {
			if subcategory == "total" || subCategoryEBUWFBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryEBUWFBBTotal)
			response.SubCategory.EBUWFBB[subcategory] = v
		}

		subCategoryEBUHBBTotal := response.SubCategory.EBUHBB["total"].Count
		for subcategory, v := range response.SubCategory.EBUHBB {
			if subcategory == "total" || subCategoryEBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryEBUHBBTotal)
			response.SubCategory.EBUHBB[subcategory] = v
		}

		//serviceCategory
		srvcCategoryCBUHBBTotal := response.ServiceCategory.CBUHBB["total"].Count
		for srvccategory, v := range response.ServiceCategory.CBUHBB {
			if srvccategory == "total" || srvcCategoryCBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(srvcCategoryCBUHBBTotal)
			response.ServiceCategory.CBUHBB[srvccategory] = v
		}

		srvcCategoryEBUHBBTotal := response.ServiceCategory.EBUHBB["total"].Count
		for srvccategory, v := range response.ServiceCategory.EBUHBB {
			if srvccategory == "total" || srvcCategoryEBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(srvcCategoryEBUHBBTotal)
			response.ServiceCategory.EBUHBB[srvccategory] = v
		}

		go func(report *models.HbbUserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(omngr.UserId, report.Product, chanel,
				isQuarterReport, month, quarter, year, nil, report, impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("product", product).
					Err(errors.New(err.Details.(string))).Msg("OMSalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("product", product).
					Msg("OMSalesCountReport : cache update success")
			}
		}(response)

		return nil, response, nil
	}

	response := &models.UserSalesCountPerPlanDB{
		Product:    product,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	response.BusinessUnit.Consumer = &models.SaleCalculation{}
	response.BusinessUnit.Corporate = &models.SaleCalculation{}
	response.BusinessUnit.Government = &models.SaleCalculation{}
	response.SubCategory.CBU = map[string]*models.SaleCalculation{"total": {}}
	response.SubCategory.EBU = map[string]*models.SaleCalculation{"total": {}}
	response.Plans.CBU = map[string]int{"total": 0}
	response.Plans.EBU = map[string]int{"total": 0}

	if len(userAgents) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMSalesCountReport : zero agents")
		return response, nil, nil
	}
	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).
		Msg("OMSalesCountReport : total agents")

	for _, eachAgent := range userAgents {
		var agentReport *models.UserSalesCountPerPlanDB

		agentReport, _, err = impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, chanel,
			isQuarterReport, month, quarter, year)
		if err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
				Msg("OMSalesCountReport : cache report failed")
			agentReport, _, err = impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, chanel, isQuarterReport,
				quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
					Err(err).Msg("OMSalesCountReport : agent report failed")
				continue
			}
			rawReport, _ := json.Marshal(agentReport)
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
				Str("record", string(rawReport)).Msg("OMSalesCountReport : agent log")

		} else if agentReport != nil {
			if !objecthelper.IsCurrentDay(agentReport.ModifiedAt) {
				agentReport, _, err = impl.updateAgentUserSalesCountReport(ctx, eachAgent, product, chanel, isQuarterReport,
					quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
						Err(err).Msg("OMSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(agentReport)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
					Str("record", string(rawReport)).Msg("OMSalesCountReport : agent log")
			}
		}

		if agentReport == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
				Err(errors.New("nil record")).Msg("OMSalesCountReport : agent report failed")
			continue
		}

		userSaleCount := agentReport

		//businessUnit
		response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
		response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
		response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

		//subcategory
		for subcategory, val := range userSaleCount.SubCategory.CBU {
			if subcategory == "total" {
				response.SubCategory.CBU["total"].Count = response.SubCategory.CBU["total"].Count + val.Count
				continue
			}
			if response.SubCategory.CBU[subcategory] == nil {
				response.SubCategory.CBU[subcategory] = &models.SaleCalculation{}
			}
			response.SubCategory.CBU[subcategory].Count += val.Count
		}

		for subcategory, val := range userSaleCount.SubCategory.EBU {
			if subcategory == "total" {
				response.SubCategory.EBU["total"].Count = response.SubCategory.EBU["total"].Count + val.Count
				continue
			}
			if response.SubCategory.EBU[subcategory] == nil {
				response.SubCategory.EBU[subcategory] = &models.SaleCalculation{}
			}
			response.SubCategory.EBU[subcategory].Count += val.Count
		}

		//plan
		for plan, count := range userSaleCount.Plans.CBU {
			response.Plans.CBU[plan] += count
		}
		for plan, count := range userSaleCount.Plans.EBU {
			response.Plans.EBU[plan] += count
		}
	}

	// OTHER CALCULATION

	//businesssUnit
	sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
		float64(response.BusinessUnit.Consumer.Count))
	if sum != 0 {
		response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
		response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
		response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
	}

	//category

	//subCategory
	subCategoryCBUTotal := response.SubCategory.CBU["total"].Count
	for subcategory, v := range response.SubCategory.CBU {
		if subcategory == "total" || subCategoryCBUTotal == 0 {
			continue
		}
		v.Percentage = float64(v.Count) * 100 / float64(subCategoryCBUTotal)
		response.SubCategory.CBU[subcategory] = v
	}

	subCategoryEBUTotal := response.SubCategory.EBU["total"].Count
	for subcategory, v := range response.SubCategory.EBU {
		if subcategory == "total" || subCategoryEBUTotal == 0 {
			continue
		}
		v.Percentage = float64(v.Count) * 100 / float64(subCategoryEBUTotal)
		response.SubCategory.EBU[subcategory] = v
	}

	go func(report *models.UserSalesCountPerPlanDB) {
		err = impl.CacheService.InsertUserSaleCountPerPlanReport(omngr.UserId, report.Product, chanel,
			isQuarterReport, month, quarter, year, report, nil, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("product", product).
				Err(errors.New(err.Details.(string))).Msg("OMSalesCountReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("product", product).
				Msg("OMSalesCountReport : cache update success")
		}
	}(response)

	return response, nil, nil
}

func (impl *UserImpl) updateManagerUserSalesCountReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB, product string,
	chanel string, isQuarterReport bool, quarter, month, year int) (*models.UserSalesCountPerPlanDB,
	*models.HbbUserSalesCountPerPlanDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if mngr == nil || mngr.UserId == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).
			Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
			Int("month", month).Int("year", year).Msg("MngrSalesCountReport : invalid info")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "updateMngrUserSalesCountReport: invalid mngr details",
		}
	}

	position, ok := constant.ManagerMap[mngr.Position]
	if !ok {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).
			Msg("MngrSalesCountReport : invalid position")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details: "updateMngrUserSalesCountReport: incorrect mngr position - userId:" + mngr.UserId +
				" position:" + mngr.Position,
		}
	}
	mngr.Position = position

	rawDetails, _ := json.Marshal(mngr)
	logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("user", string(rawDetails)).
		Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("MngrSalesCountReport : info")

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManagerAbb { // || mngr.Position == constant.TeamLeader
		report, hbbReport, err := impl.updateOutletManagerUserSaleCountReport(ctx, mngr, product, chanel, isQuarterReport,
			quarter, month, year)
		if err != nil {
			return nil, nil, err
		}

		return report, hbbReport, nil
	}

	if mngr.Position != constant.VicePresidentAbb && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).
			Msg("MngrSalesCountReport : user's manager info missing")
		return nil, nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    mngr.UserId + ": INVALID MANAGER'S MANAGER DETAILS",
			Details:    mngr.UserId + ": invalid manager's manager details",
		}
	}

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(mngr.Position, mngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).
			Err(errors.New(err.Details.(string))).Msg("MngrSalesCountReport : fetch om list failed")
		return nil, nil, err
	}

	if product == constant.HBBCategory {
		response := &models.HbbUserSalesCountPerPlanDB{
			Product:    product,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		}
		response.BusinessUnit.Consumer = &models.SaleCalculation{}
		response.BusinessUnit.Corporate = &models.SaleCalculation{}
		response.BusinessUnit.Government = &models.SaleCalculation{}
		response.Category.CBUWFBB = &models.SaleCalculation{Percentage: 100}
		response.Category.CBUHBB = &models.SaleCalculation{Percentage: 100}
		response.Category.EBUWFBB = &models.SaleCalculation{Percentage: 100}
		response.Category.EBUHBB = &models.SaleCalculation{Percentage: 100}
		response.SubCategory.CBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.SubCategory.CBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.SubCategory.EBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.SubCategory.EBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.CBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.CBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.EBUWFBB = map[string]*models.SaleCalculation{"total": {}}
		response.ServiceCategory.EBUHBB = map[string]*models.SaleCalculation{"total": {}}
		response.Plans.CBUWFBB = map[string]int{"total": 0}
		response.Plans.CBUHBB = map[string]int{"total": 0}
		response.Plans.EBUWFBB = map[string]int{"total": 0}
		response.Plans.EBUHBB = map[string]int{"total": 0}

		if len(allDistinctOM) == 0 {
			logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrSalesCountReport : zero om")
			return nil, response, nil
		}
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Int("omCount", len(allDistinctOM)).
			Msg("MngrSalesCountReport : total om")

		for _, eachAgent := range allDistinctOM {
			var agentReport *models.HbbUserSalesCountPerPlanDB

			_, agentReport, err = impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, chanel,
				isQuarterReport, month, quarter, year)
			if err != nil {
				logger.FileLogger.CtxWarn(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
					Msg("MngrSalesCountReport : cache report failed")
				_, agentReport, err = impl.updateOutletManagerUserSaleCountReport(ctx, eachAgent, product, chanel,
					isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
						Err(err).Msg("MngrSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(agentReport)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
					Str("record", string(rawReport)).Msg("MngrSalesCountReport : agent log")

			} else if agentReport != nil {
				if !objecthelper.IsCurrentDay(agentReport.ModifiedAt) {
					_, agentReport, err = impl.updateOutletManagerUserSaleCountReport(ctx, eachAgent, product, chanel,
						isQuarterReport, quarter, month, year)
					if err != nil {
						logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
							Err(err).Msg("MngrSalesCountReport : agent report failed")
						continue
					}
					rawReport, _ := json.Marshal(agentReport)
					logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
						Str("record", string(rawReport)).Msg("MngrSalesCountReport : agent log")
				}
			}

			if agentReport == nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
					Err(errors.New("nil record")).Msg("MngrSalesCountReport : om report failed")
				continue
			}

			userSaleCount := agentReport

			//businessUnit
			response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
			response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
			response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

			//category
			response.Category.CBUWFBB.Count += userSaleCount.Category.CBUWFBB.Count
			response.Category.CBUHBB.Count += userSaleCount.Category.CBUHBB.Count
			response.Category.EBUWFBB.Count += userSaleCount.Category.EBUWFBB.Count
			response.Category.EBUHBB.Count += userSaleCount.Category.EBUHBB.Count

			//subCategory
			for subcategory, val := range userSaleCount.SubCategory.CBUWFBB {
				if subcategory == "total" {
					response.SubCategory.CBUWFBB["total"].Count = response.SubCategory.CBUWFBB["total"].Count + val.Count
					continue
				}
				if response.SubCategory.CBUWFBB[subcategory] == nil {
					response.SubCategory.CBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				response.SubCategory.CBUWFBB[subcategory].Count += val.Count
			}

			for subcategory, val := range userSaleCount.SubCategory.CBUHBB {
				if subcategory == "total" {
					response.SubCategory.CBUHBB["total"].Count = response.SubCategory.CBUHBB["total"].Count + val.Count
					continue
				}
				if response.SubCategory.CBUHBB[subcategory] == nil {
					response.SubCategory.CBUHBB[subcategory] = &models.SaleCalculation{}
				}
				response.SubCategory.CBUHBB[subcategory].Count += val.Count
			}

			for subcategory, val := range userSaleCount.SubCategory.EBUWFBB {
				if subcategory == "total" {
					response.SubCategory.EBUWFBB["total"].Count = response.SubCategory.EBUWFBB["total"].Count + val.Count
					continue
				}
				if response.SubCategory.EBUWFBB[subcategory] == nil {
					response.SubCategory.EBUWFBB[subcategory] = &models.SaleCalculation{}
				}
				response.SubCategory.EBUWFBB[subcategory].Count += val.Count
			}

			for subcategory, val := range userSaleCount.SubCategory.EBUHBB {
				if subcategory == "total" {
					response.SubCategory.EBUHBB["total"].Count = response.SubCategory.EBUHBB["total"].Count + val.Count
				}
				if response.SubCategory.EBUHBB[subcategory] == nil {
					response.SubCategory.EBUHBB[subcategory] = &models.SaleCalculation{}
				}
				response.SubCategory.EBUHBB[subcategory].Count += val.Count
			}

			//serviceCategory
			for srvccategory, val := range userSaleCount.ServiceCategory.CBUWFBB {
				if srvccategory == "total" {
					response.ServiceCategory.CBUWFBB["total"].Count = response.ServiceCategory.CBUWFBB["total"].Count + val.Count
				}
			}

			for srvccategory, val := range userSaleCount.ServiceCategory.CBUHBB {
				if srvccategory == "total" {
					response.ServiceCategory.CBUHBB["total"].Count = response.ServiceCategory.CBUHBB["total"].Count + val.Count
					continue
				}
				if response.ServiceCategory.CBUHBB[srvccategory] == nil {
					response.ServiceCategory.CBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				response.ServiceCategory.CBUHBB[srvccategory].Count += val.Count
			}

			for srvccategory, val := range userSaleCount.ServiceCategory.EBUWFBB {
				if srvccategory == "total" {
					response.ServiceCategory.EBUWFBB["total"].Count = response.ServiceCategory.EBUWFBB["total"].Count + val.Count
				}
			}

			for srvccategory, val := range userSaleCount.ServiceCategory.EBUHBB {
				if srvccategory == "total" {
					response.ServiceCategory.EBUHBB["total"].Count = response.ServiceCategory.EBUHBB["total"].Count + val.Count
					continue
				}
				if response.ServiceCategory.EBUHBB[srvccategory] == nil {
					response.ServiceCategory.EBUHBB[srvccategory] = &models.SaleCalculation{}
				}
				response.ServiceCategory.EBUHBB[srvccategory].Count += val.Count
			}

			//plans
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

		// HBB CALCULATION

		//businesssUnit
		sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
			float64(response.BusinessUnit.Consumer.Count))
		if sum != 0 {
			response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
			response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
			response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
		}

		//category

		//subCategory
		subCategoryCBUWFBBTotal := response.SubCategory.CBUWFBB["total"].Count
		for subcategory, v := range response.SubCategory.CBUWFBB {
			if subcategory == "total" || subCategoryCBUWFBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryCBUWFBBTotal)
			response.SubCategory.CBUWFBB[subcategory] = v
		}

		subCategoryCBUHBBTotal := response.SubCategory.CBUHBB["total"].Count
		for subcategory, v := range response.SubCategory.CBUHBB {
			if subcategory == "total" || subCategoryCBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryCBUHBBTotal)
			response.SubCategory.CBUHBB[subcategory] = v
		}

		subCategoryEBUWFBBTotal := response.SubCategory.EBUWFBB["total"].Count
		for subcategory, v := range response.SubCategory.EBUWFBB {
			if subcategory == "total" || subCategoryEBUWFBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryEBUWFBBTotal)
			response.SubCategory.EBUWFBB[subcategory] = v
		}

		subCategoryEBUHBBTotal := response.SubCategory.EBUHBB["total"].Count
		for subcategory, v := range response.SubCategory.EBUHBB {
			if subcategory == "total" || subCategoryEBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(subCategoryEBUHBBTotal)
			response.SubCategory.EBUHBB[subcategory] = v
		}

		//serviceCategory
		srvcCategoryCBUHBBTotal := response.ServiceCategory.CBUHBB["total"].Count
		for srvccategory, v := range response.ServiceCategory.CBUHBB {
			if srvccategory == "total" || srvcCategoryCBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(srvcCategoryCBUHBBTotal)
			response.ServiceCategory.CBUHBB[srvccategory] = v
		}

		srvcCategoryEBUHBBTotal := response.ServiceCategory.EBUHBB["total"].Count
		for srvccategory, v := range response.ServiceCategory.EBUHBB {
			if srvccategory == "total" || srvcCategoryEBUHBBTotal == 0 {
				continue
			}
			v.Percentage = float64(v.Count) * 100 / float64(srvcCategoryEBUHBBTotal)
			response.ServiceCategory.EBUHBB[srvccategory] = v
		}

		go func(report *models.HbbUserSalesCountPerPlanDB) {
			err = impl.CacheService.InsertUserSaleCountPerPlanReport(mngr.UserId, report.Product, chanel,
				isQuarterReport, month, quarter, year, nil, report, impl.Config.ReportCacheTTL)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("product", product).
					Err(errors.New(err.Details.(string))).Msg("MngrSalesCountReport : force report update")
			} else {
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("product", product).
					Msg("MngrSalesCountReport : cache update success")
			}
		}(response)

		return nil, response, nil
	}

	response := &models.UserSalesCountPerPlanDB{
		Product:    product,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
	response.BusinessUnit.Consumer = &models.SaleCalculation{}
	response.BusinessUnit.Corporate = &models.SaleCalculation{}
	response.BusinessUnit.Government = &models.SaleCalculation{}
	response.SubCategory.CBU = map[string]*models.SaleCalculation{"total": {}}
	response.SubCategory.EBU = map[string]*models.SaleCalculation{"total": {}}
	response.Plans.CBU = map[string]int{"total": 0}
	response.Plans.EBU = map[string]int{"total": 0}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrSalesCountReport : zero agents")
		return response, nil, nil
	}
	logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Int("agentCount", len(allDistinctOM)).
		Msg("MngrSalesCountReport : total agents")

	for _, eachAgent := range allDistinctOM {
		var agentReport *models.UserSalesCountPerPlanDB

		agentReport, _, err = impl.CacheService.GetUserSaleCountPerPlanReport(eachAgent.UserId, product, chanel,
			isQuarterReport, month, quarter, year)
		if err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
				Msg("MngrSalesCountReport : cache report failed")
			agentReport, _, err = impl.updateOutletManagerUserSaleCountReport(ctx, eachAgent, product, chanel, isQuarterReport,
				quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
					Err(err).Msg("MngrSalesCountReport : agent report failed")
				continue
			}
			rawReport, _ := json.Marshal(agentReport)
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
				Str("record", string(rawReport)).Msg("MngrSalesCountReport : agent log")

		} else if agentReport != nil {
			if !objecthelper.IsCurrentDay(agentReport.ModifiedAt) {
				agentReport, _, err = impl.updateOutletManagerUserSaleCountReport(ctx, eachAgent, product, chanel, isQuarterReport,
					quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
						Err(err).Msg("MngrSalesCountReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(agentReport)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
					Str("record", string(rawReport)).Msg("MngrSalesCountReport : agent log")
			}
		}

		if agentReport == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
				Err(errors.New("nil record")).Msg("MngrSalesCountReport : om report failed")
			continue
		}

		userSaleCount := agentReport

		//businessUnit
		response.BusinessUnit.Government.Count += userSaleCount.BusinessUnit.Government.Count
		response.BusinessUnit.Corporate.Count += userSaleCount.BusinessUnit.Corporate.Count
		response.BusinessUnit.Consumer.Count += userSaleCount.BusinessUnit.Consumer.Count

		//category
		for subcategory, val := range userSaleCount.SubCategory.CBU {
			if subcategory == "total" {
				response.SubCategory.CBU["total"].Count = response.SubCategory.CBU["total"].Count + val.Count
				continue
			}
			if response.SubCategory.CBU[subcategory] == nil {
				response.SubCategory.CBU[subcategory] = &models.SaleCalculation{}
			}
			response.SubCategory.CBU[subcategory].Count += val.Count
		}

		for subcategory, val := range userSaleCount.SubCategory.EBU {
			if subcategory == "total" {
				response.SubCategory.EBU["total"].Count = response.SubCategory.EBU["total"].Count + val.Count
				continue
			}
			if response.SubCategory.EBU[subcategory] == nil {
				response.SubCategory.EBU[subcategory] = &models.SaleCalculation{}
			}
			response.SubCategory.EBU[subcategory].Count += val.Count
		}

		//plan
		for plan, count := range userSaleCount.Plans.CBU {
			response.Plans.CBU[plan] += count
		}
		for plan, count := range userSaleCount.Plans.EBU {
			response.Plans.EBU[plan] += count
		}
	}

	// OTHER CALCULATION

	//businesssUnit
	sum := (float64(response.BusinessUnit.Government.Count) + float64(response.BusinessUnit.Corporate.Count) +
		float64(response.BusinessUnit.Consumer.Count))
	if sum != 0 {
		response.BusinessUnit.Government.Percentage = float64(response.BusinessUnit.Government.Count) * 100 / sum
		response.BusinessUnit.Corporate.Percentage = float64(response.BusinessUnit.Corporate.Count) * 100 / sum
		response.BusinessUnit.Consumer.Percentage = float64(response.BusinessUnit.Consumer.Count) * 100 / sum
	}

	//category

	//subCategory
	subCategoryCBUTotal := response.SubCategory.CBU["total"].Count
	for subcategory, v := range response.SubCategory.CBU {
		if subcategory == "total" || subCategoryCBUTotal == 0 {
			continue
		}

		v.Percentage = float64(v.Count) * 100 / float64(subCategoryCBUTotal)
		response.SubCategory.CBU[subcategory] = v
	}

	subCategoryEBUTotal := response.SubCategory.EBU["total"].Count
	for subcategory, v := range response.SubCategory.EBU {
		if subcategory == "total" || subCategoryEBUTotal == 0 {
			continue
		}
		v.Percentage = float64(v.Count) * 100 / float64(subCategoryEBUTotal)
		response.SubCategory.EBU[subcategory] = v
	}

	go func(report *models.UserSalesCountPerPlanDB) {
		err = impl.CacheService.InsertUserSaleCountPerPlanReport(mngr.UserId, report.Product, chanel,
			isQuarterReport, month, quarter, year, report, nil, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("product", product).
				Err(errors.New(err.Details.(string))).
				Msg("MngrSalesCountReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("product", product).
				Msg("MngrSalesCountReport : cache update success")
		}
	}(response)

	return response, nil, nil
}
