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

func (impl *UserImpl) updateAgentUserMNPPlanCountReport(ctx *gin.Context, user *models.UserBasicDetailsDB, chanel string,
	isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if user == nil || user.UserId == "" || user.ManagerId == "" || user.ManagerPosition == "" {
		if user == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").
				Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).
				Int("year", year).Msg("AMNPPlanReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID AGENT DETAILS",
				Details:    "updateAgentUserMNPPlanCountReport: invalid agent details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).
			Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
			Int("month", month).Int("year", year).Msg("AMNPPlanReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    user.UserId + ": INVALID AGENT DETAILS",
			Details:    "updateAgentUserMNPPlanCountReport: invalid agent details - userId:" + user.UserId,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Str("region", user.Region).
		Str("outlet", user.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("AMNPPlanReport : info")

	if !(user.Position == constant.PromoterAbb || user.Position == constant.CSRAbb) {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("position", user.Position).
			Msg("AMNPPlanReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT AGENT POSITION - " + user.Position,
			Details: "updateAgentUserMNPPlanCountReport: incorrect agent position - userId:" + user.UserId +
				" position:" + user.Position,
		}
	}

	if user.Region == "" || user.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Str("region", user.Region).
			Str("outlet", user.Outlet).Msg("AMNPPlanReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: user.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + user.Region +
				" outlet:" + user.Outlet,
			Details: "updateAgentUserMNPPlanCountReport: agent blank region or outlet - userId: " + user.UserId +
				" region:" + user.Region + " outlet:" + user.Outlet,
		}
	}

	outletResponse, err := impl.DbService.FetchOutletDetailsViaName(user.Region, user.Outlet, quarter, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserMNPPlanCountReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).
			Msg("AMNPPlanReport : fetch outlet info failed")
		return nil, err
	}

	if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Msg("AMNPPlanReport : outlet co-ordinates invalid")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: user.UserId + ": AGENT INVALID OUTLET CO-ORDINATES - region:" + user.Region +
				" outlet:" + user.Outlet,
			Details: "updateAgentUserMNPPlanCountReport: " + user.UserId + ": AGENT invalid outlet co-ordintaes - region:" +
				user.Region + " outlet:" + user.Outlet,
		}
	}

	allUserSales, err := impl.DbService.FetchActiveUserSales(user.UserId, chanel, isQuarterReport, quarter, month, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateAgentUserMNPPlanCountReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(err).Msg("AMNPPlanReport : fetch sales failed")
		return nil, err
	}
	allUserSales.Region = user.Region
	allUserSales.Outlet = user.Outlet

	logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Int("saleGroupCount", len(allUserSales.SaleRecords)).
		Msg("AMNPPlanReport : total agents and salesGroup")

	resp := &models.MNPPlansCount{
		Total:      make(map[string]map[string]int),
		Baqati:     make(map[string]map[string]int),
		Hayyak:     make(map[string]map[string]int),
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if len(allUserSales.SaleRecords) == 0 {
		resp.Error = "NO SALES FOUND FOR OUTLET - " + user.Outlet
		return resp, nil
	}

	for grp, rcds := range allUserSales.SaleRecords {
		for _, v := range rcds {
			if v.Category == constant.UpgradeCategory {
				continue
			}

			switch grp {
			case constant.BaqatiCategory:
				plan := resp.Baqati[v.CurrentProduct]
				if plan == nil {
					plan = map[string]int{
						v.MNPDonor: v.Quantity,
					}
				} else {
					plan[v.MNPDonor] += v.Quantity
				}

				group := resp.Total[grp]
				if group == nil {
					group = map[string]int{
						v.MNPDonor: v.Quantity,
					}
				} else {
					group[v.MNPDonor] += v.Quantity
				}
			case constant.HayyakCategory:
				plan := resp.Hayyak[v.CurrentProduct]
				if plan == nil {
					plan = map[string]int{
						v.MNPDonor: v.Quantity,
					}
				} else {
					plan[v.MNPDonor] += v.Quantity
				}

				group := resp.Total[grp]
				if group == nil {
					group = map[string]int{
						v.MNPDonor: v.Quantity,
					}
				} else {
					group[v.MNPDonor] += v.Quantity
				}
			default:
				continue
			}
		}
	}

	go func(report *models.MNPPlansCount) {
		if report.Error != "" {
			return
		}
		err = impl.CacheService.InsertUserMNPPlanCountPerPlanReport(user.UserId, chanel, isQuarterReport, month,
			quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", user.UserId).Err(errors.New(err.Details.(string))).
				Msg("AMNPPlanReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", user.UserId).Msg("AMNPPlanReport : cache update success")
		}
	}(resp)

	return resp, nil
}

