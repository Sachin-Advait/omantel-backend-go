package business

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) GetOutletTotalSalesCount(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesCountPerOutlet, *errors.Error) {

	user, err := impl.PsqlDb.FetchUserBasicDetails(userID, quarter, year)
	if err != nil {
		return nil, err
	}

	_, ok := constant.ManagerMap[user.Position]
	if !ok {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Message:    user.UserId + ": INCORRECT MANAGER POSITION - " + user.Position,
			Details:    "incorrect manager position: userId-" + user.UserId + " position-" + user.Position,
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

	//TODO: remove at last
	needUpdate = true

	resp := make([]*models.ProductWiseSalesCountPerOutlet, 0)

	// for oulletManager and teamLeader
	if user.Position == constant.OutletManager { //|| user.Position == constant.TeamLeader

		outletReport, err := impl.PsqlDb.FetchOutletReport(user.Outlet, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil || (err == nil && outletReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				outletResp, err := impl.updateOutletReport(ctx, user.Region, user.Outlet, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}

				logger.FileLogger.CtxInfo(ctx).Interface("OM outletReport", outletResp).Msg("GetOutletTotalSalesCount : trace")
				//TODO: test this out for struct assignment
				outletReport = outletResp
			}
		}

		eachOutletResp := &models.ProductWiseSalesCountPerOutlet{
			OutletName: user.Outlet,
			Region:     user.Region,
			Latitude:   outletReport.Latitude,
			Longitude:  outletReport.Longitude,
			HBB:        outletReport.HBBSalesCount,
			Hayyak:     outletReport.HayyakSalesCount,
			Baqati:     outletReport.BaqatiSalesCount,
			Upgrade:    outletReport.UpgradeSalesCount,
		}

		eachOutletResp.Total = outletReport.HBBSalesCount + outletReport.HayyakSalesCount + outletReport.BaqatiSalesCount + outletReport.UpgradeSalesCount

		resp = append(resp, eachOutletResp)

		return resp, nil
	}

	// for non outlet and tl manager
	allDistinctOM, err := impl.PsqlDb.FetchAllOMReportingToGivenManager(user.Position, userID, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(allDistinctOM) == 0 {
		logger.FileLogger.CtxError(ctx).Str("position", user.Position).Msg("GetOutletSalesCount : no OM found")
		return resp, nil
	}

	for _, eachOM := range allDistinctOM {

		outletReport, err := impl.PsqlDb.FetchOutletReport(eachOM.Outlet, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil || (err == nil && outletReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				outletResp, err := impl.updateOutletReport(ctx, eachOM.Region, eachOM.Outlet, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}

				logger.FileLogger.CtxInfo(ctx).Interface("OtherMngr outletReport", outletResp).Msg("GetOutletTotalSalesCount : trace")
				//TODO: test this out for struct assignment
				outletReport = outletResp
			}
		}

		eachOutletResp := &models.ProductWiseSalesCountPerOutlet{
			OutletName: outletReport.Outlet,
			Region:     outletReport.Region,
			Latitude:   outletReport.Latitude,
			Longitude:  outletReport.Longitude,
			HBB:        outletReport.HBBSalesCount,
			Hayyak:     outletReport.HayyakSalesCount,
			Baqati:     outletReport.BaqatiSalesCount,
			Upgrade:    outletReport.UpgradeSalesCount,
		}

		eachOutletResp.Total = outletReport.HBBSalesCount + outletReport.HayyakSalesCount + outletReport.BaqatiSalesCount + outletReport.UpgradeSalesCount

		resp = append(resp, eachOutletResp)
	}

	return resp, nil
}
