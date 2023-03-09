package scheduler

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func PreLoadIntrinsicValues() {

	if err := scheduleProvide.LoadPlansFromDB(false); err != nil {
		logger.FileLogger.Error().Err(err).Msg("PreLoadIntrinsicValues : plan provider failed")
	}

	if err := scheduleProvide.LoadTargetFromDB(false); err != nil {
		logger.FileLogger.Error().Err(err).Msg("PreLoadIntrinsicValues : target provider failed")
	}

	if err := scheduleProvide.LoadMaxPayFromDB(false); err != nil {
		logger.FileLogger.Error().Err(err).Msg("PreLoadIntrinsicValues : maxPay provider failed")
	}

	if err := scheduleProvide.LoadMaxExtraPayFromDB(false); err != nil {
		logger.FileLogger.Error().Err(err).Msg("PreLoadIntrinsicValues : maxExtraPay provider failed")
	}

	if err := scheduleProvide.LoadThresholdCheckFromDB(false); err != nil {
		logger.FileLogger.Error().Err(err).Msg("PreLoadIntrinsicValues : threshold provider failed")
	}

	if err := scheduleProvide.LoadValueOfPointFromDB(false); err != nil {
		logger.FileLogger.Error().Err(err).Msg("PreLoadIntrinsicValues : vop provider failed")
	}

	if err := scheduleProvide.LoadVOCFromDB(false); err != nil {
		logger.FileLogger.Error().Err(err).Msg("PreLoadIntrinsicValues : voc provider failed")
	}
}
