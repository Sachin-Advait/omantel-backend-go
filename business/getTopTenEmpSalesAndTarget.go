package business

import (
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) GetTopTenEmpSalesAndTarget(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesTargetPerStaff, *errors.Error) {

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

	resp := make([]*models.ProductWiseSalesTargetPerStaff, 0)

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

			agentProductWiseRprt := &models.ProductWiseSalesTargetPerStaff{
				Name:                   agentName,
				HBBSales:               agentReport.HBBSalesCount,
				HBBTarget:              agentReport.HBBTargetCount,
				HBBAchievedPercent:     agentReport.HBBAchievedPercentage,
				HayyakSales:            agentReport.HayyakSalesCount,
				HayyakTarget:           agentReport.HayyakTargetCount,
				HayyakAchievedPercent:  agentReport.HayyakAchievedPercentage,
				BaqatiSales:            agentReport.BaqatiSalesCount,
				BaqatiTarget:           agentReport.BaqatiTargetCount,
				BaqatiAchievedPercent:  agentReport.BaqatiAchievedPercentage,
				UpgradeSales:           agentReport.UpgradeSalesCount,
				UpgradeTarget:          agentReport.UpgradeTargetCount,
				UpgradeAchievedPercent: agentReport.UpradeAchievedPercentage,
			}

			agentProductWiseRprt.TotalSales = agentProductWiseRprt.HBBSales + agentProductWiseRprt.HayyakSales + agentProductWiseRprt.BaqatiSales + agentProductWiseRprt.UpgradeSales
			agentProductWiseRprt.TotalTarget = agentProductWiseRprt.HBBTarget + agentProductWiseRprt.HayyakTarget + agentProductWiseRprt.BaqatiTarget + agentProductWiseRprt.UpgradeTarget
			agentProductWiseRprt.TotalAchievedPercentage = math.Round(float64(agentProductWiseRprt.TotalSales*100/int(agentProductWiseRprt.TotalTarget))*100) / 100

			resp = append(resp, agentProductWiseRprt)

			sort.Slice(resp, func(i, j int) bool {
				return resp[i].TotalAchievedPercentage > resp[j].TotalAchievedPercentage
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
		logger.FileLogger.CtxError(ctx).Str("position", user.Position).Msg("GetTopTenEmpSalesAndTarget : no Mngr found")
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

		mngrProductWiseRprt := &models.ProductWiseSalesTargetPerStaff{
			Name:                   mngrName,
			HBBSales:               mngrReport.HBBSalesCount,
			HBBTarget:              mngrReport.HBBTargetCount,
			HBBAchievedPercent:     mngrReport.HBBAchievedPercentage,
			HayyakSales:            mngrReport.HayyakSalesCount,
			HayyakTarget:           mngrReport.HayyakTargetCount,
			HayyakAchievedPercent:  mngrReport.HayyakAchievedPercentage,
			BaqatiSales:            mngrReport.BaqatiSalesCount,
			BaqatiTarget:           mngrReport.BaqatiTargetCount,
			BaqatiAchievedPercent:  mngrReport.BaqatiAchievedPercentage,
			UpgradeSales:           mngrReport.UpgradeSalesCount,
			UpgradeTarget:          mngrReport.UpgradeTargetCount,
			UpgradeAchievedPercent: mngrReport.UpradeAchievedPercentage,
		}

		mngrProductWiseRprt.TotalSales = mngrProductWiseRprt.HBBSales + mngrProductWiseRprt.HayyakSales + mngrProductWiseRprt.BaqatiSales + mngrProductWiseRprt.UpgradeSales
		mngrProductWiseRprt.TotalTarget = mngrProductWiseRprt.HBBTarget + mngrProductWiseRprt.HayyakTarget + mngrProductWiseRprt.BaqatiTarget + mngrProductWiseRprt.UpgradeTarget
		mngrProductWiseRprt.TotalAchievedPercentage = math.Round(float64(mngrProductWiseRprt.TotalSales*100/int(mngrProductWiseRprt.TotalTarget))*100) / 100

		resp = append(resp, mngrProductWiseRprt)

		sort.Slice(resp, func(i, j int) bool {
			return resp[i].TotalAchievedPercentage > resp[j].TotalAchievedPercentage
		})

		if len(resp) > 10 {
			resp = resp[:10]
		}
	}

	return resp, nil
}
