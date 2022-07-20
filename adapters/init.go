package adapters

import (
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
)

var theProvider models.IncentiveProvider

func InitIncentive(provider models.IncentiveProvider) {
	theProvider = provider
}

func InitIncentiveCron() {

	//if err := gocron.Every(1).Day().At("12:00").Do(theProvider.HandleDataInputFile); err != nil {
	//	utils.LogObj.Error("failed to initate upload pledge file task, err: " + err.Error())
	//}
	//
	//if err := gocron.Every(1).Day().At("12:00").Do(theProvider.HandlePlansInputFile); err != nil {
	//	utils.LogObj.Error("failed to initate upload pledge file task, err: " + err.Error())
	//}

	//go theProvider.HandleDataInputFile()
	//go theProvider.HandleUserMasterInputFile()
	go theProvider.HandlePlansInputFile()
	go theProvider.HandleTargetInputFile()
	go theProvider.HandleScaleOrderInputFile()
	go theProvider.HandleMaxPayInputFile()
	//go theProvider.HandleHierarchyFile()

	//<-gocron.Start()
}
