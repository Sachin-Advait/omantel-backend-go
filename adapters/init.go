package adapters

import (
	"github.com/jasonlvhit/gocron"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

var theProvider models.IncentiveProvider

func InitIncentive(provider models.IncentiveProvider) {
	theProvider = provider
}

func InitIncentiveCron() {

	if err := theProvider.LoadPlansFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadPlansFromDB: load planRate failed")
	}

	if err := theProvider.LoadTargetFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadTargetFromDB: load target failed")
	}

	if err := theProvider.LoadValueOfPointFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadValueOfPointFromDB: load VOP failed")
	}

	if err := theProvider.LoadMaxPayFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadMaxPayFromDB: load maxPay failed")
	}

	if err := theProvider.LoadMaxExtraPayFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadMaxExtraPayFromDB: load maxExtraPay failed")
	}

	if err := theProvider.LoadThresholdCheckFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadThresholdCheckFromDB: load threshold failed")
	}

	theProvider.PrintAllInMemData()

	////////

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadPlansFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadPlansFromDB: productPlans cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadTargetFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadTargetFromDB: saleTarget cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadValueOfPointFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadScaleOrderFromDB: vop cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadMaxPayFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadMaxPayFromDB: maxPay cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadMaxExtraPayFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadMaxExtraPayFromDB: maxExtraPay cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadThresholdCheckFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadThresholdCheckFromDB: threshold cron failed")
	}

	<-gocron.Start()
}
