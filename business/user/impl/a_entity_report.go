package impl

import (
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

func (impl *UserImpl) updateOutletReport(ctx *gin.Context, regionName, outletName string, chanel string,
	isQuarterReport bool, quarter, month, year int) (*models.OutletReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if regionName == "" || outletName == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    "BLANK REGION OR OUTLET - region:" + regionName + " outlet:" + outletName,
			Details:    "blank region or outlet - region:" + regionName + " outlet:" + outletName,
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("region", regionName).Str("outlet", outletName).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("update Outlet Report")

	lat, long := 0.0, 0.0
	outletResponse, err := impl.DbService.FetchOutletDetailsViaName(regionName, outletName, quarter, year)
	if err != nil {
		err.Details = fmt.Sprintf("updateOutletReport: %v", err.Details)
		logger.FileLogger.CtxError(ctx).Str("outletName", outletName).Err(errors.New(err.Details.(string))).
			Msg("OutletReport : fetch outlet info failed")
		// return nil, nil, err
	} else if outletResponse != nil {
		lat = outletResponse.Lat
		long = outletResponse.Lang
	}

	// if outletResponse.Lat == 0 || outletResponse.Lang == 0 {
	// 	return nil, &errors.Error{
	// 		StatusCode: http.StatusBadRequest,
	// 		Code:       errors.InvalidDetailsCode,
	// 		Message:    "INVALID OUTLET CO-ORDINATES - region:" + regionName + " outlet:" + outletName,
	// 		Details:    "invalid outlet co-ordintaes - region:" + regionName + " outlet:" + outletName,
	// 	}
	// }

	chanel = "" // chanel is not considered for promoters/csr sales
	hbbSales, hayyakSales, baqatiSales, upgradeSales, err := impl.DbService.FetchOutletSalesCount(regionName, outletName,
		chanel, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}

	hbbTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.HBBCategory,
		regionName, outletName)
	if internalErr != nil {
		// return nil, &errors.Error{
		// 	StatusCode: http.StatusBadRequest,
		// 	Code:       errors.InvalidDetailsCode,
		// 	Message:    internalErr.Error(),
		// 	Details:    internalErr.Error(),
		// }
	}

	hayyakTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.HayyakCategory,
		regionName, outletName)
	if internalErr != nil {
		// return nil, &errors.Error{
		// 	StatusCode: http.StatusBadRequest,
		// 	Code:       errors.InvalidDetailsCode,
		// 	Message:    internalErr.Error(),
		// 	Details:    internalErr.Error(),
		// }
	}

	baqatiTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.BaqatiCategory,
		regionName, outletName)
	if internalErr != nil {
		// return nil, &errors.Error{
		// 	StatusCode: http.StatusBadRequest,
		// 	Code:       errors.InvalidDetailsCode,
		// 	Message:    internalErr.Error(),
		// 	Details:    internalErr.Error(),
		// }
	}

	upgradeTarget, internalErr := impl.getSaleTargetPerOutlet(year, month, quarter, isQuarterReport, constant.UpgradeCategory,
		regionName, outletName)
	if internalErr != nil {
		// return nil, &errors.Error{
		// 	StatusCode: http.StatusBadRequest,
		// 	Code:       errors.InvalidDetailsCode,
		// 	Message:    internalErr.Error(),
		// 	Details:    internalErr.Error(),
		// }
	}

	outletReport := &models.OutletReportDB{
		Outlet:             outletName,
		Region:             regionName,
		Latitude:           lat,
		Longitude:          long,
		HBBSalesCount:      hbbSales,
		HayyakSalesCount:   hayyakSales,
		BaqatiSalesCount:   baqatiSales,
		UpgradeSalesCount:  upgradeSales,
		HBBTargetCount:     hbbTarget,
		HayyakTargetCount:  hayyakTarget,
		BaqatiTargetCount:  baqatiTarget,
		UpgradeTargetCount: upgradeTarget,
		CreatedAt:          time.Now(),
		ModifiedAt:         time.Now(),
	}

	go func(report *models.OutletReportDB) {
		err = impl.CacheService.InsertOutletReport(chanel, isQuarterReport, month, quarter, year, report,
			impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("outlet", outletName).Str("region", regionName).
				Err(errors.New(err.Details.(string))).Msg("OutletReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("outlet", outletName).Msg("OutletReport : cache update success")
		}

	}(outletReport)

	return outletReport, nil
}

