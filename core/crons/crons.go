package crons

import (
	"github.com/jasonlvhit/gocron"
)

type CronServices interface {
	HandleDataInputFile()
	HandlePlansInputFile()
	HandleTargetInputFile()
	HandleUserMasterInputFile()
	HandleScaleOrderInputFile()
	HandleMaxPayInputFile()
}

type CronContext struct {
	CronServices CronServices
}

func NewCronContext(cronServices CronServices) *CronContext {
	return &CronContext{CronServices: cronServices}
}

func (ctx *CronContext) StartCrons() {

	// if err := gocron.Every(1).Day().At("12:00").Do(ctx.CronServices.HandleDataInputFile); err != nil {
	// 	utils.LogObj.Error("failed to initate upload pledge file task, err: " + err.Error())
	// }

	// if err := gocron.Every(1).Day().At("12:00").Do(ctx.CronServices.HandlePlansInputFile); err != nil {
	// 	utils.LogObj.Error("failed to initate upload pledge file task, err: " + err.Error())
	// }

	go ctx.CronServices.HandleDataInputFile()
	go ctx.CronServices.HandleUserMasterInputFile()
	go ctx.CronServices.HandlePlansInputFile()
	go ctx.CronServices.HandleTargetInputFile()
	go ctx.CronServices.HandleScaleOrderInputFile()
	go ctx.CronServices.HandleMaxPayInputFile()

	<-gocron.Start()
}
