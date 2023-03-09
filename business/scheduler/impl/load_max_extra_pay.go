package impl

import (
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *SchedulerImpl) LoadMaxExtraPayFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.DbService.FetchMaxExtraPayout(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadMaxExtraPayFromDB : maxExtraPayout fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchMaxExtraPayout(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadMaxExtraPayFromDB : maxExtraPayout fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchMaxExtraPayout(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadMaxExtraPayFromDB : maxExtraPayout fetching failed")
		}
		return err
	}

	return nil
}
