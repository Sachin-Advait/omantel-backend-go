package business

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) LoadPlansFromDB(printLog bool) error {

	if err := impl.PsqlDb.FetchPlans(); err != nil {
		if printLog {
			logger.FileLogger.Error().Err(err).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return err
	}

	if printLog {
		logger.FileLogger.Info().Msg("LoadPlansFromDB : successful")
	}

	return nil
}

func (impl *IncentiveImpl) LoadTargetFromDB(printLog bool) error {

	if err := impl.PsqlDb.FetchSaleTarget(); err != nil {
		if printLog {
			logger.FileLogger.Error().Err(err).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return err
	}

	if printLog {
		logger.FileLogger.Info().Msg("LoadPlansFromDB : successful")
	}

	return nil
}

func (impl *IncentiveImpl) LoadScaleOrderFromDB(printLog bool) error {

	if err := impl.PsqlDb.FetchValueOfPoint(); err != nil {
		if printLog {
			logger.FileLogger.Error().Err(err).Msg("LoadScaleOrderFromDB : vop fetching failed")
		}
		return err
	}

	if printLog {
		logger.FileLogger.Info().Msg("LoadScaleOrderFromDB : successful")
	}

	return nil
}

func (impl *IncentiveImpl) LoadMaxPayFromDB(printLog bool) error {

	if err := impl.PsqlDb.FetchMaxPayout(); err != nil {
		if printLog {
			logger.FileLogger.Error().Err(err).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return err
	}

	if printLog {
		logger.FileLogger.Info().Msg("LoadMaxPayFromDB : successful")
	}

	return nil
}
