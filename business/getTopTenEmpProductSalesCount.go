package business

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) GetTopTenEmpProductSalesCount(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesCountPerStaff, *errors.Error) {

	user, err := impl.PsqlDb.FetchUserBasicDetails(userID)
	if err != nil {
		return nil, err
	}

	_, ok := constant.ManagerMap[user.Position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "INCORRECT USER POSITION",
			Details:    "incorrect position",
		}
	}

	var needUpdate bool

	currTime := time.Now()
	currYear, currMonth, currDay := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if isQuarterReport && year == currYear && quarter <= currQuarter {
		needUpdate = true
	} else if !isQuarterReport && year == currYear && month <= int(currMonth) {
		needUpdate = true
	}

	resp := make([]*models.ProductWiseSalesCountPerStaff, 0)

	// for oulletManager and teamLeader
	if user.Position == constant.OutletManager || user.Position == constant.TeamLeader {

		agentList, err := impl.PsqlDb.FetchAllAgentsUserOfOM(userID)
		if err != nil {
			return nil, err
		}

		if len(agentList) == 0 {
			return resp, nil
		}

		for _, agent := range agentList {

			agentReport, err := impl.PsqlDb.FetchAgentUserReport(agent.UserId, isQuarterReport, quarter, month, year)
			if needUpdate && (err != nil || (err == nil && agentReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())))) {
				_, agentReport, err = impl.updateAgentUserReport(ctx, agent, isQuarterReport, quarter, month, year)
			}

			if err != nil {
				return nil, err
			}

			agentName, err := impl.PsqlDb.FetchUserName(agentReport.UserId)
			if err != nil {
				return nil, err
			}

			agentProductWiseRprt := &models.ProductWiseSalesCountPerStaff{
				Name:    agentName,
				HBB:     agentReport.HBBSalesCount,
				Hayyak:  agentReport.HayyakSalesCount,
				Baqati:  agentReport.BaqatiSalesCount,
				Upgrade: agentReport.UpgradeSalesCount,
			}

			agentProductWiseRprt.Total = agentProductWiseRprt.HBB + agentProductWiseRprt.Hayyak + agentProductWiseRprt.Baqati + agentProductWiseRprt.Upgrade

			resp = append(resp, agentProductWiseRprt)

			sort.Slice(resp, func(i, j int) bool {
				return resp[i].Total > resp[j].Total
			})

			if len(resp) > 10 {
				resp = resp[:10]
			}
		}

		return resp, nil
	}

	// for all other non-outlet managers
	allDistinctMngr, err := impl.PsqlDb.FetchAllMngrReportingToManagerRole(user.Position, userID)
	if err != nil {
		return nil, err
	}

	if len(allDistinctMngr) == 0 {
		logger.FileLogger.CtxError(ctx).Str("position", user.Position).Msg("GetTopTenEmpProductSalesCount : no Mngr found")
		return resp, nil
	}

	for _, eachMngr := range allDistinctMngr {

		mngrReport, err := impl.PsqlDb.FetchManagerReport(eachMngr.UserId, isQuarterReport, quarter, month, year)
		if needUpdate && (err != nil || (err == nil && mngrReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())))) {
			mngrReport, err = impl.updateManagerReport(ctx, eachMngr, isQuarterReport, quarter, month, year)
		}

		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		mngrName, err := impl.PsqlDb.FetchUserName(mngrReport.UserId)
		if err != nil {
			return nil, err
		}

		mngrProductWiseRprt := &models.ProductWiseSalesCountPerStaff{
			Name:    mngrName,
			HBB:     mngrReport.HBBSalesCount,
			Hayyak:  mngrReport.HayyakSalesCount,
			Baqati:  mngrReport.BaqatiSalesCount,
			Upgrade: mngrReport.UpgradeSalesCount,
		}

		mngrProductWiseRprt.Total = mngrProductWiseRprt.HBB + mngrProductWiseRprt.Hayyak + mngrProductWiseRprt.Baqati + mngrProductWiseRprt.Upgrade

		resp = append(resp, mngrProductWiseRprt)

		sort.Slice(resp, func(i, j int) bool {
			return resp[i].Total > resp[j].Total
		})

		if len(resp) > 10 {
			resp = resp[:10]
		}
	}

	return resp, nil
}
