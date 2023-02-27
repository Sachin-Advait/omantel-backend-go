package adapters

import (
	"context"

	"github.com/jasonlvhit/gocron"
	"gitlab.com/rendezvousPoint/incentive_calculator/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

var theProvider models.IncentiveProvider

func InitIncentive(provider models.IncentiveProvider) {
	theProvider = provider
}

func InitIncentiveCron() {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.Fatal(context.TODO()).Stack().Err(err).Msg("Cron Init: error loading application config file, env - " + flags.Env())
	}

	if err := theProvider.LoadPlansFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadPlansFromDB")
	}

	if err := theProvider.LoadTargetFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadTargetFromDB")
	}

	if err := theProvider.LoadValueOfPointFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadValueOfPointFromDB")
	}

	if err := theProvider.LoadMaxPayFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadMaxPayFromDB")
	}

	if err := theProvider.LoadMaxExtraPayFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadMaxExtraPayFromDB")
	}

	if err := theProvider.LoadThresholdCheckFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadThresholdCheckFromDB")
	}

	theProvider.PrintAllInMemData()

	masterCronTime := applicationConfig.GetString(constant.MasterCronTime)

	if err := gocron.Every(1).Day().At(masterCronTime).Do(theProvider.LoadPlansFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadPlansFromDB: cron")
	}

	if err := gocron.Every(1).Day().At(masterCronTime).Do(theProvider.LoadTargetFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadTargetFromDB: cron")
	}

	if err := gocron.Every(1).Day().At(masterCronTime).Do(theProvider.LoadValueOfPointFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadValueOfPointFromDB: cron")
	}

	if err := gocron.Every(1).Day().At(masterCronTime).Do(theProvider.LoadMaxPayFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadMaxPayFromDB: cron")
	}

	if err := gocron.Every(1).Day().At(masterCronTime).Do(theProvider.LoadMaxExtraPayFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadMaxExtraPayFromDB: cron")
	}

	if err := gocron.Every(1).Day().At(masterCronTime).Do(theProvider.LoadThresholdCheckFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("LoadThresholdCheckFromDB: cron")
	}

	<-gocron.Start()
}