func (impl *UserImpl) updateOutletManagerUserMNPPlanCountReport(ctx *gin.Context, omngr *models.UserBasicDetailsDB,
	chanel string, isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if omngr == nil || omngr.UserId == "" || omngr.ManagerId == "" || omngr.ManagerPosition == "" {

		if omngr == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", "").Str("region", "").Str("outlet", "").
				Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).
				Int("year", year).Msg("OMMNPPlanCountReport : nil info")
			return nil, &errors.Error{
				StatusCode: http.StatusInternalServerError,
				Code:       errors.InternalServerErrorCode,
				Message:    "INVALID OUTLET MANAGER DETAILS",
				Details:    "updateOMMNPUserPlanReport: invalid omngr details",
			}
		}

		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
			Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
			Int("month", month).Int("year", year).Msg("OMMNPPlanCountReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    omngr.UserId + ": INVALID OUTLET MANAGER DETAILS",
			Details:    "updateOMMNPUserPlanReport: invalid omngr details - userId:" + omngr.UserId,
		}
	}

	if omngr.Position != constant.OutletManagerAbb { //!(omngr.Position == constant.OutletManager || omngr.Position == constant.TeamLeader)
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("position", omngr.Position).
			Msg("OMMNPPlanCountReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    omngr.UserId + ": INCORRECT OUTLET MANAGER POSITION - " + omngr.Position,
			Details: "updateOMMNPUserPlanReport: incorrect omngr position - userId:" + omngr.UserId +
				" position:" + omngr.Position,
		}
	}

	if omngr.Region == "" || omngr.Outlet == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
			Str("outlet", omngr.Outlet).Msg("OMMNPPlanCountReport : invalid outlet info")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message: omngr.UserId + ": AGENT BLANK REGION OR OUTLET - region:" + omngr.Region +
				" outlet:" + omngr.Outlet,
			Details: "updateOMMNPUserPlanReport: agent blank region or outlet - userId: " + omngr.UserId +
				" region:" + omngr.Region + " outlet:" + omngr.Outlet,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("region", omngr.Region).
		Str("outlet", omngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
		Int("month", month).Int("year", year).Msg("OMMNPPlanCountReport : info")

	userAgents, err := impl.DbService.FetchAllAgentsUserOfOM(omngr.UserId, quarter, year)
	if err != nil {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(errors.New(err.Details.(string))).
			Msg("OMMNPPlanCountReport : fetch agents failed")
		return nil, err
	}

	resp := &models.MNPPlansCount{
		Total:      make(map[string]map[string]int),
		Baqati:     make(map[string]map[string]int),
		Hayyak:     make(map[string]map[string]int),
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if len(userAgents) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Msg("OMMNPPlanCountReport : zero agents")
		return resp, nil
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Int("agentCount", len(userAgents)).
			Msg("OMMNPPlanCountReport : total agents")
	}

	for _, eachAgent := range userAgents {
		var (
			fromCache bool
			userResp  *models.MNPPlansCount
		)
		agentMNPReport, err := impl.CacheService.GetUserMNPPlanCountPerPlanReport(eachAgent.UserId, chanel,
			isQuarterReport, month, quarter, year)
		if err != nil || agentMNPReport == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
				Err(err).Msg("OMMNPPlanCountReport : cache report failed")
		} else {
			userResp = agentMNPReport.(*models.MNPPlansCount)
			if objecthelper.IsCurrentDay(userResp.ModifiedAt) {
				fromCache = true
			}
		}

		if !fromCache {
			agentMNPReport, err = impl.updateAgentUserMNPPlanCountReport(ctx, eachAgent, chanel, isQuarterReport,
				quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
					Err(err).Msg("OMMNPPlanCountReport : agent report failed")
				continue
			} else {
				userResp = agentMNPReport.(*models.MNPPlansCount)
			}
		}

		logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Str("agent", eachAgent.UserId).
			Interface("record", agentMNPReport).Msg("OMMNPPlanCountReport : agent log")

		for planName, planInfo := range userResp.Baqati {
			respPlan := resp.Baqati[planName]
			if respPlan == nil {
				resp.Baqati[planName] = planInfo
			} else {
				for donor, qty := range planInfo {
					respPlan[donor] += qty
				}
			}
		}

		for planName, planInfo := range userResp.Hayyak {
			respPlan := resp.Hayyak[planName]
			if respPlan == nil {
				resp.Hayyak[planName] = planInfo
			} else {
				for donor, qty := range planInfo {
					respPlan[donor] += qty
				}
			}
		}

		for grpName, grpInfo := range userResp.Total {
			respGrp := resp.Total[grpName]
			if respGrp == nil {
				resp.Total[grpName] = grpInfo
			} else {
				for donor, qty := range grpInfo {
					respGrp[donor] += qty
				}
			}
		}
	}

	go func(report *models.MNPPlansCount) {
		if report.Error != "" {
			return
		}
		err = impl.CacheService.InsertUserMNPPlanCountPerPlanReport(omngr.UserId, chanel, isQuarterReport, month,
			quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", omngr.UserId).Err(errors.New(err.Details.(string))).
				Msg("OMMNPPlanCountReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", omngr.UserId).Msg("OMMNPPlanCountReport : cache update success")
		}
	}(resp)

	return resp, nil
}

