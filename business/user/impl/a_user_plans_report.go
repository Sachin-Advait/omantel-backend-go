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

func (impl *UserImpl) updateAgentUserPlanPerCategoryReport(ctx *gin.Context, user *models.UserBasicDetailsDB, product string,
	chanel string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).
				Int("quarter", quarter).Int("month", month).Int("year", year).Msg("APlanPerCategoryReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "updateAgentUserPlanPerCategoryReport: invalid agent details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("APlanPerCategoryReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": INVALID AGENT DETAILS",
			Details:    "updateAgentUserPlanPerCategoryReport: invalid agent details - userId:" + user.UserId,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("APlanPerCategoryReport : info")

	if !(user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb) {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("position", user.Position).Msg("APlanPerCategoryReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details:    "APlanPerCategoryReport: incorrect agent position - userId:" + user.UserId + " position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).Str("outlet", user.Outlet).
			Msg("APlanPerCategoryReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + user.Region + " outlet:" + user.Outlet,
			Details:    "APlanPerCategoryReport: agent blank region or outlet - userId: " + user.UserId + " region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	allUserSales, err := impl.DbService.FetchUserSalesCountPerPlan(user.UserId, isQuarterReport, quarter, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserPlanPerCategoryReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("APlanPerCategoryReport : fetch sales failed")
		return nil, err
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("saleGroupCount", len(allUserSales)).
		Msg("APlanPerCategoryReport : total agents and salesGroup")

	response := make(map[string]interface{})
	response[product] = constant.HBBCategory
	response["CreatedAt"] = time.Now()
	response["ModifiedAt"] = time.Now()
	response["GRAND TOTAL"] = make(map[string]int)

	for _, v := range allUserSales {
		if v.ProductGroup != constant.HBBCategory {
			continue
		}

		switch {
		case product == "CBU" && v.BusinessUnit == "CONSUMER":
		case product == "EBU" && v.BusinessUnit == "CORPORATE":
		default:
			continue
		}

		//only consider HBB (i.e CBU-HBB and EBU-HBB) and ignore WFBB (i.e CBU-WFBB and EBU-WFBB)
		if !(v.LineType == "ILL SERVICE" || v.LineType == "MPLS SERVICE" || v.LineType == "P2P LEASED LINE SERVICE" ||
			v.LineType == "PSTN VOICE SERVICE") {
			continue
		}

		if response[v.SubCategory] == nil {
			response[v.SubCategory] = make(map[string]interface{})
		}
		plans := response[v.SubCategory].(map[string]interface{})

		if plans[v.CurrentPlan] == nil {
			plans[v.CurrentPlan] = make(map[string]int)
		}
		srvccategory := plans[v.CurrentPlan].(map[string]int)
		srvccategory[v.FinalServiceCategory] = srvccategory[v.FinalServiceCategory] + v.Quantity

		//total
		srvccategory["total"] = srvccategory["total"] + v.Quantity

		totalSrvcCategory := response["GRAND TOTAL"].(map[string]int)
		totalSrvcCategory[v.FinalServiceCategory] = totalSrvcCategory[v.FinalServiceCategory] + v.Quantity
		totalSrvcCategory["total"] = totalSrvcCategory["total"] + v.Quantity
	}

	go func(product string, report map[string]interface{}) {
		err = impl.CacheService.InsertUserPlanCategoryReport(user.UserId, product, chanel, isQuarterReport, month, quarter, year, report,
			impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
				Msg("APlanPerCategoryReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("product", product).Msg("APlanPerCategoryReport : cache update success")
		}
	}(product, response)

	return response, nil
}

func (impl *UserImpl) updateOutletManagerUserPlanPerCategoryReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB,
	product string, chanel string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {

		if omngr == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
				Int("month", month).Int("year", year).Msg("OMUserPlanPerCategoryReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "updateOMUserPlanPerCategoryReport: invalid omngr details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMUserPlanPerCategoryReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "updateOMUserPlanPerCategoryReport: invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if omngr.Position != constant.OutletManagerAbb { //!(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader)
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("position", omngr.Position).Msg("OMUserPlanPerCategoryReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details:    "updateOMUserPlanPerCategoryReport: incorrect omngr position - userId:" + omngr.UserId + " position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Msg("OMUserPlanPerCategoryReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + omngr.Region + " outlet:" + omngr.Outlet,
			Details:    "updateOMUserPlanPerCategoryReport: agent blank region or outlet - userId: " + omngr.UserId + " region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("OMUserPlanPerCategoryReport : info")

	userAgents, err := impl.DbService.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(errors.New(err.Details.(string))).Msg("OMUserPlanPerCategoryReport : fetch agents failed")
		return nil, err
	}

	if len(userAgents) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMUserPlanPerCategoryReport : zero agents")
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).Msg("OMUserPlanPerCategoryReport : total agents")
	}

	response := make(map[string]interface{})
	response[product] = constant.HBBCategory
	response["CreatedAt"] = time.Now()
	response["ModifiedAt"] = time.Now()
	response["GRAND TOTAL"] = make(map[string]int)

	for _, eachAgent := range userAgents {
		agentReport, err := impl.CacheService.GetUserPlanCategoryReport(eachAgent.UserId, product, chanel,
			isQuarterReport, month, quarter, year)
		if err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
				Msg("OMUserPlanPerCategoryReport : cache report failed")
			resp, err := impl.updateAgentUserPlanPerCategoryReport(ctx, eachAgent, product, chanel, isQuarterReport,
				quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Err(err).
					Msg("OMUserPlanPerCategoryReport : agent report failed")
				continue
			}
			rawReport, _ := json.Marshal(resp)
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
				Str("record", string(rawReport)).Msg("OMUserPlanPerCategoryReport : agent log")
			agentReport = resp
		} else if agentReport != nil {
			report, _ := agentReport.(map[string]interface{})
			modifiedAt, _ := report["ModifiedAt"].(time.Time)
			if objecthelper.IsCurrentDay(modifiedAt) {
				resp, err := impl.updateAgentUserPlanPerCategoryReport(ctx, eachAgent, product, chanel, isQuarterReport,
					quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
						Err(err).Msg("OMUserPlanPerCategoryReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("OMUserPlanPerCategoryReport : agent log")
				agentReport = resp
			}
		}
		if agentReport == nil {
			continue
		}
		userPlanCount := agentReport.(map[string]interface{})

		for subcategoryName, subcategoryMap := range userPlanCount {

			if subcategoryName == "GRAND TOTAL" {
				agentGrandTotal, ok := subcategoryMap.(map[string]int)
				if !ok {
					continue
				}
				for name, qty := range agentGrandTotal {
					responseGrandTotalSrvcCategory := response["GRAND TOTAL"].(map[string]int)
					responseGrandTotalSrvcCategory[name] = responseGrandTotalSrvcCategory[name] + qty
					response["GRAND TOTAL"] = responseGrandTotalSrvcCategory
				}
				continue
			}

			if response[subcategoryName] == nil {
				response[subcategoryName] = make(map[string]interface{})
			}

			planMap, _ := subcategoryMap.(map[string]interface{})
			responsePlanMap, _ := response[subcategoryName].(map[string]interface{})

			for planName, plan := range planMap {
				if responsePlanMap[planName] == nil {
					responsePlanMap[planName] = make(map[string]int)
				}

				srvccategoryMap, _ := plan.(map[string]int)
				responseSrvccategoryMap, _ := responsePlanMap[planName].(map[string]int)

				for srvcCategoryName, qty := range srvccategoryMap {
					responseSrvccategoryMap[srvcCategoryName] = responseSrvccategoryMap[srvcCategoryName] + qty
				}
			}
		}
	}

	go func(product string, report map[string]interface{}) {
		err = impl.CacheService.InsertUserPlanCategoryReport(omngr.UserId, product, chanel, isQuarterReport, month, quarter, year, report,
			impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
				Msg("OMUserPlanPerCategoryReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("product", product).Msg("OMUserPlanPerCategoryReport : cache update success")
		}
	}(product, response)

	return response, nil
}

func (impl *UserImpl) updateManagerUserPlanPerCategoryReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB, product string,
	chanel string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if mngr == nil || mngr.UserId == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).
			Int("quarter", quarter).Int("month", month).Int("year", year).Msg("MngrUserPlanPerCategoryReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "updateMngrUserPlanPerCategoryReport: invalid mngr details",
		}
	}

	position, ok := constant.ManagerMap[mngr.Position]
	if !ok {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrUserPlanPerCategoryReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details:    "updateMngrUserPlanPerCategoryReport: incorrect mngr position - userId:" + mngr.UserId + " position:" + mngr.Position,
		}
	}
	mngr.Position = position

	rawDetails, _ := json.Marshal(mngr)
	logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("user", string(rawDetails)).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("MngrUserPlanPerCategoryReport : info")

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManagerAbb { // || mngr.Position == constant.TeamLeader
		report, err := impl.updateOutletManagerUserPlanPerCategoryReport(ctx, mngr, product, chanel, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	if mngr.Position != constant.VicePresidentAbb && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Msg("MngrUserPlanPerCategoryReport : user's manager info missing")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    mngr.UserId + ": INVALID MANAGER'S MANAGER DETAILS",
			Details:    mngr.UserId + ": invalid manager's manager details",
		}
	}

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(mngr.Position, mngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).Err(errors.New(err.Details.(string))).Msg("MngrUserPlanPerCategoryReport : fetch om list failed")
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrUserPlanPerCategoryReport : zero om")
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Int("omCount", len(allDistinctOM)).Msg("MngrUserPlanPerCategoryReport : total om")
	}

	response := make(map[string]interface{})
	response[product] = constant.HBBCategory
	response["CreatedAt"] = time.Now()
	response["ModifiedAt"] = time.Now()
	response["GRAND TOTAL"] = make(map[string]int)

	for _, eachAgent := range allDistinctOM {
		agentReport, err := impl.CacheService.GetUserPlanCategoryReport(eachAgent.UserId, product, chanel,
			isQuarterReport, month, quarter, year)
		if err != nil {
			logger.FileLogger.CtxWarn(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
				Msg("MngrUserPlanPerCategoryReport : cache report failed")
			resp, err := impl.updateOutletManagerUserPlanPerCategoryReport(ctx, eachAgent, product, chanel, isQuarterReport,
				quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).
					Msg("MngrUserPlanPerCategoryReport : agent report failed")
				continue
			}
			rawReport, _ := json.Marshal(resp)
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
				Str("record", string(rawReport)).Msg("MngrUserPlanPerCategoryReport : agent log")
			agentReport = resp
		} else if agentReport != nil {
			report, _ := agentReport.(map[string]interface{})
			modifiedAt, _ := report["ModifiedAt"].(time.Time)
			if objecthelper.IsCurrentDay(modifiedAt) {
				resp, err := impl.updateOutletManagerUserPlanPerCategoryReport(ctx, eachAgent, product, chanel, isQuarterReport, quarter, month, year)
				if err != nil {
					logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Err(err).Msg("MngrUserPlanPerCategoryReport : agent report failed")
					continue
				}
				rawReport, _ := json.Marshal(resp)
				logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).Str("record", string(rawReport)).Msg("MngrUserPlanPerCategoryReport : agent log")
				agentReport = resp
			}
		}
		if agentReport == nil {
			continue
		}
		userPlanCount := agentReport.(map[string]interface{})

		for subcategoryName, subcategoryMap := range userPlanCount {

			if subcategoryName == "GRAND TOTAL" {
				agentGrandTotal, ok := subcategoryMap.(map[string]int)
				if !ok {
					continue
				}
				for name, qty := range agentGrandTotal {
					responseGrandTotalSrvcCategory := response["GRAND TOTAL"].(map[string]int)
					responseGrandTotalSrvcCategory[name] = responseGrandTotalSrvcCategory[name] + qty
					response["GRAND TOTAL"] = responseGrandTotalSrvcCategory
				}
				continue
			}

			if response[subcategoryName] == nil {
				response[subcategoryName] = make(map[string]interface{})
			}

			planMap, _ := subcategoryMap.(map[string]interface{})
			responsePlanMap, _ := response[subcategoryName].(map[string]interface{})

			for planName, plan := range planMap {
				if responsePlanMap[planName] == nil {
					responsePlanMap[planName] = make(map[string]int)
				}

				srvccategoryMap, _ := plan.(map[string]int)
				responseSrvccategoryMap, _ := responsePlanMap[planName].(map[string]int)

				for srvcCategoryName, qty := range srvccategoryMap {
					responseSrvccategoryMap[srvcCategoryName] = responseSrvccategoryMap[srvcCategoryName] + qty
				}
			}
		}
	}

	go func(product string, report map[string]interface{}) {
		err = impl.CacheService.InsertUserPlanCategoryReport(mngr.UserId, product, chanel, isQuarterReport, month, quarter, year, report,
			impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("product", product).Err(errors.New(err.Details.(string))).
				Msg("MngrUserPlanPerCategoryReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("product", product).Msg("MngrUserPlanPerCategoryReport : cache update success")
		}
	}(product, response)

	return response, nil
}
