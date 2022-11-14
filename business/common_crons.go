package business

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) LoadPlansFromDB() {

	if err := impl.PsqlDb.FetchPlans(); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadPlansFromDB : productPlans fetching failed")
		return
	}

	logger.FileLogger.Info().Msg("LoadPlansFromDB : successful")
}

func (impl *IncentiveImpl) LoadTargetFromDB() {

	if err := impl.PsqlDb.FetchSaleTarget(); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadTargetFromDB : saleTarget fetching failed")
		return
	}

	logger.FileLogger.Info().Msg("LoadPlansFromDB : successful")
}

func (impl *IncentiveImpl) LoadScaleOrderFromDB() {

	if err := impl.PsqlDb.FetchValueOfPoint(); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadScaleOrderFromDB : vop fetching failed")
		return
	}

	logger.FileLogger.Info().Msg("LoadScaleOrderFromDB : successful")
}

func (impl *IncentiveImpl) LoadMaxPayFromDB() {

	if err := impl.PsqlDb.FetchMaxPayout(); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		return
	}

	logger.FileLogger.Info().Msg("LoadMaxPayFromDB : successful")
}