func (impl *UserImpl) updateManagerUserMNPPlanCountReport(ctx *gin.Context, mngr *models.UserBasicDetailsDB, chanel string,
	isQuarterReport bool, quarter, month, year int) (interface{}, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if mngr == nil || mngr.UserId == "" {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("region", mngr.Region).
			Str("outlet", mngr.Outlet).Bool("isQuarter", isQuarterReport).Int("quarter", quarter).
			Int("month", month).Int("year", year).Msg("MngrMNPPlanCountReport : invalid info")
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    "INVALID MANAGER DETAILS",
			Details:    "updateManagerUserMNPPlanCountReport: invalid mngr details",
		}
	}

	position, ok := constant.ManagerMap[mngr.Position]
	if !ok {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).
			Msg("MngrMNPPlanCountReport : invalid position")
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    mngr.UserId + ": INCORRECT MANAGER POSITION - " + mngr.Position,
			Details: "updateManagerUserMNPPlanCountReport: incorrect mngr position - userId:" +
				mngr.UserId + " position:" + mngr.Position,
		}
	}
	mngr.Position = position

	rawDetails, _ := json.Marshal(mngr)
	logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("user", string(rawDetails)).
		Bool("isQuarter", isQuarterReport).Int("quarter", quarter).Int("month", month).Int("year", year).
		Msg("MngrMNPPlanCountReport : info")

	//for outlet manager and team lead only
	if mngr.Position == constant.OutletManagerAbb { // || mngr.Position == constant.TeamLeader
		report, err := impl.updateOutletManagerUserMNPPlanCountReport(ctx, mngr, chanel, isQuarterReport, quarter, month, year)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	if mngr.Position != constant.VicePresidentAbb && (mngr.ManagerId == "" || mngr.ManagerPosition == "") {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("position", mngr.Position).
			Msg("MngrMNPPlanCountReport : user's manager info missing")
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
			Err(errors.New(err.Details.(string))).Msg("MngrMNPPlanCountReport : fetch om list failed")
		return nil, err
	}

	resp := &models.MNPPlansCount{
		Total:      make(map[string]map[string]int),
		Baqati:     make(map[string]map[string]int),
		Hayyak:     make(map[string]map[string]int),
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Msg("MngrMNPPlanCountReport : zero om")
		return resp, nil
	} else {
		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Int("omCount", len(allDistinctOM)).
			Msg("MngrMNPPlanCountReport : total om")
	}

	for _, eachAgent := range allDistinctOM {
		var (
			fromCache bool
			userResp  *models.MNPPlansCount
		)
		agentMNPReport, err := impl.CacheService.GetUserMNPPlanCountPerPlanReport(eachAgent.UserId,
			chanel, isQuarterReport, month, quarter, year)
		if err != nil || agentMNPReport == nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
				Err(err).Msg("MngrMNPPlanCountReport : cache report failed")
		} else {
			userResp = agentMNPReport.(*models.MNPPlansCount)
			if objecthelper.IsCurrentDay(userResp.ModifiedAt) {
				fromCache = true
			}
		}

		if !fromCache {
			agentMNPReport, err = impl.updateOutletManagerUserMNPPlanCountReport(ctx, eachAgent, chanel, isQuarterReport,
				quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
					Err(err).Msg("MngrMNPPlanCountReport : agent report failed")
				continue
			} else {
				userResp = agentMNPReport.(*models.MNPPlansCount)
			}
		}

		logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Str("agent", eachAgent.UserId).
			Interface("record", agentMNPReport).Msg("MngrMNPPlanCountReport : agent log")

		for planName, planInfo := range userResp.Baqati {
			respPlan := resp.Baqati[planName]
			if respPlan == nil {
				resp.Baqati[planName] = planInfo
			} else {
				for donor, qty := range planInfo {
					respPlan[donor] += qty
				}
			}
		}

		for planName, planInfo := range userResp.Hayyak {
			respPlan := resp.Hayyak[planName]
			if respPlan == nil {
				resp.Hayyak[planName] = planInfo
			} else {
				for donor, qty := range planInfo {
					respPlan[donor] += qty
				}
			}
		}

		for grpName, grpInfo := range userResp.Total {
			respGrp := resp.Total[grpName]
			if respGrp == nil {
				resp.Total[grpName] = grpInfo
			} else {
				for donor, qty := range grpInfo {
					respGrp[donor] += qty
				}
			}
		}
	}

	go func(report *models.MNPPlansCount) {
		if report.Error != "" {
			return
		}
		err = impl.CacheService.InsertUserMNPPlanCountPerPlanReport(mngr.UserId, chanel, isQuarterReport, month,
			quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("clientID", mngr.UserId).Err(errors.New(err.Details.(string))).
				Msg("MngrMNPPlanCountReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("clientID", mngr.UserId).Msg("MngrMNPPlanCountReport : cache update success")
		}
	}(resp)

	return resp, nil
}
