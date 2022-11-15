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
		logger.FileLogger.Fatal().Err(err).Msg("InitIncentiveCron: load planRate failed")
	}

	if err := theProvider.LoadTargetFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("InitIncentiveCron: load target failed")
	}

	if err := theProvider.LoadScaleOrderFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("InitIncentiveCron: load VOP failed")
	}

	if err := theProvider.LoadMaxPayFromDB(false); err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("InitIncentiveCron: load maxPay failed")
	}

	logger.FileLogger.Info().Interface("planRate", models.PlanRates.PlanRateMap).Msg("InitIncentiveCron: load planRate success")
	logger.FileLogger.Info().Interface("target", models.SalesTrgt.SaleTargetMap).Msg("InitIncentiveCron: load target success")
	logger.FileLogger.Info().Interface("VOP", models.ValOfPoint.VopMap).Msg("InitIncentiveCron: load VOP success")
	logger.FileLogger.Info().Interface("maxPay", models.MaxPay.Payoutmap).Msg("InitIncentiveCron: load maxPay success")

	////////

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadPlansFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("InitIncentiveCron: productPlans cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadTargetFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("InitIncentiveCron: saleTarget cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadScaleOrderFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("InitIncentiveCron: vop cron failed")
	}

	if err := gocron.Every(1).Day().At("01:00").Do(theProvider.LoadMaxPayFromDB, true); err != nil {
		logger.FileLogger.Error().Err(err).Msg("InitIncentiveCron: maxPay cron failed")
	}

	<-gocron.Start()
}
