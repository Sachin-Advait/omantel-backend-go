package business

import (
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *IncentiveImpl) LoadPlansFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchPlans(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchPlans(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchPlans(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return err
	}

	return nil
}

func (impl *IncentiveImpl) LoadTargetFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchSaleTarget(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchSaleTarget(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchSaleTarget(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return err
	}

	return nil
}

func (impl *IncentiveImpl) LoadValueOfPointFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchValueOfPoint(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadValueOfPointFromDB : vop fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchValueOfPoint(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadValueOfPointFromDB : vop fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchValueOfPoint(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadValueOfPointFromDB : vop fetching failed")
		}
		return err
	}

	return nil
}

func (impl *IncentiveImpl) LoadMaxPayFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchMaxPayout(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchMaxPayout(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchMaxPayout(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return err
	}

	return nil
}

func (impl *IncentiveImpl) LoadMaxExtraPayFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchMaxExtraPayout(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadMaxExtraPayFromDB : maxExtraPayout fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchMaxExtraPayout(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadMaxExtraPayFromDB : maxExtraPayout fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchMaxExtraPayout(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadMaxExtraPayFromDB : maxExtraPayout fetching failed")
		}
		return err
	}

	return nil
}

func (impl *IncentiveImpl) LoadThresholdCheckFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.PsqlDb.FetchMaxThreshold(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadThresholdCheckFromDB : threshold fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchMaxThreshold(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadThresholdCheckFromDB : threshold fetching failed")
		}
		return err
	}

	if err := impl.PsqlDb.FetchMaxThreshold(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadThresholdCheckFromDB : threshold fetching failed")
		}
		return err
	}

	return nil
}

func (impl *IncentiveImpl) PrintAllInMemData() {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	year := currYear

	logger.FileLogger.Info().Int("year", year).Msg("Annual Data Log")
	logger.FileLogger.Info().Int("year", year).Interface("plans", models.PlanRates.PlanRateMap[year]).Msg("Load Plans")
	logger.FileLogger.Info().Int("year", year).Interface("target", models.SalesTrgt.SaleTargetMap[year]).Msg("Load Target")
	logger.FileLogger.Info().Int("year", year).Interface("vop", models.ValOfPoint.VopMap[year]).Msg("Load VOP")
	logger.FileLogger.Info().Int("year", year).Interface("maxPay", models.MaxPay.Payoutmap[year]).Msg("Load MaxPay")
	logger.FileLogger.Info().Int("year", year).Interface("maxExtraPay", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load MaxExtraPay")
	logger.FileLogger.Info().Int("year", year).Interface("threshold", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load Threshold")

	year = year - 1

	logger.FileLogger.Info().Int("year", year).Msg("Annual Data Log")
	logger.FileLogger.Info().Int("year", year).Interface("plans", models.PlanRates.PlanRateMap[year]).Msg("Load Plans")
	logger.FileLogger.Info().Int("year", year).Interface("target", models.SalesTrgt.SaleTargetMap[year]).Msg("Load Target")
	logger.FileLogger.Info().Int("year", year).Interface("vop", models.ValOfPoint.VopMap[year]).Msg("Load VOP")
	logger.FileLogger.Info().Int("year", year).Interface("maxPay", models.MaxPay.Payoutmap[year]).Msg("Load MaxPay")
	logger.FileLogger.Info().Int("year", year).Interface("maxExtraPay", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load MaxExtraPay")
	logger.FileLogger.Info().Int("year", year).Interface("threshold", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load Threshold")

	year = year + 2

	logger.FileLogger.Info().Int("year", year).Msg("Annual Data Log")
	logger.FileLogger.Info().Int("year", year).Interface("plans", models.PlanRates.PlanRateMap[year]).Msg("Load Plans")
	logger.FileLogger.Info().Int("year", year).Interface("target", models.SalesTrgt.SaleTargetMap[year]).Msg("Load Target")
	logger.FileLogger.Info().Int("year", year).Interface("vop", models.ValOfPoint.VopMap[year]).Msg("Load VOP")
	logger.FileLogger.Info().Int("year", year).Interface("maxPay", models.MaxPay.Payoutmap[year]).Msg("Load MaxPay")
	logger.FileLogger.Info().Int("year", year).Interface("maxExtraPay", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load MaxExtraPay")
	logger.FileLogger.Info().Int("year", year).Interface("threshold", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load Threshold")
}
