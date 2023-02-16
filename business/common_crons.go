package business

import (
	"errors"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) LoadPlansFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchPlans(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchPlans(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(errors.New(err.Details.(string))).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchPlans(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(errors.New(err.Details.(string))).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	// prevQuarter := impl.getQuarterFromMonth(int(currMonth)) - 1
	// if prevQuarter > 0 {
	// 	if err := impl.PsqlDb.FetchPlans(prevQuarter, currYear); err != nil {
	// 		if printLog {
	// 			logger.FileLogger.Error().Int("prevQuarter", prevQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadPlansFromDB : productPlans fetching failed")
	// 		}
	// 		return errors.New(err.Details.(string))
	// 	}
	// }

	// if printLog {
	// 	logger.FileLogger.Info().Int("prevQuarter", prevQuarter).Int("currentQuarter", currQuarter).Int("year", currYear).Msg("LoadPlansFromDB : successful")
	// }

	return nil
}

func (impl *IncentiveImpl) LoadTargetFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchSaleTarget(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchSaleTarget(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(errors.New(err.Details.(string))).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchSaleTarget(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(errors.New(err.Details.(string))).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	// prevQuarter := impl.getQuarterFromMonth(int(currMonth)) - 1
	// if prevQuarter > 0 {
	// 	if err := impl.PsqlDb.FetchSaleTarget(prevQuarter, currYear); err != nil {
	// 		if printLog {
	// 			logger.FileLogger.Error().Int("prevQuarter", prevQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadTargetFromDB : saleTarget fetching failed")
	// 		}
	// 		return errors.New(err.Details.(string))
	// 	}
	// }

	// if printLog {
	// 	logger.FileLogger.Info().Int("prevQuarter", prevQuarter).Int("currentQuarter", currQuarter).Int("year", currYear).Msg("LoadTargetFromDB : successful")
	// }

	return nil
}

func (impl *IncentiveImpl) LoadScaleOrderFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchValueOfPoint(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadScaleOrderFromDB : vop fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchValueOfPoint(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(errors.New(err.Details.(string))).Msg("LoadScaleOrderFromDB : vop fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchValueOfPoint(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(errors.New(err.Details.(string))).Msg("LoadScaleOrderFromDB : vop fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	// prevQuarter := impl.getQuarterFromMonth(int(currMonth)) - 1
	// if prevQuarter > 0 {
	// 	if err := impl.PsqlDb.FetchValueOfPoint(prevQuarter, currYear); err != nil {
	// 		if printLog {
	// 			logger.FileLogger.Error().Int("prevQuarter", prevQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadScaleOrderFromDB : vop fetching failed")
	// 		}
	// 		return errors.New(err.Details.(string))
	// 	}
	// }

	// if printLog {
	// 	logger.FileLogger.Info().Int("prevQuarter", prevQuarter).Int("currentQuarter", currQuarter).Int("year", currYear).Msg("LoadScaleOrderFromDB : successful")
	// }

	return nil
}

func (impl *IncentiveImpl) LoadMaxPayFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchMaxPayout(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchMaxPayout(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(errors.New(err.Details.(string))).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchMaxPayout(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(errors.New(err.Details.(string))).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	// prevQuarter := impl.getQuarterFromMonth(int(currMonth)) - 1
	// if prevQuarter > 0 {
	// 	if err := impl.PsqlDb.FetchMaxPayout(prevQuarter, currYear); err != nil {
	// 		if printLog {
	// 			logger.FileLogger.Error().Int("prevQuarter", prevQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
	// 		}
	// 		return errors.New(err.Details.(string))
	// 	}
	// }

	// if printLog {
	// 	logger.FileLogger.Info().Int("prevQuarter", prevQuarter).Int("currentQuarter", currQuarter).Int("year", currYear).Msg("LoadMaxPayFromDB : successful")
	// }

	return nil
}

func (impl *IncentiveImpl) LoadMaxExtraPayFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchMaxExtraPayout(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadMaxExtraPayFromDB : maxExtraPayout fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchMaxExtraPayout(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(errors.New(err.Details.(string))).Msg("LoadMaxExtraPayFromDB : maxExtraPayout fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchMaxExtraPayout(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(errors.New(err.Details.(string))).Msg("LoadMaxExtraPayFromDB : maxExtraPayout fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	return nil
}

func (impl *IncentiveImpl) LoadThresholdCheckFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchMaxThresholdCheck(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadThresholdCheckFromDB : threshold fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchMaxThresholdCheck(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(errors.New(err.Details.(string))).Msg("LoadThresholdCheckFromDB : threshold fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if err := impl.PsqlDb.FetchMaxThresholdCheck(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(errors.New(err.Details.(string))).Msg("LoadThresholdCheckFromDB : threshold fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	return nil
}

func (impl *IncentiveImpl) PrintAllInMemData() {

	for year, yearVal := range models.PlanRates.PlanRateMap {
		for quarter, val := range yearVal {
			logger.FileLogger.Info().Int("year", year).Int("quarter", quarter).Interface("planRate", val).Msg("LoadPlansFromDB: load planRate success")
		}
	}

	for year, val := range models.SalesTrgt.SaleTargetMap {
		logger.FileLogger.Info().Int("year", year).Interface("target", val).Msg("LoadTargetFromDB: load target success")
	}

	for year, yearVal := range models.ValOfPoint.VopMap {
		for quarter, val := range yearVal {
			logger.FileLogger.Info().Int("year", year).Int("quarter", quarter).Interface("VOP", val).Msg("LoadScaleOrderFromDB: load VOP success")
		}
	}

	for year, yearVal := range models.MaxPay.Payoutmap {
		for quarter, quarterVal := range yearVal {
			for position, val := range quarterVal {
				logger.FileLogger.Info().Int("year", year).Int("quarter", quarter).Str("position", position).Interface("maxPay", val).Msg("LoadMaxPayFromDB: load maxPay success")
			}
		}
	}

	for year, yearVal := range models.MaxExtraPay.ExtraPayoutmap {
		for quarter, quarterVal := range yearVal {
			for position, val := range quarterVal {
				logger.FileLogger.Info().Int("year", year).Int("quarter", quarter).Str("position", position).Interface("maxExtraPay", val).Msg("LoadMaxExtraPayFromDB: load maxExtraPay success")
			}
		}
	}

	for year, yearVal := range models.ThresholdCheck.Limit {
		for quarter, quarterVal := range yearVal {
			for position, val := range quarterVal {
				logger.FileLogger.Info().Int("year", year).Int("quarter", quarter).Str("position", position).Interface("threshold", val).Msg("LoadThresholdFromDB: load threshold success")
			}
		}
	}
}
