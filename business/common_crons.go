package business

import (
	"errors"
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) LoadPlansFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if err := impl.PsqlDb.FetchPlans(currQuarter, currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("currentQuarter", currQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	prevQuarter := impl.getQuarterFromMonth(int(currMonth)) - 1
	if prevQuarter > 0 {
		if err := impl.PsqlDb.FetchPlans(prevQuarter, currYear); err != nil {
			if printLog {
				logger.FileLogger.Error().Int("prevQuarter", prevQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadPlansFromDB : productPlans fetching failed")
			}
			return errors.New(err.Details.(string))
		}
	}

	if printLog {
		logger.FileLogger.Info().Int("prevQuarter", prevQuarter).Int("currentQuarter", currQuarter).Int("year", currYear).Msg("LoadPlansFromDB : successful")
	}

	return nil
}

func (impl *IncentiveImpl) LoadTargetFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if err := impl.PsqlDb.FetchSaleTarget(currQuarter, currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("currentQuarter", currQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	prevQuarter := impl.getQuarterFromMonth(int(currMonth)) - 1
	if prevQuarter > 0 {
		if err := impl.PsqlDb.FetchSaleTarget(prevQuarter, currYear); err != nil {
			if printLog {
				logger.FileLogger.Error().Int("prevQuarter", prevQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadTargetFromDB : saleTarget fetching failed")
			}
			return errors.New(err.Details.(string))
		}
	}

	if printLog {
		logger.FileLogger.Info().Int("prevQuarter", prevQuarter).Int("currentQuarter", currQuarter).Int("year", currYear).Msg("LoadTargetFromDB : successful")
	}

	return nil
}

func (impl *IncentiveImpl) LoadScaleOrderFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if err := impl.PsqlDb.FetchValueOfPoint(currQuarter, currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("currentQuarter", currQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadScaleOrderFromDB : vop fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	prevQuarter := impl.getQuarterFromMonth(int(currMonth)) - 1
	if prevQuarter > 0 {
		if err := impl.PsqlDb.FetchValueOfPoint(prevQuarter, currYear); err != nil {
			if printLog {
				logger.FileLogger.Error().Int("prevQuarter", prevQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadScaleOrderFromDB : vop fetching failed")
			}
			return errors.New(err.Details.(string))
		}
	}

	if printLog {
		logger.FileLogger.Info().Int("prevQuarter", prevQuarter).Int("currentQuarter", currQuarter).Int("year", currYear).Msg("LoadScaleOrderFromDB : successful")
	}

	return nil
}

func (impl *IncentiveImpl) LoadMaxPayFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if err := impl.PsqlDb.FetchMaxPayout(currQuarter, currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("currentQuarter", currQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	prevQuarter := impl.getQuarterFromMonth(int(currMonth)) - 1
	if prevQuarter > 0 {
		if err := impl.PsqlDb.FetchMaxPayout(prevQuarter, currYear); err != nil {
			if printLog {
				logger.FileLogger.Error().Int("prevQuarter", prevQuarter).Int("year", currYear).Err(errors.New(err.Details.(string))).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
			}
			return errors.New(err.Details.(string))
		}
	}

	if printLog {
		logger.FileLogger.Info().Int("prevQuarter", prevQuarter).Int("currentQuarter", currQuarter).Int("year", currYear).Msg("LoadMaxPayFromDB : successful")
	}

	return nil
}
