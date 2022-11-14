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

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadPlansFromDB); err != nil {
		logger.FileLogger.Error().Err(err).Msg("InitIncentiveCron: productPlans cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadTargetFromDB); err != nil {
		logger.FileLogger.Error().Err(err).Msg("InitIncentiveCron: saleTarget cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadScaleOrderFromDB); err != nil {
		logger.FileLogger.Error().Err(err).Msg("InitIncentiveCron: vop cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadMaxPayFromDB); err != nil {
		logger.FileLogger.Error().Err(err).Msg("InitIncentiveCron: maxPay cron failed")
	}

	<-gocron.Start()
}
