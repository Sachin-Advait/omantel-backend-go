package impl

import (
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *SchedulerImpl) LoadVOCFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.DbService.FetchVOCCriteria(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadVOCCriteriaFromDB : criteria fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchVOCCriteria(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadVOCCriteriaFromDB : criteria fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchVOCCriteria(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadVOCCriteriaFromDB : criteria fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchVOCMaxPay(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadVOCCriteriaFromDB : maxPay fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchVOCMaxPay(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadVOCCriteriaFromDB : maxPay fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchVOCMaxPay(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadVOCCriteriaFromDB : maxPay fetching failed")
		}
		return err
	}

	return nil
}
