package business

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *IncentiveImpl) GetUserReport(ctx *gin.Context, userId string, isQuarterReport bool, quarter, month, year int) (*models.UserRecord, *errors.Error) {

	userReport, err := impl.PsqlDb.FetchUserReport(userId, isQuarterReport, quarter, month, year)
	if err == nil && userReport != nil {

		currTime := time.Now()
		currYear, currMonth, currDay := currTime.Date()

		if userReport.UserId != "" && !userReport.CreateAt.Before(time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, currTime.Location())) { //entry exists and created time is today
			return &models.UserRecord{
				EmpId:                        userReport.UserId,
				ReportingManager:             userReport.ManagerId,
				Isquarterly:                  userReport.Isquarterly,
				Month:                        userReport.Month,
				Quarter:                      userReport.Quarter,
				Year:                         userReport.Year,
				Hayyak_sales_count:           userReport.Hayyak_sales_count,
				Hayyak_target:                userReport.Hayyak_target,
				Hayyak_payout:                userReport.Hayyak_payout,
				Hayyak_extra_payout:          userReport.Hayyak_extra_payout,
				Hayyak_achieved_percentage:   userReport.Hayyak_achieved_percentage,
				Hbb_sales_count:              userReport.Hbb_sales_count,
				Hbb_target:                   userReport.Hbb_target,
				Hbb_payout:                   userReport.Hbb_payout,
				Hbb_extra_payout:             userReport.Hbb_extra_payout,
				Hbb_achieved_percentage:      userReport.Hbb_achieved_percentage,
				Baqati_sales_count:           userReport.Baqati_sales_count,
				Baqati_target:                userReport.Baqati_target,
				Baqati_payout:                userReport.Baqati_payout,
				Baqati_extra_payout:          userReport.Baqati_extra_payout,
				Baqati_achieved_percentage:   userReport.Baqati_achieved_percentage,
				Upgrade_sales_count:          userReport.Upgrade_sales_count,
				Upgrade_target:               userReport.Upgrade_target,
				Upgrade_payout:               userReport.Upgrade_payout,
				Upgrade_extra_payout:         userReport.Upgrade_extra_payout,
				Uprade_achieved_percentage:   userReport.Uprade_achieved_percentage,
				Total_payout:                 userReport.Total_payout,
				HBBSaleCountPerPlanArray:     userReport.HBBSaleCountPerPlanArray,
				HayyakSaleCountPerPlanArray:  userReport.HayyakSaleCountPerPlanArray,
				BaqatiSaleCountPerPlanArray:  userReport.BaqatiSaleCountPerPlanArray,
				UpgradeSaleCountPerPlanArray: userReport.UpgradeSaleCountPerPlanArray,
			}, nil
		}
	}

	_, userReport, err = impl.updateUserReport(ctx, userId, isQuarterReport, quarter, month, year)
	if err != nil {
		return nil, err
	}

	return &models.UserRecord{
		EmpId:                        userReport.UserId,
		ReportingManager:             userReport.ManagerId,
		Isquarterly:                  userReport.Isquarterly,
		Month:                        userReport.Month,
		Quarter:                      userReport.Quarter,
		Year:                         userReport.Year,
		Hayyak_sales_count:           userReport.Hayyak_sales_count,
		Hayyak_target:                userReport.Hayyak_target,
		Hayyak_payout:                userReport.Hayyak_payout,
		Hayyak_extra_payout:          userReport.Hayyak_extra_payout,
		Hayyak_achieved_percentage:   userReport.Hayyak_achieved_percentage,
		Hbb_sales_count:              userReport.Hbb_sales_count,
		Hbb_target:                   userReport.Hbb_target,
		Hbb_payout:                   userReport.Hbb_payout,
		Hbb_extra_payout:             userReport.Hbb_extra_payout,
		Hbb_achieved_percentage:      userReport.Hbb_achieved_percentage,
		Baqati_sales_count:           userReport.Baqati_sales_count,
		Baqati_target:                userReport.Baqati_target,
		Baqati_payout:                userReport.Baqati_payout,
		Baqati_extra_payout:          userReport.Baqati_extra_payout,
		Baqati_achieved_percentage:   userReport.Baqati_achieved_percentage,
		Upgrade_sales_count:          userReport.Upgrade_sales_count,
		Upgrade_target:               userReport.Upgrade_target,
		Upgrade_payout:               userReport.Upgrade_payout,
		Upgrade_extra_payout:         userReport.Upgrade_extra_payout,
		Uprade_achieved_percentage:   userReport.Uprade_achieved_percentage,
		Total_payout:                 userReport.Total_payout,
		HBBSaleCountPerPlanArray:     userReport.HBBSaleCountPerPlanArray,
		HayyakSaleCountPerPlanArray:  userReport.HayyakSaleCountPerPlanArray,
		BaqatiSaleCountPerPlanArray:  userReport.BaqatiSaleCountPerPlanArray,
		UpgradeSaleCountPerPlanArray: userReport.UpgradeSaleCountPerPlanArray,
	}, nil
}
