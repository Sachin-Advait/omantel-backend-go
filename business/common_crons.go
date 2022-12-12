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
			logger.FileLogger.Error().Err(errors.New(err.Details.(string))).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if printLog {
		logger.FileLogger.Info().Msg("LoadPlansFromDB : successful")
	}

	return nil
}

func (impl *IncentiveImpl) LoadTargetFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if err := impl.PsqlDb.FetchSaleTarget(currQuarter, currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Err(errors.New(err.Details.(string))).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if printLog {
		logger.FileLogger.Info().Msg("LoadPlansFromDB : successful")
	}

	return nil
}

func (impl *IncentiveImpl) LoadScaleOrderFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if err := impl.PsqlDb.FetchValueOfPoint(currQuarter, currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Err(errors.New(err.Details.(string))).Msg("LoadScaleOrderFromDB : vop fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if printLog {
		logger.FileLogger.Info().Msg("LoadScaleOrderFromDB : successful")
	}

	return nil
}

func (impl *IncentiveImpl) LoadMaxPayFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, currMonth, _ := currTime.Date()
	currQuarter := impl.getQuarterFromMonth(int(currMonth))

	if err := impl.PsqlDb.FetchMaxPayout(currQuarter, currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Err(errors.New(err.Details.(string))).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return errors.New(err.Details.(string))
	}

	if printLog {
		logger.FileLogger.Info().Msg("LoadMaxPayFromDB : successful")
	}

	return nil
}
