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

func (impl *UserImpl) updateAgentUserMNPSalesCountReport(ctx *gin.Context, user *models.UserBasicDetailsDB,
	chanel string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").
				Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).
				Int("year", year).Msg("AMNPReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "updateAgentUserMNPSalesCountReport: invalid agent details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).
			Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
			Int("month", month).Int("year", year).Msg("AMNPReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": INVALID AGENT DETAILS",
			Details:    "updateAgentUserMNPSalesCountReport: invalid agent details - userId:" + user.UserId,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("region", user.Region).
		Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("AMNPReport : info")

	if !(user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb) {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("position", user.Position).
			Msg("AMNPReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details: "updateAgentUserMNPSalesCountReport: incorrect agent position - userId:" + user.UserId +
				" position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).
			Str("outlet", user.Outlet).Msg("AMNPReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: user.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + user.Region +
				" outlet:" + user.Outlet,
			Details: "updateAgentUserMNPSalesCountReport: agent blank region or outlet - userId: " + user.UserId +
				" region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	outletResponse, err := impl.DbService.FetchOutletDetailsViaName(user.Region, user.Outlet, quarter, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserMNPSalesCountReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).
			Msg("AMNPReport : fetch outlet info failed")
		return nil, err
	}

	if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Msg("AMNPReport : outlet co-ordinates invalid")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: user.UserId + ": AGENT INVALID OUTLET CO-ORDINATES - region:" + user.Region +
				" outlet:" + user.Outlet,
			Details: "updateAgentUserMNPSalesCountReport: " + user.UserId + ": AGENT invalid outlet co-ordintaes - region:" +
				user.Region + " outlet:" + user.Outlet,
		}
	}

	allUserSales, err := impl.DbService.FetchActiveUserSales(user.UserId, chanel, isQuarterReport, quarter, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserMNPSalesCountReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("AMNPReport : fetch sales failed")
		return nil, err
	}
	allUserSales.Region = user.Region
	allUserSales.Outlet = user.Outlet

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("saleGroupCount", len(allUserSales.SaleRecords)).
		Msg("AMNPReport : total agents and salesGroup")

	resp := &models.MNPSalesCount{
		Baqati: &models.MNPCategorySale{
			Types: make(map[string]*models.EachMNPSale),
		},
		Hayyak: &models.MNPCategorySale{
			Types: make(map[string]*models.EachMNPSale),
		},
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if len(allUserSales.SaleRecords) == 0 {
		resp.Baqati.Error = "NO SALES FOUND FOR OUTLET - " + user.Outlet
		return resp, nil
	}

	grpPlans := make(map[string]map[string]int)
	totalActive := make(map[string]int)

	for grp, rcds := range allUserSales.SaleRecords {
		for _, v := range rcds {
			if v.Category != constant.UpgradeCategory {
				if grpPlans[grp] == nil {
					grpPlans[grp] = make(map[string]int)
				}
				grpPlans[grp][v.MNPDonor] += v.Quantity
				totalActive[grp] += v.Quantity
			}
		}
	}

	for grp, qty := range totalActive {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("category", grp).Int("count", qty).
			Msg("AMNPReport : sale count")
	}

	for grp, mnpDonor := range grpPlans {
		switch grp {
		case constant.HayyakCategory:
			for donorName, count := range mnpDonor {
				resp.Hayyak.Count += count
				resp.Hayyak.Types[donorName] = &models.EachMNPSale{
					Count: count,
				}
			}

		case constant.BaqatiCategory:
			for donorName, count := range mnpDonor {
				resp.Baqati.Count += count
				resp.Baqati.Types[donorName] = &models.EachMNPSale{
					Count: count,
				}
			}
		}
	}

	resp.Total = resp.Hayyak.Count + resp.Baqati.Count

	if resp.Total != 0 {
		resp.Hayyak.Percentage = float64(resp.Hayyak.Count) / float64(resp.Total) * 100
		resp.Baqati.Percentage = float64(resp.Baqati.Count) / float64(resp.Total) * 100
	}

	for _, sale := range resp.Hayyak.Types {
		if resp.Hayyak.Count != 0 {
			sale.Percentage = float64(sale.Count) / float64(resp.Hayyak.Count) * 100
		}
	}

	for _, sale := range resp.Baqati.Types {
		if resp.Baqati.Count != 0 {
			sale.Percentage = float64(sale.Count) / float64(resp.Baqati.Count) * 100
		}
	}

	go func(report *models.MNPSalesCount) {
		if report.Total <= 0 {
			return
		}
		err = impl.CacheService.InsertUserMNPSaleCountPerPlanReport(user.UserId, chanel, isQuarterReport, month,
			quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).
				Msg("AMNPReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Msg("AMNPReport : cache update success")
		}
	}(resp)

	return resp, nil
}

func (impl *UserImpl) updateOutletManagerUserMNPSaleCountReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB,
	chanel string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {
		if omngr == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").
				Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).
				Int("year", year).Msg("OMMNPSalesCountReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "updateOMMNPUserSalesReport: invalid omngr details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
			Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
			Int("month", month).Int("year", year).Msg("OMMNPSalesCountReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "updateOMMNPUserSalesReport: invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if omngr.Position != constant.OutletManagerAbb { //!(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader)
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("position", omngr.Position).
			Msg("OMMNPSalesCountReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details: "updateOMMNPUserSalesReport: incorrect omngr position - userId:" + omngr.UserId +
				" position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
			Str("outlet", omngr.Outlet).Msg("OMMNPSalesCountReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + omngr.Region +
				" outlet:" + omngr.Outlet,
			Details: "updateOMMNPUserSalesReport: agent blank region or outlet - userId: " + omngr.UserId +
				" region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
		Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("OMMNPSalesCountReport : info")

	userAgents, err := impl.DbService.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(errors.New(err.Details.(string))).
			Msg("OMMNPSalesCountReport : fetch agents failed")
		return nil, err
	}

	resp := &models.MNPSalesCount{
		Baqati: &models.MNPCategorySale{
			Types: make(map[string]*models.EachMNPSale),
		},
		Hayyak: &models.MNPCategorySale{
			Types: make(map[string]*models.EachMNPSale),
		},
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if len(userAgents) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMMNPSalesCountReport : zero agents")
		return resp, nil
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).
			Msg("OMMNPSalesCountReport : total agents")
	}

	for _, eachAgent := range userAgents {
		var (
			fromCache bool
			userResp  *models.MNPSalesCount
		)
		agentMNPReport, err := impl.CacheService.GetUserMNPSaleCountPerPlanReport(eachAgent.UserId,
			chanel, isQuarterReport, month, quarter, year)
		if err != nil || agentMNPReport == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
				Err(err).Msg("OMMNPSalesCountReport : cache report failed")
		} else {
			userResp = agentMNPReport.(*models.MNPSalesCount)
			if objecthelper.IsCurrentDay(userResp.ModifiedAt) {
				fromCache = true
			}
		}

		if !fromCache {
			agentMNPReport, err = impl.updateAgentUserMNPSalesCountReport(ctx, eachAgent, chanel, isQuarterReport,
				quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
					Err(err).Msg("OMMNPSalesCountReport : agent report failed")
				continue
			} else {
				userResp = agentMNPReport.(*models.MNPSalesCount)
			}
		}

		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
			Interface("record", agentMNPReport).Msg("OMMNPSalesCountReport : agent log")

		resp.Hayyak.Count += userResp.Hayyak.Count
		for key, val := range userResp.Hayyak.Types {
			eachSale := resp.Hayyak.Types[key]
			if eachSale == nil {
				resp.Hayyak.Types[key] = &models.EachMNPSale{
					Count: val.Count,
				}
			} else {
				eachSale.Count += val.Count
			}
		}

		resp.Baqati.Count += userResp.Baqati.Count
		for key, val := range userResp.Baqati.Types {
			eachSale := resp.Baqati.Types[key]
			if eachSale == nil {
				resp.Baqati.Types[key] = &models.EachMNPSale{
					Count: val.Count,
				}
			} else {
				eachSale.Count += val.Count
			}
		}
	}

	resp.Total = resp.Hayyak.Count + resp.Baqati.Count

	if resp.Total != 0 {
		resp.Hayyak.Percentage = float64(resp.Hayyak.Count) / float64(resp.Total) * 100
		resp.Baqati.Percentage = float64(resp.Baqati.Count) / float64(resp.Total) * 100
	}

	for _, sale := range resp.Hayyak.Types {
		if resp.Hayyak.Count != 0 {
			sale.Percentage = float64(sale.Count) / float64(resp.Hayyak.Count) * 100
		}
	}

	for _, sale := range resp.Baqati.Types {
		if resp.Baqati.Count != 0 {
			sale.Percentage = float64(sale.Count) / float64(resp.Baqati.Count) * 100
		}
	}

	go func(report *models.MNPSalesCount) {
		if report.Total <= 0 {
			return
		}
		err = impl.CacheService.InsertUserMNPSaleCountPerPlanReport(omngr.UserId, chanel, isQuarterReport, month,
			quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(errors.New(err.Details.(string))).
				Msg("OMMNPReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Msg("OMMNPReport : cache update success")
		}
	}(resp)

	return resp, nil
}

func (impl *UserImpl) updateManagerUserMNPSalesCountReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB, chanel string,
	isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if mngr == nil || mngr.UserId == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).
			Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
			Int("month", month).Int("year", year).Msg("MngrMNPSalesCountReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "updateManagerUserMNPSalesCountReport: invalid mngr details",
		}
	}

	position, ok := constant.ManagerMap[mngr.Position]
	if !ok {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).
			Msg("MngrMNPSalesCountReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details: "updateManagerUserMNPSalesCountReport: incorrect mngr position - userId:" +
				mngr.UserId + " position:" + mngr.Position,
		}
	}
	mngr.Position = position

	rawDetails, _ := json.Marshal(mngr)
	logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("user", string(rawDetails)).
		Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).Int("year", year).
		Msg("MngrMNPSalesCountReport : info")

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManagerAbb { // || mngr.Position == constant.TeamLeader
		report, err := impl.updateOutletManagerUserMNPSaleCountReport(ctx, mngr, chanel, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	if mngr.Position != constant.VicePresidentAbb && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).
			Msg("MngrMNPSalesCountReport : user's manager info missing")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    mngr.UserId + ": INVALID MANAGER'S MANAGER DETAILS",
			Details:    mngr.UserId + ": invalid manager's manager details",
		}
	}

	allDistinctOM, err := impl.DbService.FetchAllOMReportingToGivenManager(mngr.Position, mngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).
			Err(errors.New(err.Details.(string))).Msg("MngrMNPSalesCountReport : fetch om list failed")
		return nil, err
	}

	resp := &models.MNPSalesCount{
		Baqati: &models.MNPCategorySale{
			Types: make(map[string]*models.EachMNPSale),
		},
		Hayyak: &models.MNPCategorySale{
			Types: make(map[string]*models.EachMNPSale),
		},
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrMNPSalesCountReport : zero om")
		return resp, nil
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Int("omCount", len(allDistinctOM)).
			Msg("MngrMNPSalesCountReport : total om")
	}

	for _, eachAgent := range allDistinctOM {
		var (
			fromCache bool
			userResp  *models.MNPSalesCount
		)
		agentMNPReport, err := impl.CacheService.GetUserMNPSaleCountPerPlanReport(eachAgent.UserId,
			chanel, isQuarterReport, month, quarter, year)
		if err != nil || agentMNPReport == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
				Err(err).Msg("MngrMNPSalesCountReport : cache report failed")
		} else {
			userResp = agentMNPReport.(*models.MNPSalesCount)
			if objecthelper.IsCurrentDay(userResp.ModifiedAt) {
				fromCache = true
			}
		}

		if !fromCache {
			agentMNPReport, err = impl.updateOutletManagerUserMNPSaleCountReport(ctx, eachAgent, chanel, isQuarterReport,
				quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
					Err(err).Msg("MngrMNPSalesCountReport : agent report failed")
				continue
			} else {
				userResp = agentMNPReport.(*models.MNPSalesCount)
			}
		}

		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
			Interface("record", agentMNPReport).Msg("MngrMNPSalesCountReport : agent log")

		resp.Hayyak.Count += userResp.Hayyak.Count
		for key, val := range userResp.Hayyak.Types {
			eachSale := resp.Hayyak.Types[key]
			if eachSale == nil {
				resp.Hayyak.Types[key] = &models.EachMNPSale{
					Count: val.Count,
				}
			} else {
				eachSale.Count += val.Count
			}
		}

		resp.Baqati.Count += userResp.Baqati.Count
		for key, val := range userResp.Baqati.Types {
			eachSale := resp.Baqati.Types[key]
			if eachSale == nil {
				resp.Baqati.Types[key] = &models.EachMNPSale{
					Count: val.Count,
				}
			} else {
				eachSale.Count += val.Count
			}
		}
	}

	resp.Total = resp.Hayyak.Count + resp.Baqati.Count

	if resp.Total != 0 {
		resp.Hayyak.Percentage = float64(resp.Hayyak.Count) / float64(resp.Total) * 100
		resp.Baqati.Percentage = float64(resp.Baqati.Count) / float64(resp.Total) * 100
	}

	for _, sale := range resp.Hayyak.Types {
		if resp.Hayyak.Count != 0 {
			sale.Percentage = float64(sale.Count) / float64(resp.Hayyak.Count) * 100
		}
	}

	for _, sale := range resp.Baqati.Types {
		if resp.Baqati.Count != 0 {
			sale.Percentage = float64(sale.Count) / float64(resp.Baqati.Count) * 100
		}
	}

	go func(report *models.MNPSalesCount) {
		if report.Total <= 0 {
			return
		}
		err = impl.CacheService.InsertUserMNPSaleCountPerPlanReport(mngr.UserId, chanel, isQuarterReport, month,
			quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Err(errors.New(err.Details.(string))).
				Msg("MngrMNPSalesCountReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Msg("MngrMNPSalesCountReport : cache update success")
		}
	}(resp)

	return resp, nil
}
