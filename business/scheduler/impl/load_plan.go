package impl

import (
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *SchedulerImpl) LoadPlansFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.DbService.FetchPlans(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchPlans(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchPlans(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadPlansFromDB : productPlans fetching failed")
		}
		return err
	}

	return nil
}
