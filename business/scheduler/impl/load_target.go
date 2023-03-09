package impl

import (
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func (impl *SchedulerImpl) LoadTargetFromDB(printLog bool) error {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	if err := impl.DbService.FetchSaleTarget(currYear); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear).Err(err).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchSaleTarget(currYear - 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear-1).Err(err).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return err
	}

	if err := impl.DbService.FetchSaleTarget(currYear + 1); err != nil {
		if printLog {
			logger.FileLogger.Error().Int("year", currYear+1).Err(err).Msg("LoadTargetFromDB : saleTarget fetching failed")
		}
		return err
	}

	return nil
}