func (impl *UserImpl) updateRegionReport(ctx *gin.Context, regionName string, chanel string, isQuarterReport bool,
	quarter, month, year int) (*models.RegionReportDB, *errors.Error) {

	if !isQuarterReport {
		quarter = getQuarterFromMonth(month)
	}

	if regionName == "" {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    "BLANK REGION NAME",
			Details:    "blank region name",
		}
	}

	logger.FileLogger.CtxInfo(ctx).Str("region", regionName).Bool("isQuarter", isQuarterReport).
		Int("quarter", quarter).Int("month", month).Int("year", year).Msg("update Outlet Report")

	allOutlets, err := impl.DbService.FetchAllOutletNames(regionName, quarter, year)
	if err != nil {
		return nil, err
	}

	if len(allOutlets) == 0 {
		return nil, &errors.Error{
			StatusCode: http.StatusBadRequest,
			Code:       errors.InvalidDetailsCode,
			Message:    "NO OUTLETS FOUND FOR REGION - " + regionName,
			Details:    "no outlet found for region - " + regionName,
		}
	}

	regionReport := &models.RegionReportDB{
		Region:     regionName,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	for _, eachOutlet := range allOutlets {

		var (
			err          *errors.Error
			outletReport *models.OutletReportDB
		)

		outletReport, err = impl.CacheService.GetOutletReport(eachOutlet, regionName, chanel, isQuarterReport, month, quarter, year)
		if !(err == nil && outletReport != nil && objecthelper.IsCurrentDay(outletReport.ModifiedAt)) {
			outletResp, err := impl.updateOutletReport(ctx, regionName, eachOutlet, chanel, isQuarterReport, quarter, month, year)
			if err != nil {
				logger.FileLogger.CtxError(ctx).Str("region", regionName).Str("outlet", eachOutlet).Bool("isQuarter", isQuarterReport).
					Int("quarter", quarter).Int("month", month).Int("year", year).Err(err).Msg("RegionReport : outlet report failed")
				//return nil, err
				continue
			}

			logger.FileLogger.CtxInfo(ctx).Interface("outletReport", outletResp).Msg("updateRegionReport : trace")
			//TODO: test this out for struct assignment
			outletReport = outletResp
		}

		regionReport.HayyakSalesCount += outletReport.HayyakSalesCount
		regionReport.HBBSalesCount += outletReport.HBBSalesCount
		regionReport.BaqatiSalesCount += outletReport.BaqatiSalesCount
		regionReport.UpgradeSalesCount += outletReport.UpgradeSalesCount

		regionReport.HayyakTargetCount += outletReport.HayyakTargetCount
		regionReport.HBBTargetCount += outletReport.HBBTargetCount
		regionReport.BaqatiTargetCount += outletReport.BaqatiTargetCount
		regionReport.UpgradeTargetCount += outletReport.UpgradeTargetCount
	}

	go func(report *models.RegionReportDB) {

		err = impl.CacheService.InsertRegionReport(chanel, isQuarterReport, month, quarter, year, report, impl.Config.ReportCacheTTL)
		if err != nil {
			logger.FileLogger.CtxError(ctx).Str("region", regionName).Err(errors.New(err.Details.(string))).Msg("RegionReport : force report update")
		} else {
			logger.FileLogger.CtxInfo(ctx).Str("region", regionName).Msg("RegionReport` : cache update success")
		}

	}(regionReport)

	return regionReport, nil
}
