package scheduler

import (
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func PrintIntrinsicValues() {

	currTime := time.Now()
	currYear, _, _ := currTime.Date()

	year := currYear

	logger.FileLogger.Info().Int("year", year).Msg("Annual Data Log")
	logger.FileLogger.Info().Int("year", year).Interface("plans", models.PlanRates.PlanRateMap[year]).Msg("Load Plans")
	logger.FileLogger.Info().Int("year", year).Interface("target", models.SalesTrgt.SaleTargetMap[year]).Msg("Load Target")
	logger.FileLogger.Info().Int("year", year).Interface("vop", models.ValOfPoint.VopMap[year]).Msg("Load VOP")
	logger.FileLogger.Info().Int("year", year).Interface("maxPay", models.MaxPay.Payoutmap[year]).Msg("Load MaxPay")
	logger.FileLogger.Info().Int("year", year).Interface("maxExtraPay", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load MaxExtraPay")
	logger.FileLogger.Info().Int("year", year).Interface("threshold", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load Threshold")
	logger.FileLogger.Info().Int("year", year).Interface("vocCriteria", models.VOCDetail.Criteria[year]).Msg("Load VOCCriteria")
	logger.FileLogger.Info().Int("year", year).Interface("vocMaxPay", models.VOCDetail.MaxPay[year]).Msg("Load VOCMaxPay")

	year = year - 1

	logger.FileLogger.Info().Int("year", year).Msg("Annual Data Log")
	logger.FileLogger.Info().Int("year", year).Interface("plans", models.PlanRates.PlanRateMap[year]).Msg("Load Plans")
	logger.FileLogger.Info().Int("year", year).Interface("target", models.SalesTrgt.SaleTargetMap[year]).Msg("Load Target")
	logger.FileLogger.Info().Int("year", year).Interface("vop", models.ValOfPoint.VopMap[year]).Msg("Load VOP")
	logger.FileLogger.Info().Int("year", year).Interface("maxPay", models.MaxPay.Payoutmap[year]).Msg("Load MaxPay")
	logger.FileLogger.Info().Int("year", year).Interface("maxExtraPay", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load MaxExtraPay")
	logger.FileLogger.Info().Int("year", year).Interface("threshold", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load Threshold")
	logger.FileLogger.Info().Int("year", year).Interface("vocCriteria", models.VOCDetail.Criteria[year]).Msg("Load VOCCriteria")
	logger.FileLogger.Info().Int("year", year).Interface("vocMaxPay", models.VOCDetail.MaxPay[year]).Msg("Load VOCMaxPay")

	year = year + 2

	logger.FileLogger.Info().Int("year", year).Msg("Annual Data Log")
	logger.FileLogger.Info().Int("year", year).Interface("plans", models.PlanRates.PlanRateMap[year]).Msg("Load Plans")
	logger.FileLogger.Info().Int("year", year).Interface("target", models.SalesTrgt.SaleTargetMap[year]).Msg("Load Target")
	logger.FileLogger.Info().Int("year", year).Interface("vop", models.ValOfPoint.VopMap[year]).Msg("Load VOP")
	logger.FileLogger.Info().Int("year", year).Interface("maxPay", models.MaxPay.Payoutmap[year]).Msg("Load MaxPay")
	logger.FileLogger.Info().Int("year", year).Interface("maxExtraPay", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load MaxExtraPay")
	logger.FileLogger.Info().Int("year", year).Interface("threshold", models.MaxExtraPay.ExtraPayoutmap[year]).Msg("Load Threshold")
	logger.FileLogger.Info().Int("year", year).Interface("vocCriteria", models.VOCDetail.Criteria[year]).Msg("Load VOCCriteria")
	logger.FileLogger.Info().Int("year", year).Interface("vocMaxPay", models.VOCDetail.MaxPay[year]).Msg("Load VOCMaxPay")

}
