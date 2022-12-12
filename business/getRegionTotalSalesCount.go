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

func (impl *IncentiveImpl) GetRegionTotalSalesCount(ctx *gin.Context, userID string, isQuarterReport bool, quarter, month, year int) ([]*models.ProductWiseSalesCountPerOutlet, *errors.Error) {

	user, err := impl.PsqlDb.FetchUserBasicDetails(userID, quarter, year)
	if err != nil {
		return nil, err
	}

	mangerPosition, ok := constant.ManagerMap[user.Position]
	if !ok || user.Position == constant.Promoter || user.Position == constant.CSR || user.Position == constant.OutletManager || user.Position == constant.TeamLeader {
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

	resp := make([]*models.ProductWiseSalesCountPerOutlet, 0)

	// for regional manager
	if user.Position == constant.Manager {

		regionReport, err := impl.PsqlDb.FetchRegionReport(user.Region, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil || (err == nil && regionReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				regionResp, err := impl.updateRegionReport(ctx, user.Region, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}

				//TODO: test this out for struct assignment
				regionReport = regionResp
			}
		}

		eachRegionResp := &models.ProductWiseSalesCountPerOutlet{
			Region:  regionReport.Region,
			HBB:     regionReport.HBBSalesCount,
			Hayyak:  regionReport.HayyakSalesCount,
			Baqati:  regionReport.BaqatiSalesCount,
			Upgrade: regionReport.UpgradeSalesCount,
		}

		eachRegionResp.Total = regionReport.HBBSalesCount + regionReport.HayyakSalesCount + regionReport.BaqatiSalesCount + regionReport.UpgradeSalesCount

		resp = append(resp, eachRegionResp)

		return resp, nil
	}

	//// for non region manager

	allDistinctMngr, err := impl.PsqlDb.FetchAllRMReportingToGivenManager(mangerPosition, userID, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(allDistinctMngr) == 0 {
		logger.FileLogger.CtxError(ctx).Str("position", mangerPosition).Msg("GetRegionSalesAndTarget : no RM found")
		return resp, nil
	}

	for _, eachMngr := range allDistinctMngr {

		regionReport, err := impl.PsqlDb.FetchRegionReport(eachMngr.Region, isQuarterReport, quarter, month, year)
		if !needUpdate {
			if err != nil {
				return nil, err
			}
		} else {
			if err != nil || (err == nil && regionReport.ModifiedAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location()))) {

				regionResp, err := impl.updateRegionReport(ctx, eachMngr.Region, isQuarterReport, quarter, month, year)
				if err != nil {
					return nil, err
				}

				//TODO: test this out for struct assignment
				regionReport = regionResp
			}
		}

		eachRegionResp := &models.ProductWiseSalesCountPerOutlet{
			Region:  regionReport.Region,
			HBB:     regionReport.HBBSalesCount,
			Hayyak:  regionReport.HayyakSalesCount,
			Baqati:  regionReport.BaqatiSalesCount,
			Upgrade: regionReport.UpgradeSalesCount,
		}

		eachRegionResp.Total = regionReport.HBBSalesCount + regionReport.HayyakSalesCount + regionReport.BaqatiSalesCount + regionReport.UpgradeSalesCount

		resp = append(resp, eachRegionResp)
	}

	return resp, nil
}
