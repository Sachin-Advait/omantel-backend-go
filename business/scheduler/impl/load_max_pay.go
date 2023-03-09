package impl

import (
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *SchedulerImpl) LoadMaxPayFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.DbService.FetchMaxPayout(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchMaxPayout(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchMaxPayout(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadMaxPayFromDB : maxPayout fetching failed")
		}
		return err
	}

	return nil
}
